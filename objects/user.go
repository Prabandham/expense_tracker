package objects

import (
	"github.com/jinzhu/gorm"
	crypt "golang.org/x/crypto/bcrypt"
)

// User model
type User struct {
	Base
	Email             string `gorm:"index;unique;size:256" json:"email" sql:"not null"`
	Name              string `gorm:"index;size:256" json:"name" sql:"not null"`
	Password          string `gorm:"-" json:"password"`
	EncryptedPassword string `json:"-"`
}

// Before we save a user we check if the password is present
// if present we will hash it and save it.
func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	if u.Password != "" {
		encryptedPassword, err := hashPassword(u.Password)

		if err != nil {
			return err
		}

		u.EncryptedPassword = encryptedPassword
	}
	return
}

func (u *User) ValidatePassword(password string) bool {
	return doPasswordsMatch(u.EncryptedPassword, password)
}

// Hash password using the Crypt hashing algorithm
// and then return the hashed password as a
// base64 encoded string
func hashPassword(password string) (string, error) {
	var passwordBytes = []byte(password)
	hashedPasswordBytes, err := crypt.GenerateFromPassword(passwordBytes, crypt.MinCost)
	return string(hashedPasswordBytes), err
}

// Check if two passwords match using crypt's CompareHashAndPassword
// which return nil on success and an error on failure.
func doPasswordsMatch(hashedPassword, currentPassword string) bool {
	err := crypt.CompareHashAndPassword(
		[]byte(hashedPassword), []byte(currentPassword))
	return err == nil
}
