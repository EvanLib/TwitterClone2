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
	ID        uint   `json:"id"`
	UserName  string `json:"username"`
	ProfileID uint   `json:"profileid"`
}

//Midleware

//Token
func CreateToken(username string, id uint, profileid uint) (string, error) {
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
		profileid,
	}

	// Creat token string
	return t.SignedString(signKey)
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
	w.Header().Set("Access-Control-Allow-Origin", "*")

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
