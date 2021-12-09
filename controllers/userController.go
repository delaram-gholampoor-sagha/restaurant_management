package controller

import "github.com/gin-gonic/gin"

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {

	}

}

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func LogIn() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func HashPassword(password string) string {
	return
}

func Verifypassword(userPassword string, providedPassword string) (bool, string) {
	return

}
