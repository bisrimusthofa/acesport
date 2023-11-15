package main

import (
	"net/http"
	"strings"

	"github.com/bisrimusthofa/acesport/auth"
	"github.com/bisrimusthofa/acesport/controller"
	"github.com/bisrimusthofa/acesport/helper"
	"github.com/bisrimusthofa/acesport/user"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/acesport?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	helper.PanicIfError(err)

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	authService := auth.NewJwtService()
	userController := controller.NewUserController(userService, authService)

	router := gin.Default()
	api := router.Group("/api/v1")

	//auth
	api.POST("/register", userController.Register)
	api.POST("/login", userController.Login)
	api.POST("/email-checker", userController.CheckEmailAvailable)

	// profile
	api.POST("/upload-avatar", authMiddleware(authService, userService), userController.UploadAvatar)

	err = router.Run()
	helper.PanicIfError(err)
}

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {

	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")

		// check apakah mengandung kata Bearer
		if !strings.Contains(header, "Bearer") {
			response := helper.APIResponse(http.StatusUnauthorized, "error", "Unauthorized", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		//ambil token
		var tokenString string
		tokenSplited := strings.Split(header, " ")
		if len(tokenSplited) == 2 {
			tokenString = tokenSplited[1]
		}

		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse(http.StatusUnauthorized, "error", "Invalid Token", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		payload, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			response := helper.APIResponse(http.StatusUnauthorized, "error", "Invalid Token", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userId := payload["user_id"].(string)

		user, err := userService.FindById(userId)
		if err != nil {
			response := helper.APIResponse(http.StatusUnauthorized, "error", "Invalid Token", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", user)
		c.Set("userId", user.Id)
	}
}
