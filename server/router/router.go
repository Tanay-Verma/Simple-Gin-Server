package router

import (
	"server/internal/user"

	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func InitRouter(userHandler *user.Handler) {
	r = gin.Default()

	r.POST("/signup", userHandler.CreateUser)
	r.POST("/login", userHandler.LoginUser)
	r.POST("/logout", userHandler.LogoutUser)
}

func Start(addr string) error {
	return r.Run(addr)
}
