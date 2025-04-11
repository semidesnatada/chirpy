package user

import "golang.org/x/crypto/bcrypt"

type User struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    password string // private field, won't be included in JSON
    loginAttempts int // private field for tracking login attempts
}

// NewUser creates a new user with a secure password
func NewUser(username, email, initialPassword string) *User {
    return &User{
        Username: username,
        Email: email,
        password: hashPassword(initialPassword), // private field initialized
        loginAttempts: 0,
    }
}

// VerifyPassword checks if the provided password is correct
func (u *User) VerifyPassword(attempt string) bool {
    if isCorrect := checkPassword(attempt, u.password); !isCorrect {
        u.loginAttempts++ // update internal state
        return false
    }
    u.loginAttempts = 0 // reset attempts on success
    return true
}

// IsLocked checks if the account is locked due to too many failed attempts
func (u *User) IsLocked() bool {
    return u.loginAttempts >= 5
}

// hashPassword converts a plaintext password into a secure hash
func hashPassword(password string) string {
    // In a real implementation, you'd handle errors properly
    hashedBytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(hashedBytes)
}

// checkPassword verifies if the attempt matches the stored password hash
func checkPassword(attempt string, hashedPassword string) bool {
    // CompareHashAndPassword returns nil on success, error on failure
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(attempt))
    return err == nil
}