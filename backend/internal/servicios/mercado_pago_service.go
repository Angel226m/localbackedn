package servicios

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sistema-toursseft/internal/entidades"
	"time"
)

// MercadoPagoService maneja la integración con Mercado Pago
type MercadoPagoService struct {
	AccessToken string
	PublicKey   string
	ApiBaseURL  string
}

// NewMercadoPagoService crea una nueva instancia del servicio de Mercado Pago
func NewMercadoPagoService() *MercadoPagoService {
	return &MercadoPagoService{
		AccessToken: os.Getenv("MERCADOPAGO_ACCESS_TOKEN"),
		PublicKey:   os.Getenv("MERCADOPAGO_PUBLIC_KEY"),
		ApiBaseURL:  "https://api.mercadopago.com",
	}
}

// PreferenceItem representa un ítem en la preferencia de Mercado Pago
type PreferenceItem struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	PictureURL  string  `json:"picture_url,omitempty"`
	CategoryID  string  `json:"category_id,omitempty"`
	Quantity    int     `json:"quantity"`
	CurrencyID  string  `json:"currency_id"`
	UnitPrice   float64 `json:"unit_price"`
}

// Payer representa al pagador en Mercado Pago
type Payer struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Email   string `json:"email"`
	Phone   struct {
		AreaCode string `json:"area_code"`
		Number   string `json:"number"`
	} `json:"phone"`
	Identification struct {
		Type   string `json:"type"`
		Number string `json:"number"`
	} `json:"identification"`
	Address struct {
		ZipCode      string `json:"zip_code"`
		StreetName   string `json:"street_name"`
		StreetNumber int    `json:"street_number"`
	} `json:"address"`
}

// BackURLs representa las URLs de redirección tras el pago
type BackURLs struct {
	Success string `json:"success"`
	Failure string `json:"failure"`
	Pending string `json:"pending"`
}

// PaymentMethods representa las configuraciones de métodos de pago
type PaymentMethods struct {
	ExcludedPaymentMethods []struct {
		ID string `json:"id"`
	} `json:"excluded_payment_methods"`
	ExcludedPaymentTypes []struct {
		ID string `json:"id"`
	} `json:"excluded_payment_types"`
	Installments int `json:"installments"`
}

// PreferenceRequest representa la solicitud para crear una preferencia
type PreferenceRequest struct {
	Items               []PreferenceItem `json:"items"`
	Payer               Payer            `json:"payer"`
	BackURLs            BackURLs         `json:"back_urls"`
	AutoReturn          string           `json:"auto_return"`
	PaymentMethods      PaymentMethods   `json:"payment_methods"`
	NotificationURL     string           `json:"notification_url"`
	ExternalReference   string           `json:"external_reference"`
	StatementDescriptor string           `json:"statement_descriptor"`
}

// PreferenceResponse representa la respuesta de Mercado Pago al crear una preferencia
type PreferenceResponse struct {
	ID               string    `json:"id"`
	InitPoint        string    `json:"init_point"`
	SandboxInitPoint string    `json:"sandbox_init_point"`
	DateCreated      time.Time `json:"date_created"`
	LastUpdated      time.Time `json:"last_updated"`
}

// PaymentNotification representa la notificación de un pago de Mercado Pago
type PaymentNotification struct {
	ID            int64     `json:"id"`
	LiveMode      bool      `json:"live_mode"`
	Type          string    `json:"type"`
	DateCreated   time.Time `json:"date_created"`
	ApplicationID int64     `json:"application_id"`
	UserID        int64     `json:"user_id"`
	Version       int       `json:"version"`
	Data          struct {
		ID string `json:"id"`
	} `json:"data"`
}

// PaymentResponse representa los detalles de un pago de Mercado Pago
type PaymentResponse struct {
	ID                int64     `json:"id"`
	DateCreated       time.Time `json:"date_created"`
	DateApproved      time.Time `json:"date_approved"`
	DateLastUpdated   time.Time `json:"date_last_updated"`
	DateOfExpiration  time.Time `json:"date_of_expiration"`
	MoneyReleaseDate  time.Time `json:"money_release_date"`
	OperationType     string    `json:"operation_type"`
	IssuerId          string    `json:"issuer_id"`
	PaymentMethodId   string    `json:"payment_method_id"`
	PaymentTypeId     string    `json:"payment_type_id"`
	Status            string    `json:"status"`
	StatusDetail      string    `json:"status_detail"`
	CurrencyId        string    `json:"currency_id"`
	Description       string    `json:"description"`
	TransactionAmount float64   `json:"transaction_amount"`
	ExternalReference string    `json:"external_reference"`
}

