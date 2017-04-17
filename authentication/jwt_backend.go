package authentication

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gorilla/context"
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

//Midleware

//Token
func CreateToken(username string, id uint) (string, error) {
	// create a signer for rsa 256
	t := jwt.New(jwt.GetSigningMethod("RS256"))

	// set our claims
	t.Claims = &UserTest{
		&jwt.StandardClaims{
			// set the expire time
			// see http://tools.ietf.org/html/draft-ietf-oauth-json-web-token-20#section-4.1.4
			ExpiresAt: time.Now().Add(time.Minute * 30).Unix(),
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

//Negroni function for middleware
func HandlerWithNext(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	fmt.Println("Called HandlerWithNext")

	err := JWTAuth(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	// If there was an error, do not call next.
	if err == nil && next != nil {
		next(w, r)
	}
}

//Extractor
func FromAuthHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", nil
	}

	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", fmt.Errorf("Authorization header format must be Bearer {token}")
	}

	return authHeaderParts[1], nil
}

//Auth user via token
//Set Claims as context
func JWTAuth(w http.ResponseWriter, r *http.Request) error {
	//Retrive header information.
	headerToken, err := FromAuthHeader(r)
	if err != nil {
		fmt.Printf("Error retrieving token : %v", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return err
	}

	//Parse Token
	token, err := jwt.ParseWithClaims(headerToken, &UserTest{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return verifyKey, nil
	})

	if err != nil {
		fmt.Println(err)
		return err
	}
	if !token.Valid {
		fmt.Println("Token not valid.")
		return fmt.Errorf("Token not valid.")
	}

	context.Set(r, "user_info", token.Claims)
	fmt.Println(token.Claims.(*UserTest).UserName)
	return nil
}
