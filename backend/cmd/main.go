/*
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sistema-toursseft/internal/config"
	"sistema-toursseft/internal/controladores"
	"sistema-toursseft/internal/repositorios"
	"sistema-toursseft/internal/rutas"
	"sistema-toursseft/internal/servicios"
	"sistema-toursseft/internal/utils"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// Cargar configuración
	cfg := config.LoadConfig()

	// Configurar modo de Gin según entorno
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Inicializar router
	// Inicializar router
	router := gin.Default()

	// Middleware de recuperación y logging
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Configurar CORS con más opciones y headers
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Content-Length",
			"Accept",
			"Authorization",
			"X-Requested-With",
			"X-CSRF-Token",
		},
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		AllowWildcard:    true,
		AllowWebSockets:  true,
	}))

	// Middleware adicional para CORS
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Health check endpoint con más información
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "ok",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"version":   "1.0.0",
			"service":   "sistema-tours-api",
		})
	})
	utils.InitValidator()

	// Conectar a la base de datos con reintentos
	db, err := connectDBWithRetry(cfg)
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}
	defer db.Close()

	log.Println("Ejecutando migraciones...")
	if err := runMigrations(db); err != nil {
		log.Fatalf("Error al ejecutar migraciones: %v", err)
	}

	// Inicializar repositorios
	usuarioRepo := repositorios.NewUsuarioRepository(db)
	embarcacionRepo := repositorios.NewEmbarcacionRepository(db)
	sedeRepo := repositorios.NewSedeRepository(db)
	tipoTourRepo := repositorios.NewTipoTourRepository(db)
	horarioTourRepo := repositorios.NewHorarioTourRepository(db)
	horarioChoferRepo := repositorios.NewHorarioChoferRepository(db)
	tourProgramadoRepo := repositorios.NewTourProgramadoRepository(db)
	metodoPagoRepo := repositorios.NewMetodoPagoRepository(db)
	tipoPasajeRepo := repositorios.NewTipoPasajeRepository(db)
	canalVentaRepo := repositorios.NewCanalVentaRepository(db)
	clienteRepo := repositorios.NewClienteRepository(db)
	reservaRepo := repositorios.NewReservaRepository(db)
	pagoRepo := repositorios.NewPagoRepository(db)
	comprobantePagoRepo := repositorios.NewComprobantePagoRepository(db)

	// Inicializar servicios
	authService := servicios.NewAuthService(usuarioRepo, cfg)
	usuarioService := servicios.NewUsuarioService(usuarioRepo)
	sedeService := servicios.NewSedeService(sedeRepo)
	embarcacionService := servicios.NewEmbarcacionService(embarcacionRepo, usuarioRepo)
	tipoTourService := servicios.NewTipoTourService(tipoTourRepo, sedeRepo)
	// En main.go, línea 64:
	horarioTourService := servicios.NewHorarioTourService(horarioTourRepo, tipoTourRepo, sedeRepo)
	// Inicializar servicios
	horarioChoferService := servicios.NewHorarioChoferService(horarioChoferRepo, usuarioRepo, sedeRepo)

	// LÍNEA ACTUALIZADA: Se agregó el parámetro sedeRepo al constructor de TourProgramadoService
	tourProgramadoService := servicios.NewTourProgramadoService(tourProgramadoRepo, tipoTourRepo, embarcacionRepo, horarioTourRepo, sedeRepo)

	metodoPagoService := servicios.NewMetodoPagoService(metodoPagoRepo, sedeRepo)
	tipoPasajeService := servicios.NewTipoPasajeService(tipoPasajeRepo, sedeRepo)
	canalVentaService := servicios.NewCanalVentaService(canalVentaRepo, sedeRepo)
	clienteService := servicios.NewClienteService(clienteRepo)

	// Actualizado para incluir sedeRepo
	reservaService := servicios.NewReservaService(
		db,
		reservaRepo,
		clienteRepo,
		tourProgramadoRepo,
		canalVentaRepo,
		tipoPasajeRepo,
		usuarioRepo,
		sedeRepo, // Añadido el repositorio de sedes
	)

	pagoService := servicios.NewPagoService(
		pagoRepo,
		reservaRepo,
		metodoPagoRepo,
		canalVentaRepo,
	)
	comprobantePagoService := servicios.NewComprobantePagoService(
		comprobantePagoRepo,
		reservaRepo,
		pagoRepo,
	)

	// Middleware global para agregar la configuración al contexto
	router.Use(func(c *gin.Context) {
		c.Set("config", cfg)
		c.Next()
	})

	// Inicializar controladores
	authController := controladores.NewAuthController(authService)
	usuarioController := controladores.NewUsuarioController(usuarioService)
	embarcacionController := controladores.NewEmbarcacionController(embarcacionService)
	tipoTourController := controladores.NewTipoTourController(tipoTourService)
	horarioTourController := controladores.NewHorarioTourController(horarioTourService)
	horarioChoferController := controladores.NewHorarioChoferController(horarioChoferService)
	tourProgramadoController := controladores.NewTourProgramadoController(tourProgramadoService)
	metodoPagoController := controladores.NewMetodoPagoController(metodoPagoService)
	tipoPasajeController := controladores.NewTipoPasajeController(tipoPasajeService)
	canalVentaController := controladores.NewCanalVentaController(canalVentaService)
	sedeController := controladores.NewSedeController(sedeService)
	clienteController := controladores.NewClienteController(clienteService, cfg)
	reservaController := controladores.NewReservaController(reservaService)
	pagoController := controladores.NewPagoController(pagoService)
	comprobantePagoController := controladores.NewComprobantePagoController(comprobantePagoService)

	// Configurar rutas
	rutas.SetupRoutes(
		router,
		cfg,
		authController,
		usuarioController,
		embarcacionController,
		tipoTourController,
		horarioTourController,
		horarioChoferController,
		tourProgramadoController,
		tipoPasajeController,
		metodoPagoController,
		canalVentaController,
		clienteController,
		reservaController,
		pagoController,
		comprobantePagoController,
		sedeController,
	)

	// Iniciar servidor
	serverAddr := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)
	log.Printf("Servidor iniciado en %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Error al iniciar servidor: %v", err)
	}
}

// connectDBWithRetry establece conexión con la base de datos PostgreSQL con reintentos
func connectDBWithRetry(cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	var db *sql.DB
	var err error

	maxRetries := 5
	retryInterval := 5 * time.Second

	for i := 0; i < maxRetries; i++ {
		log.Printf("Intentando conectar a la base de datos (intento %d/%d)...", i+1, maxRetries)

		db, err = sql.Open("postgres", dsn)
		if err != nil {
			log.Printf("Error al abrir conexión: %v. Reintentando en %s...", err, retryInterval)
			time.Sleep(retryInterval)
			continue
		}

		// Verificar conexión
		err = db.Ping()
		if err == nil {
			log.Println("Conexión exitosa a la base de datos")
			return db, nil
		}

		log.Printf("Error al verificar conexión: %v. Reintentando en %s...", err, retryInterval)
		db.Close()
		time.Sleep(retryInterval)
	}

	return nil, fmt.Errorf("no se pudo conectar a la base de datos después de %d intentos: %v", maxRetries, err)
}

// runMigrations ejecuta las migraciones de la base de datos
func runMigrations(db *sql.DB) error {
	// Verificar si la tabla sede existe
	var existsSede bool
	err := db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'sede')").Scan(&existsSede)
	if err != nil {
		return fmt.Errorf("error al verificar tabla sede: %v", err)
	}

	// Si no existe la tabla sede, ejecutamos las migraciones completas
	if !existsSede {
		log.Println("Ejecutando migraciones iniciales...")
		migrationFile, err := os.ReadFile("./migrations/crear_tablas.sql")
		if err != nil {
			return fmt.Errorf("error al leer archivo de migración: %v", err)
		}

		// Insertar datos iniciales para sede (requerido para crear usuarios)
		_, err = db.Exec(string(migrationFile))
		if err != nil {
			return fmt.Errorf("error al ejecutar migraciones: %v", err)
		}

		// Verificar si necesitamos crear datos iniciales básicos
		var countSedes int
		err = db.QueryRow("SELECT COUNT(*) FROM sede").Scan(&countSedes)
		if err != nil {
			return fmt.Errorf("error al verificar sedes existentes: %v", err)
		}

		// Si no hay sedes, insertamos una sede principal
		if countSedes == 0 {
			log.Println("Insertando sede principal...")
			_, err = db.Exec(`
				INSERT INTO sede (nombre, direccion, ciudad, pais)
				VALUES ('Sede Principal', 'Av. Principal 123', 'Lima', 'Perú')
			`)
			if err != nil {
				return fmt.Errorf("error al insertar sede principal: %v", err)
			}
		}

		// Verificar si necesitamos crear un usuario administrador
		var countAdmins int
		err = db.QueryRow("SELECT COUNT(*) FROM usuario WHERE rol = 'ADMIN'").Scan(&countAdmins)
		if err != nil {
			return fmt.Errorf("error al verificar administradores existentes: %v", err)
		}

		// Si no hay admins, insertamos uno
		/*	if countAdmins == 0 {
			log.Println("Insertando usuario administrador...")
			_, err = db.Exec(`
				INSERT INTO usuario (id_sede, nombres, apellidos, correo, telefono, direccion,
				fecha_nacimiento, rol, nacionalidad, tipo_de_documento, numero_documento, contrasena)
				VALUES (, 'Admin', 'Sistema', 'admin@sistema-tours.com', '123456789', 'Dirección Admin',
				'1990-01-01', 'ADMIN', 'Peruana', 'DNI', '12345678', '$2a$10$Lxx1J7M.A/MT6aZuIEwEoeVPnIQnAqDaJTy6cwg/K3ZGxRV7.U9b.')
			`)
			if err != nil {
				return fmt.Errorf("error al insertar usuario administrador: %v", err)
			}
		}*/
