package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/junaid9001/tripneo/auth-service/db"
	"github.com/junaid9001/tripneo/auth-service/models"
	"gorm.io/gorm"
)

var ErrEmailALreadyTaken = errors.New("email already taken")

func CreateUser(email, hashedPassword string) error {

	user := &models.User{
		Email:        email,
		PasswordHash: hashedPassword,
		CreatedAt:    time.Now(),
	}

	if err := db.DB.Create(&user).Error; err != nil {

		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrEmailALreadyTaken
		}

		return fmt.Errorf("Internal Server Error")
	}

	return nil
}
