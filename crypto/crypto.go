package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

//此处加解密、消息摘要
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5Unpadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

//AESEncryptByECB
//golang和Java中的AES模式不一样：java aes默认加密模式为ECB
//golang 默认的是CBC模式
func AESEncryptByECB(text, key string) string {
	return base64.StdEncoding.EncodeToString(EcbEncrypt([]byte(text), []byte(key)))
}

func AESDecryptByECB(text, key string) string {
	bs, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		fmt.Errorf("has err,text:%s,decodeString:%s,err:%s", text, string(bs), err)
	}
	return string(EcbDecrypt(bs, []byte(key)))
}

func EcbEncrypt(data, key []byte) []byte {
	block, _ := aes.NewCipher(key)
	data = PKCS5Padding(data, block.BlockSize())
	decrypted := make([]byte, len(data))
	size := block.BlockSize()

	for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
		block.Encrypt(decrypted[bs:be], data[bs:be])
	}
	return decrypted
}

func EcbDecrypt(data, key []byte) []byte {
	block, _ := aes.NewCipher(key)
	decrypted := make([]byte, len(data))
	size := block.BlockSize()

	for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
		block.Decrypt(decrypted[bs:be], data[bs:be])
	}

	return PKCS5Unpadding(decrypted)
}

//Sha256Sign
func Sha256Sign(noSign string) string {
	bs := []byte(noSign)
	h := sha256.New()
	h.Write(bs)
	return hex.EncodeToString(h.Sum(nil))
}

//MD5
func MD5(noSign string) string {
	md5Contain := md5.New()
	md5Contain.Write([]byte(noSign))
	return hex.EncodeToString(md5Contain.Sum(nil))
}

//HmacSha256
func ComputeHmacSha256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}
