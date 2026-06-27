package server

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProfilePageRoutesAreRegistered(t *testing.T) {
	srv := NewServer(&sql.DB{})

	tests := []struct {
		name   string
		method string
		path   string
	}{
		{name: "profile", method: http.MethodGet, path: "/api/users/user-1"},
		{name: "profile privacy", method: http.MethodPut, path: "/api/users/user-1/privacy"},
		{name: "user posts", method: http.MethodGet, path: "/api/users/user-1/posts"},
		{name: "followers", method: http.MethodGet, path: "/api/users/user-1/followers"},
		{name: "following", method: http.MethodGet, path: "/api/users/user-1/following"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()

			srv.mux.ServeHTTP(rec, req)

			if rec.Code == http.StatusNotFound {
				t.Fatalf("%s %s returned 404; route is not registered", tt.method, tt.path)
			}
		})
	}
}
