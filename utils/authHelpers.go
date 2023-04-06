package utils

import (
	"bufio"
	"math/rand"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type AuthHelpersI interface {
	GetRandomWords(numOfWords int) (string, error)
	GetRandomPassword(length int) string
	HashPasswords(password string) (string, error)
	CheckPassword(password, hash string) bool
}

type AuthHelpers struct{}

func NewAuthHelper() *AuthHelpers {
	return &AuthHelpers{}
}

func (a *AuthHelpers) GetRandomWords(numOfWords int) (string, error) {
	file, err := os.Open("words.txt")
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	var resultSlice []string
	for i := 0; i < numOfWords; i++ {
		resultSlice = append(resultSlice, lines[rand.Intn(len(lines))])
	}

	return strings.Join(resultSlice, "."), nil
}

func (a *AuthHelpers) GetRandomPassword(length int) string {
	letterRunes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_+"
	ranSrt := make([]byte, length)

	for i := 0; i < length; i++ {
		ranSrt[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(ranSrt)
}

func (a *AuthHelpers) HashPasswords(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (a *AuthHelpers) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
