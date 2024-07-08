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
	// Searching already existing zone or file
	for _, zone := range zones {
		if zone.Name == newZone.Name || zone.FileLocation == newZone.FileLocation {
			c.AbortWithError(http.StatusUnprocessableEntity, fmt.Errorf("zone with the name %s or the file %s already exists", newZone.Name, newZone.FileLocation))
			return
		}
	}
	zones = append(zones, newZone)
	if err = ws.generateNamedConfLocalAndReloadBind(zones); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

func (ws *WebServer) ModifyZone(c *gin.Context) {
	// Get zone name from path param
	zoneName := c.Param("zone")
	// Get body into ModifyZoneInput struct
	var modifyZoneInput bindmodel.ModifyZoneInput
	err := c.Bind(&modifyZoneInput)
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
	// Searching already existing zone or file
	foundZone := false
	for _, zone := range zones {
		if zone.Name == zoneName {
			zone.FileLocation = modifyZoneInput.FileLocation
			foundZone = true
		} else if zone.FileLocation == modifyZoneInput.FileLocation {
			c.AbortWithError(http.StatusUnprocessableEntity, fmt.Errorf("zone with the file %s already exists", modifyZoneInput.FileLocation))
			return
		}
	}
	if !foundZone {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("%s zone not found", zoneName))
		return
	}
	if err = ws.generateNamedConfLocalAndReloadBind(zones); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

func (ws *WebServer) DeleteZone(c *gin.Context) {
	// Get zone name from path param
	zoneName := c.Param("zone")
	// Read named.conf.local
	ws.NamedConfLocalLock.Lock()
	defer ws.NamedConfLocalLock.Unlock()
	zones, err := bindfile.NamedConfLocalParser(ws.NamedConfLocalLocation)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	var modifiedZones []bindmodel.Zone
	for _, zone := range zones {
		if zone.Name != zoneName {
			modifiedZones = append(modifiedZones, zone)
		}
	}
	if err = ws.generateNamedConfLocalAndReloadBind(modifiedZones); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

func (ws *WebServer) generateNamedConfLocalAndReloadBind(zones []bindmodel.Zone) error {
	// Generate the new file
	err := bindfile.GenerateNamedConfLocal(zones, ws.NamedConfLocalLocation)
	if err != nil {
		return err
	}
	// Reload the bind service
	return bindfile.ReloadBind(ws.CommandReloadBind)
}
