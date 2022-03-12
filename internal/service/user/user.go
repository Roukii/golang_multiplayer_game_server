package user

import (
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
	hashedPassword, err := bcrypt.GenerateFromPassword(passwordByte, bcrypt.DefaultCost)
	if err != nil {
			return nil, err
	}
	user, err := a.userDao.GetByUsernamePassword(username, hashedPassword)
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (a *UserService) Register(input UserInput) (*entity.User, error) {
	user := entity.User{}
	user.Username = input.Name
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
