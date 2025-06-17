package servicios

import (
	"errors"
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/repositorios"
	"time"
)

// PagoService maneja la lógica de negocio para pagos
type PagoService struct {
	pagoRepo       *repositorios.PagoRepository
	reservaRepo    *repositorios.ReservaRepository
	metodoPagoRepo *repositorios.MetodoPagoRepository
	canalVentaRepo *repositorios.CanalVentaRepository
	sedeRepo       *repositorios.SedeRepository // Añadido repositorio de sede
}

// NewPagoService crea una nueva instancia de PagoService
func NewPagoService(
	pagoRepo *repositorios.PagoRepository,
	reservaRepo *repositorios.ReservaRepository,
	metodoPagoRepo *repositorios.MetodoPagoRepository,
	canalVentaRepo *repositorios.CanalVentaRepository,
	sedeRepo *repositorios.SedeRepository, // Añadido repositorio de sede
) *PagoService {
	return &PagoService{
		pagoRepo:       pagoRepo,
		reservaRepo:    reservaRepo,
		metodoPagoRepo: metodoPagoRepo,
		canalVentaRepo: canalVentaRepo,
		sedeRepo:       sedeRepo, // Asignado repositorio de sede
	}
}

// Create crea un nuevo pago
func (s *PagoService) Create(pago *entidades.NuevoPagoRequest) (int, error) {
	// Verificar que la reserva existe
	reserva, err := s.reservaRepo.GetByID(pago.IDReserva)
	if err != nil {
		return 0, errors.New("la reserva especificada no existe")
	}

	// Verificar que la reserva no esté cancelada
	if reserva.Estado == "CANCELADA" {
		return 0, errors.New("no se puede registrar un pago para una reserva cancelada")
	}

	// Verificar que el método de pago existe
	_, err = s.metodoPagoRepo.GetByID(pago.IDMetodoPago)
	if err != nil {
		return 0, errors.New("el método de pago especificado no existe")
	}

	// Verificar que el canal de venta existe
	_, err = s.canalVentaRepo.GetByID(pago.IDCanal)
	if err != nil {
		return 0, errors.New("el canal de venta especificado no existe")
	}

	// Verificar que la sede existe
	_, err = s.sedeRepo.GetByID(pago.IDSede)
	if err != nil {
		return 0, errors.New("la sede especificada no existe")
	}

	// Verificar que el monto sea positivo
	if pago.Monto <= 0 {
		return 0, errors.New("el monto del pago debe ser mayor a cero")
	}

	// Verificar que el monto total pagado + el nuevo pago no exceda el total a pagar de la reserva
	totalPagado, err := s.pagoRepo.GetTotalPagadoByReserva(pago.IDReserva)
	if err != nil {
		return 0, err
	}

	if totalPagado+pago.Monto > reserva.TotalPagar {
		return 0, errors.New("el monto total pagado excedería el total a pagar de la reserva")
	}

	// Crear pago
	return s.pagoRepo.Create(pago)
}

// GetByID obtiene un pago por su ID
func (s *PagoService) GetByID(id int) (*entidades.Pago, error) {
	return s.pagoRepo.GetByID(id)
}

// Update actualiza un pago existente
func (s *PagoService) Update(id int, pago *entidades.ActualizarPagoRequest) error {
	// Verificar que el pago existe
	existingPago, err := s.pagoRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Verificar que el método de pago existe
	_, err = s.metodoPagoRepo.GetByID(pago.IDMetodoPago)
	if err != nil {
		return errors.New("el método de pago especificado no existe")
	}

	// Verificar que el canal de venta existe
	_, err = s.canalVentaRepo.GetByID(pago.IDCanal)
	if err != nil {
		return errors.New("el canal de venta especificado no existe")
	}

	// Verificar que la sede existe
	_, err = s.sedeRepo.GetByID(pago.IDSede)
	if err != nil {
		return errors.New("la sede especificada no existe")
	}

	// Verificar que el monto sea positivo
	if pago.Monto <= 0 {
		return errors.New("el monto del pago debe ser mayor a cero")
	}

	// Si cambia el monto, verificar que el total pagado no exceda el total a pagar
	if pago.Monto != existingPago.Monto {
		// Obtener el total pagado sin considerar este pago
		totalPagado, err := s.pagoRepo.GetTotalPagadoByReserva(existingPago.IDReserva)
		if err != nil {
			return err
		}

		// Restar el monto del pago actual si está procesado
		if existingPago.Estado == "PROCESADO" {
			totalPagado -= existingPago.Monto
		}

		// Verificar si el nuevo monto excedería el total a pagar
		reserva, err := s.reservaRepo.GetByID(existingPago.IDReserva)
		if err != nil {
			return err
		}

		if pago.Estado == "PROCESADO" && totalPagado+pago.Monto > reserva.TotalPagar {
			return errors.New("el monto total pagado excedería el total a pagar de la reserva")
		}
	}

	// Actualizar pago
	return s.pagoRepo.Update(id, pago)
}

