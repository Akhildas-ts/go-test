package helper

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"fmt"
	"lock/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type AuthCustomClaims struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.StandardClaims
}

func PasswordHasing(password string) (string, error) {

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		return "", errors.New("generate hash password issue")
	}

	hash := string(hashPassword)

	return hash, nil
}

func GenerateTokenUsers(userId int, userEmail string, expirationTime time.Time) (string, error) {

	// cfg, _ := config.LoadConfig()

	claims := &AuthCustomClaims{
		Id:    userId,
		Email: userEmail,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	fmt.Println("claims data", claims.Id, claims.Email)

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Println("Error generating ECDSA key:", err)
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	tokenString, err := token.SignedString(privateKey)

	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return tokenString, nil
}

func GenerateRefreshToken(user models.SignupDetailResponse) (string, error) {

	expirationTime := time.Now().Add(24 * 90 * time.Hour)
	tokeString, err := GenerateTokenUsers(user.ID, user.Email, expirationTime)
	if err != nil {
		return "", err
	}
	return tokeString, nil

}

func GenerateAccessToken(user models.SignupDetailResponse) (string, error) {

	expirationTime := time.Now().Add(15 * time.Minute)
	tokenString, err := GenerateTokenUsers(user.ID, user.Email, expirationTime)
	if err != nil {
		return "", err
	}
	return tokenString, nil

}

// func generateRandomKey(length int) (string, error) {
// 	key := make([]byte, length)
// 	_, err := rand.Read(key)

// 	if err != nil {
// 		return "", err
// 	}
// 	return hex.EncodeToString(key), nil
// }
