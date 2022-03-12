package http

import (
	"github.com/Roukii/pock_multiplayer/internal/gateway/middleware"
	"github.com/Roukii/pock_multiplayer/internal/gateway/service"
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
		h.GET("/user_world", r.userWorlds)
		h.GET("/world_list", r.worldList)
		h.POST("/join_world", r.joinWorld)
		h.DELETE("/delete_user_world", r.deleteUserWorld)
	}
}

func (wr *worldRoutes) userWorlds(c *gin.Context) {
}

func (wr *worldRoutes) worldList(c *gin.Context) {
}

func (wr *worldRoutes) joinWorld(c *gin.Context) {
}

func (wr *worldRoutes) deleteUserWorld(c *gin.Context) {
}
