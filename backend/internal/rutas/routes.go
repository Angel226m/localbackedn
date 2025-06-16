package rutas

import (
	"fmt"
	"net/http"
	"sistema-toursseft/internal/config"
	"sistema-toursseft/internal/controladores"
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/middleware"
	"sistema-toursseft/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configura todas las rutas de la API
func SetupRoutes(
	router *gin.Engine,
	config *config.Config,
	authController *controladores.AuthController,
	usuarioController *controladores.UsuarioController,
	idiomaController *controladores.IdiomaController, // NUEVO CONTROLADOR
	usuarioIdiomaController *controladores.UsuarioIdiomaController, // Asegúrate de que esté aquí

	embarcacionController *controladores.EmbarcacionController,
	tipoTourController *controladores.TipoTourController,
	galeriaTourController *controladores.GaleriaTourController,

	horarioTourController *controladores.HorarioTourController,
	horarioChoferController *controladores.HorarioChoferController,
	tourProgramadoController *controladores.TourProgramadoController,
	tipoPasajeController *controladores.TipoPasajeController,
	paquetePasajesController *controladores.PaquetePasajesController, // Nuevo controlador

	metodoPagoController *controladores.MetodoPagoController,
	canalVentaController *controladores.CanalVentaController,
	clienteController *controladores.ClienteController,
	reservaController *controladores.ReservaController, // Agregar controlador de reservas
	pagoController *controladores.PagoController,
	comprobantePagoController *controladores.ComprobantePagoController,
	sedeController *controladores.SedeController,
	instanciaTourController *controladores.InstanciaTourController, // Nuevo controlador

) {
	// Middleware global
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.ErrorMiddleware())
	router.Use(gin.Recovery())

	// Rutas públicas
	public := router.Group("/api/v1")
	{
		// Autenticación
		public.POST("/auth/login", authController.Login)
		public.POST("/auth/refresh", authController.RefreshToken)
		public.POST("/auth/logout", authController.Logout)

		// Registro de cliente
		public.POST("/clientes/registro", clienteController.Create)

		// Autenticación de clientes
		public.POST("/clientes/login", clienteController.Login)
		public.POST("/clientes/refresh", clienteController.RefreshToken)
		public.POST("/clientes/logout", clienteController.Logout)

		// En la sección de rutas públicas:
		// Tipos de tour (acceso público)
		public.GET("/tipos-tour", tipoTourController.List)
		public.GET("/tipos-tour/:id", tipoTourController.GetByID)
		public.GET("/tipos-tour/sede/:idSede", tipoTourController.ListBySede)

		// Tours programados disponibles (acceso público)
		public.GET("/tours/disponibles", tourProgramadoController.ListToursProgramadosDisponibles)
		public.GET("/tours/disponibilidad/:fecha", tourProgramadoController.GetDisponibilidadDia)
		public.GET("/tours/:id", tourProgramadoController.GetByID)
		public.GET("/tours/disponibles-en-fecha/:fecha", tourProgramadoController.GetToursDisponiblesEnFecha)
		public.GET("/tours/disponibles-en-rango", tourProgramadoController.GetToursDisponiblesEnRangoFechas)
		public.GET("/tours/verificar-disponibilidad", tourProgramadoController.VerificarDisponibilidadHorario)
		// Añadir esta línea a la configuración de rutas
		public.GET("/tours/disponibles-sin-duplicados", tourProgramadoController.GetToursDisponibles)

		// En tu configuración de rutas
		public.GET("/instancias-tour/disponibles", func(ctx *gin.Context) {
			// Crear filtro para buscar solo instancias con estado PROGRAMADO
			var filtros entidades.FiltrosInstanciaTour
			estado := "PROGRAMADO"
			filtros.Estado = &estado

			// Establecer el filtro en el contexto
			ctx.Set("filtros", filtros)

			instanciaTourController.ListByFiltros(ctx)
		})

		// Consultar instancias de tour por fecha
		public.GET("/instancias-tour/fecha/:fecha", func(ctx *gin.Context) {
			fecha := ctx.Param("fecha")
			var filtros entidades.FiltrosInstanciaTour

			filtros.FechaInicio = &fecha
			filtros.FechaFin = &fecha

			// Establecer el filtro en el contexto
			ctx.Set("filtros", filtros)

			instanciaTourController.ListByFiltros(ctx)
		})

		// Tipos de pasaje (acceso público para ver precios)
		public.GET("/tipos-pasaje", tipoPasajeController.List)
		public.GET("/tipos-pasaje/sede/:idSede", tipoPasajeController.ListBySede)
		public.GET("/tipos-pasaje/tipo-tour/:id_tipo_tour", tipoPasajeController.ListByTipoTour)

		// Galería de tours
		public.GET("/tipo-tours/:id_tipo_tour/galerias", galeriaTourController.ListByTipoTour)
		public.GET("/galerias/:id", galeriaTourController.GetByID)

		// Paquetes de pasajes (acceso público para ver precios y opciones)
		public.GET("/paquetes-pasajes", paquetePasajesController.List)
		public.GET("/paquetes-pasajes/sede/:id_sede", paquetePasajesController.ListBySede)
		public.GET("/paquetes-pasajes/tipo-tour/:id_tipo_tour", paquetePasajesController.ListByTipoTour)
		public.GET("/paquetes-pasajes/:id", paquetePasajesController.GetByID)

		// Métodos de pago (acceso público para ver opciones)
		public.GET("/metodos-pago", metodoPagoController.List)

		// Sedes (acceso público)
		public.GET("/sedes", sedeController.List)
		public.GET("/sedes/:id", sedeController.GetByID)
		public.GET("/sedes/distrito/:distrito", sedeController.GetByDistrito)
		public.GET("/sedes/pais/:pais", sedeController.GetByPais)

		// Idiomas (acceso público para ver opciones disponibles)
		public.GET("/idiomas", idiomaController.List)

		//reservas mercado pago

		public.POST("/mercadopago/reservar", reservaController.ReservarConMercadoPago)

		// Webhook para recibir notificaciones de Mercado Pago
		public.POST("/webhook/mercadopago", reservaController.WebhookMercadoPago)

		// Verificar disponibilidad de instancia
		public.GET("/instancias-tour/:idInstancia/verificar-disponibilidad", reservaController.VerificarDisponibilidadInstancia)

	}

	// Rutas protegidas (requieren autenticación)
	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(config))
	{
		// Cambiar contraseña (cualquier usuario autenticado)
		protected.GET("/auth/status", authController.CheckStatus)
		protected.POST("/auth/change-password", authController.ChangePassword)

		// Estas rutas deben estar fuera del grupo admin para mantenerlas separadas
		adminAuth := protected.Group("/auth")
		adminAuth.Use(middleware.RoleMiddleware("ADMIN"))
		{
			protected.GET("/auth/debug", func(ctx *gin.Context) {
				// Listar todo el contexto
				fmt.Println("Datos en el contexto:")
				for k, v := range ctx.Keys {
					fmt.Printf("  %s: %v\n", k, v)
				}

				ctx.JSON(http.StatusOK, gin.H{
					"message":      "Datos del contexto",
					"userID":       ctx.GetInt("userID"),
					"userRole":     ctx.GetString("userRole"),
					"sedeID":       ctx.GetInt("sedeID"),
					"adminConSede": ctx.GetBool("adminConSede"),
					"adminSinSede": ctx.GetBool("adminSinSede"),
				})
			})

			// Obtener todas las sedes disponibles para seleccionar
			adminAuth.GET("/sedes", authController.GetUserSedes)

			// Seleccionar una sede específica para la sesión
			adminAuth.POST("/select-sede", authController.SelectSede)
		}

		// Usuarios - Admin
		admin := protected.Group("/admin")
		admin.Use(middleware.RoleMiddleware("ADMIN"))
		{
			// Gestión de usuarios
			admin.POST("/usuarios", usuarioController.Create)
			admin.GET("/usuarios", usuarioController.List)
			admin.GET("/usuarios/:id", usuarioController.GetByID)
			admin.PUT("/usuarios/:id", usuarioController.Update)
			admin.DELETE("/usuarios/:id", usuarioController.Delete)
			admin.GET("/usuarios/rol/:rol", usuarioController.ListByRol)

			admin.GET("/usuarios/:id/idiomas", usuarioIdiomaController.GetIdiomasByUsuarioID)
			admin.POST("/usuarios/:id/idiomas", usuarioIdiomaController.AsignarIdioma)
			admin.DELETE("/usuarios/:id/idiomas/:idioma_id", usuarioIdiomaController.DesasignarIdioma)
			admin.PUT("/usuarios/:id/idiomas", usuarioIdiomaController.ActualizarIdiomasUsuario)
			admin.GET("/idiomas/:id/usuarios", usuarioIdiomaController.GetUsuariosByIdiomaID)

			// Gestión de idiomas
			admin.POST("/idiomas", idiomaController.Create)
			admin.GET("/idiomas", idiomaController.List)
			admin.GET("/idiomas/:id", idiomaController.GetByID)
			admin.PUT("/idiomas/:id", idiomaController.Update)
			admin.DELETE("/idiomas/:id", idiomaController.Delete)

			// Gestión de embarcaciones
			admin.POST("/embarcaciones", embarcacionController.Create)
			admin.GET("/embarcaciones", embarcacionController.List)
			admin.GET("/embarcaciones/sede/:idSede", embarcacionController.ListBySede)
			admin.GET("/embarcaciones/:id", embarcacionController.GetByID)
			admin.PUT("/embarcaciones/:id", embarcacionController.Update)
			admin.DELETE("/embarcaciones/:id", embarcacionController.Delete)

			// Gestión de tipos de tour
			admin.POST("/tipos-tour", tipoTourController.Create)
			admin.GET("/tipos-tour", tipoTourController.List)
			admin.GET("/tipos-tour/sede/:idSede", tipoTourController.ListBySede)
			admin.GET("/tipos-tour/:id", tipoTourController.GetByID)
			admin.PUT("/tipos-tour/:id", tipoTourController.Update)
			admin.DELETE("/tipos-tour/:id", tipoTourController.Delete)

			// Gestión de galería de imágenes
			admin.POST("/galerias", galeriaTourController.Create)
			admin.GET("/galerias/:id", galeriaTourController.GetByID)
			admin.PUT("/galerias/:id", galeriaTourController.Update)
			admin.DELETE("/galerias/:id", galeriaTourController.Delete)
			admin.GET("/tipo-tours/:id_tipo_tour/galerias", galeriaTourController.ListByTipoTour)

			// Gestión de horarios de tour
			admin.POST("/horarios-tour", horarioTourController.Create)
			admin.GET("/horarios-tour", horarioTourController.List)
			admin.GET("/horarios-tour/:id", horarioTourController.GetByID)
			admin.PUT("/horarios-tour/:id", horarioTourController.Update)
			admin.DELETE("/horarios-tour/:id", horarioTourController.Delete)
			admin.GET("/horarios-tour/tipo/:idTipoTour", horarioTourController.ListByTipoTour)
			admin.GET("/horarios-tour/dia/:dia", horarioTourController.ListByDia)

			// Gestión de horarios de chofer
			admin.POST("/horarios-chofer", horarioChoferController.Create)
			admin.GET("/horarios-chofer", horarioChoferController.List)
			admin.GET("/horarios-chofer/:id", horarioChoferController.GetByID)
			admin.PUT("/horarios-chofer/:id", horarioChoferController.Update)
			admin.DELETE("/horarios-chofer/:id", horarioChoferController.Delete)
			admin.GET("/horarios-chofer/chofer/:idChofer", horarioChoferController.ListByChofer)
			admin.GET("/horarios-chofer/chofer/:idChofer/activos", horarioChoferController.ListActiveByChofer)
			admin.GET("/horarios-chofer/dia/:dia", horarioChoferController.ListByDia)

			// Gestión de tours programados
			admin.POST("/tours", tourProgramadoController.Create)
			admin.GET("/tours", tourProgramadoController.List)
			admin.GET("/tours/:id", tourProgramadoController.GetByID)
			admin.PUT("/tours/:id", tourProgramadoController.Update)
			admin.DELETE("/tours/:id", tourProgramadoController.Delete)
			admin.POST("/tours/:id/chofer", tourProgramadoController.AsignarChofer)
			admin.POST("/tours/:id/estado", tourProgramadoController.CambiarEstado)
			admin.GET("/tours/programacion-semanal", tourProgramadoController.GetProgramacionSemanal)
			admin.GET("/tours/fecha/:fecha", tourProgramadoController.ListByFecha)
			admin.GET("/tours/rango-fechas", tourProgramadoController.ListByRangoFechas)
			admin.GET("/tours/estado/:estado", tourProgramadoController.ListByEstado)
			admin.GET("/tours/embarcacion/:idEmbarcacion", tourProgramadoController.ListByEmbarcacion)
			admin.GET("/tours/chofer/:idChofer", tourProgramadoController.ListByChofer)
			admin.GET("/tours/tipo-tour/:idTipoTour", tourProgramadoController.ListByTipoTour)
			admin.GET("/tours/sede/:idSede", tourProgramadoController.ListBySede)
			admin.GET("/tours/vigentes", tourProgramadoController.GetToursVigentes)
			admin.GET("/tours/disponibles-en-fecha/:fecha", tourProgramadoController.GetToursDisponiblesEnFecha)
			admin.GET("/tours/disponibles-en-rango", tourProgramadoController.GetToursDisponiblesEnRangoFechas)
			admin.GET("/tours/verificar-disponibilidad", tourProgramadoController.VerificarDisponibilidadHorario)

			// Gestión de tipos de pasaje
			admin.POST("/tipos-pasaje", tipoPasajeController.Create)
			admin.GET("/tipos-pasaje", tipoPasajeController.List)
			admin.GET("/tipos-pasaje/:id", tipoPasajeController.GetByID)
			admin.PUT("/tipos-pasaje/:id", tipoPasajeController.Update)
			admin.DELETE("/tipos-pasaje/:id", tipoPasajeController.Delete)
			admin.GET("/tipos-pasaje/sede/:idSede", tipoPasajeController.ListBySede)
			admin.GET("/tipos-pasaje/tipo-tour/:id_tipo_tour", tipoPasajeController.ListByTipoTour)

			// Gestión de paquetes de pasajes
			admin.POST("/paquetes-pasajes", paquetePasajesController.Create)
			admin.GET("/paquetes-pasajes", paquetePasajesController.List)
			admin.GET("/paquetes-pasajes/:id", paquetePasajesController.GetByID)
			admin.PUT("/paquetes-pasajes/:id", paquetePasajesController.Update)
			admin.DELETE("/paquetes-pasajes/:id", paquetePasajesController.Delete)
			admin.GET("/paquetes-pasajes/sede/:id_sede", paquetePasajesController.ListBySede)
			admin.GET("/paquetes-pasajes/tipo-tour/:id_tipo_tour", paquetePasajesController.ListByTipoTour)

			// Gestión de métodos de pago
			admin.POST("/metodos-pago", metodoPagoController.Create)
			admin.GET("/metodos-pago", metodoPagoController.List)
			admin.GET("/metodos-pago/:id", metodoPagoController.GetByID)
			admin.PUT("/metodos-pago/:id", metodoPagoController.Update)
			admin.DELETE("/metodos-pago/:id", metodoPagoController.Delete)
			admin.GET("/metodos-pago/sede/:idSede", metodoPagoController.ListBySede)

			// Gestión de canales de venta
			admin.POST("/canales-venta", canalVentaController.Create)
			admin.GET("/canales-venta", canalVentaController.List)
			admin.GET("/canales-venta/:id", canalVentaController.GetByID)
			admin.PUT("/canales-venta/:id", canalVentaController.Update)
			admin.DELETE("/canales-venta/:id", canalVentaController.Delete)
			admin.GET("/canales-venta/sede/:idSede", canalVentaController.ListBySede)

			// Gestión de clientes
			admin.GET("/clientes", clienteController.List)
			admin.GET("/clientes/:id", clienteController.GetByID)
			admin.PUT("/clientes/:id", clienteController.Update)
			admin.DELETE("/clientes/:id", clienteController.Delete)
			// Puedes usar:
			admin.GET("/clientes/buscar-documento", func(ctx *gin.Context) {
				query := ctx.Query("query")
				if query == "" {
					ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Término de búsqueda requerido", nil))
					return
				}
				// Usar el método List existente con parámetro search y type=doc
				ctx.Request.URL.RawQuery = "search=" + query + "&type=doc"
				clienteController.List(ctx)
			})
			// Gestión de pagos
			admin.POST("/pagos", pagoController.Create)
			admin.GET("/pagos", pagoController.List)
			admin.GET("/pagos/:id", pagoController.GetByID)
			admin.PUT("/pagos/:id", pagoController.Update)
			admin.DELETE("/pagos/:id", pagoController.Delete)
			admin.POST("/pagos/:id/estado", pagoController.CambiarEstado)
			admin.GET("/pagos/reserva/:idReserva", pagoController.ListByReserva)
			admin.GET("/pagos/fecha/:fecha", pagoController.ListByFecha)
			admin.GET("/pagos/estado/:estado", pagoController.ListByEstado)
			admin.GET("/pagos/reserva/:idReserva/total", pagoController.GetTotalPagadoByReserva)
			admin.GET("/pagos/cliente/:idCliente", pagoController.ListByCliente)
			admin.GET("/pagos/sede/:idSede", pagoController.ListBySede)

			// Gestión de comprobantes de pago
			admin.POST("/comprobantes", comprobantePagoController.Create)
			admin.GET("/comprobantes", comprobantePagoController.List)
			admin.GET("/comprobantes/:id", comprobantePagoController.GetByID)
			admin.GET("/comprobantes/buscar", comprobantePagoController.GetByTipoAndNumero)
			admin.PUT("/comprobantes/:id", comprobantePagoController.Update)
			admin.DELETE("/comprobantes/:id", comprobantePagoController.Delete)
			admin.POST("/comprobantes/:id/estado", comprobantePagoController.CambiarEstado)
			admin.GET("/comprobantes/reserva/:idReserva", comprobantePagoController.ListByReserva)
			admin.GET("/comprobantes/fecha/:fecha", comprobantePagoController.ListByFecha)
			admin.GET("/comprobantes/tipo/:tipo", comprobantePagoController.ListByTipo)
			admin.GET("/comprobantes/estado/:estado", comprobantePagoController.ListByEstado)
			admin.GET("/comprobantes/cliente/:idCliente", comprobantePagoController.ListByCliente)

			// Gestión de sedes
			admin.POST("/sedes", sedeController.Create)
			admin.PUT("/sedes/:id", sedeController.Update)
			admin.DELETE("/sedes/:id", sedeController.Delete)
			admin.POST("/sedes/:id/restore", sedeController.Restore)
			admin.GET("/sedes", sedeController.List)
			admin.GET("/sedes/:id", sedeController.GetByID)
			admin.GET("/sedes/distrito/:distrito", sedeController.GetByDistrito)
			admin.GET("/sedes/pais/:pais", sedeController.GetByPais)

			admin.POST("/instancias-tour", instanciaTourController.Create)
			admin.GET("/instancias-tour", instanciaTourController.List)
			admin.GET("/instancias-tour/:id", instanciaTourController.GetByID)
			admin.PUT("/instancias-tour/:id", instanciaTourController.Update)
			admin.DELETE("/instancias-tour/:id", instanciaTourController.Delete)
			admin.POST("/instancias-tour/:id/asignar-chofer", instanciaTourController.AsignarChofer)
			admin.GET("/instancias-tour/tour-programado/:id_tour_programado", instanciaTourController.ListByTourProgramado)
			admin.POST("/instancias-tour/filtrar", instanciaTourController.ListByFiltros)
			admin.POST("/instancias-tour/generar/:id_tour_programado", instanciaTourController.GenerarInstanciasDeTourProgramado)

			admin.POST("/reservas", reservaController.Create)
			admin.GET("/reservas", reservaController.List)
			admin.GET("/reservas/:id", reservaController.GetByID)
			admin.PUT("/reservas/:id", reservaController.Update)
			admin.DELETE("/reservas/:id", reservaController.Delete)
			admin.POST("/reservas/:id/estado", reservaController.CambiarEstado)
			admin.GET("/reservas/cliente/:idCliente", reservaController.ListByCliente)
			admin.GET("/reservas/instancia/:idInstancia", reservaController.ListByInstancia)
			admin.GET("/reservas/fecha/:fecha", reservaController.ListByFecha)
			admin.GET("/reservas/estado/:estado", reservaController.ListByEstado)
			admin.GET("/reservas/sede/:idSede", reservaController.ListBySede)

			// Confirmación manual de pagos con Mercado Pago
			admin.POST("/reservas/confirmar-pago", reservaController.ConfirmarPagoReserva)

		}

		// Vendedores
		vendedor := protected.Group("/vendedor")
		vendedor.Use(middleware.RoleMiddleware("ADMIN", "VENDEDOR"))
		{
			// Ver idiomas (solo lectura)
			vendedor.GET("/idiomas", idiomaController.List)

			// Ver embarcaciones (solo lectura)
			vendedor.GET("/embarcaciones", embarcacionController.List)
			vendedor.GET("/embarcaciones/:id", embarcacionController.GetByID)

			// Ver tipos de tour (solo lectura)
			vendedor.GET("/tipos-tour", tipoTourController.List)
			vendedor.GET("/tipos-tour/:id", tipoTourController.GetByID)

			// Ver galería de imágenes (solo lectura)
			vendedor.GET("/tipo-tours/:id_tipo_tour/galerias", galeriaTourController.ListByTipoTour)
			vendedor.GET("/galerias/:id", galeriaTourController.GetByID)

			// Ver horarios de tour (solo lectura)
			vendedor.GET("/horarios-tour", horarioTourController.List)
			vendedor.GET("/horarios-tour/:id", horarioTourController.GetByID)
			vendedor.GET("/horarios-tour/tipo/:idTipoTour", horarioTourController.ListByTipoTour)
			vendedor.GET("/horarios-tour/dia/:dia", horarioTourController.ListByDia)

			// Ver horarios de choferes disponibles (solo lectura)
			vendedor.GET("/horarios-chofer/dia/:dia", horarioChoferController.ListByDia)

			// Ver tours programados (solo lectura)
			vendedor.GET("/tours", tourProgramadoController.List)
			vendedor.GET("/tours/:id", tourProgramadoController.GetByID)
			vendedor.GET("/tours/fecha/:fecha", tourProgramadoController.ListByFecha)
			vendedor.GET("/tours/rango-fechas", tourProgramadoController.ListByRangoFechas)
			vendedor.GET("/tours/estado/:estado", tourProgramadoController.ListByEstado)
			vendedor.GET("/tours/sede/:idSede", tourProgramadoController.ListBySede)
			vendedor.GET("/tours/disponibles", tourProgramadoController.ListToursProgramadosDisponibles)
			vendedor.GET("/tours/disponibles-en-fecha/:fecha", tourProgramadoController.GetToursDisponiblesEnFecha)
			vendedor.GET("/tours/disponibles-en-rango", tourProgramadoController.GetToursDisponiblesEnRangoFechas)

			// Ver tipos de pasaje (solo lectura)
			vendedor.GET("/tipos-pasaje", tipoPasajeController.List)
			vendedor.GET("/tipos-pasaje/:id", tipoPasajeController.GetByID)
			vendedor.GET("/tipos-pasaje/sede/:idSede", tipoPasajeController.ListBySede)
			vendedor.GET("/tipos-pasaje/tipo-tour/:id_tipo_tour", tipoPasajeController.ListByTipoTour)

			// Ver paquetes de pasajes (solo lectura)
			vendedor.GET("/paquetes-pasajes", paquetePasajesController.List)
			vendedor.GET("/paquetes-pasajes/:id", paquetePasajesController.GetByID)
			vendedor.GET("/paquetes-pasajes/sede/:id_sede", paquetePasajesController.ListBySede)
			vendedor.GET("/paquetes-pasajes/tipo-tour/:id_tipo_tour", paquetePasajesController.ListByTipoTour)

			// Ver métodos de pago (solo lectura)
			vendedor.GET("/metodos-pago", metodoPagoController.List)
			vendedor.GET("/metodos-pago/:id", metodoPagoController.GetByID)
			vendedor.GET("/metodos-pago/sede/:idSede", metodoPagoController.ListBySede)

			// Ver canales de venta (solo lectura)
			vendedor.GET("/canales-venta", canalVentaController.List)
			vendedor.GET("/canales-venta/:id", canalVentaController.GetByID)
			vendedor.GET("/canales-venta/sede/:idSede", canalVentaController.ListBySede)

			// Gestión de clientes
			vendedor.POST("/clientes", clienteController.Create)
			vendedor.GET("/clientes", clienteController.List)
			vendedor.GET("/clientes/:id", clienteController.GetByID)
			vendedor.PUT("/clientes/:id", clienteController.Update)
			vendedor.PUT("/clientes/:id/datos-empresa", clienteController.UpdateDatosEmpresa)
			// Búsqueda rápida de clientes por documento o RUC
			vendedor.GET("/clientes/documento", clienteController.GetByDocumento)
			// Por esta implementación:
			vendedor.GET("/clientes/buscar-documento", func(ctx *gin.Context) {
				query := ctx.Query("query")
				if query == "" {
					ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Término de búsqueda requerido", nil))
					return
				}
				// Redirigir a la función List con parámetros adecuados
				ctx.Request.URL.RawQuery = "search=" + query + "&type=doc"
				clienteController.List(ctx)
			})
			// Ver instancias de tour (solo lectura)
			vendedor.GET("/instancias-tour", instanciaTourController.List)
			vendedor.GET("/instancias-tour/:id", instanciaTourController.GetByID)
			vendedor.GET("/instancias-tour/tour-programado/:id_tour_programado", instanciaTourController.ListByTourProgramado)
			vendedor.POST("/instancias-tour/filtrar", instanciaTourController.ListByFiltros)

			// Gestión de pagos (vendedor puede registrar y ver pagos)
			vendedor.POST("/pagos", pagoController.Create)
			vendedor.GET("/pagos", pagoController.List)
			vendedor.GET("/pagos/:id", pagoController.GetByID)
			vendedor.GET("/pagos/reserva/:idReserva", pagoController.ListByReserva)
			vendedor.GET("/pagos/reserva/:idReserva/total", pagoController.GetTotalPagadoByReserva)
			vendedor.GET("/pagos/sede/:idSede", pagoController.ListBySede)

			// Gestión de comprobantes (vendedor puede emitir y ver comprobantes)
			vendedor.POST("/comprobantes", comprobantePagoController.Create)
			vendedor.GET("/comprobantes", comprobantePagoController.List)
			vendedor.GET("/comprobantes/:id", comprobantePagoController.GetByID)
			vendedor.GET("/comprobantes/buscar", comprobantePagoController.GetByTipoAndNumero)
			vendedor.GET("/comprobantes/reserva/:idReserva", comprobantePagoController.ListByReserva)
			//reservas mercado pago
			vendedor.POST("/reservas", reservaController.Create)
			vendedor.GET("/reservas", reservaController.List)
			vendedor.GET("/reservas/:id", reservaController.GetByID)
			vendedor.PUT("/reservas/:id", reservaController.Update)
			vendedor.POST("/reservas/:id/estado", reservaController.CambiarEstado)
			vendedor.GET("/reservas/cliente/:idCliente", reservaController.ListByCliente)
			vendedor.GET("/reservas/instancia/:idInstancia", reservaController.ListByInstancia)
			vendedor.GET("/reservas/fecha/:fecha", reservaController.ListByFecha)
			vendedor.GET("/reservas/estado/:estado", reservaController.ListByEstado)

			// Confirmación manual de pagos con Mercado Pago
			vendedor.POST("/reservas/confirmar-pago", reservaController.ConfirmarPagoReserva)
		}

		// Choferes
		chofer := protected.Group("/chofer")
		chofer.Use(middleware.RoleMiddleware("ADMIN", "CHOFER"))
		{
			// Ver embarcaciones asignadas
			chofer.GET("/mis-embarcaciones", func(ctx *gin.Context) {
				userID := ctx.GetInt("userID")
				ctx.Request.URL.Path = "/api/v1/admin/embarcaciones/chofer/" + strconv.Itoa(userID)
				router.HandleContext(ctx)
			})

			// Ver tipos de tour (solo lectura)
			chofer.GET("/tipos-tour", tipoTourController.List)

			// Ver horarios de tour (solo lectura)
			chofer.GET("/horarios-tour", horarioTourController.List)
			chofer.GET("/horarios-tour/dia/:dia", horarioTourController.ListByDia)

			// Ver mis horarios de trabajo
			chofer.GET("/mis-horarios", horarioChoferController.GetMyActiveHorarios)
			chofer.GET("/todos-mis-horarios", func(ctx *gin.Context) {
				userID := ctx.GetInt("userID")
				ctx.Request.URL.Path = "/api/v1/admin/horarios-chofer/chofer/" + strconv.Itoa(userID)
				router.HandleContext(ctx)
			})

			// Ver mis tours programados
			chofer.GET("/mis-tours", func(ctx *gin.Context) {
				userID := ctx.GetInt("userID")
				ctx.Request.URL.Path = "/api/v1/admin/tours/chofer/" + strconv.Itoa(userID)
				router.HandleContext(ctx)
			})

			chofer.GET("/tours/:id", tourProgramadoController.GetByID)
			chofer.GET("/tours/programacion-semanal", tourProgramadoController.GetProgramacionSemanal)
			chofer.GET("/tours/fecha/:fecha", tourProgramadoController.ListByFecha)

			chofer.GET("/mis-instancias-tour", func(ctx *gin.Context) {
				userID := ctx.GetInt("userID")
				// Crear filtro para buscar instancias asignadas a este chofer
				var filtros entidades.FiltrosInstanciaTour
				idChofer := userID
				filtros.IDChofer = &idChofer

				// Establecer el filtro en el contexto
				ctx.Set("filtros", filtros)

				instanciaTourController.ListByFiltros(ctx)
			})

		}

		// Clientes
		cliente := protected.Group("/cliente")
		cliente.Use(middleware.RoleMiddleware("ADMIN", "CLIENTE"))
		{
			// Cambiar contraseña (cliente)
			cliente.POST("/change-password", clienteController.ChangePassword)

			// Ver tipos de tour disponibles (solo lectura)
			cliente.GET("/tipos-tour", tipoTourController.List)
			cliente.GET("/tipos-tour/:id", tipoTourController.GetByID)

			// Ver horarios de tour disponibles (solo lectura)
			cliente.GET("/horarios-tour/tipo/:idTipoTour", horarioTourController.ListByTipoTour)

			// Ver tours disponibles
			cliente.GET("/tours/disponibles", tourProgramadoController.ListToursProgramadosDisponibles)
			cliente.GET("/tours/disponibilidad/:fecha", tourProgramadoController.GetDisponibilidadDia)
			cliente.GET("/tours/:id", tourProgramadoController.GetByID)
			cliente.GET("/tours/disponibles-en-fecha/:fecha", tourProgramadoController.GetToursDisponiblesEnFecha)
			cliente.GET("/tours/disponibles-en-rango", tourProgramadoController.GetToursDisponiblesEnRangoFechas)
			cliente.GET("/tours/verificar-disponibilidad", tourProgramadoController.VerificarDisponibilidadHorario)

			// Ver tipos de pasaje (solo lectura)
			cliente.GET("/tipos-pasaje", tipoPasajeController.List)
			cliente.GET("/tipos-pasaje/sede/:idSede", tipoPasajeController.ListBySede)
			cliente.GET("/tipos-pasaje/tipo-tour/:id_tipo_tour", tipoPasajeController.ListByTipoTour)

			// Ver paquetes de pasajes (solo lectura)
			cliente.GET("/paquetes-pasajes", paquetePasajesController.List)
			cliente.GET("/paquetes-pasajes/:id", paquetePasajesController.GetByID)
			cliente.GET("/paquetes-pasajes/sede/:id_sede", paquetePasajesController.ListBySede)
			cliente.GET("/paquetes-pasajes/tipo-tour/:id_tipo_tour", paquetePasajesController.ListByTipoTour)

			// Ver métodos de pago (solo lectura)
			cliente.GET("/metodos-pago", metodoPagoController.List)

			// Ver canales de venta (solo lectura)
			cliente.GET("/canales-venta", canalVentaController.List)

			// Gestión del perfil propio
			cliente.GET("/mi-perfil", func(ctx *gin.Context) {
				clienteID := ctx.GetInt("userID")
				ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: strconv.Itoa(clienteID)})
				clienteController.GetByID(ctx)
			})

			cliente.PUT("/mi-perfil", func(ctx *gin.Context) {
				clienteID := ctx.GetInt("userID")
				ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: strconv.Itoa(clienteID)})
				clienteController.Update(ctx)
			})

			// Actualizar datos de empresa propios (cambiar de "datos-facturacion" a "datos-empresa")
			cliente.PUT("/mi-perfil/datos-empresa", func(ctx *gin.Context) {
				clienteID := ctx.GetInt("userID")
				ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: strconv.Itoa(clienteID)})
				clienteController.UpdateDatosEmpresa(ctx)
			})
			// Ver galería de imágenes (solo lectura)
			cliente.GET("/tipo-tours/:id_tipo_tour/galerias", galeriaTourController.ListByTipoTour)
			cliente.GET("/galerias/:id", galeriaTourController.GetByID)

			// Ver mis pagos
			cliente.GET("/mis-pagos", func(ctx *gin.Context) {
				clienteID := ctx.GetInt("userID")
				ctx.Request.URL.Path = "/api/v1/admin/pagos/cliente/" + strconv.Itoa(clienteID)
				router.HandleContext(ctx)
			})

			// Ver mis comprobantes
			cliente.GET("/mis-comprobantes", func(ctx *gin.Context) {
				clienteID := ctx.GetInt("userID")
				ctx.Request.URL.Path = "/api/v1/admin/comprobantes/cliente/" + strconv.Itoa(clienteID)
				router.HandleContext(ctx)
			})

			// Ver mis reservas
			cliente.GET("/mis-reservas", reservaController.ListMyReservas)

			// Ver detalle de una reserva específica
			cliente.GET("/mis-reservas/:id", func(ctx *gin.Context) {
				reservaID := ctx.Param("id")
				clienteID := ctx.GetInt("userID")

				// Obtener la reserva
				id, _ := strconv.Atoi(reservaID)
				reserva, err := reservaController.reservaService.GetByID(id)
				if err != nil {
					ctx.JSON(http.StatusNotFound, utils.ErrorResponse("Reserva no encontrada", err))
					return
				}

				// Verificar que la reserva pertenece al cliente
				if reserva.IDCliente != clienteID {
					ctx.JSON(http.StatusForbidden, utils.ErrorResponse("No tiene acceso a esta reserva", nil))
					return
				}

				// Mostrar la reserva
				ctx.JSON(http.StatusOK, utils.SuccessResponse("Reserva obtenida exitosamente", reserva))
			})

			// Realizar reserva (desde área de cliente)
			cliente.POST("/reservas", func(ctx *gin.Context) {
				// Asegurar que la reserva se crea para el cliente autenticado
				clienteID := ctx.GetInt("userID")

				var reservaReq entidades.NuevaReservaRequest
				if err := ctx.ShouldBindJSON(&reservaReq); err != nil {
					ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Datos inválidos", err))
					return
				}

				// Forzar el ID del cliente al usuario autenticado
				reservaReq.IDCliente = clienteID

				// Llamar al controlador para crear la reserva
				ctx.Set("reservaRequest", reservaReq)
				reservaController.Create(ctx)
			})

			// Cancelar una reserva
			cliente.POST("/mis-reservas/:id/cancelar", func(ctx *gin.Context) {
				reservaID := ctx.Param("id")
				clienteID := ctx.GetInt("userID")

				// Obtener la reserva
				id, _ := strconv.Atoi(reservaID)
				reserva, err := reservaController.reservaService.GetByID(id)
				if err != nil {
					ctx.JSON(http.StatusNotFound, utils.ErrorResponse("Reserva no encontrada", err))
					return
				}

				// Verificar que la reserva pertenece al cliente
				if reserva.IDCliente != clienteID {
					ctx.JSON(http.StatusForbidden, utils.ErrorResponse("No tiene acceso a esta reserva", nil))
					return
				}

				// Crear el request para cambiar estado
				estadoReq := entidades.CambiarEstadoReservaRequest{
					Estado: "CANCELADA",
				}

				// Añadir al contexto y llamar al controlador
				ctx.Request.Method = "POST"
				ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: reservaID})
				ctx.Set("json", estadoReq)
				if err := ctx.ShouldBindJSON(&estadoReq); err != nil {
					// Si hay error de binding, cargarlo manualmente en el contexto
					ctx.Set("estadoRequest", estadoReq)
				}

				reservaController.CambiarEstado(ctx)
			})

			// Pagar una reserva con Mercado Pago
			cliente.POST("/mis-reservas/:id/pagar", func(ctx *gin.Context) {
				reservaID := ctx.Param("id")
				clienteID := ctx.GetInt("userID")

				// Obtener la reserva
				id, _ := strconv.Atoi(reservaID)
				reserva, err := reservaController.reservaService.GetByID(id)
				if err != nil {
					ctx.JSON(http.StatusNotFound, utils.ErrorResponse("Reserva no encontrada", err))
					return
				}

				// Verificar que la reserva pertenece al cliente
				if reserva.IDCliente != clienteID {
					ctx.JSON(http.StatusForbidden, utils.ErrorResponse("No tiene acceso a esta reserva", nil))
					return
				}

				// Obtener datos del cliente
				cliente, err := clienteController.clienteService.GetByID(clienteID)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al obtener datos del cliente", err))
					return
				}

				// Crear solicitud para generar preferencia de Mercado Pago
				frontendURL := ctx.GetHeader("Origin")
				if frontendURL == "" {
					frontendURL = "https://tours-peru.com" // URL predeterminada
				}

				// Crear la preferencia de pago para esta reserva
				// Aquí tendríamos que adaptar ya que la reserva ya existe, a diferencia de ReservarConMercadoPago
				// Esta es una implementación simplificada
				response, err := reservaController.mercadoPagoService.GeneratePreferenceForExistingReserva(
					id, reserva.TotalPagar, cliente, frontendURL)

				if err != nil {
					ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error al generar preferencia de pago", err))
					return
				}

				ctx.JSON(http.StatusOK, utils.SuccessResponse("Preferencia de pago generada exitosamente", response))
			})
		}
	}
}

