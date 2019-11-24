package auth

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/models"
	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/utils"
)

//JWT authority structure
type Token struct {
	UserId int
	jwt.StandardClaims
}

func JwtAuth(claims Token, tokenStr string) (Token, map[string] interface{}) {
	JwtKeyStr, check := os.LookupEnv("secret")
	if !check {
		fmt.Println("not secret")
	}
	JwtKey := []byte(JwtKeyStr)
	tkn, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error){
		return JwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return claims, utils.Message(false, "Invalid Signature", "Unauthorized")
		}
		return claims,  utils.Message(false, err.Error(),"Bad Request")
	}
	if !tkn.Valid {
		return claims, utils.Message(false, "Invalid token","Unauthorized")
	}
	return claims, nil
}

func CreateTokenAndSetCookie(w http.ResponseWriter, user models.User) map[string] interface{} {
	//Setting session time
	expirationTime := time.Now().Add(1 * time.Minute)
	//Create token JWT
	tk := &Token{
		UserId: user.Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	//Token generation
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	//Get the complete, signed token
	JwtKeyStr, check := os.LookupEnv("secret")
	if !check {
		fmt.Println("not secret")
	}
	JwtKey := []byte(JwtKeyStr)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		return utils.Message(false, err.Error(),"Internal Server Error")
	}
	http.SetCookie(w, &http.Cookie{
		Name: "access_token",
		Value: tokenString,
		Expires: expirationTime,
	})

	// Generate refresh_token
	expirationTime = time.Now().Add(24 * time.Hour)
	rtk := &Token{
		UserId: user.Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), rtk)
	rTokenStr, err := refreshToken.SignedString(JwtKey)
	if err != nil {
		return utils.Message(false, err.Error(),"Internal Server Error")
	}
	http.SetCookie(w, &http.Cookie{
		Name: "refresh_token",
		Value: rTokenStr,
		Expires: expirationTime,
	})

	return nil
}

func RefreshToken(claims Token) (tokenStr string, resp map[string] interface{}) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	JwtKeyStr, check := os.LookupEnv("secret")
	if !check {
		fmt.Println("not secret")
	}
	JwtKey := []byte(JwtKeyStr)
	tokenStr, err := token.SignedString(JwtKey)
	if err != nil {
		return tokenStr, utils.Message(false, "refresh token error","Unauthorized")
	}
	return tokenStr, nil
}

func RefreshAccessAndRefreshToken(w http.ResponseWriter, r *http.Request, id int) map[string] interface{} {
	//getting a token from cookies
	c, err := r.Cookie("refresh_token")
	if err != nil {
		if err == http.ErrNoCookie {
			utils.Message(false, "No cookie", "Unauthorized")
		}
		return utils.Message(false,"Cookie error","Bad Request")
	}

	checkTokenStr := c.Value
	claims := Token{}
	claims, resp := JwtAuth(claims, checkTokenStr)
	if resp != nil {
		return resp
	}

	//update time
	expirationTime := time.Now().Add(1 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	// generate new token
	checkTokenStr, resp = RefreshToken(claims)
	if resp != nil {
		return resp
	}

	http.SetCookie(w, &http.Cookie{
		Name: "access_token",
		Value: checkTokenStr,
		Expires: expirationTime,
		Path: "/",
	})

	expirationTime = time.Now().Add(24 * time.Hour)
	claims.ExpiresAt = expirationTime.Unix()
	checkTokenStr, resp = RefreshToken(claims)
	if resp != nil {
		return resp
	}

	http.SetCookie(w, &http.Cookie{
		Name: "refresh_token",
		Value: checkTokenStr,
		Expires: expirationTime,
	    Path: "/",
	})

	//check id in url and id in cookies
	if claims.UserId != id {
		return utils.Message(false,"id do not match","Unauthorized")
	}
	resp = utils.Message(true, "Check token", "")
	resp["token"] = claims
	return resp
}

func CheckTokenAndRefresh(w http.ResponseWriter, r *http.Request, id int) map[string] interface{} {
	//check deadline access_token
	var resp map[string] interface{}
	c, err := r.Cookie("access_token")
	if err != nil {
		resp = RefreshAccessAndRefreshToken(w, r, id)
	} else {
		checkTokenStr := c.Value
		claims := Token{}
		claims, resp := JwtAuth(claims, checkTokenStr)
		if resp != nil {
			return resp
		}
		//check id in url and id in cookies
		if claims.UserId != id {
			return utils.Message(false,"id do not match","Unauthorized")
		}
		return utils.Message(true, "Check token", "")
	}
	return resp
}
