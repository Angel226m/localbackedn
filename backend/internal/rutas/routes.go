/*
package rutas

import (
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
	embarcacionController *controladores.EmbarcacionController,
	tipoTourController *controladores.TipoTourController,
	horarioTourController *controladores.HorarioTourController,
	horarioChoferController *controladores.HorarioChoferController,
	tourProgramadoController *controladores.TourProgramadoController,
	tipoPasajeController *controladores.TipoPasajeController,
	metodoPagoController *controladores.MetodoPagoController,
	canalVentaController *controladores.CanalVentaController,
	clienteController *controladores.ClienteController,
	reservaController *controladores.ReservaController,
	pagoController *controladores.PagoController,
	comprobantePagoController *controladores.ComprobantePagoController,
	sedeController *controladores.SedeController, // Controlador de sede
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
		public.POST("/clientes/login", clienteController.Login)

		// Tours programados disponibles (acceso público)
		public.GET("/tours/disponibles", tourProgramadoController.ListToursProgramadosDisponibles)
		public.GET("/tours/disponibilidad/:fecha", tourProgramadoController.GetDisponibilidadDia)
		public.GET("/tours/:id", tourProgramadoController.GetByID)

		// Tipos de pasaje (acceso público para ver precios)
		public.GET("/tipos-pasaje", tipoPasajeController.List)
		public.GET("/tipos-pasaje/sede/:idSede", tipoPasajeController.ListBySede) // Nueva ruta

		// Métodos de pago (acceso público para ver opciones)
		public.GET("/metodos-pago", metodoPagoController.List)

		// Canales de venta (acceso público)
		public.GET("/canales-venta", canalVentaController.List)
		public.GET("/sedes", sedeController.List)
		public.GET("/sedes/:id", sedeController.GetByID)
		public.GET("/sedes/ciudad/:ciudad", sedeController.GetByCiudad)
		public.GET("/sedes/pais/:pais", sedeController.GetByPais)
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

			// Gestión de embarcaciones
			admin.POST("/embarcaciones", embarcacionController.Create)
			admin.GET("/embarcaciones", embarcacionController.List)
			admin.GET("/embarcaciones/:id", embarcacionController.GetByID)
			admin.PUT("/embarcaciones/:id", embarcacionController.Update)
			admin.DELETE("/embarcaciones/:id", embarcacionController.Delete)
			admin.GET("/embarcaciones/chofer/:idChofer", embarcacionController.ListByChofer)

			// Gestión de tipos de tour
			admin.POST("/tipos-tour", tipoTourController.Create)
			admin.GET("/tipos-tour", tipoTourController.List)
			admin.GET("/tipos-tour/:id", tipoTourController.GetByID)
			admin.PUT("/tipos-tour/:id", tipoTourController.Update)
			admin.DELETE("/tipos-tour/:id", tipoTourController.Delete)

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
			admin.POST("/tours/:id/estado", tourProgramadoController.CambiarEstado)
			admin.GET("/tours/fecha/:fecha", tourProgramadoController.ListByFecha)
			admin.GET("/tours/rango", tourProgramadoController.ListByRangoFechas)
			admin.GET("/tours/estado/:estado", tourProgramadoController.ListByEstado)
			admin.GET("/tours/embarcacion/:idEmbarcacion", tourProgramadoController.ListByEmbarcacion)
			admin.GET("/tours/chofer/:idChofer", tourProgramadoController.ListByChofer)
			admin.GET("/tours/tipo/:idTipoTour", tourProgramadoController.ListByTipoTour)
			admin.GET("/tours/sede/:idSede", tourProgramadoController.ListBySede)

			// Gestión de tipos de pasaje
			admin.POST("/tipos-pasaje", tipoPasajeController.Create)
			admin.GET("/tipos-pasaje", tipoPasajeController.List)
			admin.GET("/tipos-pasaje/:id", tipoPasajeController.GetByID)
			admin.PUT("/tipos-pasaje/:id", tipoPasajeController.Update)
			admin.DELETE("/tipos-pasaje/:id", tipoPasajeController.Delete)
			admin.GET("/tipos-pasaje/sede/:idSede", tipoPasajeController.ListBySede) // Nueva ruta

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
			admin.POST("/reservas", reservaController.Create)
			admin.GET("/reservas", reservaController.List)
			admin.GET("/reservas/:id", reservaController.GetByID)
			admin.PUT("/reservas/:id", reservaController.Update)
			admin.DELETE("/reservas/:id", reservaController.Delete)
			admin.POST("/reservas/:id/estado", reservaController.CambiarEstado)
			admin.GET("/reservas/cliente/:idCliente", reservaController.ListByCliente)
			admin.GET("/reservas/tour/:idTourProgramado", reservaController.ListByTourProgramado)
			admin.GET("/reservas/fecha/:fecha", reservaController.ListByFecha)
			admin.GET("/reservas/estado/:estado", reservaController.ListByEstado)
			admin.GET("/reservas/sede/:idSede", reservaController.ListBySede) // Nueva ruta para listar por sede

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
			admin.GET("/pagos/sede/:idSede", pagoController.ListBySede) // ¡NUEVA! Ruta para listar pagos por sede

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
			admin.GET("/sedes", sedeController.List)                       // NUEVA - Listar todas las sedes
			admin.GET("/sedes/:id", sedeController.GetByID)                // NUEVA - Ver detalle de sede
			admin.GET("/sedes/ciudad/:ciudad", sedeController.GetByCiudad) // NUEVA - Filtrar por ciudad
			admin.GET("/sedes/pais/:pais", sedeController.GetByPais)       // NUEVA - Filtrar por país

		}

		// Vendedores
		vendedor := protected.Group("/vendedor")
		vendedor.Use(middleware.RoleMiddleware("ADMIN", "VENDEDOR"))
		{
			// Ver embarcaciones (solo lectura)
			vendedor.GET("/embarcaciones", embarcacionController.List)
			vendedor.GET("/embarcaciones/:id", embarcacionController.GetByID)

			// Ver tipos de tour (solo lectura)
			vendedor.GET("/tipos-tour", tipoTourController.List)
			vendedor.GET("/tipos-tour/:id", tipoTourController.GetByID)

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
			vendedor.GET("/tipos-pasaje/sede/:idSede", tipoPasajeController.ListBySede) // Nueva ruta

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
			vendedor.POST("/reservas", reservaController.Create)
			vendedor.GET("/reservas", reservaController.List)
			vendedor.GET("/reservas/:id", reservaController.GetByID)
			vendedor.PUT("/reservas/:id", reservaController.Update)
			vendedor.POST("/reservas/:id/estado", reservaController.CambiarEstado)
			vendedor.GET("/reservas/cliente/:idCliente", reservaController.ListByCliente)
			vendedor.GET("/reservas/tour/:idTourProgramado", reservaController.ListByTourProgramado)
			vendedor.GET("/reservas/fecha/:fecha", reservaController.ListByFecha)
			vendedor.GET("/reservas/estado/:estado", reservaController.ListByEstado)
			vendedor.GET("/reservas/sede/:idSede", reservaController.ListBySede) // Nueva ruta para listar por sede

			// Gestión de pagos (vendedor puede registrar y ver pagos)
			vendedor.POST("/pagos", pagoController.Create)
			vendedor.GET("/pagos", pagoController.List)
			vendedor.GET("/pagos/:id", pagoController.GetByID)
			vendedor.GET("/pagos/reserva/:idReserva", pagoController.ListByReserva)
			vendedor.GET("/pagos/reserva/:idReserva/total", pagoController.GetTotalPagadoByReserva)
			vendedor.GET("/pagos/sede/:idSede", pagoController.ListBySede) // ¡NUEVA! Ruta para listar pagos por sede

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
				// Obtener ID del usuario autenticado del contexto
				userID := ctx.GetInt("userID")
				// Redirigir a la ruta que obtiene embarcaciones por chofer
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
				// Obtener ID del usuario autenticado del contexto
				userID := ctx.GetInt("userID")
				// Redirigir a la ruta que obtiene horarios por chofer
				ctx.Request.URL.Path = "/api/v1/admin/horarios-chofer/chofer/" + strconv.Itoa(userID)
				router.HandleContext(ctx)
			})

			// Ver mis tours programados
			chofer.GET("/mis-tours", func(ctx *gin.Context) {
				// Obtener ID del usuario autenticado del contexto
				userID := ctx.GetInt("userID")
				// Redirigir a la ruta que obtiene tours por chofer
				ctx.Request.URL.Path = "/api/v1/admin/tours/chofer/" + strconv.Itoa(userID)
				router.HandleContext(ctx)
			})

			// Ver reservas para mis tours
			chofer.GET("/mis-tours/:idTourProgramado/reservas", reservaController.ListByTourProgramado)
		}

		// Clientes
		cliente := protected.Group("/cliente")
		cliente.Use(middleware.RoleMiddleware("ADMIN", "CLIENTE"))
		{
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
			cliente.GET("/tipos-pasaje/sede/:idSede", tipoPasajeController.ListBySede) // Nueva ruta

			// Ver métodos de pago (solo lectura)
			cliente.GET("/metodos-pago", metodoPagoController.List)

			// Ver canales de venta (solo lectura)
			cliente.GET("/canales-venta", canalVentaController.List)

			// Gestión del perfil propio
			cliente.GET("/mi-perfil", func(ctx *gin.Context) {
				// Obtener ID del cliente autenticado del contexto
				clienteID := ctx.GetInt("userID")
				// Redireccionar a la ruta que obtiene un cliente por ID
				ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: strconv.Itoa(clienteID)})
				clienteController.GetByID(ctx)
			})

			cliente.PUT("/mi-perfil", func(ctx *gin.Context) {
				// Obtener ID del cliente autenticado del contexto
				clienteID := ctx.GetInt("userID")
				// Establecer el parámetro ID en el contexto
				ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: strconv.Itoa(clienteID)})
				clienteController.Update(ctx)
			})

			// Gestión de mis reservas
			cliente.POST("/reservas", reservaController.Create)
			cliente.GET("/mis-reservas", reservaController.ListMyReservas)
			cliente.GET("/reservas/:id", reservaController.GetByID)
			cliente.POST("/reservas/:id/estado", reservaController.CambiarEstado) // Solo para cancelar

			// Ver mis pagos
			cliente.GET("/mis-pagos", func(ctx *gin.Context) {
				// Obtener ID del cliente autenticado del contexto
				clienteID := ctx.GetInt("userID")
				// Redirigir a la ruta que obtiene pagos por cliente
				ctx.Request.URL.Path = "/api/v1/admin/pagos/cliente/" + strconv.Itoa(clienteID)
				router.HandleContext(ctx)
			})

			// Ver mis comprobantes
			cliente.GET("/mis-comprobantes", func(ctx *gin.Context) {
				// Obtener ID del cliente autenticado del contexto
				clienteID := ctx.GetInt("userID")
				// Redirigir a la ruta que obtiene comprobantes por cliente
				ctx.Request.URL.Path = "/api/v1/admin/comprobantes/cliente/" + strconv.Itoa(clienteID)
				router.HandleContext(ctx)
			})
		}
	}
}
*/

