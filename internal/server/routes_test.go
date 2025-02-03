package server

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"devices_api/internal/devices"
	"devices_api/mock"

	"go.uber.org/mock/gomock"
)

func initServer(s *Server) {

}

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
	s := &Server{}
	server := httptest.NewServer(http.HandlerFunc(s.CreateDevice))
	defer server.Close()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)

	mockRepo.EXPECT().
		Create(context.Background(), devices.CreateDevice{}).
		Return(&devices.Device{Id: 1,
			Name:      "Device1",
			Brand:     "Brand1",
			State:     devices.InUse,
			CreatedAt: time.Date(2009, time.November, 10, 23, 1, 2, 0, time.UTC)},
			nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/devices", nil)
	s.db = mockRepo
	s.CreateDevice(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("io.ReadAll() error = %s; want nil", err)
	}

	want := "{\"id\":1,\"name\":\"Device1\",\"brand\":\"Brand1\",\"state\":1,\"created_at\":\"2009-11-10T23:00:00Z\" }"

	if string(body) != want {
		t.Errorf("got %v ; want %v", string(body), want)
	}

}
