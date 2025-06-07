package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

//	func GenerateKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
//		privateKey, err := rsa.GenerateKey(rand.Reader, bits)
//		if err != nil {
//			return nil, nil, err
//		}
//		return privateKey, &privateKey.PublicKey, nil
//	}
//
// Натянуть это на GenKeyPair
func GenKeyPair() error {
	// Генерация ключей
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return fmt.Errorf("failed to generate key pair: %w", err)
	}

	// Сохранение приватного ключа (PKCS1)
	privateFile, err := os.Create("/tmp/key.pem")
	if err != nil {
		return fmt.Errorf("failed to create private key file: %w", err)
	}
	defer privateFile.Close()

	privateBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	if err = pem.Encode(privateFile, privateBlock); err != nil {
		return fmt.Errorf("failed to write private key: %w", err)
	}

	// Сохранение публичного ключа (PKCS1)
	publicFile, err := os.Create("/tmp/cert.pem")
	if err != nil {
		return fmt.Errorf("failed to create public key file: %w", err)
	}
	defer publicFile.Close()

	publicBlock := &pem.Block{
		Type:  "RSA PUBLIC KEY", // Правильный тип для PKCS1
		Bytes: x509.MarshalPKCS1PublicKey(&privateKey.PublicKey),
	}
	if err = pem.Encode(publicFile, publicBlock); err != nil {
		return fmt.Errorf("failed to write public key: %w", err)
	}

	return nil
}

func LoadPrivateKey(filename string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func LoadPublicKey(filename string) (*rsa.PublicKey, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	return x509.ParsePKCS1PublicKey(block.Bytes)
}

func Encrypt(publicKey *rsa.PublicKey, data []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, publicKey, data)
}

func Decrypt(privateKey *rsa.PrivateKey, data []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, privateKey, data)
}