/*
		if countAdmins == 0 {
			log.Println("Insertando usuario administrador...")
			_, err = db.Exec(`
        INSERT INTO usuario (
            nombres, apellidos, correo, telefono, direccion,
            fecha_nacimiento, rol, nacionalidad, tipo_de_documento, numero_documento, contrasena
        )
        VALUES (
            'Admin', 'Sistema', 'admin@sistema-tours.com', '123456789', 'Dirección Admin',
            '1990-01-01', 'ADMIN', 'Peruana', 'DNI', '12345678', '$2a$10$Lxx1J7M.A/MT6aZuIEwEoeVPnIQnAqDaJTy6cwg/K3ZGxRV7.U9b.'
        )
    `)
			if err != nil {
				return fmt.Errorf("error al insertar usuario administrador: %v", err)
			}
		}

		log.Println("Migraciones iniciales completadas")
	} else {
		log.Println("Base de datos ya inicializada, saltando migraciones")
	}

	return nil
}
*/
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sistema-toursseft/internal/config"
	"sistema-toursseft/internal/controladores"
	"sistema-toursseft/internal/repositorios"
	"sistema-toursseft/internal/rutas"
	"sistema-toursseft/internal/servicios"
	"sistema-toursseft/internal/utils"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// Cargar configuración
	cfg := config.LoadConfig()

	// Configurar modo de Gin según entorno
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Inicializar router
	router := gin.Default()

	// Middleware de recuperación y logging
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Configurar CORS con más opciones y headers
	// En main.go, actualizar la configuración CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"https://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Content-Length",
			"Accept",
			"Authorization",
			"X-Requested-With",
			"X-CSRF-Token",
		},
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true, // Importante para cookies
		MaxAge:           12 * time.Hour,
		AllowWildcard:    true,
		AllowWebSockets:  true,
	}))

	// Health check endpoint con más información
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "ok",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"version":   "1.0.0",
			"service":   "sistema-tours-api",
		})
	})
	utils.InitValidator()

	// Conectar a la base de datos con reintentos
	db, err := connectDBWithRetry(cfg)
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}
	defer db.Close()

	log.Println("Ejecutando migraciones...")
	if err := runMigrations(db); err != nil {
		log.Fatalf("Error al ejecutar migraciones: %v", err)
	}

	// Inicializar repositorios
	usuarioRepo := repositorios.NewUsuarioRepository(db)
	embarcacionRepo := repositorios.NewEmbarcacionRepository(db)
	sedeRepo := repositorios.NewSedeRepository(db)
	tipoTourRepo := repositorios.NewTipoTourRepository(db)
	horarioTourRepo := repositorios.NewHorarioTourRepository(db)
	horarioChoferRepo := repositorios.NewHorarioChoferRepository(db)
	tourProgramadoRepo := repositorios.NewTourProgramadoRepository(db)
	metodoPagoRepo := repositorios.NewMetodoPagoRepository(db)
	tipoPasajeRepo := repositorios.NewTipoPasajeRepository(db)
	canalVentaRepo := repositorios.NewCanalVentaRepository(db)
	clienteRepo := repositorios.NewClienteRepository(db)
	reservaRepo := repositorios.NewReservaRepository(db)
	pagoRepo := repositorios.NewPagoRepository(db)
	comprobantePagoRepo := repositorios.NewComprobantePagoRepository(db)

	// Inicializar servicios
	authService := servicios.NewAuthService(usuarioRepo, sedeRepo, cfg)
	usuarioService := servicios.NewUsuarioService(usuarioRepo)
	sedeService := servicios.NewSedeService(sedeRepo)
	embarcacionService := servicios.NewEmbarcacionService(embarcacionRepo, usuarioRepo)
	tipoTourService := servicios.NewTipoTourService(tipoTourRepo, sedeRepo)
	horarioTourService := servicios.NewHorarioTourService(horarioTourRepo, tipoTourRepo, sedeRepo)
	horarioChoferService := servicios.NewHorarioChoferService(horarioChoferRepo, usuarioRepo, sedeRepo)
	tourProgramadoService := servicios.NewTourProgramadoService(tourProgramadoRepo, tipoTourRepo, embarcacionRepo, horarioTourRepo, sedeRepo)
	metodoPagoService := servicios.NewMetodoPagoService(metodoPagoRepo, sedeRepo)
	tipoPasajeService := servicios.NewTipoPasajeService(tipoPasajeRepo, sedeRepo)
	canalVentaService := servicios.NewCanalVentaService(canalVentaRepo, sedeRepo)
	clienteService := servicios.NewClienteService(clienteRepo, cfg)

	// Servicios de reserva
	reservaService := servicios.NewReservaService(
		db,
		reservaRepo,
		clienteRepo,
		tourProgramadoRepo,
		canalVentaRepo,
		tipoPasajeRepo,
		usuarioRepo,
		sedeRepo,
	)

	// Servicios de pago - actualizado para incluir sedeRepo
	pagoService := servicios.NewPagoService(
		pagoRepo,
		reservaRepo,
		metodoPagoRepo,
		canalVentaRepo,
		sedeRepo, // Añadido el repositorio de sedes
	)

	// Servicios de comprobante de pago - actualizado para incluir sedeRepo
	comprobantePagoService := servicios.NewComprobantePagoService(
		comprobantePagoRepo,
		reservaRepo,
		pagoRepo,
		sedeRepo, // Añadido el repositorio de sedes
	)

	// Middleware global para agregar la configuración al contexto
	router.Use(func(c *gin.Context) {
		c.Set("config", cfg)
		c.Next()
	})

	// Inicializar controladores
	authController := controladores.NewAuthController(authService)
	usuarioController := controladores.NewUsuarioController(usuarioService)
	embarcacionController := controladores.NewEmbarcacionController(embarcacionService)
	tipoTourController := controladores.NewTipoTourController(tipoTourService)
	horarioTourController := controladores.NewHorarioTourController(horarioTourService)
	horarioChoferController := controladores.NewHorarioChoferController(horarioChoferService)
	tourProgramadoController := controladores.NewTourProgramadoController(tourProgramadoService)
	metodoPagoController := controladores.NewMetodoPagoController(metodoPagoService)
	tipoPasajeController := controladores.NewTipoPasajeController(tipoPasajeService)
	canalVentaController := controladores.NewCanalVentaController(canalVentaService)
	sedeController := controladores.NewSedeController(sedeService)
	clienteController := controladores.NewClienteController(clienteService, cfg)
	reservaController := controladores.NewReservaController(reservaService)
	pagoController := controladores.NewPagoController(pagoService)
	comprobantePagoController := controladores.NewComprobantePagoController(comprobantePagoService)

	// Configurar rutas
	rutas.SetupRoutes(
		router,
		cfg,
		authController,
		usuarioController,
		embarcacionController,
		tipoTourController,
		horarioTourController,
		horarioChoferController,
		tourProgramadoController,
		tipoPasajeController,
		metodoPagoController,
		canalVentaController,
		clienteController,
		reservaController,
		pagoController,
		comprobantePagoController,
		sedeController,
	)

	// Iniciar servidor
	serverAddr := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)
	log.Printf("Servidor iniciado en %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Error al iniciar servidor: %v", err)
	}
}