/*package rutas

import (
	"fmt"
	"net/http"
	"sistema-toursseft/internal/config"
	"sistema-toursseft/internal/controladores"
	"sistema-toursseft/internal/middleware"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configura todas las rutas de la API
func SetupRoutes(

	router *gin.Engine,
	config *config.Config,
	authController *controladores.AuthController,
	usuarioController *controladores.UsuarioController,
	idiomaController *controladores.IdiomaController, // NUEVO CONTROLADOR
	usuarioIdiomaController *controladores.UsuarioIdiomaController, // Asegúrate de que esté aquí

	embarcacionController *controladores.EmbarcacionController,
	tipoTourController *controladores.TipoTourController,
	galeriaTourController *controladores.GaleriaTourController,

	horarioTourController *controladores.HorarioTourController,
	horarioChoferController *controladores.HorarioChoferController,
	tourProgramadoController *controladores.TourProgramadoController,
	tipoPasajeController *controladores.TipoPasajeController,
	paquetePasajesController *controladores.PaquetePasajesController, // Nuevo controlador

	metodoPagoController *controladores.MetodoPagoController,
	canalVentaController *controladores.CanalVentaController,
	clienteController *controladores.ClienteController,
	/*reservaController *controladores.ReservaController,*/
/*pagoController *controladores.PagoController,
	comprobantePagoController *controladores.ComprobantePagoController,
	sedeController *controladores.SedeController,

) {
	// Middleware global
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.ErrorMiddleware())
	router.Use(gin.Recovery())

	// Rutas públicas
	public := router.Group("/api/v1")
	{
		// Autenticación
		public.POST("/auth/login", authController.Login)
		public.POST("/auth/refresh", authController.RefreshToken)
		public.POST("/auth/logout", authController.Logout)

		// Registro de cliente
		public.POST("/clientes/registro", clienteController.Create)

		// Autenticación de clientes
		public.POST("/clientes/login", clienteController.Login)
		public.POST("/clientes/refresh", clienteController.RefreshToken)
		public.POST("/clientes/logout", clienteController.Logout)

		// Tours programados disponibles (acceso público)
		// Rutas públicas para TourProgramado que deben actualizarse
		public.GET("/tours/disponibles", tourProgramadoController.ListToursProgramadosDisponibles)
		public.GET("/tours/disponibilidad/:fecha", tourProgramadoController.GetDisponibilidadDia)
		public.GET("/tours/:id", tourProgramadoController.GetByID)

		public.GET("/tipos-pasaje/tipo-tour/:id_tipo_tour", tipoPasajeController.ListByTipoTour)

		// Tipos de pasaje (acceso público para ver precios)
		public.GET("/tipos-pasaje", tipoPasajeController.List)
		public.GET("/tipos-pasaje/sede/:idSede", tipoPasajeController.ListBySede)

		public.GET("/tipo-tours/:id_tipo_tour/galerias", galeriaTourController.ListByTipoTour)
		public.GET("/galerias/:id", galeriaTourController.GetByID)

		// Paquetes de pasajes (acceso público para ver precios y opciones)
		public.GET("/paquetes-pasajes", paquetePasajesController.List)
		public.GET("/paquetes-pasajes/sede/:id_sede", paquetePasajesController.ListBySede)
		public.GET("/paquetes-pasajes/tipo-tour/:id_tipo_tour", paquetePasajesController.ListByTipoTour)
		public.GET("/paquetes-pasajes/:id", paquetePasajesController.GetByID)
		// Ruta para obtener tours disponibles en un rango de fechas (para el público)
		// Tours programados disponibles (acceso público)
		// Rutas públicas para TourProgramado
		public.GET("/tours/disponibles", tourProgramadoController.ListToursProgramadosDisponibles)
		public.GET("/tours/disponibilidad/:fecha", tourProgramadoController.GetDisponibilidadDia)
		public.GET("/tours/:id", tourProgramadoController.GetByID)
		public.GET("/tours/disponibles-en-fecha/:fecha", tourProgramadoController.GetToursDisponiblesEnFecha)
		public.GET("/tours/disponibles-en-rango", tourProgramadoController.GetToursDisponiblesEnRangoFechas)
		public.GET("/tours/verificar-disponibilidad", tourProgramadoController.VerificarDisponibilidadHorario)
		// AÑADIR: Ruta para verificar disponibilidad de un horario en una fecha específica
		public.GET("/tours/verificar-disponibilidad", tourProgramadoController.VerificarDisponibilidadHorario) // Métodos de pago (acceso público para ver opciones)
		public.GET("/metodos-pago", metodoPagoController.List)

		// Canales de venta (acceso público)
		public.GET("/sedes", sedeController.List)
		public.GET("/sedes/:id", sedeController.GetByID)
		public.GET("/sedes/distrito/:distrito", sedeController.GetByDistrito) // CORREGIDO
		public.GET("/sedes/pais/:pais", sedeController.GetByPais)
		// Idiomas (acceso público para ver opciones disponibles)
		public.GET("/idiomas", idiomaController.List) // NUEVA RUTA PÚBLICA

	}

	// Rutas protegidas (requieren autenticación)
	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(config))
	{
		// Cambiar contraseña (cualquier usuario autenticado)
		protected.GET("/auth/status", authController.CheckStatus)
		protected.POST("/auth/change-password", authController.ChangePassword)

		// Estas rutas deben estar fuera del grupo admin para mantenerlas separadas
		adminAuth := protected.Group("/auth")
		adminAuth.Use(middleware.RoleMiddleware("ADMIN"))
		{
			protected.GET("/auth/debug", func(ctx *gin.Context) {
				// Listar todo el contexto
				fmt.Println("Datos en el contexto:")
				for k, v := range ctx.Keys {
					fmt.Printf("  %s: %v\n", k, v)
				}

				ctx.JSON(http.StatusOK, gin.H{
					"message":      "Datos del contexto",
					"userID":       ctx.GetInt("userID"),
					"userRole":     ctx.GetString("userRole"),
					"sedeID":       ctx.GetInt("sedeID"),
					"adminConSede": ctx.GetBool("adminConSede"),
					"adminSinSede": ctx.GetBool("adminSinSede"),
				})
			})

			// Obtener todas las sedes disponibles para seleccionar
			adminAuth.GET("/sedes", authController.GetUserSedes)

			// Seleccionar una sede específica para la sesión
			adminAuth.POST("/select-sede", authController.SelectSede)
		}

		// Usuarios - Admin
		admin := protected.Group("/admin")
		admin.Use(middleware.RoleMiddleware("ADMIN"))
		{
			// Gestión de usuarios
			admin.POST("/usuarios", usuarioController.Create)
			admin.GET("/usuarios", usuarioController.List)
			admin.GET("/usuarios/:id", usuarioController.GetByID)
			admin.PUT("/usuarios/:id", usuarioController.Update)
			admin.DELETE("/usuarios/:id", usuarioController.Delete)
			admin.GET("/usuarios/rol/:rol", usuarioController.ListByRol)

			admin.GET("/usuarios/:id/idiomas", usuarioIdiomaController.GetIdiomasByUsuarioID)
			admin.POST("/usuarios/:id/idiomas", usuarioIdiomaController.AsignarIdioma)
			admin.DELETE("/usuarios/:id/idiomas/:idioma_id", usuarioIdiomaController.DesasignarIdioma)
			admin.PUT("/usuarios/:id/idiomas", usuarioIdiomaController.ActualizarIdiomasUsuario)
			admin.GET("/idiomas/:id/usuarios", usuarioIdiomaController.GetUsuariosByIdiomaID)
			// Gestión de idiomas - NUEVAS RUTAS
			admin.POST("/idiomas", idiomaController.Create)
			admin.GET("/idiomas", idiomaController.List)
			admin.GET("/idiomas/:id", idiomaController.GetByID)
			admin.PUT("/idiomas/:id", idiomaController.Update)
			admin.DELETE("/idiomas/:id", idiomaController.Delete)

			// Gestión de embarcaciones
			admin.POST("/embarcaciones", embarcacionController.Create)
			admin.GET("/embarcaciones", embarcacionController.List)
			admin.GET("/embarcaciones/sede/:idSede", embarcacionController.ListBySede)
			admin.GET("/embarcaciones/:id", embarcacionController.GetByID)
			admin.PUT("/embarcaciones/:id", embarcacionController.Update)
			admin.DELETE("/embarcaciones/:id", embarcacionController.Delete)

			// Gestión de tipos de tour
			admin.POST("/tipos-tour", tipoTourController.Create)
			admin.GET("/tipos-tour", tipoTourController.List)
			admin.GET("/tipos-tour/sede/:idSede", tipoTourController.ListBySede)
			admin.GET("/tipos-tour/:id", tipoTourController.GetByID)
			admin.PUT("/tipos-tour/:id", tipoTourController.Update)
			admin.DELETE("/tipos-tour/:id", tipoTourController.Delete)

			// Gestión de galería de imágenes
			admin.POST("/galerias", galeriaTourController.Create)
			admin.GET("/galerias/:id", galeriaTourController.GetByID)
			admin.PUT("/galerias/:id", galeriaTourController.Update)
			admin.DELETE("/galerias/:id", galeriaTourController.Delete)
			admin.GET("/tipo-tours/:id_tipo_tour/galerias", galeriaTourController.ListByTipoTour)

			// Gestión de horarios de tour
			admin.POST("/horarios-tour", horarioTourController.Create)
			admin.GET("/horarios-tour", horarioTourController.List)
			admin.GET("/horarios-tour/:id", horarioTourController.GetByID)
			admin.PUT("/horarios-tour/:id", horarioTourController.Update)
			admin.DELETE("/horarios-tour/:id", horarioTourController.Delete)
			admin.GET("/horarios-tour/tipo/:idTipoTour", horarioTourController.ListByTipoTour)
			admin.GET("/horarios-tour/dia/:dia", horarioTourController.ListByDia)

			// Gestión de horarios de chofer
			admin.POST("/horarios-chofer", horarioChoferController.Create)
			admin.GET("/horarios-chofer", horarioChoferController.List)
			admin.GET("/horarios-chofer/:id", horarioChoferController.GetByID)
			admin.PUT("/horarios-chofer/:id", horarioChoferController.Update)
			admin.DELETE("/horarios-chofer/:id", horarioChoferController.Delete)
			admin.GET("/horarios-chofer/chofer/:idChofer", horarioChoferController.ListByChofer)
			admin.GET("/horarios-chofer/chofer/:idChofer/activos", horarioChoferController.ListActiveByChofer)
			admin.GET("/horarios-chofer/dia/:dia", horarioChoferController.ListByDia)

			// Gestión de tours programados
			// Gestión de tours programados
			admin.POST("/tours", tourProgramadoController.Create)
			admin.GET("/tours", tourProgramadoController.List)
			admin.GET("/tours/:id", tourProgramadoController.GetByID)
			admin.PUT("/tours/:id", tourProgramadoController.Update)
			admin.DELETE("/tours/:id", tourProgramadoController.Delete)
			admin.POST("/tours/:id/chofer", tourProgramadoController.AsignarChofer)
			admin.POST("/tours/:id/estado", tourProgramadoController.CambiarEstado)
			admin.GET("/tours/programacion-semanal", tourProgramadoController.GetProgramacionSemanal)

			// Rutas adicionales para tours programados con filtros
			admin.POST("/tours", tourProgramadoController.Create)
			admin.GET("/tours", tourProgramadoController.List)
			admin.GET("/tours/:id", tourProgramadoController.GetByID)
			admin.PUT("/tours/:id", tourProgramadoController.Update)
			admin.DELETE("/tours/:id", tourProgramadoController.Delete)
			admin.POST("/tours/:id/chofer", tourProgramadoController.AsignarChofer)
			admin.POST("/tours/:id/estado", tourProgramadoController.CambiarEstado)
			admin.GET("/tours/programacion-semanal", tourProgramadoController.GetProgramacionSemanal)

			// Rutas adicionales para tours programados con filtros
			admin.GET("/tours/fecha/:fecha", tourProgramadoController.ListByFecha)
			admin.GET("/tours/rango-fechas", tourProgramadoController.ListByRangoFechas)
			admin.GET("/tours/estado/:estado", tourProgramadoController.ListByEstado)
			admin.GET("/tours/embarcacion/:idEmbarcacion", tourProgramadoController.ListByEmbarcacion)
			admin.GET("/tours/chofer/:idChofer", tourProgramadoController.ListByChofer)
			admin.GET("/tours/tipo-tour/:idTipoTour", tourProgramadoController.ListByTipoTour)
			admin.GET("/tours/sede/:idSede", tourProgramadoController.ListBySede)
			admin.GET("/tours/vigentes", tourProgramadoController.GetToursVigentes)
			admin.GET("/tours/disponibles-en-fecha/:fecha", tourProgramadoController.GetToursDisponiblesEnFecha)
			admin.GET("/tours/disponibles-en-rango", tourProgramadoController.GetToursDisponiblesEnRangoFechas)
			admin.GET("/tours/verificar-disponibilidad", tourProgramadoController.VerificarDisponibilidadHorario)

			// Gestión de tipos de pasaje
			admin.POST("/tipos-pasaje", tipoPasajeController.Create)
			admin.GET("/tipos-pasaje", tipoPasajeController.List)
			admin.GET("/tipos-pasaje/:id", tipoPasajeController.GetByID)
			admin.PUT("/tipos-pasaje/:id", tipoPasajeController.Update)
			admin.DELETE("/tipos-pasaje/:id", tipoPasajeController.Delete)
			//
			admin.GET("/tipos-pasaje/sede/:idSede", tipoPasajeController.ListBySede)
			//admin.GET("/tipos-pasaje/sede/:id_sede", tipoPasajeController.ListBySede)
			admin.GET("/tipos-pasaje/tipo-tour/:id_tipo_tour", tipoPasajeController.ListByTipoTour)

			// Gestión de paquetes de pasajes
			admin.POST("/paquetes-pasajes", paquetePasajesController.Create)
			admin.GET("/paquetes-pasajes", paquetePasajesController.List)
			admin.GET("/paquetes-pasajes/:id", paquetePasajesController.GetByID)
			admin.PUT("/paquetes-pasajes/:id", paquetePasajesController.Update)
			admin.DELETE("/paquetes-pasajes/:id", paquetePasajesController.Delete)
			admin.GET("/paquetes-pasajes/sede/:id_sede", paquetePasajesController.ListBySede)
			admin.GET("/paquetes-pasajes/tipo-tour/:id_tipo_tour", paquetePasajesController.ListByTipoTour)

			// Gestión de métodos de pago
			admin.POST("/metodos-pago", metodoPagoController.Create)
			admin.GET("/metodos-pago", metodoPagoController.List)
			admin.GET("/metodos-pago/:id", metodoPagoController.GetByID)
			admin.PUT("/metodos-pago/:id", metodoPagoController.Update)
			admin.DELETE("/metodos-pago/:id", metodoPagoController.Delete)
			admin.GET("/metodos-pago/sede/:idSede", metodoPagoController.ListBySede)

			// Gestión de canales de venta
			admin.POST("/canales-venta", canalVentaController.Create)
			admin.GET("/canales-venta", canalVentaController.List)
			admin.GET("/canales-venta/:id", canalVentaController.GetByID)
			admin.PUT("/canales-venta/:id", canalVentaController.Update)
			admin.DELETE("/canales-venta/:id", canalVentaController.Delete)
			admin.GET("/canales-venta/sede/:idSede", canalVentaController.ListBySede)

			// Gestión de clientes
			admin.GET("/clientes", clienteController.List)
			admin.GET("/clientes/:id", clienteController.GetByID)
			admin.PUT("/clientes/:id", clienteController.Update)
			admin.DELETE("/clientes/:id", clienteController.Delete)

			// Gestión de reservas
			/*		admin.POST("/reservas", reservaController.Create)
					admin.GET("/reservas", reservaController.List)
					admin.GET("/reservas/:id", reservaController.GetByID)
					admin.PUT("/reservas/:id", reservaController.Update)
					admin.DELETE("/reservas/:id", reservaController.Delete)
					admin.POST("/reservas/:id/estado", reservaController.CambiarEstado)
					admin.GET("/reservas/cliente/:idCliente", reservaController.ListByCliente)
					admin.GET("/reservas/tour/:idTourProgramado", reservaController.ListByTourProgramado)
					admin.GET("/reservas/fecha/:fecha", reservaController.ListByFecha)
					admin.GET("/reservas/estado/:estado", reservaController.ListByEstado)
					admin.GET("/reservas/sede/:idSede", reservaController.ListBySede)
*/
// Gestión de pagos
/*	admin.POST("/pagos", pagoController.Create)
		admin.GET("/pagos", pagoController.List)
		admin.GET("/pagos/:id", pagoController.GetByID)
		admin.PUT("/pagos/:id", pagoController.Update)
		admin.DELETE("/pagos/:id", pagoController.Delete)
		admin.POST("/pagos/:id/estado", pagoController.CambiarEstado)
		admin.GET("/pagos/reserva/:idReserva", pagoController.ListByReserva)
		admin.GET("/pagos/fecha/:fecha", pagoController.ListByFecha)
		admin.GET("/pagos/estado/:estado", pagoController.ListByEstado)
		admin.GET("/pagos/reserva/:idReserva/total", pagoController.GetTotalPagadoByReserva)
		admin.GET("/pagos/cliente/:idCliente", pagoController.ListByCliente)
		admin.GET("/pagos/sede/:idSede", pagoController.ListBySede)

		// Gestión de comprobantes de pago
		admin.POST("/comprobantes", comprobantePagoController.Create)
		admin.GET("/comprobantes", comprobantePagoController.List)
		admin.GET("/comprobantes/:id", comprobantePagoController.GetByID)
		admin.GET("/comprobantes/buscar", comprobantePagoController.GetByTipoAndNumero)
		admin.PUT("/comprobantes/:id", comprobantePagoController.Update)
		admin.DELETE("/comprobantes/:id", comprobantePagoController.Delete)
		admin.POST("/comprobantes/:id/estado", comprobantePagoController.CambiarEstado)
		admin.GET("/comprobantes/reserva/:idReserva", comprobantePagoController.ListByReserva)
		admin.GET("/comprobantes/fecha/:fecha", comprobantePagoController.ListByFecha)
		admin.GET("/comprobantes/tipo/:tipo", comprobantePagoController.ListByTipo)
		admin.GET("/comprobantes/estado/:estado", comprobantePagoController.ListByEstado)
		admin.GET("/comprobantes/cliente/:idCliente", comprobantePagoController.ListByCliente)

		// Gestión de sedes
		admin.POST("/sedes", sedeController.Create)
		admin.PUT("/sedes/:id", sedeController.Update)
		admin.DELETE("/sedes/:id", sedeController.Delete)
		admin.POST("/sedes/:id/restore", sedeController.Restore)
		admin.GET("/sedes", sedeController.List)
		admin.GET("/sedes/:id", sedeController.GetByID)
		admin.GET("/sedes/distrito/:distrito", sedeController.GetByDistrito) // CORREGIDO
		admin.GET("/sedes/pais/:pais", sedeController.GetByPais)
	}

	// Vendedores
	vendedor := protected.Group("/vendedor")
	vendedor.Use(middleware.RoleMiddleware("ADMIN", "VENDEDOR"))
	{
		// Ver idiomas (solo lectura)
		vendedor.GET("/idiomas", idiomaController.List) // NUEVA RUTA

		// Ver embarcaciones (solo lectura)
		vendedor.GET("/embarcaciones", embarcacionController.List)
		vendedor.GET("/embarcaciones/:id", embarcacionController.GetByID)

		// Ver tipos de tour (solo lectura)
		vendedor.GET("/tipos-tour", tipoTourController.List)
		vendedor.GET("/tipos-tour/:id", tipoTourController.GetByID)

		// Ver galería de imágenes (solo lectura)
		vendedor.GET("/tipo-tours/:id_tipo_tour/galerias", galeriaTourController.ListByTipoTour)
		vendedor.GET("/galerias/:id", galeriaTourController.GetByID)

		// Ver horarios de tour (solo lectura)
		vendedor.GET("/horarios-tour", horarioTourController.List)
		vendedor.GET("/horarios-tour/:id", horarioTourController.GetByID)
		vendedor.GET("/horarios-tour/tipo/:idTipoTour", horarioTourController.ListByTipoTour)
		vendedor.GET("/horarios-tour/dia/:dia", horarioTourController.ListByDia)

		// Ver horarios de choferes disponibles (solo lectura)
		vendedor.GET("/horarios-chofer/dia/:dia", horarioChoferController.ListByDia)

		// Ver tours programados (solo lectura)
		vendedor.GET("/tours", tourProgramadoController.List)
		vendedor.GET("/tours/:id", tourProgramadoController.GetByID)
		vendedor.GET("/tours/fecha/:fecha", tourProgramadoController.ListByFecha)
		vendedor.GET("/tours/rango", tourProgramadoController.ListByRangoFechas)
		vendedor.GET("/tours/estado/:estado", tourProgramadoController.ListByEstado)
		vendedor.GET("/tours/disponibles", tourProgramadoController.ListToursProgramadosDisponibles)
		vendedor.GET("/tours/sede/:idSede", tourProgramadoController.ListBySede)

		// Ver tipos de pasaje (solo lectura)
		vendedor.GET("/tipos-pasaje", tipoPasajeController.List)
		vendedor.GET("/tipos-pasaje/:id", tipoPasajeController.GetByID)
		vendedor.GET("/tipos-pasaje/sede/:idSede", tipoPasajeController.ListBySede)
		vendedor.GET("/tipos-pasaje/tipo-tour/:id_tipo_tour", tipoPasajeController.ListByTipoTour)

		// Ver paquetes de pasajes (solo lectura)
		vendedor.GET("/paquetes-pasajes", paquetePasajesController.List)
		vendedor.GET("/paquetes-pasajes/:id", paquetePasajesController.GetByID)
		vendedor.GET("/paquetes-pasajes/sede/:id_sede", paquetePasajesController.ListBySede)
		vendedor.GET("/paquetes-pasajes/tipo-tour/:id_tipo_tour", paquetePasajesController.ListByTipoTour)

		// Ver métodos de pago (solo lectura)
		vendedor.GET("/metodos-pago", metodoPagoController.List)
		vendedor.GET("/metodos-pago/:id", metodoPagoController.GetByID)
		vendedor.GET("/metodos-pago/sede/:idSede", metodoPagoController.ListBySede)

		// Ver canales de venta (solo lectura)
		vendedor.GET("/canales-venta", canalVentaController.List)
		vendedor.GET("/canales-venta/:id", canalVentaController.GetByID)
		vendedor.GET("/canales-venta/sede/:idSede", canalVentaController.ListBySede)

		// Gestión de clientes
		vendedor.POST("/clientes", clienteController.Create)
		vendedor.GET("/clientes", clienteController.List)
		vendedor.GET("/clientes/:id", clienteController.GetByID)
		vendedor.PUT("/clientes/:id", clienteController.Update)

		// Gestión de reservas
		/*	vendedor.POST("/reservas", reservaController.Create)
			vendedor.GET("/reservas", reservaController.List)
			vendedor.GET("/reservas/:id", reservaController.GetByID)
			vendedor.PUT("/reservas/:id", reservaController.Update)
			vendedor.POST("/reservas/:id/estado", reservaController.CambiarEstado)
			vendedor.GET("/reservas/cliente/:idCliente", reservaController.ListByCliente)
			vendedor.GET("/reservas/tour/:idTourProgramado", reservaController.ListByTourProgramado)
			vendedor.GET("/reservas/fecha/:fecha", reservaController.ListByFecha)
			vendedor.GET("/reservas/estado/:estado", reservaController.ListByEstado)
			vendedor.GET("/reservas/sede/:idSede", reservaController.ListBySede)
*/
// Gestión de pagos (vendedor puede registrar y ver pagos)
/*	vendedor.POST("/pagos", pagoController.Create)
		vendedor.GET("/pagos", pagoController.List)
		vendedor.GET("/pagos/:id", pagoController.GetByID)
		vendedor.GET("/pagos/reserva/:idReserva", pagoController.ListByReserva)
		vendedor.GET("/pagos/reserva/:idReserva/total", pagoController.GetTotalPagadoByReserva)
		vendedor.GET("/pagos/sede/:idSede", pagoController.ListBySede)

		// Gestión de comprobantes (vendedor puede emitir y ver comprobantes)
		vendedor.POST("/comprobantes", comprobantePagoController.Create)
		vendedor.GET("/comprobantes", comprobantePagoController.List)
		vendedor.GET("/comprobantes/:id", comprobantePagoController.GetByID)
		vendedor.GET("/comprobantes/buscar", comprobantePagoController.GetByTipoAndNumero)
		vendedor.GET("/comprobantes/reserva/:idReserva", comprobantePagoController.ListByReserva)
	}

	// Choferes
	chofer := protected.Group("/chofer")
	chofer.Use(middleware.RoleMiddleware("ADMIN", "CHOFER"))
	{
		// Ver embarcaciones asignadas
		chofer.GET("/mis-embarcaciones", func(ctx *gin.Context) {
			userID := ctx.GetInt("userID")
			ctx.Request.URL.Path = "/api/v1/admin/embarcaciones/chofer/" + strconv.Itoa(userID)
			router.HandleContext(ctx)
		})

		// Ver tipos de tour (solo lectura)
		chofer.GET("/tipos-tour", tipoTourController.List)

		// Ver horarios de tour (solo lectura)
		chofer.GET("/horarios-tour", horarioTourController.List)
		chofer.GET("/horarios-tour/dia/:dia", horarioTourController.ListByDia)

		// Ver mis horarios de trabajo
		chofer.GET("/mis-horarios", horarioChoferController.GetMyActiveHorarios)
		chofer.GET("/todos-mis-horarios", func(ctx *gin.Context) {
			userID := ctx.GetInt("userID")
			ctx.Request.URL.Path = "/api/v1/admin/horarios-chofer/chofer/" + strconv.Itoa(userID)
			router.HandleContext(ctx)
		})

		// Ver mis tours programados
		// Rutas para chofer - TourProgramado
		chofer.GET("/mis-tours", func(ctx *gin.Context) {
			userID := ctx.GetInt("userID")
			ctx.Request.URL.Path = "/api/v1/admin/tours/chofer/" + strconv.Itoa(userID)
			router.HandleContext(ctx)
		})
		chofer.GET("/tours/:id", tourProgramadoController.GetByID)
		chofer.GET("/tours/programacion-semanal", tourProgramadoController.GetProgramacionSemanal)
		chofer.GET("/tours/fecha/:fecha", tourProgramadoController.ListByFecha)

		// Ver reservas para mis tours
		/*chofer.GET("/mis-tours/:idTourProgramado/reservas", reservaController.ListByTourProgramado)
*/
/*}

// Clientes
/*cliente := protected.Group("/cliente")
cliente.Use(middleware.RoleMiddleware("ADMIN", "CLIENTE"))
{
	// Cambiar contraseña (cliente)
	cliente.POST("/change-password", clienteController.ChangePassword)

	// Ver tipos de tour disponibles (solo lectura)
	cliente.GET("/tipos-tour", tipoTourController.List)
	cliente.GET("/tipos-tour/:id", tipoTourController.GetByID)

	// Ver horarios de tour disponibles (solo lectura)
	cliente.GET("/horarios-tour/tipo/:idTipoTour", horarioTourController.ListByTipoTour)

	// Ver tours disponibles
	cliente.GET("/tours/disponibles", tourProgramadoController.ListToursProgramadosDisponibles)
	cliente.GET("/tours/disponibilidad/:fecha", tourProgramadoController.GetDisponibilidadDia)
	cliente.GET("/tours/:id", tourProgramadoController.GetByID)

	// Ver tipos de pasaje (solo lectura)
	cliente.GET("/tipos-pasaje", tipoPasajeController.List)
	cliente.GET("/tipos-pasaje/sede/:idSede", tipoPasajeController.ListBySede)
	cliente.GET("/tipos-pasaje/tipo-tour/:id_tipo_tour", tipoPasajeController.ListByTipoTour)

	// Ver paquetes de pasajes (solo lectura)
	cliente.GET("/paquetes-pasajes", paquetePasajesController.List)
	cliente.GET("/paquetes-pasajes/:id", paquetePasajesController.GetByID)
	cliente.GET("/paquetes-pasajes/sede/:id_sede", paquetePasajesController.ListBySede)
	cliente.GET("/paquetes-pasajes/tipo-tour/:id_tipo_tour", paquetePasajesController.ListByTipoTour)

	// Ver métodos de pago (solo lectura)
	cliente.GET("/metodos-pago", metodoPagoController.List)

	// Ver canales de venta (solo lectura)
	cliente.GET("/canales-venta", canalVentaController.List)

	// Gestión del perfil propio
	cliente.GET("/mi-perfil", func(ctx *gin.Context) {
		clienteID := ctx.GetInt("userID")
		ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: strconv.Itoa(clienteID)})
		clienteController.GetByID(ctx)
	})

	cliente.PUT("/mi-perfil", func(ctx *gin.Context) {
		clienteID := ctx.GetInt("userID")
		ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: strconv.Itoa(clienteID)})
		clienteController.Update(ctx)
	})

	// Gestión de mis reservas
	/*		cliente.POST("/reservas", reservaController.Create)
			cliente.GET("/mis-reservas", reservaController.ListMyReservas)
			cliente.GET("/reservas/:id", reservaController.GetByID)
			cliente.POST("/reservas/:id/estado", reservaController.CambiarEstado)
*/
// Ver galería de imágenes (solo lectura)
/*	cliente.GET("/tipo-tours/:id_tipo_tour/galerias", galeriaTourController.ListByTipoTour)
			cliente.GET("/galerias/:id", galeriaTourController.GetByID)

			// Ver mis pagos
			cliente.GET("/mis-pagos", func(ctx *gin.Context) {
				clienteID := ctx.GetInt("userID")
				ctx.Request.URL.Path = "/api/v1/admin/pagos/cliente/" + strconv.Itoa(clienteID)
				router.HandleContext(ctx)
			})

			// Ver mis comprobantes
			cliente.GET("/mis-comprobantes", func(ctx *gin.Context) {
				clienteID := ctx.GetInt("userID")
				ctx.Request.URL.Path = "/api/v1/admin/comprobantes/cliente/" + strconv.Itoa(clienteID)
				router.HandleContext(ctx)
			})
		}
	}
}
*/
