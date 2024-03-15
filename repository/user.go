package repository

import (
	"errors"
	"fmt"
	"lock/database"
	"lock/domain"

	"gorm.io/gorm"
)

func CheckingEmailValidation(email string) (*domain.User, error) {

	var user domain.User

	result := database.DB.Where(&domain.User{Email: email}).First(&user)

	if result.Error != nil {

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	fmt.Println("user details was :", user)
	return &user, nil

}

func ChechingPhoneExist(phone string) (*domain.User, error) {
	var user domain.User
	result := database.DB.Where(&domain.User{Phone: phone}).First(&user)
	if result.Error != nil {

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &user, nil
}
