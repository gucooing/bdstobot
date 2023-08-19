package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"github.com/gucooing/bdstobot/config"
	"github.com/gucooing/bdstobot/pkg/logger"
	jsoniter "github.com/json-iterator/go"
)

var newkeymd5 []byte

type Param struct {
	Mode string `json:"mode"`
	Raw  string `json:"raw"`
}

type encryptPkt struct {
	Type   string `json:"type"`
	Params Param  `json:"params"`
}

// Encrypt 加密函数入口
func Encrypt_send(str string) []byte {
	en, err := AESBase64Encrypt(str)
	if err != nil {
		fmt.Println(err)
	}
	pkt := encryptPkt{Params: Param{Mode: "aes_cbc_pck7padding", Raw: en}, Type: "encrypted"}
	jpkt, _ := jsoniter.Marshal(pkt)
	logger.Debug().Msgf("aes加密结果：%d", string(jpkt))
	return jpkt

}

func AESBase64Encrypt(origin_data string) (base64_result string, err error) {
	keymd5 := md5.Sum([]byte(config.GetConfig().Key))
	newkeymd5 = []byte(fmt.Sprintf("%X", keymd5))
	key := newkeymd5[:16]
	iv := newkeymd5[16:32]

	var block cipher.Block
	if block, err = aes.NewCipher(key); err != nil {
		return
	}
	encrypt := cipher.NewCBCEncrypter(block, iv)
	var source = PKCS5Padding([]byte(origin_data), 16)
	var dst = make([]byte, len(source))
	encrypt.CryptBlocks(dst, source)
	base64_result = base64.StdEncoding.EncodeToString(dst)
	return
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(data []byte, blocklen int) ([]byte, error) {
	if blocklen <= 0 {
		return nil, fmt.Errorf("invalid blocklen %d", blocklen)
	}
	if len(data)%blocklen != 0 || len(data) == 0 {
		return nil, fmt.Errorf("invalid data len %d", len(data))
	}
	padlen := int(data[len(data)-1])
	if padlen > blocklen || padlen == 0 {
		return nil, fmt.Errorf("invalid padding")
	}
	pad := data[len(data)-padlen:]
	for i := 0; i < padlen; i++ {
		if pad[i] != byte(padlen) {
			return nil, fmt.Errorf("invalid padding")
		}
	}
	return data[:len(data)-padlen], nil
}
