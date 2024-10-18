package routes

import (
	"lock/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserRoutes(r *gin.RouterGroup, db *gorm.DB, uh handlers.UserHandler) *gin.RouterGroup {

	r.POST("/signup", uh.Signup)
	r.POST("/login", uh.UserLoginWithPassword)

	return r

}
