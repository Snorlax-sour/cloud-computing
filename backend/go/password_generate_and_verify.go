package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func generateSalt() string {
	b := make([]byte, 16) // 16 bytes salt, 32 byte to long!! to cause  Error generating hash and salt: bcrypt: password length exceeds 72 bytes
	_, err := rand.Read(b)
	if err != nil {
		log.Fatalf("Error generating salt: %v", err) // 若生成失敗，直接結束程式
	}
	return hex.EncodeToString(b) // 將 Salt 編碼為十六進位字串
}

func hashPassword(password string) (string, string, error) {
	// Step 1: Generate Salt
	salt := generateSalt()
	if salt == "" {
		return "", "", errors.New("generated salt is empty") // 返回明確錯誤
	}

	// Step 2: Combine Password + Salt
	passwordWithSalt := password + salt
	if passwordWithSalt == "" {
		return "", "", errors.New("password with salt is empty")
	}

	// Step 3: Generate Hash with bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordWithSalt), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error generating bcrypt hash:", err)
		return "", "", err
	}

	return string(hashedPassword), salt, nil // 返回雜湊密碼與 Salt
}
func verifyPassword(password string, storedHash string, salt string) bool {
	if password == "" || storedHash == "" || salt == "" {
		log.Println("empty input with password: ", password, " or storeHash: ", storedHash, " or salt value", salt)
		return false
	}
	passwordWithSalt := password + salt
	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(passwordWithSalt))
	return err == nil // 如果無錯誤，則密碼匹配
}