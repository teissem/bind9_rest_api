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
	ws.GinEngine.GET("/zone", ws.ZoneList)
	ws.GinEngine.POST("/zone", ws.AddZone)
	ws.GinEngine.PUT("/zone/:zone", ws.ModifyZone)
	ws.GinEngine.DELETE("/zone/:zone", ws.DeleteZone)
}
