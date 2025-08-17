package services

import (
	"crypto/rand"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"goapi/models"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/argon2"
)

func HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("password cannot be empty")
	}

	// Generate random 16-byte salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Argon2id: (time=1, memory=64MB, threads=4, keyLen=32)
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	// Return "salt$hash" (both base64)
	return fmt.Sprintf("%s$%s",
		base64.StdEncoding.EncodeToString(salt),
		base64.StdEncoding.EncodeToString(hash),
	), nil
}

// Save user in DB (signup)
func SignUp(db *sql.DB, user *models.User) error {
	// validate first
	if err := user.Validate(); err != nil {
		return err
	}

	hashed, err := HashPassword(user.Password)
	if err != nil {
		return err
	}

	// Insert into DB
	_, err = db.Exec(`INSERT INTO users (username, password) VALUES ($1, $2)`, user.Username, hashed)
	return err
}

// Verify user (signin)
func SignIn(db *sql.DB, username, password string) (bool, error, string) {
	var stored string
	err := db.QueryRow(`SELECT password FROM users WHERE username=$1`, username).Scan(&stored)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, errors.New("user not found"), ""
		}
		return false, fmt.Errorf("database query failed: %w", err), ""
	}

	// Split salt and hash
	parts := strings.Split(stored, "$")
	if len(parts) != 2 {
		return false, errors.New("invalid stored password format"), ""
	}
	saltB64, hashB64 := parts[0], parts[1]

	// Decode Base64
	salt, err := base64.StdEncoding.DecodeString(saltB64)
	if err != nil {
		return false, errors.New("invalid salt encoding"), ""
	}
	storedHash, err := base64.StdEncoding.DecodeString(hashB64)
	if err != nil {
		return false, errors.New("invalid hash encoding"), ""
	}

	// Re-hash input password
	newHash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	// Constant-time compare
	if subtle.ConstantTimeCompare(newHash, storedHash) != 1 {
		// Return immediately if passwords do not match
		return false, errors.New("invalid password"), ""
	}

	// Generate and return token ONLY on a successful password match
	token, err := JwtTokenGenerate("assaa", "mohammed", "user")
	if err != nil {
		// Return a clear error if token generation fails
		return false, fmt.Errorf("failed to generate token: %w", err), ""
	}

	// Return success and the generated token
	return true, nil, token
}

func JwtTokenGenerate(userid, username, userrole string) (string, error) {
	jwtsecret := "jwtsecretstring"

	claims := jwt.MapClaims{
		"uid":  userid,
		"user": username,
		"role": userrole,
		"exp":  jwt.NewNumericDate(time.Now().Add(15 * time.Minute)), // ✅ FIXED
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(jwtsecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
