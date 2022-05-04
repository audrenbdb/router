package router_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"router"
	"testing"
)

func TestRouter(t *testing.T) {
	t.Run("Given router has no matching /foo endpoint "+
		"When fetching /foo with GET request "+
		"Then it should error with err not found",
		func(t *testing.T) {
			r := router.New()
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/foo", nil)

			r.ServeHTTP(w, req)
			if w.Code != http.StatusNotFound {
				t.Errorf("got: %d, want: status not found", w.Code)
			}
		})

	t.Run("Given router has matching /ping endpoint for GET request "+
		"When fetching /ping with POST request "+
		"Then it should error with method not found", func(t *testing.T) {
		r := router.New()
		r.GET("/ping", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "pong")
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/ping", nil)

		r.ServeHTTP(w, req)
		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("got: %d, want: method not allowed", w.Code)
		}
	})

	t.Run("Given router has matching /ping endpoint for GET request "+
		"When fetching /ping with GET request "+
		"Then pong should be written in request body", func(t *testing.T) {
		r := router.New()
		r.GET("/ping", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "pong")
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/ping", nil)
		r.ServeHTTP(w, req)
		if body := w.Body.String(); body != "pong" {
			t.Errorf("got: %s, want: pong", body)
		}
	})

}
