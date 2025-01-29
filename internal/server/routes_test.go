package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	repo "devices_api/internal/devices/postgres"

	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	s := &Server{}
	server := httptest.NewServer(http.HandlerFunc(s.HelloWorldHandler))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	defer resp.Body.Close()

	// Assertions
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}

	expected := "{\"message\":\"Hello World\"}"
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	if expected != string(body) {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}
}

func TestCreateDevice(t *testing.T) {
	s := &Server{
		port: 8080,
		db:   repo.NewRepository(),
	}
	server := httptest.NewServer(http.HandlerFunc(s.CreateDevice))
	defer server.Close()

	type device struct {
		Name  string `json:"name"`
		Brand string `json:"brand"`
		State int    `json:"state"`
	}

	d := device{Name: "abc123", Brand: "abc", State: 0}
	dJson, err := json.Marshal(d)
	if err != nil {
		t.Errorf("POST /devices/new: %s", err)
	}

	respBody := bytes.NewReader(dJson)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/devices/new", respBody)

	s.CreateDevice(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	assert.Equal(t, body, dJson)

}
