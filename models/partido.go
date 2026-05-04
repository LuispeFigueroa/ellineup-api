package models

type Partido struct {
	ID             int    `json:"id"`
	DivisionID     int    `json:"division_id"`
	EquipoLocalID  int    `json:"equipo_local_id"`
	EquipoVisitaID int    `json:"equipo_visita_id"`
	CarrerasLocal  int    `json:"carreras_local"`
	CarrerasVisita int    `json:"carreras_visita"`
	Campo          string `json:"campo"`
	Fecha          string `json:"fecha"`
	Estado         string `json:"estado"`
}
