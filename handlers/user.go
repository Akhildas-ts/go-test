package handlers

import (
	"lock/models"
	"lock/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func SignUp(c *gin.Context) {

	var usersign models.SignupDetail

	if err := c.ShouldBindJSON(&usersign); err != nil {

		errRes := response.ClientResponse(http.StatusBadRequest, "field are provied wrong formate ", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return

	}

	if err := validator.New().Struct(usersign); err != nil {

		errRes := response.ClientResponse(404, "they are  not in formate", nil, err.Error())

		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	usercreate, err := usecase.UsersingUp(usersign)

	if err != nil {
		errRes := response.ClientResponse(http.StatusBadGateway, "user signup format error ", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusCreated, "User sign up succsed", usercreate, nil)
	c.JSON(http.StatusCreated, successRes)

}

