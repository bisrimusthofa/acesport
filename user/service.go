package user

import (
	"errors"
	"time"

	"github.com/bisrimusthofa/acesport/helper"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(input RegisterUserInput) (User, error)
	Login(input LoginInput) (User, error)
	IsEmailAvailable(input CheckEmailInput) (bool, error)
	SaveAvatar(id string, filePath string) (User, error)
	FindById(id string) (User, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository: repository}
}

func (s *service) Register(input RegisterUserInput) (User, error) {
	user := User{
		Id:        uuid.New().String(),
		Name:      input.Name,
		Email:     input.Email,
		Phone:     input.Phone,
		Role:      "Reseller",
		CreatedAt: time.Now(),
	}

	password, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	helper.PanicIfError(err)
	user.Password = string(password)

	dataUser, err := s.repository.Save(user)
	if err != nil {
		return user, err
	}

	return dataUser, nil
}

func (s *service) Login(input LoginInput) (User, error) {
	email := input.Email
	password := input.Password

	// find by email
	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return user, err
	}

	if user.Id == "" {
		return user, errors.New("Email atau Password tidak valid")
	}

	// matching password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return user, errors.New("Email atau Password tidak valid")
	}

	return user, nil
}

func (s *service) IsEmailAvailable(input CheckEmailInput) (bool, error) {
	email := input.Email

	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return false, err
	}

	if user.Id == "" {
		return true, nil
	}

	return false, nil
}

func (s *service) SaveAvatar(id string, filePath string) (User, error) {
	user, err := s.repository.FindById(id)
	if err != nil {
		return user, err
	}

	user.Avatar = filePath

	if user.Id == "" {
		return user, errors.New("User not Found")
	}
	
	updatedUser, err := s.repository.Update(user)
	if err != nil {
		return updatedUser, err
	}

	return updatedUser, nil
}

func (s *service) FindById(id string) (User, error) {
	user, err := s.repository.FindById(id)
	if err != nil {
		return user, err
	}

	if user.Id == "" {
		return user, errors.New("User not found")
	}

	return user, nil
}
