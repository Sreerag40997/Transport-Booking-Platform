package service

import (
	"errors"

	"github.com/junaid9001/tripneo/auth-service/repository"
	"github.com/junaid9001/tripneo/auth-service/utils"
)

var EmailAlreadyTaken = errors.New("email already taken")

func CreateUser(email, password string) error {
	hashedPass, err := utils.GenerateHashedPassword(password)
	if err != nil {
		return errors.New("internal server error")
	}
	err = repository.CreateUser(email, hashedPass)
	if err != nil {
		if errors.Is(err, repository.ErrEmailALreadyTaken) {
			return EmailAlreadyTaken
		}
		return errors.New("internal server error")
	}
	return nil
}
