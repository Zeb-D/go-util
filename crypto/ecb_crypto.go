package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

func AesECBDecrypt(crypted, key []byte) (origData []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("err is:", err)
		return
	}
	blockMode := NewECBDecrypter(block)
	origData = make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return
}

func AesECBEncrypt(src, key []byte) (crypted []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		err = fmt.Errorf("failed to new cipher, error: %v", err.Error())
		return
	}
	if len(src) == 0 {
		err = fmt.Errorf("plain content empty")
		return
	}
	ecb := NewECBEncrypter(block)
	content := src
	content = PKCS5Padding(content, block.BlockSize())
	crypted = make([]byte, len(content))
	ecb.CryptBlocks(crypted, content)
	return
}

//
//func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
//	padding := blockSize - len(ciphertext)%blockSize
//	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
//	return append(ciphertext, padtext...)
//}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

type ecb struct {
	b         cipher.Block
	blockSize int
}

func newECB(b cipher.Block) *ecb {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

type ecbEncrypter ecb

// NewECBEncrypter returns a BlockMode which encrypts in electronic code book
// mode, using the given Block.
func NewECBEncrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbEncrypter)(newECB(b))
}
func (x *ecbEncrypter) BlockSize() int {
	return x.blockSize
}
func (x *ecbEncrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("aes_crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("aes_crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Encrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

type ecbDecrypter ecb

// NewECBDecrypter returns a BlockMode which decrypts in electronic code book
// mode, using the given Block.
func NewECBDecrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbDecrypter)(newECB(b))
}
func (x *ecbDecrypter) BlockSize() int {
	return x.blockSize
}
func (x *ecbDecrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("aes_crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("aes_crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Decrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}
