package repository

import (
	"errors"
	"fmt"
	"lock/database"
	"lock/domain"
	"lock/models"

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

func SignupInsert(user models.SignupDetail) (models.SignupDetailResponse, error) {

	var signupRes models.SignupDetailResponse

	err := database.DB.Raw("INSERT INTO users(firstname, lastname, email, phone, password) VALUES (?,?,?,?,?) RETURNING id,firstname,lastname,email,phone", user.FirstName, user.LastName, user.Email, user.Phone, user.Password).Scan(&signupRes).Error

	if err != nil {
		fmt.Println(err, "asd")
		return models.SignupDetailResponse{}, err
	}

	fmt.Println("signup inserted data's are :", signupRes)
	return signupRes, nil

}
