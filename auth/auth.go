package auth

import (
	"context"
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/earqq/gqlgen-easybill/db"
	"github.com/earqq/gqlgen-easybill/graph/model"
	"gopkg.in/mgo.v2/bson"
)

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)
var userCtxKey = &contextKey{"person"}

type contextKey struct {
	Name string
}

func init() {
	privateBytes, err := ioutil.ReadFile("/var/go/src/easybill/private.rsa")
	if err != nil {
		log.Fatal("No se puede leer llave privada")
	}
	publicBytes, err := ioutil.ReadFile("/var/go/src/easybill/public.rsa.pub")
	if err != nil {
		log.Fatal("No se puedo leer llave pública")
	}
	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateBytes)
	if err != nil {
		log.Fatal("No se pudo parsear llave privada")
	}
	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicBytes)
	if err != nil {
		log.Fatal("No se pudo parsear llave pública")
	}
}

func GenerateJWT(key string) string {
	claims := model.Claim{
		PrivateKey: key,
		StandardClaims: jwt.StandardClaims{
			Issuer: "Login token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	result, _ := token.SignedString(privateKey)
	return result
}

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := request.ParseFromRequestWithClaims(r, request.OAuth2Extractor, &model.Claim{}, func(token *jwt.Token) (interface{}, error) {
				return publicKey, nil
			})
			if err == nil && token.Valid {
				var person model.Person
				tokenString := TokenFromHttpRequest(r)
				privateKeyString := PrivateKeyFromToken(tokenString)
				peopleDB := db.GetCollection("people")
				_ = peopleDB.Find(bson.M{"private_key": privateKeyString}).Select(bson.M{"_id": 0, "password": 0}).One(&person)
				// put it in context
				ctx := context.WithValue(r.Context(), userCtxKey, &person)
				// and call the next with our new context
				r = r.WithContext(ctx)
			}
			next.ServeHTTP(w, r)

		})
	}
}

func TokenFromHttpRequest(r *http.Request) string {
	reqToken := r.Header.Get("Authorization")
	var tokenString string
	splitToken := strings.Split(reqToken, "Bearer ")
	if len(splitToken) > 1 {
		tokenString = splitToken[1]
	}
	return tokenString
}
func JwtDecode(token string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, &model.Claim{}, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
}
func PrivateKeyFromToken(tokenString string) string {

	token, err := JwtDecode(tokenString)
	if err != nil {
		fmt.Println(err)
		return "1"
	}
	if claims, ok := token.Claims.(*model.Claim); ok && token.Valid {
		if claims == nil {
			return "2 "
		}
		return claims.PrivateKey
	} else {
		return "3"
	}
}
func ForContext(ctx context.Context) *model.Person {
	raw, _ := ctx.Value(userCtxKey).(*model.Person)
	return raw
}
func GetAuthFromContext(ctx context.Context) *model.Person {
	return ForContext(ctx)
}
