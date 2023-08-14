package pgk

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/gucooing/bdstobot/config"
	"io"
)

// Encrypt 加密函数
func Encrypt(plaintext []byte) ([]byte, error) {
	key := []byte(config.GetConfig().Key)
	aeskey := key[:16]
	fmt.Println("aeskey:", aeskey)
	block, err := aes.NewCipher(aeskey)
	if err != nil {
		return nil, err
	}

	// PKCS7Padding填充
	blockSize := block.BlockSize()
	plaintext = padding(plaintext, block.BlockSize())

	// 初始化向量IV，长度与blockSize相同
	ciphertext := make([]byte, blockSize+len(plaintext))
	iv := key[len(key)-16:]
	fmt.Println("aesiv:", iv)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	// 执行加密操作
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[blockSize:], plaintext)

	return ciphertext, nil
}

// Decrypt 解密函数
func Decrypt(ciphertext []byte) ([]byte, error) {
	key := []byte(config.GetConfig().Key)
	aeskey := key[:16]
	block, err := aes.NewCipher(aeskey)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	if len(ciphertext) < blockSize {
		return nil, errors.New("密文长度错误")
	}

	iv := key[len(key)-16:]
	ciphertext = ciphertext[blockSize:]

	// 执行解密操作
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	// PKCS7Padding去除填充
	plaintext, err := unpadding(ciphertext)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// PKCS7Padding填充
func padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

// PKCS7Padding去除填充
func unpadding(src []byte) ([]byte, error) {
	length := len(src)
	unpadding := int(src[length-1])
	if unpadding > length {
		return nil, errors.New("PKCS7填充错误")
	}
	return src[:(length - unpadding)], nil
}
