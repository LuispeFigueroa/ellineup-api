package handlers

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"github.com/LuispeFigueroa/ellineup-api/models"
	"github.com/gin-gonic/gin"
)

// Get /divisiones/:id/equipos
// Obtener todos los equipos de una division
func GetEquipos(c *gin.Context) {
	divisionID := c.Param("id")
	q := c.Query("q")

	var rows *sql.Rows
	var err error

	if q != "" {
		rows, err = DB.Query(
			"SELECT id, division_id, nombre, logo_url FROM equipos WHERE division_id = $1 AND LOWER(nombre) LIKE LOWER($2)",
			divisionID, "%"+q+"%")
	} else {
		rows, err = DB.Query(
			"SELECT id, division_id, nombre, logo_url FROM equipos WHERE division_id = $1", divisionID)
	}

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
	//validacion de parametros
	if e.Nombre == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El nombre del equipo es requerido"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos invalidos"})
		return
	}

	// Validacion de parametros
	if e.Nombre == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El nombre del equipo es requerido"})
		return
	}

	result, err := DB.Exec(
		"UPDATE equipos SET nombre=$1, logo_url=$2 WHERE id=$3",
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

// POST /equipos/:id/imagen
func UploadLogoEquipo(c *gin.Context) {
	id := c.Param("id")

	file, err := c.FormFile("imagen")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se recibió ninguna imagen"})
		return
	}

	ext := filepath.Ext(file.Filename)
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
	if !allowed[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Solo se permiten imágenes jpg, png o webp"})
		return
	}

	if file.Size > 1<<20 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "La imagen no puede superar 1MB"})
		return
	}

	// Leer el archivo y convertirlo a Base64
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al leer la imagen"})
		return
	}
	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar la imagen"})
		return
	}

	mimeType := "image/jpeg"
	if ext == ".png" {
		mimeType = "image/png"
	} else if ext == ".webp" {
		mimeType = "image/webp"
	}

	base64Str := fmt.Sprintf("data:%s;base64,%s", mimeType, base64.StdEncoding.EncodeToString(data))

	_, err = DB.Exec("UPDATE equipos SET logo_url=$1 WHERE id=$2", base64Str, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"logo_url": base64Str})
}
