package usecase

import (
	"errors"
	"lock/domain"
	"lock/mocks"
	"lock/models"
	"testing"

	"github.com/golang/mock/gomock"
)

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

// package usecase_test

// import (
// 	"errors"
// 	"fmt"
// 	"lock/models" // Adjust this import based on your project structure
// 	"lock/repository"
// 	"lock/usecase" // Adjust this import based on your project structure
// 	"reflect"
// 	"testing"

// 	"github.com/DATA-DOG/go-sqlmock"
// 	"github.com/stretchr/testify/assert"
// 	"gorm.io/driver/postgres"
// 	"gorm.io/driver/sqlite"
// 	"gorm.io/gorm"
// )

// func TestUserUsecase_UsersingUp(t *testing.T) {
// 	type args struct {
// 		user models.SignupDetail
// 	}
// 	tests := []struct {
// 		name       string
// 		args       args
// 		beforeTest func(mock sqlmock.Sqlmock)
// 		want       *models.TokenUser
// 		wantErr    bool
// 	}{
// 		{
// 			name: "fail email validation",
// 			args: args{
// 				user: models.SignupDetail{
// 					Email:    "jane@gmail.com",
// 					Password: "password",
// 					Phone:    "1234567890",
// 				},
// 			},
// 			beforeTest: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectQuery("SELECT \\* FROM \"users\" WHERE \"users\"\\.\"email\" = \\$1 AND \"users\"\\.\"deleted_at\" IS NULL ORDER BY \"users\"\\.\"id\" LIMIT \\$2").
// 					WithArgs("jane@gmail.com", 1). // Pass both arguments
// 					WillReturnRows(sqlmock.NewRows([]string{"email"}).AddRow("jane@gmail.com"))
// 			},
// 			wantErr: true,
// 			want:    &models.TokenUser{},
// 		},

// 		{
// 			name: "email already exists",
// 			args: args{
// 				user: models.SignupDetail{
// 					Email:    "jane@example.com",
// 					Password: "password",
// 					Phone:    "1234567890",
// 				},
// 			},
// 			beforeTest: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectQuery("SELECT \\* FROM \"users\" WHERE \"users\"\\.\"email\" = \\$1 AND \"users\"\\.\"deleted_at\" IS NULL ORDER BY \"users\"\\.\"id\" LIMIT \\$2").
// 					WithArgs("jane@example.com", 1).
// 					WillReturnError(errors.New("error checking email"))
// 			},
// 			wantErr: true,
// 			want:    &models.TokenUser{},
// 		},

// 		// {
// 		// 	name: "success user signup",
// 		// 	args: args{
// 		// 		user: models.SignupDetail{
// 		// 			Email:    "john@example.com",
// 		// 			Password: "password",
// 		// 			Phone:    "1234567890",
// 		// 		},
// 		// 	},

// 		// 	beforeTest: func(mock sqlmock.Sqlmock) {
// 		// 		// Mocking the email check
// 		// 		mock.ExpectQuery("SELECT \\* FROM \"users\" WHERE \"users\"\\.\"email\" = \\$1 AND \"users\"\\.\"deleted_at\" IS NULL ORDER BY \"users\"\\.\"id\" LIMIT \\$2").
// 		// 			WithArgs("john@example.com", 1).
// 		// 			WillReturnRows(sqlmock.NewRows([]string{"email"})) // No existing email

// 		// 		// Mocking the phone check
// 		// 		mock.ExpectQuery("SELECT phone FROM \"users\" WHERE phone = \\$1").
// 		// 			WithArgs("1234567890").
// 		// 			WillReturnRows(sqlmock.NewRows([]string{"phone"})) // No existing phone

// 		// 		// Mocking the insert operation
// 		// 		mock.ExpectExec("INSERT INTO users \\(firstname, lastname, email, phone, password\\) VALUES \\(\\?, \\?, \\?, \\?, \\?\\) RETURNING id, firstname, lastname, email, phone").
// 		// 			WithArgs("John", "Doe", "john@example.com", "1234567890", "hashedPassword").
// 		// 			WillReturnResult(sqlmock.NewResult(1, 1)) // Successful insert

