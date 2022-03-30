package user

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/Roukii/pock_multiplayer/internal/gateway/dao"
	"github.com/Roukii/pock_multiplayer/internal/gateway/entity"
	"github.com/gin-gonic/gin"
)

type UserService struct {
	userDao      *dao.UserDao
	connexionDao *dao.ConnexionDao
}

// New -.
func New(r *dao.UserDao, c *dao.ConnexionDao) *UserService {
	return &UserService{
		userDao:      r,
		connexionDao: c,
	}
}

func (a *UserService) GetById(userId string) (*entity.User, error) {
	user, err := a.userDao.GetById(userId)
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (a *UserService) Login(username string, password string, c *gin.Context) (*entity.User, error) {
	passwordByte := []byte(password)
	user, err := a.userDao.GetByEmail(username)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword(user.Password, passwordByte)
	if err != nil {
		return nil, err
	}
	a.connexionDao.SaveOrUpdate(&entity.Connexion{UserId: user.UUID, Ip: c.ClientIP(), UserAgent: c.GetHeader("User-Agent")})
	return &user, err
}

func (a *UserService) Register(input UserInput) (*entity.User, error) {
	user := entity.User{}
	user.Email = input.Email
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
