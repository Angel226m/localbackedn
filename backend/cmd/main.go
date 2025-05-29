/*package main

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
	idiomaRepo := repositorios.NewIdiomaRepository(db)
	embarcacionRepo := repositorios.NewEmbarcacionRepository(db)
	usuarioIdiomaRepo := repositorios.NewUsuarioIdiomaRepository(db) // Nuevo repositorio

	sedeRepo := repositorios.NewSedeRepository(db)
	tipoTourRepo := repositorios.NewTipoTourRepository(db)
	tipoTourGaleriaRepo := repositorios.NewTipoTourGaleriaRepository(db)
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
	//usuarioService := servicios.NewUsuarioService(usuarioRepo)
	usuarioService := servicios.NewUsuarioService(usuarioRepo, usuarioIdiomaRepo)                         // Modificado para incluir usuarioIdiomaRepo
	usuarioIdiomaService := servicios.NewUsuarioIdiomaService(usuarioIdiomaRepo, idiomaRepo, usuarioRepo) // Nuevo servicio

	idiomaService := servicios.NewIdiomaService(idiomaRepo)
	sedeService := servicios.NewSedeService(sedeRepo)
	embarcacionService := servicios.NewEmbarcacionService(embarcacionRepo, sedeRepo)
	tipoTourService := servicios.NewTipoTourService(tipoTourRepo, sedeRepo, idiomaRepo)
	tipoTourGaleriaService := servicios.NewTipoTourGaleriaService(tipoTourGaleriaRepo, tipoTourRepo)
	horarioTourService := servicios.NewHorarioTourService(horarioTourRepo, tipoTourRepo, sedeRepo)
	horarioChoferService := servicios.NewHorarioChoferService(horarioChoferRepo, usuarioRepo, sedeRepo)

	// 🔧 LÍNEA CORREGIDA - Verifica el orden de parámetros en tu constructor TourProgramadoService
	tourProgramadoService := servicios.NewTourProgramadoService(
		tourProgramadoRepo,
		tipoTourRepo,
		embarcacionRepo,
		horarioTourRepo,
		sedeRepo,
		usuarioRepo, // *repositorios.UsuarioRepository <- FALTA ESTE

	)

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

	// Servicios de pago
	pagoService := servicios.NewPagoService(
		pagoRepo,
		reservaRepo,
		metodoPagoRepo,
		canalVentaRepo,
		sedeRepo,
	)

	// Servicios de comprobante de pago
	comprobantePagoService := servicios.NewComprobantePagoService(
		comprobantePagoRepo,
		reservaRepo,
		pagoRepo,
		sedeRepo,
	)

	// Middleware global para agregar la configuración al contexto
	router.Use(func(c *gin.Context) {
		c.Set("config", cfg)
		c.Next()
	})

	// Inicializar controladores
	authController := controladores.NewAuthController(authService)
	usuarioController := controladores.NewUsuarioController(usuarioService)
	idiomaController := controladores.NewIdiomaController(idiomaService)
	usuarioIdiomaController := controladores.NewUsuarioIdiomaController(usuarioIdiomaService) // Nuevo controlador
	embarcacionController := controladores.NewEmbarcacionController(embarcacionService)
	tipoTourController := controladores.NewTipoTourController(tipoTourService)
	tipoTourGaleriaController := controladores.NewTipoTourGaleriaController(tipoTourGaleriaService)
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
		idiomaController,
		usuarioIdiomaController,
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
		tipoTourGaleriaController,
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
	idiomaRepo := repositorios.NewIdiomaRepository(db)
	embarcacionRepo := repositorios.NewEmbarcacionRepository(db)
	usuarioIdiomaRepo := repositorios.NewUsuarioIdiomaRepository(db) // Nuevo repositorio

	sedeRepo := repositorios.NewSedeRepository(db)
	tipoTourRepo := repositorios.NewTipoTourRepository(db)

	horarioTourRepo := repositorios.NewHorarioTourRepository(db)
	horarioChoferRepo := repositorios.NewHorarioChoferRepository(db)
	tourProgramadoRepo := repositorios.NewTourProgramadoRepository(db)
	metodoPagoRepo := repositorios.NewMetodoPagoRepository(db)
	tipoPasajeRepo := repositorios.NewTipoPasajeRepository(db)
	paquetePasajesRepo := repositorios.NewPaquetePasajesRepository(db)

	canalVentaRepo := repositorios.NewCanalVentaRepository(db)
	clienteRepo := repositorios.NewClienteRepository(db)
	reservaRepo := repositorios.NewReservaRepository(db)
	pagoRepo := repositorios.NewPagoRepository(db)
	comprobantePagoRepo := repositorios.NewComprobantePagoRepository(db)

	// Inicializar servicios
	authService := servicios.NewAuthService(usuarioRepo, sedeRepo, cfg)
	//usuarioService := servicios.NewUsuarioService(usuarioRepo)
	usuarioService := servicios.NewUsuarioService(usuarioRepo, usuarioIdiomaRepo)                         // Modificado para incluir usuarioIdiomaRepo
	usuarioIdiomaService := servicios.NewUsuarioIdiomaService(usuarioIdiomaRepo, idiomaRepo, usuarioRepo) // Nuevo servicio

	idiomaService := servicios.NewIdiomaService(idiomaRepo)
	sedeService := servicios.NewSedeService(sedeRepo)
	embarcacionService := servicios.NewEmbarcacionService(embarcacionRepo, sedeRepo)
	tipoTourService := servicios.NewTipoTourService(tipoTourRepo, sedeRepo)
	paquetePasajesService := servicios.NewPaquetePasajesService(paquetePasajesRepo, sedeRepo, tipoTourRepo)

	horarioTourService := servicios.NewHorarioTourService(horarioTourRepo, tipoTourRepo, sedeRepo)
	horarioChoferService := servicios.NewHorarioChoferService(horarioChoferRepo, usuarioRepo, sedeRepo)

	// 🔧 LÍNEA CORREGIDA - Verifica el orden de parámetros en tu constructor TourProgramadoService
	tourProgramadoService := servicios.NewTourProgramadoService(
		tourProgramadoRepo,
		tipoTourRepo,
		embarcacionRepo,
		horarioTourRepo,
		sedeRepo,
		usuarioRepo, // *repositorios.UsuarioRepository <- FALTA ESTE
	)

	metodoPagoService := servicios.NewMetodoPagoService(metodoPagoRepo, sedeRepo)
	tipoPasajeService := servicios.NewTipoPasajeService(tipoPasajeRepo, sedeRepo, tipoTourRepo)
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

	// Servicios de pago
	pagoService := servicios.NewPagoService(
		pagoRepo,
		reservaRepo,
		metodoPagoRepo,
		canalVentaRepo,
		sedeRepo,
	)

	// Servicios de comprobante de pago
	comprobantePagoService := servicios.NewComprobantePagoService(
		comprobantePagoRepo,
		reservaRepo,
		pagoRepo,
		sedeRepo,
	)

	// Middleware global para agregar la configuración al contexto
	router.Use(func(c *gin.Context) {
		c.Set("config", cfg)
		c.Next()
	})

	// Inicializar controladores
	authController := controladores.NewAuthController(authService)
	usuarioController := controladores.NewUsuarioController(usuarioService)
	idiomaController := controladores.NewIdiomaController(idiomaService)
	usuarioIdiomaController := controladores.NewUsuarioIdiomaController(usuarioIdiomaService) // Nuevo controlador
	embarcacionController := controladores.NewEmbarcacionController(embarcacionService)
	tipoTourController := controladores.NewTipoTourController(tipoTourService)
	horarioTourController := controladores.NewHorarioTourController(horarioTourService)
	horarioChoferController := controladores.NewHorarioChoferController(horarioChoferService)
	tourProgramadoController := controladores.NewTourProgramadoController(tourProgramadoService)
	metodoPagoController := controladores.NewMetodoPagoController(metodoPagoService)
	tipoPasajeController := controladores.NewTipoPasajeController(tipoPasajeService)
	paquetePasajesController := controladores.NewPaquetePasajesController(paquetePasajesService)

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
		idiomaController,
		usuarioIdiomaController,
		embarcacionController,
		tipoTourController,
		horarioTourController,
		horarioChoferController,
		tourProgramadoController,
		tipoPasajeController,
		paquetePasajesController, // Nuevo controlador
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
