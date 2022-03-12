package app

import (
	"github.com/Roukii/pock_multiplayer/internal/service"
	"github.com/Roukii/pock_multiplayer/pkg/logger"
	"github.com/Roukii/pock_multiplayer/pkg/postgres"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func RunGateway() {
	l := logger.New("logger.level")

	// HTTP Server
	r := gin.New()

	// TODO : add option in env
	db, err := postgres.New(&postgres.PostgresAuth{
		Host:     "localhost",
		User:     "User",
		Password: "postgres",
		Dbname:   "dbname",
		Port:     "3500",
		Sslmode:  "disable",
		TimeZone: "Europe/Paris",
	})
	if err != nil {
		l.Fatal("Couldn't start database : ", err.Error())
	}
	service.New(db, l)

	err = r.Run(viper.GetString("http.port"))
	if err != nil {
		l.Fatal("Couldn't start server : ", err.Error())
	}
}
