/*
Copyright (C) 2024 Cloud Temple

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/glamour"
	"github.com/chrlesur/aiyou.golib"
	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile      string
	email        string
	password     string
	assistantID  string
	baseURL      string
	debug        bool
	markdownMode bool
	quietMode    bool

	// Couleurs pour l'interface
	assistantColor = color.New(color.FgCyan, color.Bold)
	userColor      = color.New(color.FgGreen)
	errorColor     = color.New(color.FgRed)
	infoColor      = color.New(color.FgYellow)
)

// Configuration structure
type Config struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	AssistantID string `json:"assistant_id"`
	BaseURL     string `json:"base_url"`
	Debug       bool   `json:"debug"`
}

// chatLogger gère les logs pendant la conversation
type chatLogger struct {
	*log.Logger
	enabled bool
}

func newChatLogger() *chatLogger {
	return &chatLogger{
		Logger:  log.New(os.Stderr, "", 0),
		enabled: false,
	}
}

func (l *chatLogger) Printf(format string, v ...interface{}) {
	if l.enabled {
		l.Logger.Printf(format, v...)
	}
}

var rootCmd = &cobra.Command{
	Use:   "aiyou-cli",
	Short: "Un client CLI pour interagir avec AI.YOU",
	Long: `Un client en ligne de commande complet pour interagir avec AI.YOU.
Il permet de discuter avec des assistants, gérer les conversations,
et sauvegarder la configuration.`,
}

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Démarrer une conversation avec un assistant",
	Run:   startChat,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lister les assistants disponibles",
	Run:   listAssistants,
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Gérer la configuration",
	Run:   manageConfig,
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "fichier de configuration (par défaut $HOME/.aiyou.yaml)")
	rootCmd.PersistentFlags().StringVar(&email, "email", "", "Email pour l'authentification")
	rootCmd.PersistentFlags().StringVar(&password, "password", "", "Mot de passe pour l'authentification")
	rootCmd.PersistentFlags().StringVar(&assistantID, "assistant", "", "ID de l'assistant à utiliser")
	rootCmd.PersistentFlags().StringVar(&baseURL, "url", "https://ai.dragonflygroup.fr", "URL de base de l'API")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Active les logs de debug")
	rootCmd.PersistentFlags().BoolVar(&markdownMode, "markdown", true, "Active le rendu markdown")
	rootCmd.PersistentFlags().BoolVar(&quietMode, "quiet", false, "Désactive les messages de statut")

	rootCmd.AddCommand(chatCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(configCmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			errorColor.Fprintf(os.Stderr, "Erreur lors de la recherche du répertoire utilisateur: %v\n", err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".aiyou")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil && !quietMode {
		infoColor.Fprintf(os.Stderr, "Utilisation du fichier de configuration: %s\n", viper.ConfigFileUsed())
	}
}

func showSpinner(done chan bool) {
	spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	i := 0
	for {
		select {
		case <-done:
			fmt.Fprintf(os.Stderr, "\r\033[K")
			return
		default:
			if !quietMode {
				fmt.Fprintf(os.Stderr, "\r%s Réflexion en cours...", spinner[i])
			}
			i = (i + 1) % len(spinner)
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func formatCodeBlocks(markdown string) string {
	codeBlockRegex := regexp.MustCompile("```([a-zA-Z]*)\n([^`]+)```")
	return codeBlockRegex.ReplaceAllStringFunc(markdown, func(block string) string {
		matches := codeBlockRegex.FindStringSubmatch(block)
		if len(matches) >= 3 {
			language := matches[1]
			code := matches[2]
			return fmt.Sprintf("```%s\n%s```", language, strings.TrimSpace(code))
		}
		return block
	})
}

func createClient() (*aiyou.Client, error) {
	clientEmail := email
	if clientEmail == "" {
		clientEmail = viper.GetString("email")
	}
	clientPassword := password
	if clientPassword == "" {
		clientPassword = viper.GetString("password")
	}
	clientBaseURL := baseURL
	if clientBaseURL == "" {
		clientBaseURL = viper.GetString("base_url")
	}

	if clientEmail == "" || clientPassword == "" {
		return nil, fmt.Errorf("email et mot de passe requis")
	}

	logger := aiyou.NewDefaultLogger(os.Stderr)
	if debug {
		logger.SetLevel(aiyou.DEBUG)
	} else if quietMode {
		logger.SetLevel(aiyou.ERROR) // En mode quiet, on n'affiche que les erreurs
	} else {
		logger.SetLevel(aiyou.INFO)
	}

	return aiyou.NewClient(
		clientEmail,
		clientPassword,
		aiyou.WithLogger(logger),
		aiyou.WithBaseURL(clientBaseURL),
	)
}

func listAssistants(cmd *cobra.Command, args []string) {
	client, err := createClient()
	if err != nil {
		errorColor.Fprintf(os.Stderr, "Erreur lors de la création du client: %v\n", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	response, err := client.GetUserAssistants(ctx)
	if err != nil {
		errorColor.Fprintf(os.Stderr, "Erreur lors de la récupération des assistants: %v\n", err)
		return
	}

	if !quietMode {
		infoColor.Fprintf(os.Stderr, "\nAssistants disponibles (%d):\n\n", response.TotalItems)
	}
	for _, assistant := range response.Members {
		fmt.Fprintf(os.Stderr, "=== Assistant ===\n")
		fmt.Fprintf(os.Stderr, "ID: %s\n", assistant.AssistantID)
		fmt.Fprintf(os.Stderr, "Nom: %s\n", assistant.Name)
		if assistant.Model != "" {
			fmt.Fprintf(os.Stderr, "Modèle: %s\n", assistant.Model)
		}
		fmt.Fprintln(os.Stderr, strings.Repeat("-", 50))
	}
}

func startChat(cmd *cobra.Command, args []string) {
	client, err := createClient()
	if err != nil {
		errorColor.Fprintf(os.Stderr, "Erreur lors de la création du client: %v\n", err)
		return
	}

	currentAssistantID := assistantID
	if currentAssistantID == "" {
		currentAssistantID = viper.GetString("assistant_id")
	}
	if currentAssistantID == "" {
		errorColor.Fprintln(os.Stderr, "ID d'assistant requis. Utilisez --assistant ou listez les assistants disponibles avec la commande 'list'")
		return
	}

	// Configuration du rendu markdown
	markdownRenderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(100),
	)
	if err != nil {
		errorColor.Fprintf(os.Stderr, "Erreur lors de l'initialisation du rendu markdown: %v\n", err)
		return
	}

	// Configuration de readline
	rl, err := readline.NewEx(&readline.Config{
		Prompt:            "\033[32m➜\033[0m ",
		HistoryFile:       filepath.Join(os.TempDir(), ".aiyou_history"),
		AutoComplete:      readline.NewPrefixCompleter(),
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit",
		HistorySearchFold: true,
		FuncFilterInputRune: func(r rune) (rune, bool) {
			if r == readline.CharCtrlZ {
				return r, false
			}
			return r, true
		},
		FuncOnWidthChanged: func(f func()) {},
	})
	if err != nil {
		errorColor.Fprintf(os.Stderr, "Erreur lors de l'initialisation de readline: %v\n", err)
		return
	}
	defer rl.Close()

	if !quietMode {
		fmt.Fprintf(os.Stderr, "\nChat démarré avec l'assistant %s\n", currentAssistantID)
		fmt.Fprintln(os.Stderr, "Tapez /help pour voir les commandes disponibles")
		fmt.Fprintln(os.Stderr, strings.Repeat("-", 50))
	}

	var conversation []aiyou.Message
	for {
		line, err := rl.Readline()
		if err == readline.ErrInterrupt {
			continue
		} else if err == io.EOF || strings.ToLower(line) == "exit" || strings.ToLower(line) == "/exit" {
			fmt.Fprintln(os.Stderr, "\nAu revoir!")
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "/") {
			cmd := strings.TrimSpace(strings.ToLower(line))
			switch cmd {
			case "/help":
				fmt.Fprintln(os.Stderr, "\nCommandes disponibles:")
				fmt.Fprintln(os.Stderr, "/clear - Réinitialiser la conversation")
				fmt.Fprintln(os.Stderr, "/history - Afficher l'historique")
				fmt.Fprintln(os.Stderr, "/save - Sauvegarder la conversation")
				fmt.Fprintln(os.Stderr, "/exit - Quitter")
				continue
			case "/clear":
				conversation = nil
				fmt.Fprintln(os.Stderr, "Conversation réinitialisée.")
				continue
			case "/history":
				fmt.Fprintln(os.Stderr, "\nHistorique de la conversation:")
				for _, msg := range conversation {
					if msg.Role == "user" {
						userColor.Fprintf(os.Stderr, "Vous: %s\n", msg.Content[0].Text)
					} else {
						assistantColor.Fprintf(os.Stderr, "Assistant: %s\n", msg.Content[0].Text)
					}
				}
				continue
			case "/save":
				filename := fmt.Sprintf("conversation_%s.json", time.Now().Format("20060102_150405"))
				data, err := json.MarshalIndent(conversation, "", " ")
				if err != nil {
					errorColor.Fprintf(os.Stderr, "Erreur lors de la sauvegarde: %v\n", err)
				} else if err := os.WriteFile(filename, data, 0644); err != nil {
					errorColor.Fprintf(os.Stderr, "Erreur lors de l'écriture: %v\n", err)
				} else {
					infoColor.Fprintf(os.Stderr, "Conversation sauvegardée dans %s\n", filename)
				}
				continue
			}
		}

		userColor.Fprintf(os.Stderr, "\nVous: %s\n", line)

		userMsg := aiyou.NewTextMessage("user", line)
		conversation = append(conversation, userMsg)

		spinnerDone := make(chan bool)
		go showSpinner(spinnerDone)

		ctx := context.Background()
		req := aiyou.ChatCompletionRequest{
			Messages:    conversation,
			AssistantID: currentAssistantID,
			Stream:      true,
		}

		stream, err := client.ChatCompletionStream(ctx, req)
		if err != nil {
			close(spinnerDone)
			errorColor.Fprintf(os.Stderr, "Erreur: %v\n", err)
			continue
		}

		close(spinnerDone)

		fmt.Fprintf(os.Stderr, "\n")
		assistantColor.Fprintf(os.Stderr, "Assistant: ")

		var fullResponse strings.Builder
		for {
			chunk, err := stream.ReadChunk()
			if err == io.EOF {
				break
			}
			if err != nil {
				errorColor.Fprintf(os.Stderr, "\nErreur lors de la lecture de la réponse: %v\n", err)
				break
			}

			if chunk != nil && len(chunk.Choices) > 0 {
				choice := chunk.Choices[0]
				if choice.Delta != nil && choice.Delta.Content != "" {
					if markdownMode {
						// Pour le markdown, on accumule d'abord la réponse
						fullResponse.WriteString(choice.Delta.Content)
					} else {
						// Sans markdown, on affiche directement en streaming
						fmt.Fprintf(os.Stderr, "%s", choice.Delta.Content)
					}
				}
			}
		}

		if markdownMode {
			// Rendu markdown une seule fois à la fin
			formattedResponse := formatCodeBlocks(fullResponse.String())
			rendered, err := markdownRenderer.Render(formattedResponse)
			if err != nil {
				errorColor.Fprintf(os.Stderr, "Erreur lors du rendu markdown: %v\n", err)
				fmt.Fprintln(os.Stderr, fullResponse.String())
			} else {
				fmt.Fprint(os.Stderr, rendered)
			}
		} else {
			// Déjà affiché en streaming
			fmt.Fprintf(os.Stderr, "\n")
		}

		// Ajouter la réponse à la conversation
		assistantMsg := aiyou.NewTextMessage("assistant", fullResponse.String())
		conversation = append(conversation, assistantMsg)
	}
}

func manageConfig(cmd *cobra.Command, args []string) {
	config := Config{
		Email:       email,
		Password:    password,
		AssistantID: assistantID,
		BaseURL:     baseURL,
		Debug:       debug,
	}

	home, err := os.UserHomeDir()
	if err != nil {
		errorColor.Fprintf(os.Stderr, "Erreur lors de la recherche du répertoire utilisateur: %v\n", err)
		return
	}

	configFile := filepath.Join(home, ".aiyou.json")
	data, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		errorColor.Fprintf(os.Stderr, "Erreur lors de la sérialisation de la configuration: %v\n", err)
		return
	}

	if err := os.WriteFile(configFile, data, 0600); err != nil {
		errorColor.Fprintf(os.Stderr, "Erreur lors de l'écriture de la configuration: %v\n", err)
		return
	}

	if !quietMode {
		infoColor.Fprintf(os.Stderr, "Configuration sauvegardée dans %s\n", configFile)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		errorColor.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
