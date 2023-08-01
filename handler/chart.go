package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rickyromansyah2045/halocat-backend-go/chart"
	"github.com/rickyromansyah2045/halocat-backend-go/helper"
)

type chartHandler struct {
	chartSvc chart.Service
}

func NewChartHandler(
	chartService chart.Service,
) *chartHandler {
	return &chartHandler{
		chartSvc: chartService,
	}
}

func (handler *chartHandler) GetChart(ctx *gin.Context) {
	chartName := ctx.Query("chart")
	year := ctx.DefaultQuery("year", "")

	if chartName == "" {
		response := helper.APIResponseError(http.StatusBadRequest, "Get chart failed!", "Invalid Request!")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	chart, err := handler.chartSvc.GetChart(chartName, year)

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Get chart failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.APIResponse(http.StatusOK, "Get chart successfully!", chart)

	ctx.JSON(http.StatusOK, response)
}
