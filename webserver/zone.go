package webserver

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"teissem.fr/bind9_rest_api/bindfile"
	"teissem.fr/bind9_rest_api/bindmodel"
)

func (ws *WebServer) ZoneList(c *gin.Context) {
	/* ZoneList returns the list of the zone currently declared in BIND9 configuration
	Calling this endpoint lock the named.conf.local file */
	ws.NamedConfLocalLock.Lock()
	defer ws.NamedConfLocalLock.Unlock()
	zones, err := bindfile.NamedConfLocalParser(ws.NamedConfLocalLocation)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, zones)
}

func (ws *WebServer) AddZone(c *gin.Context) {
	// Get body into Zone struct
	var newZone bindmodel.Zone
	err := c.Bind(&newZone)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	// Read named.conf.local
	ws.NamedConfLocalLock.Lock()
	defer ws.NamedConfLocalLock.Unlock()
	zones, err := bindfile.NamedConfLocalParser(ws.NamedConfLocalLocation)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	//
	for _, zone := range zones {
		if zone.Name == newZone.Name || zone.FileLocation == newZone.FileLocation {
			c.AbortWithError(http.StatusUnprocessableEntity, fmt.Errorf("zone with the name %s or the file %s already exists", newZone.Name, newZone.FileLocation))
			return
		}
	}
	zones = append(zones, newZone)
	// Generate the new file
	err = bindfile.GenerateNamedConfLocal(zones, ws.NamedConfLocalLocation)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	// Reload the bind service
	err = bindfile.ReloadBind(ws.CommandReloadBind)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.Status(http.StatusOK)
}
