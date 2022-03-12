package http

import (
	"github.com/Roukii/pock_multiplayer/internal/middleware"
	"github.com/Roukii/pock_multiplayer/pkg/logger"
	"github.com/gin-gonic/gin"

)

type worldRoutes struct {
	l logger.Interface
}

func newWorldRoutes(handler *gin.RouterGroup, l logger.Interface) {
	r := &worldRoutes{l}

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
