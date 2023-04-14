package router

import (
	"github.com/gin-gonic/gin"
	"log"
	Controller "store_monitoring/controller"
)

func InitRouter() {
	router := gin.Default()
	controller := &Controller.Controller{}
	router.PUT("/trigger_report", controller.TriggerReport())
	router.GET("/get_report/:reportId", controller.GetReport())
	err := router.Run(":8080")
	if err != nil {
		log.Printf("Error while running on port, Error = " + err.Error())
	}
}