// 		// 		// Mocking the retrieval of the user after insertion
// 		// 		mock.ExpectQuery("SELECT id, firstname, lastname, email, phone FROM users WHERE id = \\$1").
// 		// 			WithArgs(1).
// 		// 			WillReturnRows(sqlmock.NewRows([]string{"id", "firstname", "lastname", "email", "phone"}).
// 		// 				AddRow(1, "John", "Doe", "john@example.com", "1234567890")) // Mock response for the inserted user
// 		// 	},

// 		// 	want: &models.TokenUser{
// 		// 		Users: models.SignupDetailResponse{
// 		// 			Email: "john@example.com",
// 		// 			Phone: "1234567890",
// 		// 		},
// 		// 		AccesToken:  "mockAccessToken",  // Replace with actual token generation mock
// 		// 		RefresToken: "mockRefreshToken", // Replace with actual token generation mock
// 		// 	},
// 		// },
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Create a new sqlmock database
// 			db, mock, err := sqlmock.New()
// 			assert.Nil(t, err, "an error '%s' not expected when opening mock database", err)

// 			gormDB, err := gorm.Open(postgres.New(postgres.Config{
// 				Conn: db,
// 			}), &gorm.Config{})
// 			assert.Nil(t, err, "an error '%s' not expected when opening gorm database", err)

// 			// test.buildStub(mock)
// 			ur := repository.UserRepoImpl{DB: gormDB}
// 			uc := &usecase.UserUsecase{UserRepo: &ur}

// 			if tt.beforeTest != nil {
// 				tt.beforeTest(mock)
// 			}

// 			got, err := uc.UsersingUp(tt.args.user)

// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("UserUsecase.UsersingUp() error = %v, wantErr = %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("UserUsecase.UsersingUp() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestMserUsecase_UsersingUp(t *testing.T) {
// 	type args struct {
// 		user models.SignupDetail
// 	}
// 	tests := []struct {
// 		name       string
// 		args       args
// 		beforeTest func(mock sqlmock.Sqlmock)
// 		want       *models.TokenUser
// 		wantErr    bool
// 	}{
// 		{
// 			name: "fail email validation",
// 			args: args{
// 				user: models.SignupDetail{
// 					Email:    "jane@example.com",
// 					Password: "password",
// 					Phone:    "1234567890",
// 				},
// 			},
// 			beforeTest: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectQuery("SELECT \\* FROM \"users\" WHERE \"users\"\\.\"email\" = \\$1 AND \"users\"\\.\"deleted_at\" IS NULL ORDER BY \"users\"\\.\"id\" LIMIT \\$2").
// 					WithArgs("jane@example.com", 1).
// 					WillReturnRows(sqlmock.NewRows([]string{"email"}).AddRow("jane@example.com")) // Indicating email already exists
// 			},
// 			wantErr: true,
// 			want:    &models.TokenUser{},
// 		},

// 		{
// 			name: "email already exists",
// 			args: args{
// 				user: models.SignupDetail{
// 					Email:    "jane@example.com",
// 					Password: "password",
// 					Phone:    "1234567890",
// 				},
// 			},
// 			beforeTest: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectQuery("SELECT \\* FROM \"users\" WHERE \"users\"\\.\"email\" = \\$1 AND \"users\"\\.\"deleted_at\" IS NULL ORDER BY \"users\"\\.\"id\" LIMIT \\$2").
// 					WithArgs("jane@example.com", 1).
// 					WillReturnError(errors.New("error checking email"))
// 			},
// 			wantErr: true,
// 			want:    &models.TokenUser{},
// 		},

// 		{
// 			name: "success user signup",
// 			args: args{
// 				user: models.SignupDetail{
// 					Email:    "john@example.com",
// 					Password: "password",
// 					Phone:    "1234567890",
// 				},
// 			},

// 			beforeTest: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"."email" = \$1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT \$2`).
// 					WithArgs("john@example.com", 1).
// 					WillReturnRows(sqlmock.NewRows([]string{"email"})) // No existing email

// 				// Mocking the phone check
// 				mock.ExpectQuery(`SELECT phone FROM "users" WHERE phone = \$1`).
// 					WithArgs("1234567890").
// 					WillReturnRows(sqlmock.NewRows([]string{"phone"})) // No existing phone

