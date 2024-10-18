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

type UserUsecase struct {
	UserRepo *repository.UserRepoImpl
}

func UseruseCase(repo repository.UserRepoImpl) *UserUsecase {
	return &UserUsecase{UserRepo: &repo}
}

func (uc *UserUsecase) UserLogged(user models.LoginDetails) (*models.TokenUser, error) {

	_, err := mail.ParseAddress(user.Email)
	if err != nil {

		return &models.TokenUser{}, errors.New("EMAIL SHOULD BE CORRECT FORMAT ")

	}

	email, err := uc.UserRepo.CheckingEmailValidation(user.Email)

	if err != nil {
		return &models.TokenUser{}, errors.New("SERVER ERROR  from : checking-email-validation")

	}

	if email == nil {
		return &models.TokenUser{}, errors.New("eroor email not found ")

	}

	userDetails, err := uc.UserRepo.FindUserDetailByEmail(user)
	if err != nil {
		return &models.TokenUser{}, err

	}

	// CHECKING THE HASSED PASSWORD

	err = bcrypt.CompareHashAndPassword([]byte(userDetails.Password), []byte(user.Passoword))
	if err != nil {
		fmt.Println("HEyyyy")
		return &models.TokenUser{}, errors.New("hassed password not matching")
	}
	var user_details models.SignupDetailResponse

	err = copier.Copy(&user_details, &userDetails)
	if err != nil {

		return &models.TokenUser{}, errors.New("eorro in copier")
	}

	// TOKEN....

	accessToken, err := helper.GenerateAccessToken(user_details)
	if err != nil {
		return &models.TokenUser{}, errors.New("could not create accesstoken due to internal error")
	}
	refreshToken, err := helper.GenerateRefreshToken(user_details)
	if err != nil {
		return &models.TokenUser{}, errors.New("counldnt create refresh token due to error")
	}

	return &models.TokenUser{
		Users:       user_details,
		AccesToken:  accessToken,
		RefresToken: refreshToken,
	}, nil
}

func (uc *UserUsecase) UsersingUp(user models.SignupDetail) (*models.TokenUser, error) {

	email, err := uc.UserRepo.CheckingEmailValidation(user.Email)

	if err != nil {

		return &models.TokenUser{}, errors.New("error with the singup server")

	}

	if email != nil {

		return &models.TokenUser{}, errors.New("email is already exisit ")
	}

	phone, err := uc.UserRepo.CheckingPhoneExists(user.Phone)

	if err != nil {

		return &models.TokenUser{}, errors.New("p_server have issue ")

	}
	if phone != nil {
		return &models.TokenUser{}, errors.New("phone number is already exist ")

	}

	//   Passoword Hash

	hashedPassword, err := helper.PasswordHasing(user.Password)

	if err != nil {
		return &models.TokenUser{}, errors.New("hash_server have issue")
	}

	user.Password = hashedPassword

	dataInsert, err := uc.UserRepo.SignupInsert(user)
	if err != nil {
		return &models.TokenUser{}, errors.New("could not add User ")
	}

	fmt.Println("data inserted in signup :", dataInsert)

	// CREATING A JWT TOKEN FOR THE NEW USER\\

	accessToken, err := helper.GenerateAccessToken(dataInsert)

	if err != nil {
		return &models.TokenUser{}, errors.New("can't create a acces token")
	}

	refershToken, err := helper.GenerateRefreshToken(dataInsert)

	if err != nil {
		return &models.TokenUser{}, errors.New("can't create a Refersh token")

	}

	return &models.TokenUser{
		Users:       dataInsert,
		AccesToken:  accessToken,
		RefresToken: refershToken,
	}, nil

}
