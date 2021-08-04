package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"main/app/parkmap"
	"main/app/parkmap/geo"
	"net/http"
)

type Payload struct {
	Latitude          float64 `form:"lat,string,omitempty" json:"lat,string,omitempty" xml:"lat,string,omitempty" binding:"required"`
	Longitude         float64 `form:"lon,string,omitempty" json:"lon,string,omitempty" xml:"lon,string,omitempty" binding:"required"`
	MaximumDistanceKm float64 `form:"dist,string,omitempty" json:"dist,string,omitempty" xml:"dist,string,omitempty" default:"0.015"`
	IsDisabled        bool    `form:"is_disabled" json:"is_disabled" xml:"is_disabled" default:"false"`
	Hours             int     `form:"hours" json:"hours" xml:"hours" binding:"required"`
}

type ParkingMapController struct{}

var mappia = parkmap.NewParkingMap()

func (p ParkingMapController) GetClosest(c *gin.Context) {
	var payload Payload
	if err := c.ShouldBind(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": fmt.Sprintf("%v", http.StatusBadRequest),
			"error":  err.Error(),
		})
		return
	}
	me := geo.NewPoint(payload.Latitude, payload.Longitude)
	err, parking := mappia.Closest(me, true, payload.IsDisabled, payload.MaximumDistanceKm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": fmt.Sprintf("%v", http.StatusBadRequest),
			"error":  fmt.Sprintf("%s", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  fmt.Sprintf("%v", http.StatusOK),
		"code":    parking.ParkingZoneNumber,
		"address": parking.Address,
		"hours":   payload.Hours,
	})
}
