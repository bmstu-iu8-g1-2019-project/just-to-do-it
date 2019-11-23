package auth

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/models"
	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/utils"
)

//TODO: refactor this awesome way to store the secret
var JwtKey = []byte("secret")

//JWT authority structure
type Token struct {
	UserId int
	jwt.StandardClaims
}


func JwtAuth(claims Token, tokenStr string) (Token, map[string] interface{}) {
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
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		return utils.Message(false, "","Internal Server Error")
	}
	http.SetCookie(w, &http.Cookie{
		Name: "token",
		Value: tokenString,
		Expires: expirationTime,
	})
	return nil
}

func RefreshToken(claims Token) (tokenStr string, resp map[string] interface{}) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(JwtKey)
	if err != nil {
		return tokenStr, utils.Message(false, "refresh token error","Unauthorized")
	}
	return tokenStr, nil
}

func CheckTokenAndRefresh(w http.ResponseWriter, r *http.Request, id int) map[string] interface{}{
	//getting a token from cookies
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			return utils.Message(false, "No cookie","Unauthorized")
		}
		return utils.Message(false,"Cookie error","Bad Request")
	}

	tknStr := c.Value
	claims := Token{}
	//check token and setting parameters
	claims, resp := JwtAuth(claims, tknStr)
	if resp != nil {
		return resp
	}

	//refresh token if expiration time < 30sec
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) < 30 * time.Second {
		//update time
		expirationTime := time.Now().Add(1 * time.Minute)
		claims.ExpiresAt = expirationTime.Unix()
		// generate new token
		tknStr, resp = RefreshToken(claims)
		if resp != nil {
			return resp
		}

		http.SetCookie(w, &http.Cookie{
			Name: "token",
			Value: tknStr,
			Expires: expirationTime,
		})
	}

	//check id in url and id in cookies
	if claims.UserId != id {
		return utils.Message(false,"id do not match","Unauthorized")
	}
	resp = utils.Message(true, "Check token", "")
	resp["token"] = claims
	return resp
}
