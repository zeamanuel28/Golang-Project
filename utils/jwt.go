package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5" // Import the new JWT library
)

// Define a struct for custom JWT claims (payload)
type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// TODO: Store this in an environment variable or a secure configuration
// NEVER hardcode this in a production application!
var jwtSecret = []byte("your_super_secret_jwt_key_that_is_at_least_32_bytes_long")

// GenerateToken generates a new JWT for the given user ID
func GenerateToken(userID uint, role string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Token expires in 24 hours

	claims := &Claims{
		UserID: userID,
		Role:   role, // Include the role in the token

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "user_authentication",
		},
	}

	// Create the token with the algorithm and claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT string and returns the claims if valid
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid // Return a specific error if method is wrong
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token") // Use errors.New("invalid token") from standard library
	}

	return claims, nil
}
