package handlers

import (
	"database/sql"
	"net/http"

	"github.com/LuispeFigueroa/ellineup-api/models"
	"github.com/gin-gonic/gin"
)

var DB *sql.DB

// GET /divisiones
func GetDivisiones(c *gin.Context) {
	q := c.Query("q")

	var rows *sql.Rows
	var err error

	if q != "" {
		rows, err = DB.Query("SELECT id, nombre, temporada FROM divisiones WHERE LOWER(nombre) LIKE LOWER($1)", "%"+q+"%")
	} else {
		rows, err = DB.Query("SELECT id, nombre, temporada FROM divisiones")
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

// GET /divisiones/:id
func GetDivision(c *gin.Context) {
	id := c.Param("id")
	var d models.Division

	err := DB.QueryRow("SELECT id, nombre, temporada FROM divisiones WHERE id = $1", id).
		Scan(&d.ID, &d.Nombre, &d.Temporada)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Division no encontrada"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, d)
}

// POST /divisiones
func CreateDivision(c *gin.Context) {
	var d models.Division

	if err := c.ShouldBindJSON(&d); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos invalidos"})
		return
	}

	// Validaciones
	if d.Nombre == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El nombre de la division es requerido"})
		return
	}
	if d.Temporada == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "La temporada es requerida"})
		return
	}

	err := DB.QueryRow(
		"INSERT INTO divisiones (nombre, temporada) VALUES ($1, $2) RETURNING id",
		d.Nombre, d.Temporada,
	).Scan(&d.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, d)
}

// PUT /divisiones/:id
func UpdateDivision(c *gin.Context) {
	id := c.Param("id")
	var d models.Division

	if err := c.ShouldBindJSON(&d); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos invalidos"})
		return
	}

	// Validaciones
	if d.Nombre == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El nombre de la division es requerido"})
		return
	}
	if d.Temporada == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "La temporada es requerida"})
		return
	}

	result, err := DB.Exec(
		"UPDATE divisiones SET nombre=$1, temporada=$2 WHERE id=$3",
		d.Nombre, d.Temporada, id,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Division no encontrada"})
		return
	}

	c.JSON(http.StatusOK, d)
}

// DELETE /divisiones/:id
func DeleteDivision(c *gin.Context) {
	id := c.Param("id")

	result, err := DB.Exec("DELETE FROM divisiones WHERE id=$1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Division no encontrada"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
