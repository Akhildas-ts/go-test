package interfaces

import (
	"lock/domain"
	"lock/models"
)

type UserRepo interface {
	CheckingEmailValidation(email string) (*domain.User, error)
	CheckingPhoneExists(phone string) (*domain.User, error)
	SignupInsert(user models.SignupDetail) (models.SignupDetailResponse, error)
}
