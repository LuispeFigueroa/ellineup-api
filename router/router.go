package router

import (
	"github.com/LuispeFigueroa/ellineup-api/handlers"
	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) {
	// Divisiones
	r.GET("/divisiones", handlers.GetDivisiones)
	r.GET("/divisiones/:id", handlers.GetDivision)
	r.POST("/divisiones", handlers.CreateDivision)
	r.PUT("/divisiones/:id", handlers.UpdateDivision)
	r.DELETE("/divisiones/:id", handlers.DeleteDivision)
}
