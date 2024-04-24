package mySM9

import (
	"SM9_Client/common"
	"SM9_Client/server/pkg"
	cipher2 "crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/emmansun/gmsm/kdf"
	"github.com/emmansun/gmsm/padding"
	"github.com/emmansun/gmsm/sm3"
	"github.com/emmansun/gmsm/sm4"
	"github.com/emmansun/gmsm/sm9"
	"log"
)

// Encrypt SM9加密
func Encrypt(publicKey []byte, plainText []byte) ([]byte, error) {
	masterKey := new(sm9.EncryptMasterPublicKey)
	err := masterKey.UnmarshalASN1(publicKey)
	if err != nil {
		return nil, err
	}
	hid := byte(0x01)
	encData, err := sm9.Encrypt(rand.Reader, masterKey, []byte(pkg.Pkg.PKGID), hid, plainText)
	if err != nil {
		return nil, err
	}
	return encData, nil
}

// Decrypt SM9解密
func Decrypt() {

}

// Sign SM9加签
func Sign(plainText []byte) ([]byte, error) {

	hid := byte(0x01)
	masterKey, err := sm9.GenerateSignMasterKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	userKey, err := masterKey.GenerateUserKey([]byte(pkg.Pkg.PKGID), hid)
	if err != nil {
		return nil, err
	}
	SignData, err := sm9.SignASN1(rand.Reader, userKey, plainText)
	if err != nil {
		return nil, err
	}

	// 将&userKey.SignMaterPublicKey与SignCipherSID打包
	temp, err := masterKey.Public().MarshalASN1()
	if err != nil {
		return nil, err
	}
	data := map[string][]byte{
		"SignCipherSID": SignData,
		"SignPublicKey": temp,
	}
	res, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// EncryptSm4 SM4加密
func EncryptSm4(key []byte, message []byte) ([]byte, error) {
	pkcs7 := padding.NewPKCS7Padding(sm4.BlockSize)
	c, err := sm4.NewCipher(key)
	if err != nil {
		fmt.Printf(": NewCipher(bytes) = %s", err)
		return nil, err
	}
	encrypter := cipher2.NewCBCEncrypter(c, common.CK)
	plainText := pkcs7.Pad(message)
	encData := make([]byte, len(plainText))
	encrypter.CryptBlocks(encData, plainText)
	return encData, nil
}

// DecryptSm4 SM4解密
func DecryptSm4(key []byte, plainText []byte) ([]byte, error) {
	c, err := sm4.NewCipher(key)
	if err != nil {
		log.Println("NewCipher(fail)", err)
		return nil, err
	}
	decrypter := cipher2.NewCBCDecrypter(c, key)
	decData := make([]byte, len(plainText))
	decrypter.CryptBlocks(decData, plainText)
	return decData, nil
}

// Kdf 生成CK和IK
func Kdf(key []byte) ([]byte, []byte) {
	res := kdf.Kdf(sm3.New(), key, 32)
	return res[:len(res)/2], res[len(res)/2:]
}

// HmacData 哈希数据
func HmacData(plainText []byte, key []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(plainText)
	return mac.Sum(nil)
}
