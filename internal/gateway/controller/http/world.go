package http

import (
	"net/http"

	"github.com/Roukii/pock_multiplayer/internal/gateway/middleware"
	"github.com/Roukii/pock_multiplayer/internal/gateway/service"
	"github.com/Roukii/pock_multiplayer/internal/gateway/service/world"
	"github.com/gin-gonic/gin"
)

type worldRoutes struct {
	services *service.Service
}

func newWorldRoutes(handler *gin.RouterGroup, services *service.Service) {
	r := &worldRoutes{services}

	h := handler.Group("/world")
	h.Use(middleware.CheckTokenJWT)
	{
		h.GET("/user", r.userWorlds)
		h.GET("/list", r.worldList)
		h.POST("/join", r.joinWorld)
		h.DELETE("/delete_user", r.deleteUserWorld)
	}
}

func (wr *worldRoutes) userWorlds(c *gin.Context) {
	userId := c.GetString("uuid")
	worlds, err := wr.services.WorldService.GetUserWorlds(userId)
	if err != nil {
		wr.services.Logger.Error("userWorlds failed : ", err.Error())
		c.JSON(http.StatusBadRequest, "There are no worlds for this user")
		return
	}
	c.JSON(http.StatusOK, worlds)
}

func (wr *worldRoutes) worldList(c *gin.Context) {
	worlds, err := wr.services.WorldService.GetAvailableWorld()
	if err != nil {
		wr.services.Logger.Error("WorldList failed: ", err.Error())
		c.JSON(http.StatusInternalServerError, "There is no world found")
		return
	}
	c.JSON(http.StatusOK, worlds)
}

func (wr *worldRoutes) joinWorld(c *gin.Context) {
	var w world.WorldInput
	userId := c.GetString("uuid")
	if err := c.ShouldBindJSON(&w); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	err := wr.services.UserWorldAffdService.JoinWorld(userId, w.Name)
	if err != nil {
		wr.services.Logger.Error("joinWorld failed: ", err.Error())
		c.JSON(http.StatusBadRequest, "Invalid json provided")
		return
	}
	c.JSON(http.StatusOK, true)
}

func (wr *worldRoutes) deleteUserWorld(c *gin.Context) {
	var w world.WorldInput
	userId := c.GetString("uuid")
	if err := c.ShouldBindJSON(&w); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	err := wr.services.UserWorldAffdService.DeleteWorld(userId, w.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid json provided")
		return
	}
	c.JSON(http.StatusOK, true)
}
