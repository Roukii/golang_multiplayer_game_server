package app

import (
	"github.com/Roukii/pock_multiplayer/internal/controller/http"
	"github.com/Roukii/pock_multiplayer/internal/service"
	"github.com/Roukii/pock_multiplayer/pkg/logger"
	"github.com/Roukii/pock_multiplayer/pkg/postgres"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func RunGateway() {
	l := logger.New("logger.level")

	r := gin.New()

	db, err := postgres.New(&postgres.PostgresAuth{
		Host:     viper.GetString("database.host"),
		User:     viper.GetString("database.user"),
		Password: viper.GetString("database.password"),
		Dbname:   viper.GetString("database.dbname"),
		Port:     viper.GetString("database.port"),
		Sslmode:  viper.GetString("database.sslmode"),
		TimeZone: viper.GetString("database.timeZone"),
	})
	if err != nil {
		l.Fatal("Couldn't start database : ", err.Error())
	}
	services := service.New(db, l)

	http.NewRouter(r, services)
	err = r.Run(viper.GetString("http.port"))
	if err != nil {
		l.Fatal("Couldn't start server : ", err.Error())
	}
}
