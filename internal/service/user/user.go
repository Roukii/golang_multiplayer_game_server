package user

import (
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

func (a *UserService) getById(userId string) (entity.User, error) {
	user, err := a.userDao.GetById(userId)
	if err != nil {
		return entity.User{}, err
	}
	return user, err
}

func (a *UserService) register(input UserInput) (entity.User, error) {
	user := entity.User{}
	user.Username = input.Name
	user.Password = input.Password
	err := a.userDao.SaveOrUpdate(&user)
	if err != nil {
		return entity.User{}, err
	}
	return user, err
}
