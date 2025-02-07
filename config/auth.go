package config

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var GoogleOauthConfig = &oauth2.Config{
	ClientID:     "117760943520-9ifkhadv9p8jnlklp2sp23nvp883aovd.apps.googleusercontent.com",
	ClientSecret: "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",
	RedirectURL:  "http://localhost:8080/auth/google/callback",
	Scopes:       []string{"email", "profile"},
	Endpoint:     google.Endpoint,
}

var JwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateJWT(userID string, email string) (string, error) {
	claims := jwt.MapClaims{
		"id":    userID,
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecret)
}
