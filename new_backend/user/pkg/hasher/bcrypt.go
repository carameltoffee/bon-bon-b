package hasher

import (
	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct {
	cost int
}

func NewBcryptHasher(cost ...int) *BcryptHasher {
	c := bcrypt.DefaultCost
	if len(cost) > 0 {
		c = cost[0]
	}
	return &BcryptHasher{cost: c}
}

func (b *BcryptHasher) Hash(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), b.cost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func (b *BcryptHasher) Compare(password, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}
