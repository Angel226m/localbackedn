package servicios

import (
	"database/sql"
	"errors"
	"sistema-tours/internal/entidades"
	"sistema-tours/internal/repositorios"
	"time"
)

// ReservaService maneja la lógica de negocio para reservas
type ReservaService struct {
	db                 *sql.DB
	reservaRepo        *repositorios.ReservaRepository
	clienteRepo        *repositorios.ClienteRepository
	tourProgramadoRepo *repositorios.TourProgramadoRepository
	canalVentaRepo     *repositorios.CanalVentaRepository
	tipoPasajeRepo     *repositorios.TipoPasajeRepository
	usuarioRepo        *repositorios.UsuarioRepository
}

// NewReservaService crea una nueva instancia de ReservaService
func NewReservaService(
	db *sql.DB,
	reservaRepo *repositorios.ReservaRepository,
	clienteRepo *repositorios.ClienteRepository,
	tourProgramadoRepo *repositorios.TourProgramadoRepository,
	canalVentaRepo *repositorios.CanalVentaRepository,
	tipoPasajeRepo *repositorios.TipoPasajeRepository,
	usuarioRepo *repositorios.UsuarioRepository,
) *ReservaService {
	return &ReservaService{
		db:                 db,
		reservaRepo:        reservaRepo,
		clienteRepo:        clienteRepo,
		tourProgramadoRepo: tourProgramadoRepo,
		canalVentaRepo:     canalVentaRepo,
		tipoPasajeRepo:     tipoPasajeRepo,
		usuarioRepo:        usuarioRepo,
	}
}

