package servicios

/*
import (
	"database/sql"
	"errors"
	"fmt"
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/repositorios"
	"time"
)

// ReservaService maneja la lógica de negocio para reservas
// Coordina las operaciones entre el repositorio y las reglas de negocio
type ReservaService struct {
	db                 *sql.DB
	reservaRepo        *repositorios.ReservaRepository
	clienteRepo        *repositorios.ClienteRepository
	tourProgramadoRepo *repositorios.TourProgramadoRepository
	canalVentaRepo     *repositorios.CanalVentaRepository
	tipoPasajeRepo     *repositorios.TipoPasajeRepository
	usuarioRepo        *repositorios.UsuarioRepository
	sedeRepo           *repositorios.SedeRepository // Añadido repositorio de sedes
}

// NewReservaService crea una nueva instancia de ReservaService
// Inicializa el servicio con todas las dependencias necesarias
func NewReservaService(
	db *sql.DB,
	reservaRepo *repositorios.ReservaRepository,
	clienteRepo *repositorios.ClienteRepository,
	tourProgramadoRepo *repositorios.TourProgramadoRepository,
	canalVentaRepo *repositorios.CanalVentaRepository,
	tipoPasajeRepo *repositorios.TipoPasajeRepository,
	usuarioRepo *repositorios.UsuarioRepository,
	sedeRepo *repositorios.SedeRepository, // Añadido repositorio de sedes
) *ReservaService {
	return &ReservaService{
		db:                 db,
		reservaRepo:        reservaRepo,
		clienteRepo:        clienteRepo,
		tourProgramadoRepo: tourProgramadoRepo,
		canalVentaRepo:     canalVentaRepo,
		tipoPasajeRepo:     tipoPasajeRepo,
		usuarioRepo:        usuarioRepo,
		sedeRepo:           sedeRepo, // Asignado nuevo repositorio
	}
}

// Create crea una nueva reserva
// Valida todos los datos y realiza las operaciones necesarias en la base de datos
// Retorna el ID de la reserva creada o un error si falla
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

	// Verificar que la sede existe
	_, err = s.sedeRepo.GetByID(reserva.IDSede)
	if err != nil {
		return 0, errors.New("la sede especificada no existe")
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
// Retorna la reserva completa con todos sus datos relacionados
func (s *ReservaService) GetByID(id int) (*entidades.Reserva, error) {
	return s.reservaRepo.GetByID(id)
}

// Update actualiza una reserva existente
// Valida todos los datos y actualiza la información en la base de datos
// Maneja la lógica de cambios en el cupo de pasajeros si es necesario
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

	// Verificar que la sede existe
	_, err = s.sedeRepo.GetByID(reserva.IDSede)
	if err != nil {
		return errors.New("la sede especificada no existe")
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
// Actualiza el estado y maneja la lógica de negocio relacionada con el cambio
// Por ejemplo, libera cupos si se cancela una reserva
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

// Delete realiza una eliminación lógica de una reserva
// Verifica restricciones como pagos o comprobantes asociados
// Actualiza el cupo disponible en el tour si es necesario
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

	// Eliminar reserva (lógicamente)
	return s.reservaRepo.Delete(id)
}

// List obtiene todas las reservas activas del sistema
// Retorna un slice de reservas con toda su información relacionada
func (s *ReservaService) List() ([]*entidades.Reserva, error) {
	return s.reservaRepo.List()
}

// ListByCliente lista todas las reservas de un cliente específico
// Verifica primero que el cliente exista
func (s *ReservaService) ListByCliente(idCliente int) ([]*entidades.Reserva, error) {
	// Verificar que el cliente existe
	_, err := s.clienteRepo.GetByID(idCliente)
	if err != nil {
		return nil, errors.New("el cliente especificado no existe")
	}

	return s.reservaRepo.ListByCliente(idCliente)
}

// ListByTourProgramado lista todas las reservas para un tour programado específico
// Verifica primero que el tour programado exista
func (s *ReservaService) ListByTourProgramado(idTourProgramado int) ([]*entidades.Reserva, error) {
	// Verificar que el tour programado existe
	_, err := s.tourProgramadoRepo.GetByID(idTourProgramado)
	if err != nil {
		return nil, errors.New("el tour programado especificado no existe")
	}

	return s.reservaRepo.ListByTourProgramado(idTourProgramado)
}

// ListByFecha lista todas las reservas para una fecha específica
// Útil para ver todas las reservas de un día determinado
func (s *ReservaService) ListByFecha(fecha time.Time) ([]*entidades.Reserva, error) {
	return s.reservaRepo.ListByFecha(fecha)
}

// ListByEstado lista todas las reservas por estado específico (RESERVADO, CANCELADA)
// Verifica que el estado sea válido antes de ejecutar la consulta
func (s *ReservaService) ListByEstado(estado string) ([]*entidades.Reserva, error) {
	// Verificar que el estado es válido
	if estado != "RESERVADO" && estado != "CANCELADA" {
		return nil, errors.New("estado de reserva inválido")
	}

	return s.reservaRepo.ListByEstado(estado)
}

// ListBySede lista todas las reservas de una sede específica
// Verifica primero que la sede exista
// ListBySede lista todas las reservas de una sede específica
// Verifica primero que la sede exista
func (s *ReservaService) ListBySede(idSede int) ([]*entidades.Reserva, error) {
	// Verificar que la sede existe
	sede, err := s.sedeRepo.GetByID(idSede)
	if err != nil {
		return nil, errors.New("la sede especificada no existe")
	}

	// Verificar que la sede no está eliminada
	if sede.Eliminado {
		return nil, errors.New("la sede especificada está eliminada")
	}

	// Convertir idSede a puntero para pasarlo al repositorio
	sedeID := idSede // Crear una variable local para obtener su dirección

	// Obtener las reservas de la sede
	reservas, err := s.reservaRepo.ListBySede(&sedeID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener las reservas de la sede: %v", err)
	}

	return reservas, nil
}

// Agregar un nuevo método para listar todas las reservas (para ADMIN)
func (s *ReservaService) ListAllReservas() ([]*entidades.Reserva, error) {
	// Pasar nil para obtener todas las reservas sin filtrar por sede
	return s.reservaRepo.ListBySede(nil)
}
*/
