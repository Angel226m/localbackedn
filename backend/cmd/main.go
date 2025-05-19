/*
package main

import (

	"database/sql"
	"fmt"
	"log"
	"os"
	"sistema-tours/internal/config"
	"sistema-tours/internal/controladores"
	"sistema-tours/internal/repositorios"
	"sistema-tours/internal/rutas"
	"sistema-tours/internal/servicios"
	"sistema-tours/internal/utils"

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

		// Inicializar validador
		utils.InitValidator()

		// Conectar a la base de datos
		db, err := connectDB(cfg)
		if err != nil {
			log.Fatalf("Error al conectar a la base de datos: %v", err)
		}
		defer db.Close()

		// Inicializar repositorios
		usuarioRepo := repositorios.NewUsuarioRepository(db)
		embarcacionRepo := repositorios.NewEmbarcacionRepository(db)
		// Otros repositorios...

		// Inicializar servicios
		authService := servicios.NewAuthService(usuarioRepo, cfg)
		usuarioService := servicios.NewUsuarioService(usuarioRepo)
		embarcacionService := servicios.NewEmbarcacionService(embarcacionRepo, usuarioRepo)
		// Otros servicios...

		// Inicializar controladores
		authController := controladores.NewAuthController(authService)
		usuarioController := controladores.NewUsuarioController(usuarioService)
		embarcacionController := controladores.NewEmbarcacionController(embarcacionService)
		// Otros controladores...

		// Configurar rutas
		rutas.SetupRoutes(
			router,
			cfg,
			authController,
			usuarioController,
			embarcacionController,
			// Otros controladores...
		)

		// Iniciar servidor
		serverAddr := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)
		log.Printf("Servidor iniciado en %s", serverAddr)
		if err := router.Run(serverAddr); err != nil {
			log.Fatalf("Error al iniciar servidor: %v", err)
		}
	}

// connectDB establece conexión con la base de datos PostgreSQL

	func connectDB(cfg *config.Config) (*sql.DB, error) {
		dsn := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
		)

		db, err := sql.Open("postgres", dsn)
		if err != nil {
			return nil, err
		}

		// Verificar conexión
		if err := db.Ping(); err != nil {
			return nil, err
		}

		log.Println("Conexión exitosa a la base de datos")
		return db, nil
	}

	func runMigrations(db *sql.DB) error {
		migrationFile, err := os.ReadFile("./migrations/crear_tablas.sql")
		if err != nil {
			return err
		}

		_, err = db.Exec(string(migrationFile))
		if err != nil {
			return err
		}

		log.Println("Migraciones ejecutadas exitosamente")
		return nil
	}
*/
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sistema-tours/internal/config"
	"sistema-tours/internal/controladores"
	"sistema-tours/internal/repositorios"
	"sistema-tours/internal/rutas"
	"sistema-tours/internal/servicios"
	"sistema-tours/internal/utils"
	"time"

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

	// Inicializar validador
	utils.InitValidator()

	// Conectar a la base de datos con reintentos
	db, err := connectDBWithRetry(cfg)
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}
	defer db.Close()

	// Inicializar repositorios
	usuarioRepo := repositorios.NewUsuarioRepository(db)
	embarcacionRepo := repositorios.NewEmbarcacionRepository(db)
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
	embarcacionService := servicios.NewEmbarcacionService(embarcacionRepo, usuarioRepo)
	tipoTourService := servicios.NewTipoTourService(tipoTourRepo)
	horarioTourService := servicios.NewHorarioTourService(horarioTourRepo, tipoTourRepo)
	horarioChoferService := servicios.NewHorarioChoferService(horarioChoferRepo, usuarioRepo)
	tourProgramadoService := servicios.NewTourProgramadoService(tourProgramadoRepo, tipoTourRepo, embarcacionRepo, horarioTourRepo)
	metodoPagoService := servicios.NewMetodoPagoService(metodoPagoRepo)
	tipoPasajeService := servicios.NewTipoPasajeService(tipoPasajeRepo)
	canalVentaService := servicios.NewCanalVentaService(canalVentaRepo)
	clienteService := servicios.NewClienteService(clienteRepo)
	reservaService := servicios.NewReservaService(
		db,
		reservaRepo,
		clienteRepo,
		tourProgramadoRepo,
		canalVentaRepo,
		tipoPasajeRepo,
		usuarioRepo,
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
	clienteController := controladores.NewClienteController(clienteService, cfg) // Pasamos cfg como segundo parámetro
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

// connectDB establece conexión con la base de datos PostgreSQL (función original sin reintentos)
func connectDB(cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Verificar conexión
	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Conexión exitosa a la base de datos")
	return db, nil
}

func runMigrations(db *sql.DB) error {
	migrationFile, err := os.ReadFile("./migrations/crear_tablas.sql")
	if err != nil {
		return err
	}

	_, err = db.Exec(string(migrationFile))
	if err != nil {
		return err
	}

	log.Println("Migraciones ejecutadas exitosamente")
	return nil
}
