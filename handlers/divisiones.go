package handlers

import (
	"database/sql"
	"net/http"

	"github.com/LuispeFigueroa/ellineup-api/models"
	"github.com/gin-gonic/gin"
)

var DB *sql.DB

// Get para obtener las divisiones /divisiones
func GetDivisiones(c *gin.Context) {
	rows, err := DB.Query("SELECT id, nombre, temporada FROM divisiones")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener todas las divisiones"})
		return
	}
	defer rows.Close()

	divisiones := []models.Division{}
	for rows.Next() {
		var d models.Division
		rows.Scan(&d.ID, &d.Nombre, &d.Temporada)
		divisiones = append(divisiones, d)
	}

	c.JSON(http.StatusOK, divisiones)
}

// Get divisiones por id /divisiones/id
func GetDivision(c *gin.Context) {
	id := c.Param("id")
	var d models.Division

	err := DB.QueryRow("SELECT id,nombre, temporada FROM divisiones WHERE id = $1", id).
		Scan(&d.ID, &d.Nombre, &d.Temporada)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Error al obtener la division"})
		return
	}
	c.JSON(http.StatusOK, d)
}

//POST para agregar una divisio /divisiones

func CreateDivision(c *gin.Context) {
	var d models.Division

	if err := c.ShouldBindJSON(&d); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos nos validos"})
		return
	}

	err := DB.QueryRow(
		"INSERT INTO divisiones (nombre, temporada) VALUES ($1, $2) RETURNIN id",
		d.Nombre, d.Temporada,
	).Scan(&d.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creando la division"})
		return
	}
	c.JSON(http.StatusCreated, d)
}
