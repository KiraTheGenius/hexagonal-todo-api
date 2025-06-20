package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

type TodoResponse struct {
	ID          string  `json:"id"`
	Description string  `json:"description"`
	DueDate     string  `json:"dueDate"`
	FileID      *string `json:"fileId"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
}

type ListTodosResponse struct {
	Todos      []TodoResponse `json:"todos"`
	Pagination interface{}    `json:"pagination"`
}

var baseURL = "http://localhost:8080"

func TestIntegration_TodoCRUD(t *testing.T) {
	if os.Getenv("SKIP_INTEGRATION") == "1" {
		t.Skip("Skipping integration test")
	}

	// 1. Health check
	r, err := http.Get(baseURL + "/health")
	if err != nil {
		t.Fatalf("health check failed: %v", err)
	}
	if r.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", r.StatusCode)
	}
	r.Body.Close()

	// 2. Create todo
	due := time.Now().Add(48 * time.Hour).UTC().Format(time.RFC3339)
	desc := "integration test todo " + time.Now().Format("150405.000")
	createBody := map[string]interface{}{
		"description": desc,
		"dueDate":     due,
	}
	bodyBytes, _ := json.Marshal(createBody)
	r, err = http.Post(baseURL+"/todo", "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatalf("create todo failed: %v", err)
	}
	if r.StatusCode != 201 {
		b, _ := io.ReadAll(r.Body)
		t.Fatalf("expected 201, got %d, body: %s", r.StatusCode, string(b))
	}
	var created TodoResponse
	json.NewDecoder(r.Body).Decode(&created)
	r.Body.Close()
	if created.Description != desc {
		t.Fatalf("expected description %q, got %q", desc, created.Description)
	}

	// 3. List todos
	r, err = http.Get(baseURL + "/todo?limit=5")
	if err != nil {
		t.Fatalf("list todos failed: %v", err)
	}
	if r.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", r.StatusCode)
	}
	var list ListTodosResponse
	json.NewDecoder(r.Body).Decode(&list)
	r.Body.Close()
	found := false
	for _, todo := range list.Todos {
		if todo.ID == created.ID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("created todo not found in list")
	}

	// 4. Update todo
	newDesc := desc + " updated"
	updateBody := map[string]interface{}{
		"description": newDesc,
	}
	bodyBytes, _ = json.Marshal(updateBody)
	req, _ := http.NewRequest(http.MethodPut, baseURL+"/todo/"+created.ID, bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	r, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("update todo failed: %v", err)
	}
	if r.StatusCode != 200 {
		b, _ := io.ReadAll(r.Body)
		t.Fatalf("expected 200, got %d, body: %s", r.StatusCode, string(b))
	}
	var updated TodoResponse
	json.NewDecoder(r.Body).Decode(&updated)
	r.Body.Close()
	if updated.Description != newDesc {
		t.Fatalf("expected updated description %q, got %q", newDesc, updated.Description)
	}

	// 5. Delete todo
	req, _ = http.NewRequest(http.MethodDelete, baseURL+"/todo/"+created.ID, nil)
	r, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("delete todo failed: %v", err)
	}
	if r.StatusCode != 204 {
		t.Fatalf("expected 204, got %d", r.StatusCode)
	}
	r.Body.Close()
}
