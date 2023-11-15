package controller

import (
	"net/http"

	"github.com/bisrimusthofa/acesport/auth"
	"github.com/bisrimusthofa/acesport/helper"
	"github.com/bisrimusthofa/acesport/user"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService user.Service
	authService auth.Service
}

func NewUserController(userService user.Service, authService auth.Service) *UserController {
	return &UserController{
		userService: userService,
		authService: authService,
	}
}

func (controller *UserController) Register(c *gin.Context) {
	var input user.RegisterUserInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.ErrorValidationResponse(err)
		errorMessages := gin.H{"errors": errors}

		response := helper.APIResponse(
			http.StatusUnprocessableEntity,
			"error",
			"The data given was invalid",
			errorMessages,
		)
		c.JSON(http.StatusBadRequest, response)
		panic(err)
	}

	newUser, err := controller.userService.Register(input)
	if err != nil {
		response := helper.APIResponse(
			http.StatusBadRequest,
			"error",
			"Failed register account",
			nil,
		)
		c.JSON(http.StatusBadRequest, response)
		panic(err)
	}

	dataFormated := user.FormatUser(newUser)

	// token
	token, err := controller.authService.GenerateToken(dataFormated.Id)

	dataFormated.Token = token

	response := helper.APIResponse(
		http.StatusCreated,
		"success",
		"Account has been created",
		dataFormated,
	)
	c.JSON(http.StatusOK, response)
}

func (controller *UserController) Login(c *gin.Context) {
	var input user.LoginInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.ErrorValidationResponse(err)
		errorMessages := gin.H{"errors": errors}

		response := helper.APIResponse(
			http.StatusUnprocessableEntity,
			"error",
			"The data given was invalid",
			errorMessages,
		)
		c.JSON(http.StatusBadRequest, response)
		panic(err)
	}

	dataUser, err := controller.userService.Login(input)
	if err != nil {
		errorMessages := gin.H{"errors": err.Error()}
		response := helper.APIResponse(
			http.StatusBadRequest,
			"error",
			"The data given was invalid",
			errorMessages,
		)
		c.JSON(http.StatusBadRequest, response)
		panic(err)
	}

	dataFormated := user.FormatUser(dataUser)

	// token
	token, err := controller.authService.GenerateToken(dataFormated.Id)

	dataFormated.Token = token

	response := helper.APIResponse(
		http.StatusAccepted,
		"success",
		"Login Successfully",
		dataFormated,
	)

	c.JSON(http.StatusAccepted, response)
}

func (controller *UserController) CheckEmailAvailable(c *gin.Context) {
	var input user.CheckEmailInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.ErrorValidationResponse(err)
		errorMessages := gin.H{"errors": errors}

		response := helper.APIResponse(
			http.StatusUnprocessableEntity,
			"error",
			"The email was used",
			errorMessages,
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	isEmailAvailable, err := controller.userService.IsEmailAvailable(input)
	if err != nil {
		errorMessages := gin.H{"errors": "Server Error"}

		response := helper.APIResponse(
			http.StatusUnprocessableEntity,
			"error",
			"Email checking failed",
			errorMessages,
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{
		"is_available": isEmailAvailable,
	}

	metaMessage := "Email is not available"

	if isEmailAvailable {
		metaMessage = "Email is available"
	}

	response := helper.APIResponse(http.StatusOK, "success", metaMessage, data)

	c.JSON(http.StatusOK, response)
}

func (controller *UserController) UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {
		errorMessages := gin.H{"errors": false}

		response := helper.APIResponse(
			http.StatusUnprocessableEntity,
			"error",
			"Failed upload avatar",
			errorMessages,
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// simpan file
	avatarPath := user.AvatarPath + file.Filename
	err = c.SaveUploadedFile(file, user.AvatarPath+file.Filename)
	if err != nil {
		errorMessages := gin.H{"errors": err.Error()}

		response := helper.APIResponse(
			http.StatusUnprocessableEntity,
			"error",
			"Failed upload avatar",
			errorMessages,
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = controller.userService.SaveAvatar("b2f426f3-83fa-421d-9a58-c2c9a4521d2f", avatarPath)
	if err != nil {
		errorMessages := gin.H{"errors": err.Error()}

		response := helper.APIResponse(
			http.StatusUnprocessableEntity,
			"error",
			"Failed upload avatar",
			errorMessages,
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	dataResponse := gin.H{"is_uploaded": true}
	response := helper.APIResponse(
		http.StatusUnprocessableEntity,
		"error",
		"Success upload avatar",
		dataResponse,
	)
	c.JSON(http.StatusOK, response)
}
