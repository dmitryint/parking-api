package api

import (
	"github.com/gin-gonic/gin"
	"main/app/api/controllers"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	v1 := router.Group("v1")
	{
		userGroup := v1.Group("location")
		{
			parkmap := new(controllers.ParkingMapController)
			userGroup.POST("/", parkmap.GetClosest)
		}
	}
	return router
}
