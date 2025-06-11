package controllers

import (
	"gocheck/models"
	"gocheck/services"
	"gocheck/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserController handles user-related HTTP requests
type UserController struct {
	userService *services.UserService
}

// NewUserController creates a new UserController
func NewUserController(db *gorm.DB) *UserController {
	return &UserController{
		userService: services.NewUserService(db),
	}
}

// CreateUser handles user registration and generates a token upon successful creation.
func (uc *UserController) CreateUser(c *gin.Context) {
	var user models.User // User model (or UserCreateInput DTO if you chose that approach)
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Important: If models.User.Password still has `json:"-"`, this field
	// will be empty here. You need to ensure the Password from the request
	// is properly bound. If you're not using a DTO, you must change
	// `json:"-"` to `json:"password"` in models.User for `CreateUser` to work
	// correctly with the incoming password. Or, as in the previous example,
	// manually extract it from a map if you absolutely refuse DTOs and `json` tag.
	// Assuming models.User's Password field has `json:"password"` now, OR
	// you are using a separate DTO (like UserCreateInput) with json:"password".
	// If you are sticking to `json:"-"` and no DTO for create, you would need
	// to bind to a map[string]interface{} first, then manually construct `models.User`
	// like we did for Login. For now, I'm assuming a bound 'user.Password'.

	if err := uc.userService.CreateUser(&user); err != nil {
		// Handle specific database errors if needed (e.g., duplicate entry)
		// For example, if you imported pq for PostgreSQL specific errors:
		// if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
		//     c.JSON(http.StatusConflict, gin.H{"error": "Username or email already exists"})
		//     return
		// }
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// --- NEW: Generate Token after successful user creation ---
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Respond with success message, user info (non-sensitive), and the token
	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user": gin.H{
			"ID":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			// Do NOT include Password here, even if it were possible
		},
		"token": token, // Include the generated token
	})
}

// GetUserByID handles fetching a user by ID
func (uc *UserController) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := uc.userService.GetUserByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetAllUsers handles fetching all users
func (uc *UserController) GetAllUsers(c *gin.Context) {
	users, err := uc.userService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// UpdateUser handles updating an existing user
func (uc *UserController) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.ID = uint(id) // Ensure the ID from the URL is used

	updatedUser, err := uc.userService.UpdateUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

// DeleteUser handles deleting a user
func (uc *UserController) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := uc.userService.DeleteUser(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusNoContent, nil) // 204 No Content for successful deletion
}

// / Login handles user login requests and generates a token upon successful authentication.
func (uc *UserController) Login(c *gin.Context) {
	// Bind incoming JSON to a generic map for email and password
	var loginData map[string]string
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	email, ok := loginData["email"]
	if !ok || email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}

	password, ok := loginData["password"]
	if !ok || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password is required"})
		return
	}

	// Authenticate the user using the UserService
	authenticatedUser, err := uc.userService.AuthenticateUser(email, password)
	if err != nil {
		// Return a generic "Invalid credentials" error for security,
		// whether it's email not found or password mismatch.
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// --- NEW: Generate Token after successful authentication ---
	token, err := utils.GenerateToken(authenticatedUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Login successful - return message, user info, and the token
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user": gin.H{ // Return only safe user data, NOT the password
			"ID":       authenticatedUser.ID,
			"username": authenticatedUser.Username,
			"email":    authenticatedUser.Email,
		},
		"token": token, // Include the generated token
	})
}
