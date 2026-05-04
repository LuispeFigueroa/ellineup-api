package handlers

import (
	"database/sql"
	"net/http"

	"github.com/LuispeFigueroa/ellineup-api/models"
	"github.com/gin-gonic/gin"
)

//GET /equipos/:id/jugadores
//obtener todos los jugadores de un equipo

func GetJugadores(c *gin.Context) {
	equipoID := c.Param("id")

	rows, err := DB.Query("SELECT id, equipo_id, nombre, numero, posicion FROM jugadores WHERE equipo_id = $1", equipoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	jugadores := []models.Jugador{}
	for rows.Next() {
		var j models.Jugador
		rows.Scan(&j.ID, &j.EquipoID, &j.Nombre, &j.Numero, &j.Posicion)
		jugadores = append(jugadores, j)
	}
	c.JSON(http.StatusOK, jugadores)
}

// GET /jugadores/id
//buscar a un jugador en especifico por su ID

func GetJugador(c *gin.Context) {
	id := c.Param("id")
	var j models.Jugador

	err := DB.QueryRow("SELECT id, equipo_id, nombre, numero, posicion FROM jugadores WHERE id = $1", id).
		Scan(&j.ID, &j.EquipoID, &j.Nombre, &j.Numero, &j.Posicion)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Jugador no encontrado"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, j)
}

//POST /equipos/id/Jugadores
//Crear un nuevo jugador para un equip

func CreateJugador(c *gin.Context) {
	equipoID := c.Param("id")
	var j models.Jugador

	if err := c.ShouldBindJSON(&j); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos invalidos"})
		return
	}

	err := DB.QueryRow(
		"INSERT INTO jugadores (equipo_id, nombre, numero, posicion) VALUES ($1, $2, $3, $4) RETURNING id",
		equipoID, j.Nombre, j.Numero, j.Posicion,
	).Scan(&j.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, j)
}

// PUT /jugadores/:id
// modificar un jugador especifico
func UpdateJugador(c *gin.Context) {
	id := c.Param("id")
	var j models.Jugador

	if err := c.ShouldBindJSON(&j); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos no validos"})
		return
	}

	result, err := DB.Exec(
		"UPDATE jugadores SET nombre=$1, numero=$2, posicion=$3 WHERE id=$4",
		j.Nombre, j.Numero, j.Posicion, id,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Jugador no encontrado"})
		return
	}

	c.JSON(http.StatusOK, j)
}

// DELETE /jugadores/:id
// eliminar a un jugador
func DeleteJugador(c *gin.Context) {
	id := c.Param("id")

	result, err := DB.Exec("DELETE FROM jugadores WHERE id=$1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Jugador no encontrado"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
