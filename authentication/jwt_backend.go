package authentication

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

const (
	privKeyPath = "authentication/app.rsa"
	pubKeyPath  = "authentication/app.rsa.pub"
)

var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

func init() {
	signBytes, err := ioutil.ReadFile(privKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		log.Fatal(err)
	}

	verifyBytes, err := ioutil.ReadFile(pubKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		log.Fatal(err)
	}

}

type UserTest struct {
	*jwt.StandardClaims
	ID       uint   `json:"id"`
	UserName string `json:"username"`
}

//Token
func CreateToken(username string, id uint) (string, error) {
	// create a signer for rsa 256
	t := jwt.New(jwt.GetSigningMethod("RS256"))

	// set our claims
	t.Claims = &UserTest{
		&jwt.StandardClaims{
			// set the expire time
			// see http://tools.ietf.org/html/draft-ietf-oauth-json-web-token-20#section-4.1.4
			ExpiresAt: time.Now().Add(time.Minute * 1).Unix(),
		},
		id,
		username,
	}

	// Creat token string
	return t.SignedString(signKey)
}
func AuthHandler(w http.ResponseWriter, r *http.Request) {
	// make sure its post
	// if r.Method != "POST" {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	fmt.Fprintln(w, "No POST", r.Method)
	// 	return
	// }
	//
	// user := r.FormValue("user")
	// pass := r.FormValue("pass")
	//
	// log.Printf("Authenticate: user[%s] pass[%s]\n", user, pass)
	//
	// // check values
	// if user != "test" || pass != "known" {
	// 	w.WriteHeader(http.StatusForbidden)
	// 	fmt.Fprintln(w, "Wrong info")
	// 	return
	// }
	//
	// tokenString, err := CreateToken(user)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	fmt.Fprintln(w, "Sorry, error while Signing Token!")
	// 	log.Printf("Token Signing error: %v\n", err)
	// 	return
	// }

	// w.Header().Set("Content-Type", "application/jwt")
	// w.WriteHeader(http.StatusOK)
	// fmt.Fprintln(w, tokenString)
}
func RestrictedHandler(w http.ResponseWriter, r *http.Request) {
	// Get token from request
	token, err := request.ParseFromRequestWithClaims(r, request.OAuth2Extractor, &UserTest{}, func(token *jwt.Token) (interface{}, error) {
		// since we only use the one private key to sign the tokens,
		// we also only use its public counter part to verify

		return verifyKey, nil
	})

	// If the token is missing or invalid, return error
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "Invalid token:", err)
		return
	}

	// Token is valid
	fmt.Fprintln(w, "Welcome,", token.Claims.(*UserTest).UserName)
	return
}
