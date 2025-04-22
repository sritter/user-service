package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}

type User struct {
	ID    int    `json:"id" example:"1"`
	Name  string `json:"name" example:"John Doe" binding:"required"`
	Email string `json:"email" example:"john@example.com" binding:"required,email"`
}

var users = []User{}
var nextID = 1

// CreateUser godoc
// @Summary Create a user
// @Tags users
// @Accept json
// @Produce json
// @Param user body User true "User info"
// @Success 201 {object} User
// @Router /users [post]
func CreateUser(c *gin.Context) {
	var u User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u.ID = nextID
	nextID++
	users = append(users, u)
	c.JSON(http.StatusCreated, u)
}

// ListUsers godoc
// @Summary List users
// @Tags users
// @Produce json
// @Success 200 {array} User
// @Router /users [get]
func ListUsers(c *gin.Context) {
	c.JSON(http.StatusOK, users)
}

// GetUser godoc
// @Summary Get a user by ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} User
// @Failure 404 {object} ErrorResponse
// @Router /users/{id} [get]
func GetUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	for _, u := range users {
		if u.ID == id {
			c.JSON(http.StatusOK, u)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
}

// UpdateUser godoc
// @Summary Update a user
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body User true "User info"
// @Success 200 {object} User
// @Failure 404 {object} ErrorResponse
// @Router /users/{id} [put]
func UpdateUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var updated User
	if err := c.ShouldBindJSON(&updated); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for i, u := range users {
		if u.ID == id {
			updated.ID = id
			users[i] = updated
			c.JSON(http.StatusOK, updated)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
}

// DeleteUser godoc
// @Summary Delete a user
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 204 {object} nil
// @Failure 404 {object} ErrorResponse
// @Router /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	for i, u := range users {
		if u.ID == id {
			users = append(users[:i], users[i+1:]...)
			c.Status(http.StatusNoContent)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
}
