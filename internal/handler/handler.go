// Package handler wires HTTP handlers, serves static files and the HTML frontend.
package handler

import (
	v1 "Hades/internal/handler/v1"
	"Hades/internal/service"
	"html/template"
	"net/http"

	"github.com/wb-go/wbf/ginext"
)

const templatePath = "web/templates/index.html"

// NewHandler creates and returns an http.Handler with all routes configured.
// It sets up static file serving, API v1 group, and the root HTML page.
func NewHandler(service service.Service) http.Handler {

	handler := ginext.New("")

	handler.Use(ginext.Recovery())
	handler.Static("/static", "./web/static")

	apiV1 := handler.Group("/api/v1")
	handlerV1 := v1.NewHandler(service)

	apiV1.GET("/items", handlerV1.GetItems)
	apiV1.POST("/items", handlerV1.CreateItem)
	apiV1.PUT("/items/:id", handlerV1.UpdateItem)
	apiV1.DELETE("/items/:id", handlerV1.DeleteItem)

	apiV1.GET("/analytics", handlerV1.GetAnalytics)

	handler.GET("/", homePage(template.Must(template.ParseFiles(templatePath))))

	return handler

}

// homePage renders the main HTML template for the root route.
func homePage(t *template.Template) ginext.HandlerFunc {
	return func(c *ginext.Context) {
		if err := t.Execute(c.Writer, nil); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, ginext.H{"error": "Failed to render page"})
		}
	}
}
