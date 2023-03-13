package model

import "golang.org/x/crypto/bcrypt"

type User struct {
	Login        string `json:"login" db:"login"`
	Password     string `json:"password"`
	PasswordHash string `json:"password_hash,omitempty" db:"password_hash"`
}

func (u *User) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	if err != nil {
		return err
	}
	u.PasswordHash = string(bytes)
	return nil
}

func (u *User) CheckPasswordHash(registeredUser *User) bool {
	err := bcrypt.CompareHashAndPassword([]byte(registeredUser.PasswordHash), []byte(u.Password))
	return err == nil
}
