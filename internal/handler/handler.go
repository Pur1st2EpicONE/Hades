package handler

import (
	v1 "Hades/internal/handler/v1"
	"Hades/internal/service"
	"net/http"

	"github.com/wb-go/wbf/ginext"
)

func NewHandler(service service.Service) http.Handler {

	handler := ginext.New("")

	handler.Use(ginext.Recovery())

	apiV1 := handler.Group("/api/v1")
	handlerV1 := v1.NewHandler(service)

	apiV1.GET("/items", handlerV1.GetItems)
	apiV1.POST("/items", handlerV1.CreateItem)
	apiV1.PUT("/items/:id", handlerV1.UpdateItem)
	apiV1.DELETE("/items/:id", handlerV1.DeleteItem)

	apiV1.GET("/analytics", handlerV1.GetAnalytics)

	return handler

}
