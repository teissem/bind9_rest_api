package webserver

import (
	"sync"

	"github.com/gin-gonic/gin"
)

type WebServer struct {
	GinEngine              *gin.Engine
	NamedConfLocalLocation string
	CommandReloadBind      string
	NamedConfLocalLock     sync.Mutex
	ZonesLock              map[string]sync.Mutex
}

func InitWebServer(namedConfLocalLocation string, commandReloadBind string) *WebServer {
	ws := WebServer{
		GinEngine:              gin.Default(),
		NamedConfLocalLocation: namedConfLocalLocation,
		CommandReloadBind:      commandReloadBind,
	}
	ws.AddRoutes()
	return &ws
}

func (ws *WebServer) AddRoutes() {
	api := ws.GinEngine.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.GET("/zone", ws.ZoneList)
			v1.POST("/zone", ws.AddZone)
			v1.PUT("/zone/:zone", ws.ModifyZone)
			v1.DELETE("/zone/:zone", ws.DeleteZone)
		}
	}
}
