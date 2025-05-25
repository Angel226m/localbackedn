package servicios

import (
	"errors"
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/repositorios"
	"time"
)

// ComprobantePagoService maneja la lógica de negocio para comprobantes de pago
type ComprobantePagoService struct {
	comprobantePagoRepo *repositorios.ComprobantePagoRepository
	reservaRepo         *repositorios.ReservaRepository
	pagoRepo            *repositorios.PagoRepository
	sedeRepo            *repositorios.SedeRepository // Añadido repositorio de sede
}

// NewComprobantePagoService crea una nueva instancia de ComprobantePagoService
func NewComprobantePagoService(
	comprobantePagoRepo *repositorios.ComprobantePagoRepository,
	reservaRepo *repositorios.ReservaRepository,
	pagoRepo *repositorios.PagoRepository,
	sedeRepo *repositorios.SedeRepository, // Añadido repositorio de sede
) *ComprobantePagoService {
	return &ComprobantePagoService{
		comprobantePagoRepo: comprobantePagoRepo,
		reservaRepo:         reservaRepo,
		pagoRepo:            pagoRepo,
		sedeRepo:            sedeRepo, // Asignado repositorio de sede
	}
}

// Create crea un nuevo comprobante de pago
func (s *ComprobantePagoService) Create(comprobante *entidades.NuevoComprobantePagoRequest) (int, error) {
	// Verificar que la reserva existe
	reserva, err := s.reservaRepo.GetByID(comprobante.IDReserva)
	if err != nil {
		return 0, errors.New("la reserva especificada no existe")
	}

	// Verificar que la reserva no esté cancelada
	if reserva.Estado == "CANCELADA" {
		return 0, errors.New("no se puede emitir un comprobante para una reserva cancelada")
	}

	// Verificar que la sede existe
	_, err = s.sedeRepo.GetByID(comprobante.IDSede)
	if err != nil {
		return 0, errors.New("la sede especificada no existe")
	}

	// Verificar que el número de comprobante no exista para este tipo
	existing, err := s.comprobantePagoRepo.GetByTipoAndNumero(comprobante.Tipo, comprobante.NumeroComprobante)
	if err == nil && existing != nil {
		return 0, errors.New("ya existe un comprobante con este tipo y número")
	}

	// Verificar que los montos sean correctos
	if comprobante.Total != comprobante.Subtotal+comprobante.IGV {
		return 0, errors.New("el total debe ser igual a subtotal + IGV")
	}

	// Verificar que el total no exceda el total a pagar de la reserva
	if comprobante.Total > reserva.TotalPagar {
		return 0, errors.New("el total del comprobante excede el total a pagar de la reserva")
	}

	// Verificar que haya pagos suficientes para cubrir el total del comprobante
	totalPagado, err := s.pagoRepo.GetTotalPagadoByReserva(comprobante.IDReserva)
	if err != nil {
		return 0, err
	}

	if totalPagado < comprobante.Total {
		return 0, errors.New("no hay pagos suficientes para cubrir el total del comprobante")
	}

	// Crear comprobante de pago
	return s.comprobantePagoRepo.Create(comprobante)
}

// GetByID obtiene un comprobante de pago por su ID
func (s *ComprobantePagoService) GetByID(id int) (*entidades.ComprobantePago, error) {
	return s.comprobantePagoRepo.GetByID(id)
}

// GetByTipoAndNumero obtiene un comprobante de pago por su tipo y número
func (s *ComprobantePagoService) GetByTipoAndNumero(tipo, numero string) (*entidades.ComprobantePago, error) {
	// Verificar tipo válido
	if tipo != "BOLETA" && tipo != "FACTURA" {
		return nil, errors.New("tipo inválido, debe ser BOLETA o FACTURA")
	}

	// Obtener comprobante
	return s.comprobantePagoRepo.GetByTipoAndNumero(tipo, numero)
}

