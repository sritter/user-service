package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTest() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	return r
}

func TestCreateUser(t *testing.T) {
	r := setupTest()
	r.POST("/users", CreateUser)

	tests := []struct {
		name       string
		user       User
		wantStatus int
	}{
		{
			name: "valid user",
			user: User{
				Name:  "John Doe",
				Email: "john@example.com",
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "invalid user - missing name",
			user: User{
				Email: "john@example.com",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset users slice before each test
			users = []User{}
			nextID = 1

			body, _ := json.Marshal(tt.user)
			req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantStatus == http.StatusCreated {
				var response User
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.user.Name, response.Name)
				assert.Equal(t, tt.user.Email, response.Email)
				assert.Equal(t, 1, response.ID)
			}
		})
	}
}

func TestListUsers(t *testing.T) {
	r := setupTest()
	r.GET("/users", ListUsers)

	// Reset and populate users slice
	users = []User{
		{ID: 1, Name: "John Doe", Email: "john@example.com"},
		{ID: 2, Name: "Jane Doe", Email: "jane@example.com"},
	}
	nextID = 3

	req := httptest.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, users, response)
}

func TestGetUser(t *testing.T) {
	r := setupTest()
	r.GET("/users/:id", GetUser)

	// Reset and populate users slice
	users = []User{
		{ID: 1, Name: "John Doe", Email: "john@example.com"},
	}
	nextID = 2

	tests := []struct {
		name       string
		userID     string
		wantStatus int
		wantUser   *User
	}{
		{
			name:       "existing user",
			userID:     "1",
			wantStatus: http.StatusOK,
			wantUser:   &users[0],
		},
		{
			name:       "non-existent user",
			userID:     "999",
			wantStatus: http.StatusNotFound,
			wantUser:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/users/"+tt.userID, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantUser != nil {
				var response User
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantUser.ID, response.ID)
				assert.Equal(t, tt.wantUser.Name, response.Name)
				assert.Equal(t, tt.wantUser.Email, response.Email)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	r := setupTest()
	r.PUT("/users/:id", UpdateUser)

	// Reset and populate users slice
	users = []User{
		{ID: 1, Name: "John Doe", Email: "john@example.com"},
	}
	nextID = 2

	tests := []struct {
		name       string
		userID     string
		updateUser User
		wantStatus int
	}{
		{
			name:   "valid update",
			userID: "1",
			updateUser: User{
				Name:  "John Updated",
				Email: "john.updated@example.com",
			},
			wantStatus: http.StatusOK,
		},
		{
			name:   "non-existent user",
			userID: "999",
			updateUser: User{
				Name:  "John Updated",
				Email: "john.updated@example.com",
			},
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.updateUser)
			req := httptest.NewRequest("PUT", "/users/"+tt.userID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantStatus == http.StatusOK {
				var response User
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.updateUser.Name, response.Name)
				assert.Equal(t, tt.updateUser.Email, response.Email)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	r := setupTest()
	r.DELETE("/users/:id", DeleteUser)

	tests := []struct {
		name       string
		userID     string
		wantStatus int
	}{
		{
			name:       "existing user",
			userID:     "1",
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "non-existent user",
			userID:     "999",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset and populate users slice before each test
			users = []User{
				{ID: 1, Name: "John Doe", Email: "john@example.com"},
			}
			nextID = 2

			req := httptest.NewRequest("DELETE", "/users/"+tt.userID, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantStatus == http.StatusNoContent {
				assert.Empty(t, users)
			}
		})
	}
}