// connectDBWithRetry establece conexión con la base de datos PostgreSQL con reintentos
func connectDBWithRetry(cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	var db *sql.DB
	var err error

	maxRetries := 5
	retryInterval := 5 * time.Second

	for i := 0; i < maxRetries; i++ {
		log.Printf("Intentando conectar a la base de datos (intento %d/%d)...", i+1, maxRetries)

		db, err = sql.Open("postgres", dsn)
		if err != nil {
			log.Printf("Error al abrir conexión: %v. Reintentando en %s...", err, retryInterval)
			time.Sleep(retryInterval)
			continue
		}

		// Verificar conexión
		err = db.Ping()
		if err == nil {
			log.Println("Conexión exitosa a la base de datos")
			return db, nil
		}

		log.Printf("Error al verificar conexión: %v. Reintentando en %s...", err, retryInterval)
		db.Close()
		time.Sleep(retryInterval)
	}

	return nil, fmt.Errorf("no se pudo conectar a la base de datos después de %d intentos: %v", maxRetries, err)
}

// runMigrations ejecuta las migraciones de la base de datos
func runMigrations(db *sql.DB) error {
	// Verificar si la tabla sede existe
	var existsSede bool
	err := db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'sede')").Scan(&existsSede)
	if err != nil {
		return fmt.Errorf("error al verificar tabla sede: %v", err)
	}

	// Si no existe la tabla sede, ejecutamos las migraciones completas
	if !existsSede {
		log.Println("Ejecutando migraciones iniciales...")
		migrationFile, err := os.ReadFile("./migrations/crear_tablas.sql")
		if err != nil {
			return fmt.Errorf("error al leer archivo de migración: %v", err)
		}

		// Insertar datos iniciales para sede (requerido para crear usuarios)
		_, err = db.Exec(string(migrationFile))
		if err != nil {
			return fmt.Errorf("error al ejecutar migraciones: %v", err)
		}

		log.Println("Migraciones iniciales completadas")
	} else {
		log.Println("Base de datos ya inicializada, saltando migraciones")
	}

	return nil
}