// Update actualiza un comprobante de pago existente
func (s *ComprobantePagoService) Update(id int, comprobante *entidades.ActualizarComprobantePagoRequest) error {
	// Verificar que el comprobante existe
	existingComprobante, err := s.comprobantePagoRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Verificar que la sede existe
	_, err = s.sedeRepo.GetByID(comprobante.IDSede)
	if err != nil {
		return errors.New("la sede especificada no existe")
	}

	// Si cambia el tipo o número, verificar que no exista otro comprobante con esos datos
	if comprobante.Tipo != existingComprobante.Tipo || comprobante.NumeroComprobante != existingComprobante.NumeroComprobante {
		existing, err := s.comprobantePagoRepo.GetByTipoAndNumero(comprobante.Tipo, comprobante.NumeroComprobante)
		if err == nil && existing != nil && existing.ID != id {
			return errors.New("ya existe otro comprobante con este tipo y número")
		}
	}

	// Verificar que los montos sean correctos
	if comprobante.Total != comprobante.Subtotal+comprobante.IGV {
		return errors.New("el total debe ser igual a subtotal + IGV")
	}

	// Verificar que el total no exceda el total a pagar de la reserva
	reserva, err := s.reservaRepo.GetByID(existingComprobante.IDReserva)
	if err != nil {
		return err
	}

	if comprobante.Total > reserva.TotalPagar {
		return errors.New("el total del comprobante excede el total a pagar de la reserva")
	}

	// Verificar que haya pagos suficientes para cubrir el total del comprobante
	totalPagado, err := s.pagoRepo.GetTotalPagadoByReserva(existingComprobante.IDReserva)
	if err != nil {
		return err
	}

	if totalPagado < comprobante.Total {
		return errors.New("no hay pagos suficientes para cubrir el total del comprobante")
	}

	// Actualizar comprobante
	return s.comprobantePagoRepo.Update(id, comprobante)
}

// CambiarEstado cambia el estado de un comprobante de pago
func (s *ComprobantePagoService) CambiarEstado(id int, estado string) error {
	// Verificar estado válido
	if estado != "EMITIDO" && estado != "ANULADO" {
		return errors.New("estado inválido, debe ser EMITIDO o ANULADO")
	}

	// Verificar que el comprobante existe
	comprobante, err := s.comprobantePagoRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Si ya tiene ese estado, no hacer nada
	if comprobante.Estado == estado {
		return nil
	}

	// Cambiar estado
	return s.comprobantePagoRepo.UpdateEstado(id, estado)
}

// Delete elimina un comprobante de pago
func (s *ComprobantePagoService) Delete(id int) error {
	// Verificar que el comprobante existe
	comprobante, err := s.comprobantePagoRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Verificar que el comprobante esté anulado
	if comprobante.Estado != "ANULADO" {
		return errors.New("no se puede eliminar un comprobante que no está anulado")
	}

	// Eliminar comprobante
	return s.comprobantePagoRepo.Delete(id)
}

// List lista todos los comprobantes de pago
func (s *ComprobantePagoService) List() ([]*entidades.ComprobantePago, error) {
	return s.comprobantePagoRepo.List()
}

// ListByReserva lista todos los comprobantes de pago de una reserva específica
func (s *ComprobantePagoService) ListByReserva(idReserva int) ([]*entidades.ComprobantePago, error) {
	// Verificar que la reserva existe
	_, err := s.reservaRepo.GetByID(idReserva)
	if err != nil {
		return nil, errors.New("la reserva especificada no existe")
	}

	// Listar comprobantes por reserva
	return s.comprobantePagoRepo.ListByReserva(idReserva)
}

// ListByFecha lista todos los comprobantes de pago de una fecha específica
func (s *ComprobantePagoService) ListByFecha(fecha time.Time) ([]*entidades.ComprobantePago, error) {
	return s.comprobantePagoRepo.ListByFecha(fecha)
}

// ListByTipo lista todos los comprobantes de pago de un tipo específico
func (s *ComprobantePagoService) ListByTipo(tipo string) ([]*entidades.ComprobantePago, error) {
	// Verificar tipo válido
	if tipo != "BOLETA" && tipo != "FACTURA" {
		return nil, errors.New("tipo inválido, debe ser BOLETA o FACTURA")
	}

	// Listar comprobantes por tipo
	return s.comprobantePagoRepo.ListByTipo(tipo)
}

// ListByEstado lista todos los comprobantes de pago con un estado específico
func (s *ComprobantePagoService) ListByEstado(estado string) ([]*entidades.ComprobantePago, error) {
	// Verificar estado válido
	if estado != "EMITIDO" && estado != "ANULADO" {
		return nil, errors.New("estado inválido, debe ser EMITIDO o ANULADO")
	}

	// Listar comprobantes por estado
	return s.comprobantePagoRepo.ListByEstado(estado)
}

// ListByCliente lista todos los comprobantes relacionados con un cliente específico
func (s *ComprobantePagoService) ListByCliente(idCliente int) ([]*entidades.ComprobantePago, error) {
	// Verificar que el cliente existe
	// Esta verificación depende de cómo tienes organizados tus servicios

	// Listar comprobantes por cliente
	return s.comprobantePagoRepo.ListByCliente(idCliente)
}

// ListBySede lista todos los comprobantes de pago de una sede específica
func (s *ComprobantePagoService) ListBySede(idSede int) ([]*entidades.ComprobantePago, error) {
	// Verificar que la sede existe
	_, err := s.sedeRepo.GetByID(idSede)
	if err != nil {
		return nil, errors.New("la sede especificada no existe")
	}

	// Listar comprobantes por sede
	return s.comprobantePagoRepo.ListBySede(idSede)
}
