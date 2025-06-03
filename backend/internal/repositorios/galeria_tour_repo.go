package repositorios

import (
	"database/sql"
	"fmt"
	"sistema-toursseft/internal/entidades"
)

type GaleriaTourRepo struct {
	DB *sql.DB
}

func NewGaleriaTourRepo(db *sql.DB) *GaleriaTourRepo {
	return &GaleriaTourRepo{DB: db}
}

func (r *GaleriaTourRepo) Crear(galeria *entidades.GaleriaTour) (int, error) {
	query := `
		INSERT INTO galeria_tour (id_tipo_tour, url_imagen, descripcion, orden)
		VALUES ($1, $2, $3, $4)
		RETURNING id_galeria
	`
	var id int
	err := r.DB.QueryRow(query, galeria.IDTipoTour, galeria.URLImagen, galeria.Descripcion, galeria.Orden).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error al crear galería de tour: %v", err)
	}
	return id, nil
}

func (r *GaleriaTourRepo) ObtenerPorID(id int) (*entidades.GaleriaTour, error) {
	query := `
		SELECT id_galeria, id_tipo_tour, url_imagen, descripcion, orden, fecha_creacion, eliminado
		FROM galeria_tour
		WHERE id_galeria = $1 AND eliminado = false
	`
	var galeria entidades.GaleriaTour
	err := r.DB.QueryRow(query, id).Scan(
		&galeria.ID, &galeria.IDTipoTour, &galeria.URLImagen,
		&galeria.Descripcion, &galeria.Orden, &galeria.FechaCreacion, &galeria.Eliminado,
	)
	if err != nil {
		return nil, fmt.Errorf("error al obtener galería de tour: %v", err)
	}
	return &galeria, nil
}

func (r *GaleriaTourRepo) ListarPorTipoTour(idTipoTour int) ([]*entidades.GaleriaTour, error) {
	query := `
		SELECT id_galeria, id_tipo_tour, url_imagen, descripcion, orden, fecha_creacion, eliminado
		FROM galeria_tour
		WHERE id_tipo_tour = $1 AND eliminado = false
		ORDER BY orden ASC
	`
	rows, err := r.DB.Query(query, idTipoTour)
	if err != nil {
		return nil, fmt.Errorf("error al listar galerías de tour: %v", err)
	}
	defer rows.Close()

	var galerias []*entidades.GaleriaTour
	for rows.Next() {
		var galeria entidades.GaleriaTour
		err := rows.Scan(
			&galeria.ID, &galeria.IDTipoTour, &galeria.URLImagen,
			&galeria.Descripcion, &galeria.Orden, &galeria.FechaCreacion, &galeria.Eliminado,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear galería de tour: %v", err)
		}
		galerias = append(galerias, &galeria)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar galerías de tour: %v", err)
	}
	return galerias, nil
}

func (r *GaleriaTourRepo) Actualizar(galeria *entidades.GaleriaTour) error {
	query := `
		UPDATE galeria_tour
		SET url_imagen = $1, descripcion = $2, orden = $3
		WHERE id_galeria = $4 AND eliminado = false
	`
	_, err := r.DB.Exec(query, galeria.URLImagen, galeria.Descripcion, galeria.Orden, galeria.ID)
	if err != nil {
		return fmt.Errorf("error al actualizar galería de tour: %v", err)
	}
	return nil
}

func (r *GaleriaTourRepo) Eliminar(id int) error {
	query := `UPDATE galeria_tour SET eliminado = true WHERE id_galeria = $1`
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error al eliminar galería de tour: %v", err)
	}
	return nil
}

func (r *GaleriaTourRepo) EliminarPorTipoTour(idTipoTour int) error {
	query := `UPDATE galeria_tour SET eliminado = true WHERE id_tipo_tour = $1`
	_, err := r.DB.Exec(query, idTipoTour)
	if err != nil {
		return fmt.Errorf("error al eliminar galerías por tipo de tour: %v", err)
	}
	return nil
}
