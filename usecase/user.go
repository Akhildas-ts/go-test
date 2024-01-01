package usecase

import (
	"errors"
	"fmt"
	"lock/helper"
	"lock/models"
	"lock/repository"
	"net/mail"

	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
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
	fmt.Println("inseted data are", dataInsert)
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
	}, nil

}

func LoginUser(login models.LoginDetails) (*models.TokenUser, error) {

	_, err := mail.ParseAddress(login.Email)

	if err != nil {
		return &models.TokenUser{}, errors.New("email should be correct formate")
	}

	email, err := repository.CheckingEmailValidation(login.Email)

	if err != nil {
		return &models.TokenUser{}, err
	}

	if email == nil {
		return &models.TokenUser{}, errors.New("email is not found")
	}

	userDetails, err := repository.FindUserDetailByEmail(login)

	if err != nil {
		return &models.TokenUser{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(userDetails.Password), []byte(login.Passoword))
	if err != nil {

		return &models.TokenUser{}, errors.New("password not matching")

	}

	var user_details models.SignupDetailResponse
	fmt.Println("user detials is ",userDetails)

	err = copier.Copy(&user_details, userDetails)

	if err != nil{
		return &models.TokenUser{},err
	}
	fmt.Println("user_detials ",user_details)

	accessToken, err := helper.GenerateAccessToken(user_details)

	if err != nil {
		return &models.TokenUser{}, err
	}

	refreshToken, err := helper.GenerateRefreshToken(user_details)

	if err != nil {
		return &models.TokenUser{}, err
	}

	return &models.TokenUser{
		Users:       user_details,
		AccesToken:  accessToken,
		RefresToken: refreshToken,
	}, nil

}
