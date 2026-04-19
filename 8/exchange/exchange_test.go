package exchange

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGetRate_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"base":"USD","target":"EUR","rate":1.23}`))
	}))
	defer ts.Close()

	svc := NewExchangeService(ts.URL)
	rate, err := svc.GetRate("USD", "EUR")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if rate != 1.23 {
		t.Errorf("expected rate 1.23, got %f", rate)
	}
}

func TestGetRate_APIBusinessError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"invalid currency pair"}`))
	}))
	defer ts.Close()

	svc := NewExchangeService(ts.URL)
	_, err := svc.GetRate("INVALID", "PAIR")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "invalid currency pair") {
		t.Errorf("expected api error message, got: %v", err)
	}
}

func TestGetRate_MalformedJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`not valid json at all`))
	}))
	defer ts.Close()

	svc := NewExchangeService(ts.URL)
	_, err := svc.GetRate("USD", "EUR")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "decode error") {
		t.Errorf("expected decode error, got: %v", err)
	}
}

func TestGetRate_Timeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(300 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	svc := NewExchangeService(ts.URL)
	svc.Client.Timeout = 50 * time.Millisecond

	_, err := svc.GetRate("USD", "EUR")
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
	if !strings.Contains(err.Error(), "network error") {
		t.Errorf("expected network error wrapping timeout, got: %v", err)
	}
}

func TestGetRate_ServerPanic(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("simulated server panic")
	}))
	defer ts.Close()

	svc := NewExchangeService(ts.URL)
	_, err := svc.GetRate("USD", "EUR")
	if err == nil {
		t.Fatal("expected error from panicking server, got nil")
	}
}

func TestGetRate_InternalServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"internal server error"}`))
	}))
	defer ts.Close()

	svc := NewExchangeService(ts.URL)
	_, err := svc.GetRate("USD", "EUR")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "internal server error") {
		t.Errorf("expected api error with server error message, got: %v", err)
	}
}

func TestGetRate_EmptyBody(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	svc := NewExchangeService(ts.URL)
	_, err := svc.GetRate("USD", "EUR")
	if err == nil {
		t.Fatal("expected error for empty body, got nil")
	}
	if !strings.Contains(err.Error(), "decode error") {
		t.Errorf("expected decode error for empty body, got: %v", err)
	}
}