// Create crea una nueva reserva
func (s *ReservaService) Create(reserva *entidades.NuevaReservaRequest) (int, error) {
	// Verificar que el cliente existe
	_, err := s.clienteRepo.GetByID(reserva.IDCliente)
	if err != nil {
		return 0, errors.New("el cliente especificado no existe")
	}

	// Verificar que el tour programado existe
	tourProgramado, err := s.tourProgramadoRepo.GetByID(reserva.IDTourProgramado)
	if err != nil {
		return 0, errors.New("el tour programado especificado no existe")
	}

	// Verificar que el tour programado está en estado PROGRAMADO
	if tourProgramado.Estado != "PROGRAMADO" {
		return 0, errors.New("no se puede reservar en un tour que no está programado")
	}

	// Verificar que el canal de venta existe
	_, err = s.canalVentaRepo.GetByID(reserva.IDCanal)
	if err != nil {
		return 0, errors.New("el canal de venta especificado no existe")
	}

	// Si se especifica un vendedor, verificar que existe y es vendedor
	if reserva.IDVendedor != nil {
		usuario, err := s.usuarioRepo.GetByID(*reserva.IDVendedor)
		if err != nil {
			return 0, errors.New("el vendedor especificado no existe")
		}
		if usuario.Rol != "VENDEDOR" && usuario.Rol != "ADMIN" {
			return 0, errors.New("el usuario especificado no es un vendedor")
		}
	}

	// Verificar que los tipos de pasaje existen
	totalPasajeros := 0
	for _, pasaje := range reserva.CantidadPasajes {
		_, err := s.tipoPasajeRepo.GetByID(pasaje.IDTipoPasaje)
		if err != nil {
			return 0, errors.New("uno de los tipos de pasaje especificados no existe")
		}
		totalPasajeros += pasaje.Cantidad
	}

	// Verificar disponibilidad de cupo
	if totalPasajeros > tourProgramado.CupoDisponible {
		return 0, errors.New("no hay suficiente cupo disponible para la cantidad de pasajeros solicitada")
	}

	// Iniciar transacción
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Crear reserva
	id, err := s.reservaRepo.Create(tx, reserva)
	if err != nil {
		return 0, err
	}

	// Actualizar cupo disponible del tour programado
	nuevoCupo := tourProgramado.CupoDisponible - totalPasajeros
	err = s.tourProgramadoRepo.UpdateCupoDisponible(reserva.IDTourProgramado, nuevoCupo)
	if err != nil {
		return 0, err
	}

	// Commit de la transacción
	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetByID obtiene una reserva por su ID
func (s *ReservaService) GetByID(id int) (*entidades.Reserva, error) {
	return s.reservaRepo.GetByID(id)
}

// Update actualiza una reserva existente
func (s *ReservaService) Update(id int, reserva *entidades.ActualizarReservaRequest) error {
	// Verificar que la reserva existe
	existingReserva, err := s.reservaRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Verificar que el cliente existe
	_, err = s.clienteRepo.GetByID(reserva.IDCliente)
	if err != nil {
		return errors.New("el cliente especificado no existe")
	}

	// Verificar que el tour programado existe
	tourProgramado, err := s.tourProgramadoRepo.GetByID(reserva.IDTourProgramado)
	if err != nil {
		return errors.New("el tour programado especificado no existe")
	}

	// Verificar que el canal de venta existe
	_, err = s.canalVentaRepo.GetByID(reserva.IDCanal)
	if err != nil {
		return errors.New("el canal de venta especificado no existe")
	}

	// Si se especifica un vendedor, verificar que existe y es vendedor
	if reserva.IDVendedor != nil {
		usuario, err := s.usuarioRepo.GetByID(*reserva.IDVendedor)
		if err != nil {
			return errors.New("el vendedor especificado no existe")
		}
		if usuario.Rol != "VENDEDOR" && usuario.Rol != "ADMIN" {
			return errors.New("el usuario especificado no es un vendedor")
		}
	}

	// Verificar que los tipos de pasaje existen
	totalPasajerosNuevo := 0
	for _, pasaje := range reserva.CantidadPasajes {
		_, err := s.tipoPasajeRepo.GetByID(pasaje.IDTipoPasaje)
		if err != nil {
			return errors.New("uno de los tipos de pasaje especificados no existe")
		}
		totalPasajerosNuevo += pasaje.Cantidad
	}

	// Obtener la cantidad actual de pasajeros en la reserva
	totalPasajerosActual, err := s.reservaRepo.GetCantidadPasajerosByReserva(id)
	if err != nil {
		return err
	}

	// Calcular diferencia de pasajeros
	diferenciaPasajeros := totalPasajerosNuevo - totalPasajerosActual

	// Si es el mismo tour programado, verificar disponibilidad de cupo considerando la diferencia
	if reserva.IDTourProgramado == existingReserva.IDTourProgramado {
		if diferenciaPasajeros > 0 && diferenciaPasajeros > tourProgramado.CupoDisponible {
			return errors.New("no hay suficiente cupo disponible para aumentar la cantidad de pasajeros")
		}
	} else {
		// Si es otro tour programado, verificar disponibilidad total
		if totalPasajerosNuevo > tourProgramado.CupoDisponible {
			return errors.New("no hay suficiente cupo disponible en el nuevo tour programado")
		}
	}

	// Iniciar transacción
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Actualizar reserva
	err = s.reservaRepo.Update(tx, id, reserva)
	if err != nil {
		return err
	}

	// Si cambió el tour programado, actualizar cupos de ambos tours
	if reserva.IDTourProgramado != existingReserva.IDTourProgramado {
		// Liberar cupo en el tour anterior
		tourAnterior, err := s.tourProgramadoRepo.GetByID(existingReserva.IDTourProgramado)
		if err != nil {
			return err
		}
		nuevoCupoAnterior := tourAnterior.CupoDisponible + totalPasajerosActual
		err = s.tourProgramadoRepo.UpdateCupoDisponible(existingReserva.IDTourProgramado, nuevoCupoAnterior)
		if err != nil {
			return err
		}

		// Reservar cupo en el nuevo tour
		nuevoCupoNuevo := tourProgramado.CupoDisponible - totalPasajerosNuevo
		err = s.tourProgramadoRepo.UpdateCupoDisponible(reserva.IDTourProgramado, nuevoCupoNuevo)
		if err != nil {
			return err
		}
	} else if diferenciaPasajeros != 0 {
		// Si es el mismo tour pero cambió la cantidad de pasajeros, actualizar cupo
		nuevoCupo := tourProgramado.CupoDisponible - diferenciaPasajeros
		err = s.tourProgramadoRepo.UpdateCupoDisponible(reserva.IDTourProgramado, nuevoCupo)
		if err != nil {
			return err
		}
	}

	// Commit de la transacción
	return tx.Commit()
}

// CambiarEstado cambia el estado de una reserva
func (s *ReservaService) CambiarEstado(id int, estado string) error {
	// Verificar que la reserva existe
	reserva, err := s.reservaRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Verificar que el estado es válido
	if estado != "RESERVADO" && estado != "CANCELADA" {
		return errors.New("estado de reserva inválido")
	}

	// Si se está cancelando una reserva, liberar el cupo
	if estado == "CANCELADA" && reserva.Estado != "CANCELADA" {
		// Obtener la cantidad de pasajeros en la reserva
		totalPasajeros, err := s.reservaRepo.GetCantidadPasajerosByReserva(id)
		if err != nil {
			return err
		}

		// Liberar cupo en el tour programado
		tourProgramado, err := s.tourProgramadoRepo.GetByID(reserva.IDTourProgramado)
		if err != nil {
			return err
		}
		nuevoCupo := tourProgramado.CupoDisponible + totalPasajeros
		err = s.tourProgramadoRepo.UpdateCupoDisponible(reserva.IDTourProgramado, nuevoCupo)
		if err != nil {
			return err
		}
	}

	// Si se está reactivando una reserva cancelada, verificar disponibilidad y reservar cupo
	if estado == "RESERVADO" && reserva.Estado == "CANCELADA" {
		// Obtener la cantidad de pasajeros en la reserva
		totalPasajeros, err := s.reservaRepo.GetCantidadPasajerosByReserva(id)
		if err != nil {
			return err
		}

		// Verificar disponibilidad de cupo
		tourProgramado, err := s.tourProgramadoRepo.GetByID(reserva.IDTourProgramado)
		if err != nil {
			return err
		}
		if totalPasajeros > tourProgramado.CupoDisponible {
			return errors.New("no hay suficiente cupo disponible para reactivar la reserva")
		}

		// Reservar cupo
		nuevoCupo := tourProgramado.CupoDisponible - totalPasajeros
		err = s.tourProgramadoRepo.UpdateCupoDisponible(reserva.IDTourProgramado, nuevoCupo)
		if err != nil {
			return err
		}
	}

	// Actualizar estado de la reserva
	return s.reservaRepo.UpdateEstado(id, estado)
}

// Delete elimina una reserva
func (s *ReservaService) Delete(id int) error {
	// Verificar que la reserva existe
	reserva, err := s.reservaRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Si la reserva no está cancelada, liberar el cupo
	if reserva.Estado != "CANCELADA" {
		// Obtener la cantidad de pasajeros en la reserva
		totalPasajeros, err := s.reservaRepo.GetCantidadPasajerosByReserva(id)
		if err != nil {
			return err
		}

		// Liberar cupo en el tour programado
		tourProgramado, err := s.tourProgramadoRepo.GetByID(reserva.IDTourProgramado)
		if err != nil {
			return err
		}
		nuevoCupo := tourProgramado.CupoDisponible + totalPasajeros
		err = s.tourProgramadoRepo.UpdateCupoDisponible(reserva.IDTourProgramado, nuevoCupo)
		if err != nil {
			return err
		}
	}

	// Eliminar reserva
	return s.reservaRepo.Delete(id)
}

// List lista todas las reservas
func (s *ReservaService) List() ([]*entidades.Reserva, error) {
	return s.reservaRepo.List()
}

// ListByCliente lista todas las reservas de un cliente
func (s *ReservaService) ListByCliente(idCliente int) ([]*entidades.Reserva, error) {
	// Verificar que el cliente existe
	_, err := s.clienteRepo.GetByID(idCliente)
	if err != nil {
		return nil, errors.New("el cliente especificado no existe")
	}

	return s.reservaRepo.ListByCliente(idCliente)
}

// ListByTourProgramado lista todas las reservas para un tour programado
func (s *ReservaService) ListByTourProgramado(idTourProgramado int) ([]*entidades.Reserva, error) {
	// Verificar que el tour programado existe
	_, err := s.tourProgramadoRepo.GetByID(idTourProgramado)
	if err != nil {
		return nil, errors.New("el tour programado especificado no existe")
	}

	return s.reservaRepo.ListByTourProgramado(idTourProgramado)
}

// ListByFecha lista todas las reservas para una fecha específica
func (s *ReservaService) ListByFecha(fecha time.Time) ([]*entidades.Reserva, error) {
	return s.reservaRepo.ListByFecha(fecha)
}

// ListByEstado lista todas las reservas por estado
func (s *ReservaService) ListByEstado(estado string) ([]*entidades.Reserva, error) {
	// Verificar que el estado es válido
	if estado != "RESERVADO" && estado != "CANCELADA" {
		return nil, errors.New("estado de reserva inválido")
	}

	return s.reservaRepo.ListByEstado(estado)
}
