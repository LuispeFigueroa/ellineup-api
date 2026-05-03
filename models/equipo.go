package models

type Equipo struct {
	ID         int    `json:"id"`
	DivisionID int    `json:"division_id"`
	Nombre     string `json:"nombre"`
	LogoURL    string `json:"logo_url"`
}
