package dto

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin/binding"
)

func TestUsersAdminUpdateStatusRequest(t *testing.T) {
	for _, tt := range []struct {
		name  string
		body  string
		valid bool
		want  int8
	}{
		{name: "disable", body: `{"status":0}`, valid: true, want: 0},
		{name: "enable", body: `{"status":1}`, valid: true, want: 1},
		{name: "missing", body: `{}`},
		{name: "null", body: `{"status":null}`},
		{name: "out of range", body: `{"status":2}`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(tt.body))
			var req UsersAdminUpdateStatusRequest
			err := binding.JSON.Bind(request, &req)
			if !tt.valid {
				if err == nil {
					t.Fatal("invalid status was accepted")
				}
				return
			}
			if err != nil {
				t.Fatalf("bind status: %v", err)
			}
			if req.Status == nil || *req.Status != tt.want {
				t.Fatalf("status must bind to %d", tt.want)
			}
		})
	}
}
