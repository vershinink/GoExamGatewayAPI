package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLatest(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /news/latest", Latest())
	srv := httptest.NewServer(mux)
	defer srv.Close()

	req := httptest.NewRequest(http.MethodGet, "/news/latest?page=1", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Latest() error = unexpected status code, got = %v, want = %v", rr.Code, http.StatusOK)
	}

	if rr.Body.Len() == 0 {
		t.Errorf("Latest() error = empty body")
	}
}

func TestFilter(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /news/filter", Filter())
	srv := httptest.NewServer(mux)
	defer srv.Close()

	req := httptest.NewRequest(http.MethodGet, "/news/filter", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Filter() error = unexpected status code, got = %v, want = %v", rr.Code, http.StatusOK)
	}

	if rr.Body.Len() == 0 {
		t.Errorf("Filter() error = empty body")
	}
}

func TestDetailed(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /news/detailed/{id}", Detailed())
	srv := httptest.NewServer(mux)
	defer srv.Close()

	req := httptest.NewRequest(http.MethodGet, "/news/detailed/news02", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Detailed() error = unexpected status code, got = %v, want = %v", rr.Code, http.StatusOK)
	}

	if rr.Body.Len() == 0 {
		t.Errorf("Detailed() error = empty body")
	}
}

func TestAddComment(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /news/comment", AddComment())
	srv := httptest.NewServer(mux)
	defer srv.Close()

	b, err := json.Marshal(CommentsNews1)
	if err != nil {
		t.Fatalf("cannot encode comment")
	}
	req := httptest.NewRequest(http.MethodPost, "/news/comment", bytes.NewReader(b))
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Errorf("Latest() error = unexpected status code, got = %v, want = %v", rr.Code, http.StatusCreated)
	}
}
