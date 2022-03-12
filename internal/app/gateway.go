package app

import (
	"github.com/Roukii/pock_multiplayer/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func RunGateway() {
	l := logger.New("logger.level")

	// HTTP Server
	r := gin.New()
	
	err := r.Run(viper.GetString("http.port"))
	if err != nil {
		l.Fatal("Couldn't start server : ", err.Error())
	}
}
