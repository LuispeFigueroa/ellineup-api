package models

type Jugador struct {
	ID       int    `json:"id"`
	EquipoID int    `json:"equipo_id"`
	Nombre   string `json:"nombre"`
	Numero   int    `json:"numero"`
	Posicion string `json:"posicion"`
}
