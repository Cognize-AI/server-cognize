package models

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"

	"github.com/Cognize-AI/server-cognize/config"
	"gorm.io/gorm"
)

type Key struct {
	gorm.Model
	Name   string
	Value  string `gorm:"unique"`
	Hash   string
	UserID uint `gorm:"index"`

	User User `gorm:"foreignKey:UserID;references:ID"`
}

var encryptionKey []byte

func init() {
	_config, err := config.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	encryptionKey = []byte(_config.EncSecret)
}

func (k *Key) BeforeCreate(tx *gorm.DB) (err error) {
	hash := sha256.Sum256([]byte(k.Value))
	k.Hash = hex.EncodeToString(hash[:])

	encrypted, err := encryptValue(k.Value)
	if err != nil {
		return err
	}
	k.Value = encrypted
	return
}

func (k *Key) BeforeUpdate(tx *gorm.DB) (err error) {
	encrypted, err := encryptValue(k.Value)
	if err != nil {
		return err
	}
	k.Value = encrypted
	return
}

func (k *Key) AfterFind(tx *gorm.DB) (err error) {
	decrypted, err := decryptValue(k.Value)
	if err != nil {
		return err
	}
	k.Value = decrypted
	return
}

func (k *Key) AfterCreate(tx *gorm.DB) (err error) {
	decrypted, err := decryptValue(k.Value)
	if err != nil {
		return err
	}
	k.Value = decrypted
	return nil
}

func encryptValue(plainText string) (string, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plainText))
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plainText))

	return hex.EncodeToString(ciphertext), nil
}

func decryptValue(encryptedHex string) (string, error) {
	ciphertext, err := hex.DecodeString(encryptedHex)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}
