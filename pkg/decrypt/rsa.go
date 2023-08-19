package decrypt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"github.com/gucooing/bdstobot/pkg/logger"
	"io/ioutil"
	"os"
)

var ec2b []byte

func Protoxor(content, sign string) string {
	newcontent, _ := base64.StdEncoding.DecodeString(content)
	newsign, _ := base64.StdEncoding.DecodeString(sign)
	rsadata := Rsaen(newcontent, newsign)
	logger.Debug("rsa解密结果:", rsadata)
	// 读取 EC2B 客户端首包密钥
	var err error
	ec2b, err = os.ReadFile("data/ec2b.bin")
	if err != nil {
		logger.Error("读取ec2b错误:", err)
		return "读取ec2b错误"
	}
	logger.Debug("使用的ec2b是:", base64.StdEncoding.EncodeToString(ec2b))
	//解密结果的异或
	newrsadata, _ := base64.StdEncoding.DecodeString(rsadata)
	Xorec2b(newrsadata)
	logger.Debug("逆异或的结果是:", string(newrsadata))
	newnewrsadata := base64.StdEncoding.EncodeToString(newrsadata)
	return newnewrsadata
}

func Xorec2b(data []byte) {
	for i := 0; i < len(data); i++ {
		data[i] ^= ec2b[i%4096]
	}
}

func Rsaen(content, sign []byte) string {
	// 读取私钥文件
	privateKeyFile, err := ioutil.ReadFile("data/private.pem")
	if err != nil {
		logger.Error("读取私钥文件失败:", err)
		return "读取私钥文件失败"
	}
	// 解析私钥
	block, _ := pem.Decode(privateKeyFile)
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		logger.Error("解析私钥失败:", err)
		return "解析私钥失败,密钥格式错误"
	}

	// 读取公钥文件
	publicKeyFile, err := ioutil.ReadFile("data/public.pem")
	if err != nil {
		logger.Error("读取公钥文件失败:", err)
		return "读取公钥文件失败"
	}
	// 解析公钥
	block, _ = pem.Decode(publicKeyFile)
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		logger.Error("解析公钥失败:", err)
		return "解析公钥失败，格式错误"
	}
	publicKey := publicKeyInterface.(*rsa.PublicKey)

	// 使用私钥进行解密
	decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey.(*rsa.PrivateKey), content)
	if err != nil {
		logger.Error("解密失败:", err)
		return "解密失败"
	}
	logger.Debug("私钥解密结果是:", decrypted)

	// 计算待验证签名的哈希值
	hashed := sha256.Sum256(decrypted)
	logger.Debug("sing的哈希值是:", hashed)

	// 使用公钥进行验证签名
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], sign)
	if err != nil {
		logger.Warn("签名验证失败:", err)
		return "签名验证失败"
	}
	data := base64.StdEncoding.EncodeToString(decrypted)
	return data
}
