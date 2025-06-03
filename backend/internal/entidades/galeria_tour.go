package entidades

import "time"

type GaleriaTour struct {
	ID            int       `json:"id_galeria" db:"id_galeria"`
	IDTipoTour    int       `json:"id_tipo_tour" db:"id_tipo_tour"`
	URLImagen     string    `json:"url_imagen" db:"url_imagen"`
	Descripcion   string    `json:"descripcion" db:"descripcion"`
	Orden         int       `json:"orden" db:"orden"`
	FechaCreacion time.Time `json:"fecha_creacion" db:"fecha_creacion"`
	Eliminado     bool      `json:"eliminado" db:"eliminado"`
}

type GaleriaTourRequest struct {
	IDTipoTour  int    `json:"id_tipo_tour" validate:"required"`
	URLImagen   string `json:"url_imagen" validate:"required"`
	Descripcion string `json:"descripcion"`
	Orden       int    `json:"orden"`
}

type GaleriaTourUpdateRequest struct {
	URLImagen   string `json:"url_imagen" validate:"required"`
	Descripcion string `json:"descripcion"`
	Orden       int    `json:"orden"`
}
