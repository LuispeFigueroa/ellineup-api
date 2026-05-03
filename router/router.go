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

	//Equipos
	r.GET("/divisiones/:id/equipos", handlers.GetEquipos)
	r.GET("/equipos/:id", handlers.GetEquipo)
	r.POST("/divisiones/:id/equipos", handlers.CreateEquipo)
	r.PUT("/equipos/:id", handlers.UpdateEquipo)
	r.DELETE("/equipos/:id", handlers.DeleteEquipo)

}
