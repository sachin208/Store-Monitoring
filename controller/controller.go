package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"store_monitoring/db/file"
	"store_monitoring/db/repository"
	"store_monitoring/service"
)

type Controller struct{}

func (c *Controller) TriggerReport() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		serviceInstance := service.Service{}
		fileInstance := file.File{}
		response := serviceInstance.CalculateLastUpAndDownTime()
		reportId := fileInstance.PrepareCsvForReport(response)
		if reportId == "" {
			c.JSON(http.StatusInternalServerError, reportId)
		} else {
			c.JSON(http.StatusOK, reportId)
		}
	}
	return fn
}

func (c *Controller) GetReport() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		reportId := c.Param("reportId")
		reportCompletionStatusRepository := repository.ReportCompletionStatusRepository{}
		response, err := reportCompletionStatusRepository.GetCsvFilePathByReportId(reportId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			if response.Status == "COMPLETED" {
				c.JSON(http.StatusOK, response)
			} else {
				c.JSON(http.StatusOK, "Running")
			}
		}
	}
	return fn
}
