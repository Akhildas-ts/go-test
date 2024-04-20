package handlers

import (
	"lock/models"
	"lock/response"
	"lock/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func SignUp(c *gin.Context) {


	// signup 
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
	usercreate, err := usecase.UsersignUp(usersign)

	if err != nil {
		errRes := response.ClientResponse(http.StatusBadGateway, "user signup format error ", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusCreated, "User sign up succsed", usercreate, nil)
	c.JSON(http.StatusCreated, successRes)

}

func Login(c *gin.Context) {

	var userLogin models.LoginDetails

	if err := c.ShouldBindJSON(&userLogin); err != nil {

		errRes := response.ClientResponse(http.StatusBadRequest, "json formte was incorrect", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	if err := validator.New().Struct(userLogin); err != nil {
		erres := response.ClientResponse(http.StatusBadGateway, "Login field was wrong formate ahn", nil, err.Error())
		c.JSON(http.StatusBadGateway, erres)
		return
	}

	login, err := usecase.LoginUser(userLogin)

	if err != nil {

		errRes := response.ClientResponse(http.StatusBadRequest, "login error", nil, err.Error())

		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	succesRes := response.ClientResponse(http.StatusOK, "login succesfully", login, nil)
	c.JSON(http.StatusOK, succesRes)

}

func SelectApp(c *gin.Context) {

	

}
