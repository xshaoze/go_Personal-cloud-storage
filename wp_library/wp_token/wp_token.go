package wp_token

import (
	"encoding/base64"
	"log"

	"github.com/dgrijalva/jwt-go"
)

var header = map[string]interface{}{
	"alg":    "HS256",
	"typ":    "JWT",
	"msgtyp": "JSON",
}

var password = []byte("|cqLcZ_23~hVzGUi8$SljZ-eXj-Fe/@4")

func JwtEncryption(data map[string]interface{}, key []byte) string {
	claims := jwt.MapClaims{}
	for key, value := range data {
		claims[key] = value
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header = header
	tokenString, err := token.SignedString(key)
	if err != nil {
		log.Println("wp_token.go 的 JwtEncryption()出现错误:", err)
		return ""
	}

	token1, err := encrypt([]byte(tokenString), password)
	if err != nil {
		log.Fatalln(err)
	}

	return token1
}

func JwtDecryption(data string, key []byte) map[string]interface{} {
	data1, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		log.Println("jwt解密失败", err)
		return nil
	}
	dataEncode, _ := decrypt(data1, password)

	token, err := jwt.Parse(string(dataEncode), func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		log.Println("wp_token.go 的 JwtDecryption() 出现错误(1):", err)
		return nil
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		log.Println("wp_token.go 的 JwtDecryption() 出现错误(2):", "invalid token")
		return nil
	}
	result := make(map[string]interface{})
	for key, value := range claims {
		result[key] = value
	}

	return result
}
