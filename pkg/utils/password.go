package utils

import "golang.org/x/crypto/bcrypt"

type Password struct {
	password string
}

func NewPassword(p string) *Password {
	return &Password{
		password: p,
	}
}

func (p *Password) Hash() (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(p.password), 14)

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (p *Password) Compare(h string) bool {
	return bcrypt.CompareHashAndPassword([]byte(h), []byte(p.password)) == nil
}
