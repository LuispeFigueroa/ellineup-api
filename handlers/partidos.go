package handlers

import (
	"database/sql"
	"net/http"

	"github.com/LuispeFigueroa/ellineup-api/models"
	"github.com/gin-gonic/gin"
)

//GET /divisiones/id/partidos
//obtener todos los partidos de una division

func GetPartidos(c *gin.Context) {
	divisionID := c.Param("id")

	rows, err := DB.Query(`
		SELECT id, division_id, equipo_local_id, equipo_visita_id, 
		carreras_local, carreras_visita, campo, fecha, estado 
		FROM partidos WHERE division_id = $1 ORDER BY fecha DESC`, divisionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	partidos := []models.Partido{}
	for rows.Next() {
		var p models.Partido
		rows.Scan(&p.ID, &p.DivisionID, &p.EquipoLocalID, &p.EquipoVisitaID,
			&p.CarrerasLocal, &p.CarrerasVisita, &p.Campo, &p.Fecha, &p.Estado)
		partidos = append(partidos, p)
	}
	c.JSON(http.StatusOK, partidos)

}

//Get /partidos/id
//obtener un partido especifico

func GetPartido(c *gin.Context) {
	id := c.Param("id")
	var p models.Partido

	err := DB.QueryRow(`
		SELECT id, division_id, equipo_local_id, equipo_visita_id,
		carreras_local, carreras_visita, campo, fecha, estado
		FROM partidos WHERE id = $1`, id).
		Scan(&p.ID, &p.DivisionID, &p.EquipoLocalID, &p.EquipoVisitaID,
			&p.CarrerasLocal, &p.CarrerasVisita, &p.Campo, &p.Fecha, &p.Estado)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "No se encontro el partido"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}

//POST /divisiones/id/partidos
//PUBLICAR resultados de un partido

func CreatePartido(c *gin.Context) {
	divisionID := c.Param("id")
	var p models.Partido

	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos no validos"})
		return
	}
	if p.Estado == "" {
		p.Estado = "programado"
	}

	err := DB.QueryRow(`
		INSERT INTO partidos (division_id, equipo_local_id, equipo_visita_id, carreras_local, carreras_visita, campo, fecha, estado)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
		divisionID, p.EquipoLocalID, p.EquipoVisitaID,
		p.CarrerasLocal, p.CarrerasVisita, p.Campo, p.Fecha, p.Estado,
	).Scan(&p.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, p)
}

// PUT /partido/id
// modificar resultados de un paritod
func UpdatePartido(c *gin.Context) {
	id := c.Param("id")
	var p models.Partido

	if err := c.ShouldBindBodyWithJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos no validos"})
		return
	}

	result, err := DB.Exec(`
		UPDATE partidos SET equipo_local_id=$1, equipo_visita_id=$2,
		carreras_local=$3, carreras_visita=$4, campo=$5, fecha=$6, estado=$7
		WHERE id=$8`,
		p.EquipoLocalID, p.EquipoVisitaID,
		p.CarrerasLocal, p.CarrerasVisita, p.Campo, p.Fecha, p.Estado, id,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No se encontro el partido"})
		return
	}
	c.JSON(http.StatusOK, p)
}

// DELETE /partidos/id
// elminar un partido
func DeletePartido(c *gin.Context) {
	id := c.Param("id")

	result, err := DB.Exec("DELETE FROM partidos WHERE id=$1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Partido no encontrado"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// GET /divisiones/:id/standings
func GetStandings(c *gin.Context) {
	divisionID := c.Param("id")

	rows, err := DB.Query(`
		SELECT 
			e.id,
			e.nombre,
			COUNT(p.id) as juegos_jugados,
			SUM(CASE 
				WHEN p.equipo_local_id = e.id AND p.carreras_local > p.carreras_visita THEN 1
				WHEN p.equipo_visita_id = e.id AND p.carreras_visita > p.carreras_local THEN 1
				ELSE 0 END) as ganados,
			SUM(CASE 
				WHEN p.equipo_local_id = e.id AND p.carreras_local < p.carreras_visita THEN 1
				WHEN p.equipo_visita_id = e.id AND p.carreras_visita < p.carreras_local THEN 1
				ELSE 0 END) as perdidos,
			SUM(CASE 
				WHEN p.equipo_local_id = e.id THEN p.carreras_local
				WHEN p.equipo_visita_id = e.id THEN p.carreras_visita
				ELSE 0 END) as carreras_anotadas
		FROM equipos e
		LEFT JOIN partidos p ON (p.equipo_local_id = e.id OR p.equipo_visita_id = e.id)
			AND p.estado = 'final'
			AND p.division_id = $1
		WHERE e.division_id = $1
		GROUP BY e.id, e.nombre
		ORDER BY ganados DESC, carreras_anotadas DESC
	`, divisionID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	type Standing struct {
		ID               int    `json:"id"`
		Nombre           string `json:"nombre"`
		JuegosJugados    int    `json:"juegos_jugados"`
		Ganados          int    `json:"ganados"`
		Perdidos         int    `json:"perdidos"`
		CarrerasAnotadas int    `json:"carreras_anotadas"`
	}

	standings := []Standing{}
	for rows.Next() {
		var s Standing
		rows.Scan(&s.ID, &s.Nombre, &s.JuegosJugados, &s.Ganados, &s.Perdidos, &s.CarrerasAnotadas)
		standings = append(standings, s)
	}

	c.JSON(http.StatusOK, standings)
}
