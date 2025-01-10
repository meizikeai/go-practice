package tool

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
)

type Secret struct {
	secret map[string]string
}

func NewSecret() *Secret {
	return &Secret{
		secret: map[string]string{
			"bf5d": "Av7Cc1pfGdKQkJpcJig1Hg==",
		},
	}
}

// 加密 - 服务端加密
func (s *Secret) HandleServiceEncrypt(ak, p string) string {
	result := ""

	cipher := s.HandleEncrypt(ak, p)
	ciphertext, err := hex.DecodeString(cipher)

	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	result = base64.StdEncoding.EncodeToString(ciphertext)
	// fmt.Println(result)

	return result
}

// 解码 - 服务端解码
func (s *Secret) HandleServiceDecrypt(c string) string {
	result := ""

	cipher, err := base64.StdEncoding.DecodeString(c)

	if err != nil {
		fmt.Println(err.Error())
		return result
	}

	ciphertext := hex.EncodeToString(cipher)
	result = s.HandleDecrypt(ciphertext)

	return result
}

func (s *Secret) getSecretKey(k string) []byte {
	result, _ := base64.StdEncoding.DecodeString(s.secret[k])
	return result
}

// PKCS7 padding of data
func (s *Secret) pkcs7Padding(plaintext []byte, blockSize int) []byte {
	padding := blockSize - len(plaintext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(plaintext, padText...)
}

// PKCS7 unpadding of data
func (s *Secret) pkcs7UnPadding(plaintext []byte) []byte {
	length := len(plaintext)
	unpadding := int(plaintext[length-1])
	return plaintext[:(length - unpadding)]
}

// encrypt data using AES-128-CBC
func (s *Secret) EncryptAes128CBC(text string, key []byte) ([]byte, error) {
	// fill in the input data
	plaintext := s.pkcs7Padding([]byte(text), aes.BlockSize)

	// generate initial vector
	iv := generateInitVector(aes.BlockSize)

	// create aes block cipher
	block, err := aes.NewCipher(key)

	if err != nil {
		return nil, err
	}

	// encrypt using CBC mode
	ciphertext := make([]byte, len(plaintext))
	stream := cipher.NewCBCEncrypter(block, iv)

	stream.CryptBlocks(ciphertext, plaintext)

	// append iv to encrypted data
	ciphertext = append(iv, ciphertext...)

	return ciphertext, nil
}

// decrypt data using AES-128-CBC
func (s *Secret) DecryptAes128CBC(ciphertext, key []byte) ([]byte, error) {
	// create aes block cipher
	block, err := aes.NewCipher(key)

	if err != nil {
		return nil, err
	}

	// extract initial vector from the ciphertext
	iv := ciphertext[:block.BlockSize()]
	ciphertext = ciphertext[block.BlockSize():]

	// decrypt using CBC mode
	stream := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))

	stream.CryptBlocks(plaintext, ciphertext)

	// PKCS7 unpadding
	plaintext = s.pkcs7UnPadding(plaintext)

	return plaintext, nil
}

func (s *Secret) HandleEncrypt(ak, p string) string {
	result := ""

	key := s.getSecretKey(ak)
	cipher, err := s.EncryptAes128CBC(p, key)

	if err != nil {
		fmt.Println(err.Error())
		return result
	}

	pos := hexToDec(ak[0:1])

	iv := hex.EncodeToString(cipher[:16])
	ciphertext := hex.EncodeToString(cipher[16:])

	result = fmt.Sprintf("%s%s%s%s", ak, ciphertext[0:pos], iv, ciphertext[pos:])
	// fmt.Println(result)

	return result
}

func (s *Secret) HandleDecrypt(c string) string {
	result := ""

	ak := c[0:4]
	enc := c[4:]

	pos := hexToDec(ak[0:1])

	iv, _ := hex.DecodeString(enc[pos : pos+32])
	cipher, _ := hex.DecodeString(fmt.Sprintf("%s%s", enc[0:pos], enc[pos+32:]))

	ciphertext := make([]byte, 0)
	ciphertext = append(ciphertext, iv...)
	ciphertext = append(ciphertext, cipher...)

	key := s.getSecretKey(ak)
	plaintext, err := s.DecryptAes128CBC(ciphertext, key)

	if err != nil {
		fmt.Println(err.Error())
		return result
	}

	result = string(plaintext)
	// fmt.Println(result)

	return result
}

// 生成随机的 iv (初始向量)
func generateInitVector(size int) []byte {
	iv := make([]byte, size)

	_, err := rand.Read(iv)

	if err != nil {
		return nil
	}

	return iv
}

func hexToDec(value string) int {
	dec, err := strconv.ParseInt(value, 16, 64)

	if err != nil {
		return 0
	}

	return int(dec)
}
