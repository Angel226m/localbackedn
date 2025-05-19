package servicios

import (
	"errors"
	"sistema-tours/internal/entidades"
	"sistema-tours/internal/repositorios"
	"sistema-tours/internal/utils"
)

// ClienteService maneja la lógica de negocio para clientes
type ClienteService struct {
	clienteRepo *repositorios.ClienteRepository
}

// NewClienteService crea una nueva instancia de ClienteService
func NewClienteService(clienteRepo *repositorios.ClienteRepository) *ClienteService {
	return &ClienteService{
		clienteRepo: clienteRepo,
	}
}

// Create crea un nuevo cliente
func (s *ClienteService) Create(cliente *entidades.NuevoClienteRequest) (int, error) {
	// Verificar si ya existe un cliente con el mismo documento
	existing, err := s.clienteRepo.GetByDocumento(cliente.TipoDocumento, cliente.NumeroDocumento)
	if err == nil && existing != nil {
		return 0, errors.New("ya existe un cliente con este tipo y número de documento")
	}

	// Verificar si ya existe un cliente con el mismo correo (si se proporcionó)
	if cliente.Correo != "" {
		existingByEmail, err := s.clienteRepo.GetByCorreo(cliente.Correo)
		if err == nil && existingByEmail != nil {
			return 0, errors.New("ya existe un cliente con este correo electrónico")
		}
	}

	// Si se proporcionó una contraseña, hashearla
	if cliente.Contrasena != "" {
		hashedPassword, err := utils.HashPassword(cliente.Contrasena)
		if err != nil {
			return 0, errors.New("error al procesar la contraseña")
		}
		cliente.Contrasena = hashedPassword
	}

	// Crear cliente
	return s.clienteRepo.Create(cliente)
}

// GetByID obtiene un cliente por su ID
func (s *ClienteService) GetByID(id int) (*entidades.Cliente, error) {
	return s.clienteRepo.GetByID(id)
}

// GetByDocumento obtiene un cliente por tipo y número de documento
func (s *ClienteService) GetByDocumento(tipoDocumento, numeroDocumento string) (*entidades.Cliente, error) {
	return s.clienteRepo.GetByDocumento(tipoDocumento, numeroDocumento)
}

// GetByCorreo obtiene un cliente por su correo electrónico
func (s *ClienteService) GetByCorreo(correo string) (*entidades.Cliente, error) {
	return s.clienteRepo.GetByCorreo(correo)
}

// Update actualiza un cliente existente
func (s *ClienteService) Update(id int, cliente *entidades.ActualizarClienteRequest) error {
	// Verificar que el cliente existe
	existing, err := s.clienteRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Verificar si ya existe otro cliente con el mismo documento
	if cliente.TipoDocumento != existing.TipoDocumento || cliente.NumeroDocumento != existing.NumeroDocumento {
		existingDoc, err := s.clienteRepo.GetByDocumento(cliente.TipoDocumento, cliente.NumeroDocumento)
		if err == nil && existingDoc != nil && existingDoc.ID != id {
			return errors.New("ya existe otro cliente con este tipo y número de documento")
		}
	}

	// Verificar si ya existe otro cliente con el mismo correo (si se proporcionó)
	if cliente.Correo != "" && cliente.Correo != existing.Correo {
		existingByEmail, err := s.clienteRepo.GetByCorreo(cliente.Correo)
		if err == nil && existingByEmail != nil && existingByEmail.ID != id {
			return errors.New("ya existe otro cliente con este correo electrónico")
		}
	}

	// Actualizar cliente
	return s.clienteRepo.Update(id, cliente)
}

// UpdatePassword actualiza la contraseña de un cliente
func (s *ClienteService) UpdatePassword(id int, contrasenaActual, nuevaContrasena string) error {
	// Verificar que el cliente existe
	cliente, err := s.clienteRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Obtener contraseña actual del cliente
	passwordHash, err := s.clienteRepo.GetPasswordByCorreo(cliente.Correo)
	if err != nil {
		return err
	}

	// Verificar contraseña actual
	if !utils.CheckPasswordHash(contrasenaActual, passwordHash) {
		return errors.New("contraseña actual incorrecta")
	}

	// Hashear nueva contraseña
	hashedPassword, err := utils.HashPassword(nuevaContrasena)
	if err != nil {
		return errors.New("error al procesar la nueva contraseña")
	}

	// Actualizar contraseña
	return s.clienteRepo.UpdatePassword(id, hashedPassword)
}

// Delete elimina un cliente
func (s *ClienteService) Delete(id int) error {
	// Verificar que el cliente existe
	_, err := s.clienteRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Eliminar cliente
	return s.clienteRepo.Delete(id)
}

// List lista todos los clientes
func (s *ClienteService) List() ([]*entidades.Cliente, error) {
	return s.clienteRepo.List()
}

// SearchByName busca clientes por nombre o apellido
func (s *ClienteService) SearchByName(query string) ([]*entidades.Cliente, error) {
	if query == "" {
		return s.clienteRepo.List()
	}
	return s.clienteRepo.SearchByName(query)
}

// Login realiza el login de un cliente
func (s *ClienteService) Login(correo, contrasena string) (*entidades.Cliente, error) {
	// Verificar que existe un cliente con ese correo
	cliente, err := s.clienteRepo.GetByCorreo(correo)
	if err != nil {
		return nil, errors.New("correo electrónico o contraseña incorrectos")
	}

	// Obtener contraseña hash
	passwordHash, err := s.clienteRepo.GetPasswordByCorreo(correo)
	if err != nil {
		return nil, errors.New("error al verificar credenciales")
	}

	// Verificar contraseña
	if !utils.CheckPasswordHash(contrasena, passwordHash) {
		return nil, errors.New("correo electrónico o contraseña incorrectos")
	}

	return cliente, nil
}
