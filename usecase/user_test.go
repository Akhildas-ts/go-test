package usecase

import (
	"errors"
	"lock/domain"
	"lock/mocks"
	"lock/models"
	"testing"

	"github.com/golang/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestUserLogged(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepo(ctrl)
	uc := UserUsecase{UserRepo: mockRepo}

	testUser := models.LoginDetails{
		Email:     "akhil@gmail.com",
		Passoword: "akhil@123",
	}

	t.Run("successful login", func(t *testing.T) {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("akhil@123"), bcrypt.DefaultCost)

		mockRepo.EXPECT().CheckingEmailValidation(testUser.Email).Return(&domain.User{Email: testUser.Email}, nil)
		mockRepo.EXPECT().FindUserDetailByEmail(testUser).Return(models.UserLoginResponse{
			Email:    testUser.Email,
			Password: string(hashedPassword), // Use the generated hashed password here
		}, nil)

		// Simulate bcrypt comparison with the correct password
		if err := bcrypt.CompareHashAndPassword([]byte(string(hashedPassword)), []byte(testUser.Passoword)); err != nil {
			t.Fatalf("passwords do not match: %v", err)
		}

		tokenUser, err := uc.UserLogged(testUser)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if tokenUser.Users.Email != testUser.Email {
			t.Errorf("expected email %v, got %v", testUser.Email, tokenUser.Users.Email)
		}
	})

	t.Run("invalid email format", func(t *testing.T) {
		invalidUser := models.LoginDetails{
			Email:     "invalid-email",
			Passoword: "somepassword",
		}
		_, err := uc.UserLogged(invalidUser)
		if err == nil || err.Error() != "EMAIL SHOULD BE CORRECT FORMAT " {
			t.Fatalf("expected error 'EMAIL SHOULD BE CORRECT FORMAT ', got %v", err)
		}
	})

	t.Run("email not found", func(t *testing.T) {
		mockRepo.EXPECT().CheckingEmailValidation(testUser.Email).Return(nil, nil)

		_, err := uc.UserLogged(testUser)
		if err == nil || err.Error() != "eroor email not found " {
			t.Fatalf("expected error 'eroor email not found ', got %v", err)
		}
	})

	t.Run("incorrect password", func(t *testing.T) {
		mockRepo.EXPECT().CheckingEmailValidation(testUser.Email).Return(&domain.User{Email: testUser.Email}, nil)
		mockRepo.EXPECT().FindUserDetailByEmail(testUser).Return(models.UserLoginResponse{
			Email:    testUser.Email,
			Password: "$2a$10$hashedpassword", // Use a valid hashed password for testing
		}, nil)

		// Simulate bcrypt comparison to fail
		err := bcrypt.CompareHashAndPassword([]byte("$2a$10$hashedpassword"), []byte("wrongpassword"))
		if err == nil {
			t.Fatalf("expected error, got none")
		}

		_, err = uc.UserLogged(testUser)
		if err == nil || err.Error() != "hassed password not matching" {
			t.Fatalf("expected error 'hassed password not matching', got %v", err)
		}
	})

	t.Run("error during finding user detail", func(t *testing.T) {
		mockRepo.EXPECT().CheckingEmailValidation(testUser.Email).Return(&domain.User{Email: testUser.Email}, nil)
		mockRepo.EXPECT().FindUserDetailByEmail(testUser).Return(models.UserLoginResponse{}, errors.New("error fetching user details"))

		_, err := uc.UserLogged(testUser)
		if err == nil || err.Error() != "error fetching user details" {
			t.Fatalf("expected error 'error fetching user details', got %v", err)
		}
	})
}

func TestUsersingUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepo(ctrl)
	uc := UserUsecase{UserRepo: mockRepo}

	testUser := models.SignupDetail{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Password:  "securepassword",
		Phone:     "1234567890",
	}

	t.Run("successful signup", func(t *testing.T) {
		mockRepo.EXPECT().CheckingEmailValidation(testUser.Email).Return(nil, nil)
		mockRepo.EXPECT().CheckingPhoneExists(testUser.Phone).Return(nil, nil)
		mockRepo.EXPECT().SignupInsert(gomock.Any()).Return(models.SignupDetailResponse{ID: 1, FirstName: "John", LastName: "Doe", Email: "john.doe@example.com", Phone: "1234567890"}, nil)

		tokenUser, err := uc.UsersingUp(testUser)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if tokenUser.Users.Email != testUser.Email {
			t.Errorf("expected email %v, got %v", testUser.Email, tokenUser.Users.Email)
		}
	})

	t.Run("email already exists", func(t *testing.T) {
		// Create a dummy domain.User to simulate an existing user
		existingUser := &domain.User{
			Email: testUser.Email,
		}

		// Mock the CheckingEmailValidation method to return the existing user
		mockRepo.EXPECT().CheckingEmailValidation(testUser.Email).Return(existingUser, nil)

		_, err := uc.UsersingUp(testUser)
		if err == nil || err.Error() != "email is already exisit " {
			t.Fatalf("expected error 'email is already exisit ', got %v", err)
		}
	})

	t.Run("phone already exists", func(t *testing.T) {
		mockRepo.EXPECT().CheckingEmailValidation(testUser.Email).Return(nil, nil)

		// Create a dummy domain.User to simulate an existing user with the same phone
		existingUser := &domain.User{
			Phone: testUser.Phone,
		}

		// Mock the CheckingPhoneExists method to return the existing user
		mockRepo.EXPECT().CheckingPhoneExists(testUser.Phone).Return(existingUser, nil)

		_, err := uc.UsersingUp(testUser)
		if err == nil || err.Error() != "phone number is already exist " {
			t.Fatalf("expected error 'phone number is already exist ', got %v", err)
		}
	})

	t.Run("error during signup insert", func(t *testing.T) {
		mockRepo.EXPECT().CheckingEmailValidation(testUser.Email).Return(nil, nil)
		mockRepo.EXPECT().CheckingPhoneExists(testUser.Phone).Return(nil, nil)
		mockRepo.EXPECT().SignupInsert(gomock.Any()).Return(models.SignupDetailResponse{}, errors.New("could not add User "))

		_, err := uc.UsersingUp(testUser)
		if err == nil || err.Error() != "could not add User " {
			t.Fatalf("expected error 'could not add User ', got %v", err)
		}
	})
}
