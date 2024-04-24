package main

import (
	"crypto/cipher"
	"encoding/json"
	"fmt"
	"github.com/emmansun/gmsm/padding"
	"github.com/emmansun/gmsm/sm4"
	"github.com/golang-jwt/jwt/v4"
	"math/rand"
	"strconv"
	"time"
)

var jwtKey = []byte("gyuachn9asd08O)@!(#HOKACXZBIU")

type Claims struct {
	UserName string `json:"userName"`
	jwt.RegisteredClaims
}

func GeneToken(userName string) (tokenString string, err error) {
	claim := Claims{
		UserName: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(3 * time.Hour * time.Duration(1))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err = token.SignedString(jwtKey)
	return
}

func getcipher() []byte {
	// 生成加密密文
	rand.Seed(time.Now().UnixNano())
	r2 := strconv.FormatInt(rand.Int63n(9999999999999999), 10)
	uid := "UID001"
	token, err := GeneToken("UID001")
	if err != nil {
		fmt.Println("token fail", err)
	}
	cipherMap := map[string]string{
		"R2":    r2,
		"UID":   uid,
		"TOKEN": token,
	}
	data, err := json.Marshal(cipherMap)
	if err != nil {
		fmt.Println("json fail", err)
	}
	fmt.Println("序列化后的结果 = ", data, len(data))
	return data
}

func main() {

	cipherData := getcipher()

	pkcs7 := padding.NewPKCS7Padding(sm4.BlockSize)

	r1 := []byte(strconv.FormatInt(rand.Int63n(9999999999999999), 10))
	c, err := sm4.NewCipher(r1)
	if err != nil {
		fmt.Printf(": NewCipher(bytes) = %s", err)
	}

	encrypter := cipher.NewCBCEncrypter(c, r1)
	plainText := pkcs7.Pad(cipherData)
	data := make([]byte, len(plainText))
	copy(data, plainText)

	encrypter.CryptBlocks(data, plainText)
	fmt.Println("加密后的内容: ", data, len(data), len(plainText))

	// 	解密
	decrypter := cipher.NewCBCDecrypter(c, r1)
	decData := make([]byte, len(data))
	copy(decData, data)

	decrypter.CryptBlocks(decData, data)
	fmt.Println("解密:", decData, len(decData))
	fmt.Println(string(decData))

	// 	反序列化
	res := make(map[string]string)
	err = json.Unmarshal(decData[:len(cipherData)], &res)
	if err != nil {
		fmt.Println("un err", err)
	}
	fmt.Println("res = ", res)

}
