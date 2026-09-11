package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
)

func GenerateRsaKeyPair() (privateKeyPem string, publicKeyPem string, err error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", err
	}
	privBytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return "", "", err
	}
	privBlock := pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privBytes,
	}
	privateKeyPem = string(pem.EncodeToMemory(&privBlock))
	pubBytes, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		return "", "", err
	}
	pubBlock := pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	}
	publicKeyPem = string(pem.EncodeToMemory(&pubBlock))
	return privateKeyPem, publicKeyPem, nil
}

func EncryptWithPublicKey(data string, pubKey string) {
	block, _ := pem.Decode([]byte(pubKey))
	if block == nil {
		log.Fatal(errors.New("failed to parse PEM block containing the public key"))
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}
	rsaKey, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		log.Fatal(errors.New("key type is not RSA public key"))
	}
	result, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaKey, []byte(data), nil)
	if err != nil {
		log.Fatal(err)
	}
	encoded := base64.StdEncoding.EncodeToString(result)
	fmt.Println(encoded)
}

func DecryptWithPrivateKey(data []byte, privateKeyPem string) ([]byte, error) {
	block, _ := pem.Decode([]byte(privateKeyPem))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the private key")
	}

	priInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	priKey, ok := priInterface.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("key type is not RSA private key")
	}

	// 使用 OAEP 填充模式解密
	return rsa.DecryptPKCS1v15(rand.Reader, priKey, data)
}
