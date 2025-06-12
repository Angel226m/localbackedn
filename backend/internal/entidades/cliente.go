// Cliente representa la estructura de un cliente en el sistema
package entidades

// Cliente representa a un cliente en el sistema, pudiendo ser persona natural o empresa
type Cliente struct {
	ID              int    `json:"id_cliente" db:"id_cliente"`
	TipoDocumento   string `json:"tipo_documento" db:"tipo_documento"` // DNI, CE, Pasaporte, RUC
	NumeroDocumento string `json:"numero_documento" db:"numero_documento"`
	// Campos para persona natural (obligatorios si tipo_documento es DNI, CE o Pasaporte)
	Nombres   string `json:"nombres,omitempty" db:"nombres"`
	Apellidos string `json:"apellidos,omitempty" db:"apellidos"`
	// Campos para empresas (obligatorios si tipo_documento es RUC)
	RazonSocial     string `json:"razon_social,omitempty" db:"razon_social"`
	DireccionFiscal string `json:"direccion_fiscal,omitempty" db:"direccion_fiscal"`
	// Campos de contacto
	Correo        string `json:"correo" db:"correo"`
	NumeroCelular string `json:"numero_celular" db:"numero_celular"`
	// Datos de sistema
	Contrasena     string `json:"-" db:"contrasena"`
	NombreCompleto string `json:"nombre_completo,omitempty" db:"-"` // Campo calculado
	Eliminado      bool   `json:"eliminado" db:"eliminado"`
}

// NuevoClienteRequest representa los datos necesarios para crear un nuevo cliente
type NuevoClienteRequest struct {
	TipoDocumento   string `json:"tipo_documento" validate:"required,oneof=DNI CE Pasaporte RUC"`
	NumeroDocumento string `json:"numero_documento" validate:"required"`
	// Campos para persona natural
	Nombres   string `json:"nombres" validate:"required_if=TipoDocumento DNI CE Pasaporte"`
	Apellidos string `json:"apellidos" validate:"required_if=TipoDocumento DNI CE Pasaporte"`
	// Campos para empresas
	RazonSocial     string `json:"razon_social" validate:"required_if=TipoDocumento RUC"`
	DireccionFiscal string `json:"direccion_fiscal" validate:"required_if=TipoDocumento RUC"`
	// Campos de contacto
	Correo        string `json:"correo" validate:"required,email"`
	NumeroCelular string `json:"numero_celular" validate:"required"`
	// Datos de sistema
	Contrasena string `json:"contrasena,omitempty" validate:"omitempty,min=6"`
}

// ActualizarClienteRequest representa los datos para actualizar un cliente
type ActualizarClienteRequest struct {
	TipoDocumento   string `json:"tipo_documento" validate:"required,oneof=DNI CE Pasaporte RUC"`
	NumeroDocumento string `json:"numero_documento" validate:"required"`
	// Campos para persona natural
	Nombres   string `json:"nombres" validate:"required_if=TipoDocumento DNI CE Pasaporte"`
	Apellidos string `json:"apellidos" validate:"required_if=TipoDocumento DNI CE Pasaporte"`
	// Campos para empresas
	RazonSocial     string `json:"razon_social" validate:"required_if=TipoDocumento RUC"`
	DireccionFiscal string `json:"direccion_fiscal" validate:"required_if=TipoDocumento RUC"`
	// Campos de contacto
	Correo        string `json:"correo" validate:"required,email"`
	NumeroCelular string `json:"numero_celular" validate:"required"`
}

// ActualizarDatosEmpresaRequest representa los datos para actualizar información de empresa
type ActualizarDatosEmpresaRequest struct {
	RazonSocial     string `json:"razon_social" validate:"required"`
	DireccionFiscal string `json:"direccion_fiscal" validate:"required"`
}
