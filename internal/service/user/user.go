package user

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/Roukii/pock_multiplayer/internal/dao"
	"github.com/Roukii/pock_multiplayer/internal/entity"
)

type UserService struct {
	userDao *dao.UserDao
}

// New -.
func New(r *dao.UserDao) *UserService {
	return &UserService{
		userDao: r,
	}
}

func (a *UserService) GetById(userId string) (*entity.User, error) {
	user, err := a.userDao.GetById(userId)
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (a *UserService) Login(username string, password string) (*entity.User, error) {
	passwordByte := []byte(password)
	user, err := a.userDao.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword(user.Password, passwordByte)
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (a *UserService) Register(input UserInput) (*entity.User, error) {
	fmt.Println("register pass : " + input.Password)
	user := entity.User{}
	user.Username = input.Username
	passwordByte := []byte(input.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(passwordByte, bcrypt.DefaultCost)
	if err != nil {
			return nil, err
	}
	user.Password = hashedPassword

	err = a.userDao.SaveOrUpdate(&user)
	if err != nil {
		return nil, err
	}
	return &user, err
}
