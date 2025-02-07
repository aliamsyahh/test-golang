// service/authService.go

package service

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Secret key for JWT signing.
var jwtKey = []byte("685439f3e9484bdd6f4d1246e8eb7940d29d020800bde59524dbbfdb70479962")
var sessionExpiration = 24 * time.Hour

// Claims struct to hold the JWT claims.
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenerateJWT creates a new JWT token for a user.
func GenerateJWT(username string) (string, error) {
	// Set expiration time for token.
	expirationTime := time.Now().Add(sessionExpiration)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	// Create token with claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// ValidateJWT checks the validity of the provided JWT.
func ValidateJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, err
	}
	return claims, nil
}
