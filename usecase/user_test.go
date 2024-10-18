package usecase_test

import (
	"context"
	"fmt"
	"lock/models"
	"lock/repository"
	"reflect"
	"regexp"

	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Test_userRepo_Create(t *testing.T) {

	type args struct {
		ctx   context.Context
		input models.SignupDetail
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(sqlmock.Sqlmock)
		want       models.SignupDetailResponse
		wantErr    bool
	}{
		{
			name: "fail create user",
			args: args{
				ctx:   context.TODO(),
				input: models.SignupDetail{Email: "john@example.com", FirstName: "John", LastName: "Doe", Phone: "923539327", Password: "akhil@123"},
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				mockSQL.
					ExpectQuery(regexp.QuoteMeta(`
						INSERT INTO users (email, firstname, lastname, phone, password)
						VALUES ($1, $2, $3, $4, $5) RETURNING id, firstname, lastname, email, phone`,
					)).
					WithArgs("john@example.com", "John", "Doe", "923539327", "akhil@123").
					WillReturnError(fmt.Errorf("failed to insert user"))
			},
			wantErr: true,
			want:    models.SignupDetailResponse{}, // Use empty struct instead of nil
		},

		{
			name: "success create user",
			args: args{
				ctx:   context.TODO(),
				input: models.SignupDetail{Email: "john@example.com", FirstName: "John", LastName: "Doe", Phone: "923539327", Password: "akhil@123"},
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				mockSQL.
					ExpectQuery(regexp.QuoteMeta(`
						INSERT INTO users (email, firstname, lastname, phone, password)
						VALUES ($1, $2, $3, $4, $5) RETURNING id, firstname, lastname, email, phone`,
					)).
					WithArgs("john@example.com", "John", "Doe", "923539327", "akhil@123").
					WillReturnRows(sqlmock.NewRows([]string{"id", "firstname", "lastname", "email", "phone"}).
						AddRow(1, "John", "Doe", "john@example.com", "923539327"))
			},
			want:    models.SignupDetailResponse{ID: 1, Email: "john@example.com", FirstName: "John", LastName: "Doe"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			sqlDB, mock, _ := sqlmock.New()

			gormDB, err := gorm.Open(postgres.New(postgres.Config{
				Conn: sqlDB,
			}), &gorm.Config{})

			if err != nil {
				t.Fatal("Error found of connection of database ", err)
			}

			u := &repository.UserRepoImpl{
				DB: gormDB,
			}

			if tt.beforeTest != nil {
				tt.beforeTest(mock)
			}

			got, err := u.SignupInsert(tt.args.input)
			fmt.Println("BEFORE SIGNUP INSERT ")

			if (err != nil) != tt.wantErr {
				fmt.Println("ERROR FROM USER REPO ")
				t.Errorf("userRepo.Create() error = %v, wantErr =%v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				fmt.Println("ERROR FROM USER REPO  2")

				t.Errorf("userRepo.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}
