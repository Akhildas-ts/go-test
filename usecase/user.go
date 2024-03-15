package usecase

import (
	"errors"
	"lock/models"
	"lock/repository"
)

func UsersignUp(user models.SignupDetail) (*models.TokenUser, error) {

	email, err := repository.CheckingEmailValidation(user.Email)

	if err != nil {
		return &models.TokenUser{}, errors.New("error with the singup server")
	}

	if email != nil {
		return &models.TokenUser{}, errors.New("email is already exist ")
	}

	phone, err := repository.ChechingPhoneExist(user.Phone)

	if err != nil {
		return &models.TokenUser{}, errors.New("error with the phone number")
	}

	if phone != nil {
		return &models.TokenUser{}, errors.New("phone number is already exist")
	}

	hashPassword,err := helper.Passwor

}