// CreatePreference crea una preferencia de pago en Mercado Pago
func (s *MercadoPagoService) CreatePreference(
	tourNombre string,
	monto float64,
	idReserva int,
	cliente *entidades.Cliente,
	frontendURL string,
) (*PreferenceResponse, error) {
	// Construir la solicitud de preferencia
	preferenceURL := fmt.Sprintf("%s/checkout/preferences", s.ApiBaseURL)

	// Crear item para la preferencia
	items := []PreferenceItem{
		{
			ID:          fmt.Sprintf("TOUR-%d", idReserva),
			Title:       fmt.Sprintf("Reserva: %s", tourNombre),
			Description: "Reserva de tour en Tours Perú",
			Quantity:    1,
			CurrencyID:  "PEN", // Soles peruanos
			UnitPrice:   monto,
		},
	}

	// Configurar información del pagador
	payer := Payer{
		Name:    cliente.Nombres,
		Surname: cliente.Apellidos,
		Email:   cliente.Correo,
	}

	// Configurar número de teléfono si está disponible
	if cliente.NumeroCelular != "" {
		// Suponiendo que el número es algo como +51987654321
		payer.Phone.AreaCode = "51" // Código de país para Perú
		payer.Phone.Number = cliente.NumeroCelular
	}

	// Configurar documento de identidad si está disponible
	if cliente.NumeroDocumento != "" {
		payer.Identification.Type = "DNI" // Para Perú generalmente es DNI
		payer.Identification.Number = cliente.NumeroDocumento
	}

	// URLs de redirección después del pago
	backURLs := BackURLs{
		Success: fmt.Sprintf("%s/reserva-exitosa", frontendURL),
		Failure: fmt.Sprintf("%s/pago-fallido", frontendURL),
		Pending: fmt.Sprintf("%s/pago-pendiente", frontendURL),
	}

	// Configuración de métodos de pago
	paymentMethods := PaymentMethods{
		Installments: 1, // Solo pago en una cuota
	}

	// Crear la solicitud completa
	preferenceReq := PreferenceRequest{
		Items:               items,
		Payer:               payer,
		BackURLs:            backURLs,
		AutoReturn:          "approved",
		PaymentMethods:      paymentMethods,
		NotificationURL:     fmt.Sprintf("%s/api/webhook/mercadopago", frontendURL),
		ExternalReference:   fmt.Sprintf("RESERVA-%d", idReserva),
		StatementDescriptor: "TOURS PERU",
	}

	// Convertir la solicitud a JSON
	jsonData, err := json.Marshal(preferenceReq)
	if err != nil {
		return nil, err
	}

	// Crear la solicitud HTTP
	req, err := http.NewRequest("POST", preferenceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	// Configurar headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.AccessToken))

	// Realizar la solicitud
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Leer respuesta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Verificar código de respuesta
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error al crear preferencia: %s - código: %d", string(body), resp.StatusCode)
	}

	// Deserializar respuesta
	var preferenceResp PreferenceResponse
	err = json.Unmarshal(body, &preferenceResp)
	if err != nil {
		return nil, err
	}

	return &preferenceResp, nil
}

// GetPaymentInfo obtiene la información de un pago específico
func (s *MercadoPagoService) GetPaymentInfo(paymentId string) (*PaymentResponse, error) {
	// Construir URL para obtener detalles del pago
	paymentURL := fmt.Sprintf("%s/v1/payments/%s", s.ApiBaseURL, paymentId)

	// Crear la solicitud HTTP
	req, err := http.NewRequest("GET", paymentURL, nil)
	if err != nil {
		return nil, err
	}

	// Configurar headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.AccessToken))

	// Realizar la solicitud
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Leer respuesta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Verificar código de respuesta
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error al obtener información del pago: %s - código: %d", string(body), resp.StatusCode)
	}

	// Deserializar respuesta
	var paymentResp PaymentResponse
	err = json.Unmarshal(body, &paymentResp)
	if err != nil {
		return nil, err
	}

	return &paymentResp, nil
}

// MapMercadoPagoStatusToInternal mapea los estados de Mercado Pago a estados internos del sistema
func (s *MercadoPagoService) MapMercadoPagoStatusToInternal(mpStatus string) string {
	switch mpStatus {
	case "approved":
		return "PROCESADO"
	case "refunded", "cancelled", "rejected":
		return "ANULADO"
	case "pending", "in_process", "authorized":
		return "PENDIENTE"
	default:
		return "PENDIENTE"
	}
}

// ProcessPaymentWebhook procesa la notificación de webhook de Mercado Pago
func (s *MercadoPagoService) ProcessPaymentWebhook(notification *PaymentNotification) (*PaymentResponse, error) {
	if notification.Type != "payment" {
		return nil, errors.New("tipo de notificación no soportado")
	}

	// Obtener información detallada del pago
	paymentInfo, err := s.GetPaymentInfo(notification.Data.ID)
	if err != nil {
		return nil, err
	}

	return paymentInfo, nil
}

// GeneratePreferenceForExistingReserva genera una preferencia de pago para una reserva existente
func (s *MercadoPagoService) GeneratePreferenceForExistingReserva(
	idReserva int,
	monto float64,
	cliente *entidades.Cliente,
	frontendURL string,
) (*PreferenceResponse, error) {
	// Podemos reutilizar el método CreatePreference, pero necesitamos un nombre para el tour
	// En un caso real, obtendrías el nombre del tour desde la reserva
	tourNombre := "Reserva de Tour"

	return s.CreatePreference(tourNombre, monto, idReserva, cliente, frontendURL)
}