package rutas

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
	embarcacionController *controladores.EmbarcacionController,
	tipoTourController *controladores.TipoTourController,
	horarioTourController *controladores.HorarioTourController,
	horarioChoferController *controladores.HorarioChoferController,
	tourProgramadoController *controladores.TourProgramadoController,
	tipoPasajeController *controladores.TipoPasajeController,
	metodoPagoController *controladores.MetodoPagoController,
	canalVentaController *controladores.CanalVentaController,
	clienteController *controladores.ClienteController,
	reservaController *controladores.ReservaController,
	pagoController *controladores.PagoController,
	comprobantePagoController *controladores.ComprobantePagoController,
	sedeController *controladores.SedeController, // Controlador de sede
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
		public.POST("/clientes/refresh", clienteController.RefreshToken) // Nueva ruta para refresh token de clientes
		public.POST("/clientes/logout", clienteController.Logout)        // Nueva ruta para logout de clientes

		// Tours programados disponibles (acceso público)
		public.GET("/tours/disponibles", tourProgramadoController.ListToursProgramadosDisponibles)
		public.GET("/tours/disponibilidad/:fecha", tourProgramadoController.GetDisponibilidadDia)
		public.GET("/tours/:id", tourProgramadoController.GetByID)

		// Tipos de pasaje (acceso público para ver precios)
		public.GET("/tipos-pasaje", tipoPasajeController.List)
		public.GET("/tipos-pasaje/sede/:idSede", tipoPasajeController.ListBySede) // Nueva ruta

		// Métodos de pago (acceso público para ver opciones)
		public.GET("/metodos-pago", metodoPagoController.List)

		// Canales de venta (acceso público)
		public.GET("/canales-venta", canalVentaController.List)
		public.GET("/sedes", sedeController.List)
		public.GET("/sedes/:id", sedeController.GetByID)
		public.GET("/sedes/ciudad/:ciudad", sedeController.GetByCiudad)
		public.GET("/sedes/pais/:pais", sedeController.GetByPais)
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

			// Gestión de embarcaciones
			admin.POST("/embarcaciones", embarcacionController.Create)
			admin.GET("/embarcaciones", embarcacionController.List)
			admin.GET("/embarcaciones/:id", embarcacionController.GetByID)
			admin.PUT("/embarcaciones/:id", embarcacionController.Update)
			admin.DELETE("/embarcaciones/:id", embarcacionController.Delete)
			admin.GET("/embarcaciones/chofer/:idChofer", embarcacionController.ListByChofer)

			// Gestión de tipos de tour
			admin.POST("/tipos-tour", tipoTourController.Create)
			admin.GET("/tipos-tour", tipoTourController.List)
			admin.GET("/tipos-tour/:id", tipoTourController.GetByID)
			admin.PUT("/tipos-tour/:id", tipoTourController.Update)
			admin.DELETE("/tipos-tour/:id", tipoTourController.Delete)

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
			admin.POST("/tours/:id/estado", tourProgramadoController.CambiarEstado)
			admin.GET("/tours/fecha/:fecha", tourProgramadoController.ListByFecha)
			admin.GET("/tours/rango", tourProgramadoController.ListByRangoFechas)
			admin.GET("/tours/estado/:estado", tourProgramadoController.ListByEstado)
			admin.GET("/tours/embarcacion/:idEmbarcacion", tourProgramadoController.ListByEmbarcacion)
			admin.GET("/tours/chofer/:idChofer", tourProgramadoController.ListByChofer)
			admin.GET("/tours/tipo/:idTipoTour", tourProgramadoController.ListByTipoTour)
			admin.GET("/tours/sede/:idSede", tourProgramadoController.ListBySede)

			// Gestión de tipos de pasaje
			admin.POST("/tipos-pasaje", tipoPasajeController.Create)
			admin.GET("/tipos-pasaje", tipoPasajeController.List)
			admin.GET("/tipos-pasaje/:id", tipoPasajeController.GetByID)
			admin.PUT("/tipos-pasaje/:id", tipoPasajeController.Update)
			admin.DELETE("/tipos-pasaje/:id", tipoPasajeController.Delete)
			admin.GET("/tipos-pasaje/sede/:idSede", tipoPasajeController.ListBySede) // Nueva ruta

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
			admin.POST("/reservas", reservaController.Create)
			admin.GET("/reservas", reservaController.List)
			admin.GET("/reservas/:id", reservaController.GetByID)
			admin.PUT("/reservas/:id", reservaController.Update)
			admin.DELETE("/reservas/:id", reservaController.Delete)
			admin.POST("/reservas/:id/estado", reservaController.CambiarEstado)
			admin.GET("/reservas/cliente/:idCliente", reservaController.ListByCliente)
			admin.GET("/reservas/tour/:idTourProgramado", reservaController.ListByTourProgramado)
			admin.GET("/reservas/fecha/:fecha", reservaController.ListByFecha)
			admin.GET("/reservas/estado/:estado", reservaController.ListByEstado)
			admin.GET("/reservas/sede/:idSede", reservaController.ListBySede) // Nueva ruta para listar por sede

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
			admin.GET("/pagos/sede/:idSede", pagoController.ListBySede) // ¡NUEVA! Ruta para listar pagos por sede

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
			admin.GET("/sedes", sedeController.List)                       // NUEVA - Listar todas las sedes
			admin.GET("/sedes/:id", sedeController.GetByID)                // NUEVA - Ver detalle de sede
			admin.GET("/sedes/ciudad/:ciudad", sedeController.GetByCiudad) // NUEVA - Filtrar por ciudad
			admin.GET("/sedes/pais/:pais", sedeController.GetByPais)       // NUEVA - Filtrar por país

		}

		// Vendedores
		vendedor := protected.Group("/vendedor")
		vendedor.Use(middleware.RoleMiddleware("ADMIN", "VENDEDOR"))
		{
			// Ver embarcaciones (solo lectura)
			vendedor.GET("/embarcaciones", embarcacionController.List)
			vendedor.GET("/embarcaciones/:id", embarcacionController.GetByID)

			// Ver tipos de tour (solo lectura)
			vendedor.GET("/tipos-tour", tipoTourController.List)
			vendedor.GET("/tipos-tour/:id", tipoTourController.GetByID)

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
			vendedor.GET("/tipos-pasaje/sede/:idSede", tipoPasajeController.ListBySede) // Nueva ruta

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
			vendedor.POST("/reservas", reservaController.Create)
			vendedor.GET("/reservas", reservaController.List)
			vendedor.GET("/reservas/:id", reservaController.GetByID)
			vendedor.PUT("/reservas/:id", reservaController.Update)
			vendedor.POST("/reservas/:id/estado", reservaController.CambiarEstado)
			vendedor.GET("/reservas/cliente/:idCliente", reservaController.ListByCliente)
			vendedor.GET("/reservas/tour/:idTourProgramado", reservaController.ListByTourProgramado)
			vendedor.GET("/reservas/fecha/:fecha", reservaController.ListByFecha)
			vendedor.GET("/reservas/estado/:estado", reservaController.ListByEstado)
			vendedor.GET("/reservas/sede/:idSede", reservaController.ListBySede) // Nueva ruta para listar por sede

			// Gestión de pagos (vendedor puede registrar y ver pagos)
			vendedor.POST("/pagos", pagoController.Create)
			vendedor.GET("/pagos", pagoController.List)
			vendedor.GET("/pagos/:id", pagoController.GetByID)
			vendedor.GET("/pagos/reserva/:idReserva", pagoController.ListByReserva)
			vendedor.GET("/pagos/reserva/:idReserva/total", pagoController.GetTotalPagadoByReserva)
			vendedor.GET("/pagos/sede/:idSede", pagoController.ListBySede) // ¡NUEVA! Ruta para listar pagos por sede

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
				// Obtener ID del usuario autenticado del contexto
				userID := ctx.GetInt("userID")
				// Redirigir a la ruta que obtiene embarcaciones por chofer
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
				// Obtener ID del usuario autenticado del contexto
				userID := ctx.GetInt("userID")
				// Redirigir a la ruta que obtiene horarios por chofer
				ctx.Request.URL.Path = "/api/v1/admin/horarios-chofer/chofer/" + strconv.Itoa(userID)
				router.HandleContext(ctx)
			})

			// Ver mis tours programados
			chofer.GET("/mis-tours", func(ctx *gin.Context) {
				// Obtener ID del usuario autenticado del contexto
				userID := ctx.GetInt("userID")
				// Redirigir a la ruta que obtiene tours por chofer
				ctx.Request.URL.Path = "/api/v1/admin/tours/chofer/" + strconv.Itoa(userID)
				router.HandleContext(ctx)
			})

			// Ver reservas para mis tours
			chofer.GET("/mis-tours/:idTourProgramado/reservas", reservaController.ListByTourProgramado)
		}

		// Clientes
		cliente := protected.Group("/cliente")
		cliente.Use(middleware.RoleMiddleware("ADMIN", "CLIENTE"))
		{
			// Cambiar contraseña (cliente)
			cliente.POST("/change-password", clienteController.ChangePassword) // Nueva ruta para cambio de contraseña

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
			cliente.GET("/tipos-pasaje/sede/:idSede", tipoPasajeController.ListBySede) // Nueva ruta

			// Ver métodos de pago (solo lectura)
			cliente.GET("/metodos-pago", metodoPagoController.List)

			// Ver canales de venta (solo lectura)
			cliente.GET("/canales-venta", canalVentaController.List)

			// Gestión del perfil propio
			cliente.GET("/mi-perfil", func(ctx *gin.Context) {
				// Obtener ID del cliente autenticado del contexto
				clienteID := ctx.GetInt("userID")
				// Redireccionar a la ruta que obtiene un cliente por ID
				ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: strconv.Itoa(clienteID)})
				clienteController.GetByID(ctx)
			})

			cliente.PUT("/mi-perfil", func(ctx *gin.Context) {
				// Obtener ID del cliente autenticado del contexto
				clienteID := ctx.GetInt("userID")
				// Establecer el parámetro ID en el contexto
				ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: strconv.Itoa(clienteID)})
				clienteController.Update(ctx)
			})

			// Gestión de mis reservas
			cliente.POST("/reservas", reservaController.Create)
			cliente.GET("/mis-reservas", reservaController.ListMyReservas)
			cliente.GET("/reservas/:id", reservaController.GetByID)
			cliente.POST("/reservas/:id/estado", reservaController.CambiarEstado) // Solo para cancelar

			// Ver mis pagos
			cliente.GET("/mis-pagos", func(ctx *gin.Context) {
				// Obtener ID del cliente autenticado del contexto
				clienteID := ctx.GetInt("userID")
				// Redirigir a la ruta que obtiene pagos por cliente
				ctx.Request.URL.Path = "/api/v1/admin/pagos/cliente/" + strconv.Itoa(clienteID)
				router.HandleContext(ctx)
			})

			// Ver mis comprobantes
			cliente.GET("/mis-comprobantes", func(ctx *gin.Context) {
				// Obtener ID del cliente autenticado del contexto
				clienteID := ctx.GetInt("userID")
				// Redirigir a la ruta que obtiene comprobantes por cliente
				ctx.Request.URL.Path = "/api/v1/admin/comprobantes/cliente/" + strconv.Itoa(clienteID)
				router.HandleContext(ctx)
			})
		}
	}
}
