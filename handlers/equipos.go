package handlers

import (
	"database/sql"
	"net/http"

	"github.com/LuispeFigueroa/ellineup-api/models"
	"github.com/gin-gonic/gin"
)

// Get /divisiones/:id/equipos
// Obtener todos los equipos de una division
func GetEquipos(c *gin.Context) {
	divisionID := c.Param("id")

	rows, err := DB.Query("SELECT id, division_id, nombre, logo_url FROM equipos WHERE division_id = $1", divisionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	equipos := []models.Equipo{}
	for rows.Next() {
		var e models.Equipo
		rows.Scan(&e.ID, &e.DivisionID, &e.Nombre, &e.LogoURL)
		equipos = append(equipos, e)
	}

	c.JSON(http.StatusOK, equipos)
}

// GET /equipos/:id
// obtener un euqipo especifico por id
func GetEquipo(c *gin.Context) {
	id := c.Param("id")
	var e models.Equipo

	err := DB.QueryRow("SELECT id, division_id, nombre, logo_url FROM equipos WHERE id = $1", id).Scan(&e.ID, &e.DivisionID, &e.Nombre, &e.LogoURL)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, e)
}

// POST /divisiones/:id/equipos
// craer un nuevo equipo en una division
func CreateEquipo(c *gin.Context) {
	divisionID := c.Param("id")
	var e models.Equipo

	if err := c.ShouldBindJSON(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos invalidos"})
		return
	}
	err := DB.QueryRow(
		"INSERT INTO equipos (division_id, nombre, logo_url) VALUES ($1, $2, $3) RETURNING id",
		divisionID, e.Nombre, e.LogoURL,
	).Scan(&e.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, e)
}

// PUT /equipos/id
// MODIFICAR un equipo especifico
func UpdateEquipo(c *gin.Context) {
	id := c.Param("id")
	var e models.Equipo

	if err := c.ShouldBindJSON(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos no Validos"})
		return
	}
	result, err := DB.Exec(
		"UPDATE equipos SET nombre= $1, logo_url=$2 WHERE id=$3",
		e.Nombre, e.LogoURL, id,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Equipo no encontrado"})
		return
	}

	c.JSON(http.StatusOK, e)
}

// DELETE /EQUIPOS/id
// eliminar a un equipo
func DeleteEquipo(c *gin.Context) {
	id := c.Param("id")

	result, err := DB.Exec("DELETE FROM equipos WHERE id=$1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Equipo no encontrado"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
