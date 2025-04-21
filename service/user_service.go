package service

import (
	"errors"
	"tracker/middleware"
	"tracker/models"
	"tracker/repo"
	"tracker/utils"
)

type UserService struct {
	Repo repo.UserRepository
}

func(s *UserService) RegisterUser(user *models.User) error {
	// check if user exists already
	_, err := s.Repo.GetUserByEmail(user.Email)
	if err != nil {
		return errors.New("email already in use")
	}

	// Hash the password

	hashPass, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashPass

	// Calling create user method

	err = s.Repo.CreateUser(user)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) LoginUser(loginRequest models.LoginUser) (string, error) {
	user, err := s.Repo.GetUserByEmail(loginRequest.Email)
	if err != nil {
		return "", err
	}

	// compare password
	err = utils.ComparePassword(user.Password, loginRequest.Password)
	if err != nil {
		return "", err
	}

	// generate token
	token, err := middleware.GenerateJWT(user.ID, user.Role)
	if err != nil {
		return "", err
	}
	return token, nil
}