-- Tabla de sedes
CREATE TABLE sede (
    id_sede SERIAL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    direccion VARCHAR(255) NOT NULL,
    telefono VARCHAR(20),
    correo VARCHAR(100),
    ciudad VARCHAR(100) NOT NULL,
    provincia VARCHAR(100),
    pais VARCHAR(100) NOT NULL,
    eliminado BOOLEAN DEFAULT FALSE
);

 
-- Tabla de usuarios
CREATE TABLE usuario (
    id_usuario SERIAL PRIMARY KEY,
    id_sede INT,
    nombres VARCHAR(100) NOT NULL,
    apellidos VARCHAR(100) NOT NULL,
    correo VARCHAR(100) UNIQUE,
    telefono VARCHAR(20),
    direccion VARCHAR(255),
    fecha_nacimiento DATE,
    rol VARCHAR(20) NOT NULL,
    nacionalidad VARCHAR(50),
    tipo_de_documento VARCHAR(50) NOT NULL,
    numero_documento VARCHAR(20) NOT NULL,
    fecha_registro TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    contrasena VARCHAR(255),
    eliminado BOOLEAN DEFAULT FALSE,
    UNIQUE (numero_documento),
    FOREIGN KEY (id_sede) REFERENCES sede(id_sede),
    CONSTRAINT check_user_sede CHECK (
        (rol = 'ADMIN' AND id_sede IS NULL) OR
        (rol != 'ADMIN' AND id_sede IS NOT NULL)
    ),
    CONSTRAINT check_valid_rol CHECK (
        rol IN ('ADMIN', 'VENDEDOR', 'CHOFER', 'CLIENTE')
    )
);


-- Tabla de embarcaciones (actualizada con estado)
CREATE TABLE embarcacion (
    id_embarcacion SERIAL PRIMARY KEY,
    id_sede INT NOT NULL,
    nombre VARCHAR(100) NOT NULL,
    capacidad INT NOT NULL,
    descripcion VARCHAR(255),
    eliminado BOOLEAN DEFAULT FALSE,
    id_usuario INT NOT NULL,
    estado VARCHAR(20) NOT NULL DEFAULT 'DISPONIBLE',
    FOREIGN KEY (id_usuario) REFERENCES usuario(id_usuario),
    FOREIGN KEY (id_sede) REFERENCES sede(id_sede),
    CHECK (estado IN ('DISPONIBLE', 'OCUPADA', 'MANTENIMIENTO', 'FUERA_DE_SERVICIO'))
);

-- Tabla de tipo de tour
CREATE TABLE tipo_tour (
    id_tipo_tour SERIAL PRIMARY KEY,
    id_sede INT NOT NULL,
    nombre VARCHAR(100) NOT NULL,
    descripcion VARCHAR(255),
    duracion_minutos INT NOT NULL,
    precio_base DECIMAL(10,2) NOT NULL,
    cantidad_pasajeros INT NOT NULL,
    url_imagen VARCHAR(255),
    eliminado BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (id_sede) REFERENCES sede(id_sede)
);

-- Tabla de horario de tour
CREATE TABLE horario_tour (
    id_horario SERIAL PRIMARY KEY,
    id_tipo_tour INT NOT NULL,
    id_sede INT NOT NULL,
    hora_inicio TIME NOT NULL,
    hora_fin TIME NOT NULL,
    disponible_lunes BOOLEAN DEFAULT FALSE,
    disponible_martes BOOLEAN DEFAULT FALSE,
    disponible_miercoles BOOLEAN DEFAULT FALSE,
    disponible_jueves BOOLEAN DEFAULT FALSE,
    disponible_viernes BOOLEAN DEFAULT FALSE,
    disponible_sabado BOOLEAN DEFAULT FALSE,
    disponible_domingo BOOLEAN DEFAULT FALSE,
    eliminado BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (id_tipo_tour) REFERENCES tipo_tour(id_tipo_tour),
    FOREIGN KEY (id_sede) REFERENCES sede(id_sede)
);

-- Tabla de horario de choferes
CREATE TABLE horario_chofer (
    id_horario_chofer SERIAL PRIMARY KEY,
    id_usuario INT NOT NULL,
    id_sede INT NOT NULL,
    hora_inicio TIME NOT NULL,
    hora_fin TIME NOT NULL,
    disponible_lunes BOOLEAN DEFAULT FALSE,
    disponible_martes BOOLEAN DEFAULT FALSE,
    disponible_miercoles BOOLEAN DEFAULT FALSE,
    disponible_jueves BOOLEAN DEFAULT FALSE,
    disponible_viernes BOOLEAN DEFAULT FALSE,
    disponible_sabado BOOLEAN DEFAULT FALSE,
    disponible_domingo BOOLEAN DEFAULT FALSE,
    fecha_inicio DATE NOT NULL,
    fecha_fin DATE,
    eliminado BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (id_usuario) REFERENCES usuario(id_usuario),
    FOREIGN KEY (id_sede) REFERENCES sede(id_sede)
);

