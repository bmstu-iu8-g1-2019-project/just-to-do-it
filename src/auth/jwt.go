package auth

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/dgrijalva/jwt-go"
)

func CreateAccessToken(user_id int) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = user_id
	claims["exp"] = time.Now().Add(time.Minute * 1).Unix() //Token expires after 1 hour
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	JwtSecret, _ := os.LookupEnv("secret")
	JwtKey := []byte(JwtSecret)
	return token.SignedString(JwtKey)
}

func CreateRefreshToken(user_id int) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = user_id
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() //Token expires after 1 hour
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	JwtSecret, _ := os.LookupEnv("secret")
	JwtKey := []byte(JwtSecret)
	return token.SignedString(JwtKey)
}

func ExtractAccessToken(r *http.Request) string {
	c, err := r.Cookie("access_token")
	if err != nil {
		return ""
	}
	token := c.Value
	if token != "" {
		return token
	}
	return ""
}

func ExtractRefreshToken(r *http.Request) string {
	c, err := r.Cookie("refresh_token")
	if err != nil {
		return ""
	}
	token := c.Value
	if token != "" {
		return token
	}
	return ""
}

func TokenValid(w http.ResponseWriter, r *http.Request, userId int) error {
	//check acc token
	tokenString := ExtractAccessToken(r)
	// if acc token has expired
	if tokenString == "" {
		err := RefTokenValid(r, userId)
		if err != nil {
			return err
		}
		err = RefreshTokens(w, r)
		if err != nil {
			return err
		}
		return nil
	}

	JwtSecret, _ := os.LookupEnv("secret")
	JwtKey := []byte(JwtSecret)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v ", token.Header["alg"])
		}
		return JwtKey, nil
	})
	if err != nil {
		return err
	}
	claims := token.Claims.(jwt.MapClaims)
	userIdFromToken := int(claims["user_id"].(float64))
	if userId != userIdFromToken {
		return fmt.Errorf("Id don't match ")
	}
	return nil
}

func RefTokenValid(r *http.Request, userId int) error {
	//check refresh token
	tokenString := ExtractRefreshToken(r)
	JwtSecret, _ := os.LookupEnv("secret")
	JwtKey := []byte(JwtSecret)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v ", token.Header["alg"])
		}
		return JwtKey, nil
	})
	if err != nil {
		return err
	}
	claims := token.Claims.(jwt.MapClaims)
	userIdFromToken := int(claims["user_id"].(float64))
	if userId != userIdFromToken {
		return fmt.Errorf("Id don't match ")
	}
	return nil
}

func RefreshTokens (w http.ResponseWriter, r *http.Request) error {
	refTokenStr := ExtractRefreshToken(r)
	JwtSecret, _ := os.LookupEnv("secret")
	JwtKey := []byte(JwtSecret)
	token, err := jwt.Parse(refTokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v ", token.Header["alg"])
		}
		return JwtKey, nil
	})
	if err != nil {
		return err
	}
	claims := token.Claims.(jwt.MapClaims)
	userId := int(claims["user_id"].(float64))
	accToken, err := CreateAccessToken(userId)
	if err != nil {
		return err
	}
	SetCookieForAccToken(w, accToken)
	refTokenStr, err = CreateRefreshToken(userId)
	if err != nil {
		return err
	}
	SetCookieForRefToken(w, refTokenStr)
	return nil
}

func SetCookieForAccToken (w http.ResponseWriter, token string) {
	expirationTime := time.Now().Add(1 * time.Minute)
	http.SetCookie(w, &http.Cookie{
		Name: "access_token",
		Value: token,
		Expires: expirationTime,
		Path: "/",
	})
}

func SetCookieForRefToken (w http.ResponseWriter, token string) {
	expirationTime := time.Now().Add(24 * time.Hour)
	http.SetCookie(w, &http.Cookie{
		Name: "refresh_token",
		Value: token,
		Expires: expirationTime,
		Path: "/",
	})
}

func CheckUser(w http.ResponseWriter, r *http.Request) (int, error) {
	paramFromURL := mux.Vars(r)
	userId, err := strconv.Atoi(paramFromURL["id"])
	if err != nil {
		return 0, err
	}
	//проверка и в случае таймута рефреш токена
	err = TokenValid(w, r, userId)
	if err != nil {
		return 0, err
	}
	return userId, nil
}