// CambiarEstado cambia el estado de un pago
func (s *PagoService) CambiarEstado(id int, estado string) error {
	// Verificar estado válido
	if estado != "PROCESADO" && estado != "ANULADO" {
		return errors.New("estado inválido, debe ser PROCESADO o ANULADO")
	}

	// Verificar que el pago existe
	pago, err := s.pagoRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Si ya tiene ese estado, no hacer nada
	if pago.Estado == estado {
		return nil
	}

	// Si está anulando el pago, verificar que no tenga comprobantes asociados
	if estado == "ANULADO" {
		// Lógica para verificar comprobantes
		// (Esto dependería de la estructura de tu aplicación)
	}

	// Cambiar estado
	return s.pagoRepo.UpdateEstado(id, estado)
}

// Delete elimina un pago
func (s *PagoService) Delete(id int) error {
	// Verificar que el pago existe
	_, err := s.pagoRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Eliminar pago
	return s.pagoRepo.Delete(id)
}

// List lista todos los pagos
func (s *PagoService) List() ([]*entidades.Pago, error) {
	return s.pagoRepo.List()
}

// ListByReserva lista todos los pagos de una reserva específica
func (s *PagoService) ListByReserva(idReserva int) ([]*entidades.Pago, error) {
	// Verificar que la reserva existe
	_, err := s.reservaRepo.GetByID(idReserva)
	if err != nil {
		return nil, errors.New("la reserva especificada no existe")
	}

	// Listar pagos por reserva
	return s.pagoRepo.ListByReserva(idReserva)
}

// ListByFecha lista todos los pagos de una fecha específica
func (s *PagoService) ListByFecha(fecha time.Time) ([]*entidades.Pago, error) {
	return s.pagoRepo.ListByFecha(fecha)
}

// GetTotalPagadoByReserva obtiene el total pagado de una reserva específica
func (s *PagoService) GetTotalPagadoByReserva(idReserva int) (float64, error) {
	// Verificar que la reserva existe
	_, err := s.reservaRepo.GetByID(idReserva)
	if err != nil {
		return 0, errors.New("la reserva especificada no existe")
	}

	// Obtener total pagado
	return s.pagoRepo.GetTotalPagadoByReserva(idReserva)
}

// ListByEstado lista todos los pagos con un estado específico
func (s *PagoService) ListByEstado(estado string) ([]*entidades.Pago, error) {
	// Verificar estado válido
	if estado != "PROCESADO" && estado != "ANULADO" {
		return nil, errors.New("estado inválido, debe ser PROCESADO o ANULADO")
	}

	// Listar pagos por estado
	return s.pagoRepo.ListByEstado(estado)
}

// ListByCliente lista todos los pagos relacionados con un cliente específico
func (s *PagoService) ListByCliente(idCliente int) ([]*entidades.Pago, error) {
	// Verificar que el cliente existe
	// Esto dependería de cómo tienes organizados tus servicios

	// Listar pagos por cliente
	return s.pagoRepo.ListByCliente(idCliente)
}

// ListBySede lista todos los pagos de una sede específica
func (s *PagoService) ListBySede(idSede int) ([]*entidades.Pago, error) {
	// Verificar que la sede existe
	_, err := s.sedeRepo.GetByID(idSede)
	if err != nil {
		return nil, errors.New("la sede especificada no existe")
	}

	// Listar pagos por sede
	return s.pagoRepo.ListBySede(idSede)
}

// UpdateEstado actualiza el estado de un pago
func (s *PagoService) UpdateEstado(id int, estado string) error {
	// Verificar que el pago existe
	pago, err := s.pagoRepo.GetByID(id)
	if err != nil {
		return errors.New("pago no encontrado")
	}

	// Validar que el estado sea uno de los permitidos
	estadosPermitidos := map[string]bool{
		"PENDIENTE": true,
		"PROCESADO": true,
		"RECHAZADO": true,
		"DEVUELTO":  true,
		"CANCELADO": true,
	}

	if !estadosPermitidos[estado] {
		return errors.New("estado de pago no válido")
	}

	// Si ya tiene ese estado, no hacer nada
	if pago.Estado == estado {
		return nil
	}

	// Actualizar estado en la base de datos
	return s.pagoRepo.UpdateEstado(id, estado)
}
