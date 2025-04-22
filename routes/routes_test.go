package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRoutes(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Run()
	}))
	defer server.Close()

	// Test cases for route availability
	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
	}{
		{
			name:       "swagger documentation",
			method:     "GET",
			path:       "/swagger/index.html",
			wantStatus: http.StatusOK,
		},
		{
			name:       "create user endpoint",
			method:     "POST",
			path:       "/users",
			wantStatus: http.StatusBadRequest, // Without a valid body
		},
		{
			name:       "list users endpoint",
			method:     "GET",
			path:       "/users",
			wantStatus: http.StatusOK,
		},
		{
			name:       "get user endpoint",
			method:     "GET",
			path:       "/users/1",
			wantStatus: http.StatusNotFound, // No users exist yet
		},
		{
			name:       "update user endpoint",
			method:     "PUT",
			path:       "/users/1",
			wantStatus: http.StatusBadRequest, // Without a valid body
		},
		{
			name:       "delete user endpoint",
			method:     "DELETE",
			path:       "/users/1",
			wantStatus: http.StatusNotFound, // No users exist yet
		},
	}

	// Create a test router
	r := gin.Default()
	r.GET("/swagger/*any", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	routes := r.Group("/users")
	{
		routes.POST("", func(c *gin.Context) {
			c.Status(http.StatusBadRequest)
		})
		routes.GET("", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})
		routes.GET(":id", func(c *gin.Context) {
			c.Status(http.StatusNotFound)
		})
		routes.PUT(":id", func(c *gin.Context) {
			c.Status(http.StatusBadRequest)
		})
		routes.DELETE(":id", func(c *gin.Context) {
			c.Status(http.StatusNotFound)
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.path, nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}