// 				// Mocking the insert operation
// 				mock.ExpectExec(`INSERT INTO "users" \("firstname", "lastname", "email", "phone", "password"\) VALUES \(\$1, \$2, \$3, \$4, \$5\) RETURNING "id", "firstname", "lastname", "email", "phone"`).
// 					WithArgs("John", "Doe", "john@example.com", "1234567890", "hashedPassword").
// 					WillReturnResult(sqlmock.NewResult(1, 1)) // Successful insert

// 				// Mocking the retrieval of the user after insertion
// 				mock.ExpectQuery(`SELECT "id", "firstname", "lastname", "email", "phone" FROM "users" WHERE "id" = \$1`).
// 					WithArgs(1).
// 					WillReturnRows(sqlmock.NewRows([]string{"id", "firstname", "lastname", "email", "phone"}).
// 						AddRow(1, "John", "Doe", "john@example.com", "1234567890")) // Mock response for the inserted user
// 			},
// 			want: &models.TokenUser{
// 				Users: models.SignupDetailResponse{
// 					Email: "john@example.com",
// 					Phone: "1234567890",
// 				},
// 				AccesToken:  "mockAccessToken",  // Replace with actual token generation mock if needed
// 				RefresToken: "mockRefreshToken", // Replace with actual token generation mock if needed
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Create a new sqlmock database
// 			db, mock, err := sqlmock.New()
// 			assert.Nil(t, err, "an error '%s' not expected when opening mock database", err)

// 			gormDB, err := gorm.Open(postgres.New(postgres.Config{
// 				Conn: db,
// 			}), &gorm.Config{})
// 			assert.Nil(t, err, "an error '%s' not expected when opening gorm database", err)

// 			// Create the user repository and use case
// 			ur := repository.UserRepoImpl{DB: gormDB}
// 			uc := &usecase.UserUsecase{UserRepo: &ur}

// 			if tt.beforeTest != nil {

// 				fmt.Println("Before test >>>>/n")
// 				tt.beforeTest(mock) // Setup the mock expectations
// 			}

// 			got, err := uc.UsersingUp(tt.args.user) // Call the method under test

// 			fmt.Println("gottt", got, err)
// 			if (err != nil) != tt.wantErr {

// 				fmt.Println(" ERRRR", err)
// 				fmt.Println("wanttt ERRRR", tt.wantErr)
// 				t.Errorf("UserUsecase.UsersingUp() error = %v, wantErr = %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("UserUsecase.UsersingUp() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestCheckingEmailValidation(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	assert.NoError(t, err)
// 	defer db.Close()

// 	// Open gorm.DB with the mocked SQL database
// 	gormDB, err := gorm.Open(sqlite.New(sqlite.Config{Conn: db}), &gorm.Config{})
// 	assert.NoError(t, err)

// 	userRepo := repository.NewUserRepo(gormDB)

// 	// Mocking the SQLite version check
// 	mock.ExpectQuery(`SELECT sqlite_version\(\)`).
// 		WillReturnRows(sqlmock.NewRows([]string{"sqlite_version"}).AddRow("3.36.0")) // Mock version response

// 	// Mocking the email check
// 	mock.ExpectQuery(`SELECT * FROM "users" WHERE "users"."email" = \$1 AND "users"."deleted_at" IS NULL`).
// 		WithArgs("john@example.com").
// 		WillReturnRows(sqlmock.NewRows([]string{"email", "id", "firstname", "lastname", "phone"}).
// 			AddRow("john@example.com", 1, "John", "Doe", "1234567890")) // Mocked existing user with all fields

// 	// Execute method
// 	user, err := userRepo.CheckingEmailValidation("john@example.com")

// 	// Assert that no error occurred
// 	assert.NoError(t, err)

// 	// Assert that the user is not nil
// 	assert.NotNil(t, user)

// 	// Assert that the returned email is as expected
// 	assert.Equal(t, "john@example.com", user.Email)

// 	// Ensure all expectations were met
// 	err = mock.ExpectationsWereMet()
// 	assert.NoError(t, err)
// }
