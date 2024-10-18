package repository

import (
	"errors"
	"lock/database"
	"lock/domain"
	"lock/models"
)

func AdminLogin(admin models.AdminLogin) (domain.Admin, error) {
	var admindomain domain.Admin

	if err := database.DB.Raw("select * from users where email= ? and isadmin= true ", admin.Email).Scan(&admindomain).Error; err != nil {
		return domain.Admin{}, errors.New("admin email is not available on database")
	}

	return admindomain, nil
}
