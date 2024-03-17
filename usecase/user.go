package usecase

import (
	"errors"
	"fmt"
	"lock/helper"
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

	hashPassword, err := helper.PasswordHasing(user.Password)

	if err != nil {

		return &models.TokenUser{}, errors.New("hash password issue")
	}

	user.Password = hashPassword

	dataInsert, err := repository.SignupInsert(user)

	if err != nil {
		return &models.TokenUser{}, errors.New("cloud not add user")
	}
	fmt.Println("inseted data are",dataInsert)
	refresh, err := helper.GenerateRefreshToken(dataInsert)
	
		if err != nil {
			return &models.TokenUser{}, err
		}
	accessToken, err := helper.GenerateAccessToken(dataInsert)

	if err != nil {
		return &models.TokenUser{}, errors.New("issue from acces token ")

	}


	return &models.TokenUser{
		Users:       dataInsert,
		AccesToken:  accessToken,
		RefresToken: refresh,
	},nil

}
