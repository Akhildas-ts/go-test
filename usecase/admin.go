package usecase

import (
	"errors"
	"lock/domain"
	"lock/helper"
	"lock/models"
	"lock/repository"

	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
)

func AdminLogin(adminmodel models.AdminLogin) (domain.TokenAdmin, error) {

	// GETTING DETAILS FROM ADMIN WITH EMAIL

	AdminDetail, err := repository.AdminLogin(adminmodel)
	if err != nil {

		return domain.TokenAdmin{}, errors.New("given mail formate have")

	}

	// password := AdminDetail.Password

	if AdminDetail.Password == "" {
		return domain.TokenAdmin{}, errors.New("error from admin password")

	}

	err = bcrypt.CompareHashAndPassword([]byte(AdminDetail.Password), []byte(adminmodel.Password))

	if err != nil {

		return domain.TokenAdmin{}, models.PasswordIsNotCorrect

	}

	var AdminDetailsResponse models.AdminDetailsResponse

	err = copier.Copy(&AdminDetailsResponse, &AdminDetail)

	if err != nil {
		return domain.TokenAdmin{}, err
	}

	tokenString, err := helper.GenerateTokenAdmin(AdminDetailsResponse)

	if err != nil {

		return domain.TokenAdmin{}, errors.New("demo : repository adminlogin")

	}

	return domain.TokenAdmin{
		Admin: AdminDetailsResponse,
		Token: tokenString,
	}, nil

}