-- Tabla de tours programados
CREATE TABLE tour_programado (
    id_tour_programado SERIAL PRIMARY KEY,
    id_tipo_tour INT NOT NULL,
    id_embarcacion INT NOT NULL,
    id_horario INT NOT NULL,
    id_sede INT NOT NULL,
    fecha DATE NOT NULL,
    cupo_maximo INT NOT NULL,
    cupo_disponible INT NOT NULL,
    estado VARCHAR(20) DEFAULT 'PROGRAMADO',
    eliminado BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (id_tipo_tour) REFERENCES tipo_tour(id_tipo_tour),
    FOREIGN KEY (id_embarcacion) REFERENCES embarcacion(id_embarcacion),
    FOREIGN KEY (id_horario) REFERENCES horario_tour(id_horario),
    FOREIGN KEY (id_sede) REFERENCES sede(id_sede),
    UNIQUE (id_embarcacion, fecha, id_horario)
);

-- Tabla de métodos de pago
CREATE TABLE metodo_pago (
    id_metodo_pago SERIAL PRIMARY KEY,
    id_sede INT NOT NULL,
    nombre VARCHAR(50) NOT NULL,
    descripcion VARCHAR(255),
    eliminado BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (id_sede) REFERENCES sede(id_sede)
);

-- Tabla de canales de venta
CREATE TABLE canal_venta (
    id_canal SERIAL PRIMARY KEY,
    id_sede INT NOT NULL,
    nombre VARCHAR(50) NOT NULL,
    descripcion VARCHAR(255),
    eliminado BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (id_sede) REFERENCES sede(id_sede)
);

-- Tabla de cliente
CREATE TABLE cliente (
    id_cliente SERIAL PRIMARY KEY,
    tipo_documento VARCHAR(50) NOT NULL,
    numero_documento VARCHAR(20) NOT NULL,
    nombres VARCHAR(100) NOT NULL,
    apellidos VARCHAR(100) NOT NULL,
    correo VARCHAR(100),
    contrasena VARCHAR(255),
    eliminado BOOLEAN DEFAULT FALSE
);

-- Tabla de reservas
CREATE TABLE reserva (
    id_reserva SERIAL PRIMARY KEY,
    id_vendedor INT,
    id_cliente INT NOT NULL,
    id_tour_programado INT NOT NULL,
    id_canal INT NOT NULL,
    id_sede INT NOT NULL,
    fecha_reserva TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    total_pagar DECIMAL(10,2) NOT NULL,
    notas TEXT,
    estado VARCHAR(20) DEFAULT 'RESERVADO',
    eliminado BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (id_vendedor) REFERENCES usuario(id_usuario),
    FOREIGN KEY (id_cliente) REFERENCES cliente(id_cliente),
    FOREIGN KEY (id_tour_programado) REFERENCES tour_programado(id_tour_programado),
    FOREIGN KEY (id_canal) REFERENCES canal_venta(id_canal),
    FOREIGN KEY (id_sede) REFERENCES sede(id_sede)
);

-- Tabla de tipo de pasaje
CREATE TABLE tipo_pasaje (
    id_tipo_pasaje SERIAL PRIMARY KEY,
    id_sede INT NOT NULL,
    nombre VARCHAR(100) NOT NULL,
    costo DECIMAL(10,2) NOT NULL,
    edad VARCHAR(50),
    eliminado BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (id_sede) REFERENCES sede(id_sede)
);

-- Tabla intermedia pasajes_cantidad
CREATE TABLE pasajes_cantidad (
    id_pasajes_cantidad SERIAL PRIMARY KEY,
    id_reserva INT NOT NULL,
    id_tipo_pasaje INT NOT NULL,
    cantidad INT NOT NULL,
    eliminado BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (id_reserva) REFERENCES reserva(id_reserva),
    FOREIGN KEY (id_tipo_pasaje) REFERENCES tipo_pasaje(id_tipo_pasaje),
    UNIQUE (id_reserva, id_tipo_pasaje)
);

-- Tabla de pagos
CREATE TABLE pago (
    id_pago SERIAL PRIMARY KEY,
    id_reserva INT NOT NULL,
    id_metodo_pago INT NOT NULL,
    id_canal INT NOT NULL,
    id_sede INT NOT NULL,
    monto DECIMAL(10,2) NOT NULL,
    fecha_pago TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    comprobante VARCHAR(100),
    estado VARCHAR(20) DEFAULT 'PROCESADO',
    eliminado BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (id_reserva) REFERENCES reserva(id_reserva),
    FOREIGN KEY (id_metodo_pago) REFERENCES metodo_pago(id_metodo_pago),
    FOREIGN KEY (id_canal) REFERENCES canal_venta(id_canal),
    FOREIGN KEY (id_sede) REFERENCES sede(id_sede)
);

-- Tabla de comprobantes de pago
CREATE TABLE comprobante_pago (
    id_comprobante SERIAL PRIMARY KEY,
    id_reserva INT NOT NULL,
    id_sede INT NOT NULL,
    tipo VARCHAR(20) NOT NULL,
    numero_comprobante VARCHAR(20) NOT NULL,
    fecha_emision TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    subtotal DECIMAL(10,2) NOT NULL,
    igv DECIMAL(10,2) NOT NULL,
    total DECIMAL(10,2) NOT NULL,
    estado VARCHAR(20) DEFAULT 'EMITIDO',
    eliminado BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (id_reserva) REFERENCES reserva(id_reserva),
    FOREIGN KEY (id_sede) REFERENCES sede(id_sede),
    UNIQUE (tipo, numero_comprobante)
);