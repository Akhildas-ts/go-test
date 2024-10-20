package repository

import (
	"errors"
	"fmt"
	"lock/database"
	"lock/domain"
	"lock/models"

	"gorm.io/gorm"
)

type UserRepo interface {
	CheckingEmailValidation(email string) (*domain.User, error)
	CheckingPhoneExists(phone string) (*domain.User, error)
	SignupInsert(user models.SignupDetail) (models.SignupDetailResponse, error)
	FindUserDetailByEmail(user models.LoginDetails) (models.UserLoginResponse, error)
}

type UserRepoImpl struct {
	DB *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepoImpl {
	return &UserRepoImpl{DB: db}
}

func (ur *UserRepoImpl) CheckingEmailValidation(email string) (*domain.User, error) {

	var user domain.User

	result := ur.DB.Where(&domain.User{Email: email}).First(&user)
	fmt.Printf("Executed SQL: %s\n", result.Statement.SQL.String())

	if result.Error != nil {

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {

			return nil, nil

		}
		return nil, result.Error
	}
	fmt.Println("user details was :", user)
	return &user, nil

}

func (ur *UserRepoImpl) CheckingPhoneExists(phone string) (*domain.User, error) {

	var user domain.User

	result := database.DB.Where(&domain.User{Phone: phone}).First(&user)

	if result.Error != nil {

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			fmt.Println("NOT FOUND PHONE NUMBER IN DB  >>> ")
			return nil, nil
		}

		return nil, result.Error
	}
	return &user, nil

}

func (ur *UserRepoImpl) SignupInsert(user models.SignupDetail) (models.SignupDetailResponse, error) {

	var signupRes models.SignupDetailResponse

	err := database.DB.Raw("INSERT INTO users(firstname, lastname, email, phone, password) VALUES (?,?,?,?,?) RETURNING id,firstname,lastname,email,phone", user.FirstName, user.LastName, user.Email, user.Phone, user.Password).Scan(&signupRes).Error

	if err != nil {
		fmt.Println("ERROR FROM INSERTING DATA", err)
		return models.SignupDetailResponse{}, err
	}

	fmt.Println("AScle inserted data's are :", signupRes)
	return signupRes, nil
}

func (ur *UserRepoImpl) FindUserDetailByEmail(user models.LoginDetails) (models.UserLoginResponse, error) {

	var UserDetails models.UserLoginResponse

	err := database.DB.Raw(
		`select * from users where email = ? and blocked = false`, user.Email).Scan(&UserDetails).Error

	if err != nil {
		return models.UserLoginResponse{}, errors.New("got an error fron ! searching users by email")

	}

	return UserDetails, nil
}
