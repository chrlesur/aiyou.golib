// File: pkg/aiyou/models_test.go
package aiyou

import (
 "context"
 "encoding/json"
 "net/http"
 "net/http/httptest"
 "testing"
 "time"
)

func TestCreateModel(t *testing.T) {
 server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
 if r.URL.Path == "/api/login" {
 w.WriteHeader(http.StatusOK)
 json.NewEncoder(w).Encode(LoginResponse{
 Token: "test_token",
 ExpiresAt: time.Now().Add(time.Hour),
 })
 return
 }

 if r.URL.Path == "/api/v1/models" && r.Method == "POST" {
 var req ModelRequest
 if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
 t.Errorf("Failed to decode request body: %v", err)
 w.WriteHeader(http.StatusBadRequest)
 return
 }

 response := ModelResponse{
 Model: Model{
 ID: "model_123",
 Name: req.Name,
 Description: req.Description,
 Version: "1.0",
  CreatedAt: time.Now(),
 UpdatedAt: time.Now(),
 Properties: req.Properties,
 },
 }
 w.WriteHeader(http.StatusOK)
 json.NewEncoder(w).Encode(response)
 }
 }))
 defer server.Close()

 client, err := NewClient("test@example.com", "password", WithBaseURL(server.URL))
 if err != nil {
 t.Fatalf("Failed to create client: %v", err)
 }

 req := ModelRequest{
 Name: "Test Model",
 Description: "A test model",
 Properties: ModelProperties{
 MaxTokens: 4096,
 Temperature: 0.7,
 Provider: "OpenAI",
 Capabilities: []string{"chat", "completion"},
 },
 }

 resp, err := client.CreateModel(context.Background(), req)
 if err != nil {
 t.Fatalf("CreateModel failed: %v", err)
 }

 if resp.Model.Name != req.Name {
 t.Errorf("Expected model name %s, got %s", req.Name, resp.Model.Name)
 }
}

func TestGetModels(t *testing.T) {
 // Test similaire à CreateModel pour la méthode GetModels
 // À implémenter
}