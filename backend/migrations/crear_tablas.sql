-- Tabla sede
CREATE TABLE sede (
    id_sede SERIAL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    direccion VARCHAR(255) NOT NULL,
    telefono VARCHAR(20),
    correo VARCHAR(100),
    distrito VARCHAR(100) NOT NULL,
    provincia VARCHAR(100),
    pais VARCHAR(100) NOT NULL,
    image_url VARCHAR(255),
    eliminado BOOLEAN DEFAULT FALSE
);
CREATE INDEX idx_sede_nombre ON sede(nombre);
CREATE INDEX idx_sede_distrito ON sede(distrito);
CREATE INDEX idx_sede_eliminado ON sede(eliminado);

-- Tabla idioma
CREATE TABLE idioma (
    id_idioma SERIAL PRIMARY KEY,  
    nombre VARCHAR(50) NOT NULL UNIQUE,
    eliminado BOOLEAN DEFAULT false
);
CREATE INDEX idx_idioma_nombre ON idioma(nombre);

-- Tabla usuario (sin la columna id_idioma)
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
    FOREIGN KEY (id_sede) REFERENCES sede(id_sede) ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT check_user_sede CHECK (
        (rol = 'ADMIN' AND id_sede IS NULL) OR
        (rol != 'ADMIN' AND id_sede IS NOT NULL)
    ),
    CONSTRAINT check_valid_rol CHECK (
        rol IN ('ADMIN', 'VENDEDOR', 'CHOFER')
    )
);
CREATE INDEX idx_usuario_sede ON usuario(id_sede);
CREATE INDEX idx_usuario_rol ON usuario(rol);
CREATE INDEX idx_usuario_nombres_apellidos ON usuario(nombres, apellidos);
CREATE INDEX idx_usuario_documento ON usuario(tipo_de_documento, numero_documento);
CREATE INDEX idx_usuario_eliminado ON usuario(eliminado);

-- NUEVA Tabla usuario_idioma (relación muchos a muchos)
CREATE TABLE usuario_idioma (
    id_usuario_idioma SERIAL PRIMARY KEY,
    id_usuario INT NOT NULL,
    id_idioma INT NOT NULL,
    nivel VARCHAR(20), -- Opcional: básico, intermedio, avanzado, nativo
    eliminado BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (id_usuario) REFERENCES usuario(id_usuario) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (id_idioma) REFERENCES idioma(id_idioma) ON UPDATE CASCADE ON DELETE RESTRICT,
    UNIQUE (id_usuario, id_idioma) -- Un usuario no puede tener el mismo idioma duplicado
);
CREATE INDEX idx_usuario_idioma_usuario ON usuario_idioma(id_usuario);
CREATE INDEX idx_usuario_idioma_idioma ON usuario_idioma(id_idioma);
CREATE INDEX idx_usuario_idioma_eliminado ON usuario_idioma(eliminado);

-- Tabla embarcacion
CREATE TABLE embarcacion (
    id_embarcacion SERIAL PRIMARY KEY,
    id_sede INT NOT NULL,
    nombre VARCHAR(100) NOT NULL,
    capacidad INT NOT NULL,
    descripcion VARCHAR(255),
    eliminado BOOLEAN DEFAULT FALSE,
    estado VARCHAR(20) NOT NULL DEFAULT 'DISPONIBLE',
    FOREIGN KEY (id_sede) REFERENCES sede(id_sede) ON UPDATE CASCADE ON DELETE RESTRICT,
    CHECK (estado IN ('DISPONIBLE', 'OCUPADA', 'MANTENIMIENTO', 'FUERA_DE_SERVICIO'))
);
CREATE INDEX idx_embarcacion_sede ON embarcacion(id_sede);
CREATE INDEX idx_embarcacion_estado ON embarcacion(estado);
CREATE INDEX idx_embarcacion_eliminado ON embarcacion(eliminado);

-- Tabla tipo_tour (también actualizada para usar la nueva relación con idioma)
-- Tabla tipo_tour simplificada
CREATE TABLE tipo_tour (
    id_tipo_tour SERIAL PRIMARY KEY,
    id_sede INT NOT NULL,
    nombre VARCHAR(100) NOT NULL,
    descripcion TEXT,
    duracion_minutos INT NOT NULL,
    url_imagen VARCHAR(255),
    eliminado BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (id_sede) REFERENCES sede(id_sede) ON UPDATE CASCADE ON DELETE RESTRICT
);
CREATE INDEX idx_tipo_tour_sede ON tipo_tour(id_sede);
CREATE INDEX idx_tipo_tour_eliminado ON tipo_tour(eliminado);

-- Tabla para la galería de imágenes de tours
CREATE TABLE galeria_tour (
    id_galeria SERIAL PRIMARY KEY,
    id_tipo_tour INT NOT NULL, -- Nueva relación con tipo_tour
    url_imagen VARCHAR(255) NOT NULL,
    descripcion TEXT,
    orden INT DEFAULT 0,
    fecha_creacion TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    eliminado BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (id_tipo_tour) REFERENCES tipo_tour(id_tipo_tour) ON UPDATE CASCADE ON DELETE CASCADE
);
CREATE INDEX idx_galeria_tour_tipo_tour ON galeria_tour(id_tipo_tour);
CREATE INDEX idx_galeria_tour_orden ON galeria_tour(orden);
CREATE INDEX idx_galeria_tour_eliminado ON galeria_tour(eliminado);


-- Tabla horario_tour
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
    FOREIGN KEY (id_tipo_tour) REFERENCES tipo_tour(id_tipo_tour) ON UPDATE CASCADE ON DELETE RESTRICT,
    FOREIGN KEY (id_sede) REFERENCES sede(id_sede) ON UPDATE CASCADE ON DELETE RESTRICT
);
CREATE INDEX idx_horario_tour_tipo_tour ON horario_tour(id_tipo_tour);
CREATE INDEX idx_horario_tour_sede ON horario_tour(id_sede);
CREATE INDEX idx_horario_tour_hora_inicio ON horario_tour(hora_inicio);
CREATE INDEX idx_horario_tour_eliminado ON horario_tour(eliminado);

-- Tabla horario_chofer
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
    FOREIGN KEY (id_usuario) REFERENCES usuario(id_usuario) ON UPDATE CASCADE ON DELETE RESTRICT,
    FOREIGN KEY (id_sede) REFERENCES sede(id_sede) ON UPDATE CASCADE ON DELETE RESTRICT
);
CREATE INDEX idx_horario_chofer_usuario ON horario_chofer(id_usuario);
CREATE INDEX idx_horario_chofer_sede ON horario_chofer(id_sede);
CREATE INDEX idx_horario_chofer_fecha ON horario_chofer(fecha_inicio, fecha_fin);
CREATE INDEX idx_horario_chofer_eliminado ON horario_chofer(eliminado);

-- Tabla tour_programado
CREATE TABLE tour_programado (
    id_tour_programado SERIAL PRIMARY KEY,
    id_tipo_tour INT NOT NULL,
    id_embarcacion INT NOT NULL,
    id_horario INT NOT NULL,
    id_sede INT NOT NULL,
    id_chofer INT,
    fecha DATE NOT NULL,
    cupo_maximo INT NOT NULL,
    cupo_disponible INT NOT NULL,
    estado VARCHAR(20) DEFAULT 'PROGRAMADO',
    eliminado BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (id_tipo_tour) REFERENCES tipo_tour(id_tipo_tour) ON UPDATE CASCADE ON DELETE RESTRICT,
    FOREIGN KEY (id_embarcacion) REFERENCES embarcacion(id_embarcacion) ON UPDATE CASCADE ON DELETE RESTRICT,
    FOREIGN KEY (id_horario) REFERENCES horario_tour(id_horario) ON UPDATE CASCADE ON DELETE RESTRICT,
    FOREIGN KEY (id_sede) REFERENCES sede(id_sede) ON UPDATE CASCADE ON DELETE RESTRICT,
    FOREIGN KEY (id_chofer) REFERENCES usuario(id_usuario) ON UPDATE CASCADE ON DELETE RESTRICT,
    UNIQUE (id_embarcacion, fecha, id_horario),
    CONSTRAINT check_chofer_rol CHECK (
        id_chofer IS NULL OR 
        EXISTS (SELECT 1 FROM usuario WHERE id_usuario = id_chofer AND rol = 'CHOFER')
    )
);
CREATE INDEX idx_tour_programado_tipo_tour ON tour_programado(id_tipo_tour);
CREATE INDEX idx_tour_programado_embarcacion ON tour_programado(id_embarcacion);
CREATE INDEX idx_tour_programado_horario ON tour_programado(id_horario);
CREATE INDEX idx_tour_programado_sede ON tour_programado(id_sede);
CREATE INDEX idx_tour_programado_chofer ON tour_programado(id_chofer);
CREATE INDEX idx_tour_programado_fecha ON tour_programado(fecha);
CREATE INDEX idx_tour_programado_estado ON tour_programado(estado);
CREATE INDEX idx_tour_programado_eliminado ON tour_programado(eliminado);


 
-- Tabla metodo_pago
CREATE TABLE metodo_pago (
    id_metodo_pago SERIAL PRIMARY KEY,
    id_sede INT NOT NULL,
    nombre VARCHAR(50) NOT NULL,
    descripcion VARCHAR(255),
    eliminado BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (id_sede) REFERENCES sede(id_sede) ON UPDATE CASCADE ON DELETE RESTRICT
);
CREATE INDEX idx_metodo_pago_sede ON metodo_pago(id_sede);
CREATE INDEX idx_metodo_pago_eliminado ON metodo_pago(eliminado);

-- Tabla canal_venta
CREATE TABLE canal_venta (
    id_canal SERIAL PRIMARY KEY,
    id_sede INT NOT NULL,
    nombre VARCHAR(50) NOT NULL,
    descripcion VARCHAR(255),
    eliminado BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (id_sede) REFERENCES sede(id_sede) ON UPDATE CASCADE ON DELETE RESTRICT
);
CREATE INDEX idx_canal_venta_sede ON canal_venta(id_sede);
CREATE INDEX idx_canal_venta_eliminado ON canal_venta(eliminado);

-- Tabla cliente
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
CREATE INDEX idx_cliente_documento ON cliente(tipo_documento, numero_documento);
CREATE INDEX idx_cliente_nombres_apellidos ON cliente(nombres, apellidos);
CREATE INDEX idx_cliente_correo ON cliente(correo);
CREATE INDEX idx_cliente_eliminado ON cliente(eliminado);

-- Tabla tipo_pasaje (C  relacionado a tipo_tour)
CREATE TABLE tipo_pasaje (
    id_tipo_pasaje SERIAL PRIMARY KEY,
    id_sede INT NOT NULL,
    id_tipo_tour INT NOT NULL,
    nombre VARCHAR(100) NOT NULL,
    costo DECIMAL(10,2) NOT NULL,
    edad VARCHAR(50),
     eliminado BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (id_sede) REFERENCES sede(id_sede) ON UPDATE CASCADE ON DELETE RESTRICT,
    FOREIGN KEY (id_tipo_tour) REFERENCES tipo_tour(id_tipo_tour) ON UPDATE CASCADE ON DELETE RESTRICT
);
CREATE INDEX idx_tipo_pasaje_sede ON tipo_pasaje(id_sede);
CREATE INDEX idx_tipo_pasaje_tipo_tour ON tipo_pasaje(id_tipo_tour);
 CREATE INDEX idx_tipo_pasaje_eliminado ON tipo_pasaje(eliminado);

-- Tabla paquete_pasajes (relacionado a tipo_tour)
CREATE TABLE paquete_pasajes (
    id_paquete SERIAL PRIMARY KEY,
    id_sede INT NOT NULL,
    id_tipo_tour INT NOT NULL,
    nombre VARCHAR(100) NOT NULL,
    descripcion TEXT,
    precio_total DECIMAL(10,2) NOT NULL,
    cantidad_total INT NOT NULL,  -- Nueva columna: número total de pasajes incluidos
    eliminado BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (id_sede) REFERENCES sede(id_sede) ON UPDATE CASCADE ON DELETE RESTRICT,
    FOREIGN KEY (id_tipo_tour) REFERENCES tipo_tour(id_tipo_tour) ON UPDATE CASCADE ON DELETE RESTRICT
);
CREATE INDEX idx_paquete_pasajes_sede ON paquete_pasajes(id_sede);
CREATE INDEX idx_paquete_pasajes_tipo_tour ON paquete_pasajes(id_tipo_tour);
 CREATE INDEX idx_paquete_pasajes_eliminado ON paquete_pasajes(eliminado);

-- Tabla paquete_pasaje_detalle
CREATE TABLE paquete_pasaje_detalle (
    id_paquete_detalle SERIAL PRIMARY KEY,
    id_paquete INT NOT NULL,
    id_reserva INT NOT NULL,
    cantidad INT NOT NULL,
    eliminado BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (id_paquete) REFERENCES paquete_pasajes(id_paquete) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (id_reserva) REFERENCES reserva(id_reserva) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE INDEX idx_paquete_pasaje_detalle_paquete ON paquete_pasaje_detalle(id_paquete);
CREATE INDEX idx_paquete_pasaje_detalle_reserva ON paquete_pasaje_detalle(id_reserva);
CREATE INDEX idx_paquete_pasaje_detalle_eliminado ON paquete_pasaje_detalle(eliminado);

-- Tabla reserva (SIN relación directa a pasajes_cantidad)
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
    FOREIGN KEY (id_vendedor) REFERENCES usuario(id_usuario) ON UPDATE CASCADE ON DELETE RESTRICT,
    FOREIGN KEY (id_cliente) REFERENCES cliente(id_cliente) ON UPDATE CASCADE ON DELETE RESTRICT,
    FOREIGN KEY (id_tour_programado) REFERENCES tour_programado(id_tour_programado) ON UPDATE CASCADE ON DELETE RESTRICT,
    FOREIGN KEY (id_canal) REFERENCES canal_venta(id_canal) ON UPDATE CASCADE ON DELETE RESTRICT,
    FOREIGN KEY (id_sede) REFERENCES sede(id_sede) ON UPDATE CASCADE ON DELETE RESTRICT,
    FOREIGN KEY (id_paquete) REFERENCES paquete_pasajes(id_paquete) ON UPDATE CASCADE ON DELETE RESTRICT

);
CREATE INDEX idx_reserva_vendedor ON reserva(id_vendedor);
CREATE INDEX idx_reserva_cliente ON reserva(id_cliente);
CREATE INDEX idx_reserva_tour_programado ON reserva(id_tour_programado);
CREATE INDEX idx_reserva_canal ON reserva(id_canal);
CREATE INDEX idx_reserva_sede ON reserva(id_sede);
CREATE INDEX idx_reserva_paquete ON reserva(id_paquete);
CREATE INDEX idx_reserva_fecha ON reserva(fecha_reserva);
CREATE INDEX idx_reserva_estado ON reserva(estado);
CREATE INDEX idx_reserva_eliminado ON reserva(eliminado);

-- Tabla pasajes_cantidad (CON relación OPCIONAL a reserva)
CREATE TABLE pasajes_cantidad (
    id_pasajes_cantidad SERIAL PRIMARY KEY,
    id_reserva INT, -- OPCIONAL: puede ser NULL cuando no está asociado a una reserva específica
    id_tipo_pasaje INT NOT NULL,
    cantidad INT NOT NULL,
    eliminado BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (id_reserva) REFERENCES reserva(id_reserva) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (id_tipo_pasaje) REFERENCES tipo_pasaje(id_tipo_pasaje) ON UPDATE CASCADE ON DELETE RESTRICT
);
CREATE INDEX idx_pasajes_cantidad_reserva ON pasajes_cantidad(id_reserva);
CREATE INDEX idx_pasajes_cantidad_tipo_pasaje ON pasajes_cantidad(id_tipo_pasaje);
CREATE INDEX idx_pasajes_cantidad_eliminado ON pasajes_cantidad(eliminado);

-- Tabla pago
CREATE TABLE pago (
    id_pago SERIAL PRIMARY KEY,
    id_reserva INT NOT NULL,
    id_metodo_pago INT NOT NULL,
    id_canal INT NOT NULL,
    id_sede INT NOT NULL,
    monto DECIMAL(10,2) NOT NULL,
    fecha_pago TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    estado VARCHAR(20) DEFAULT 'PROCESADO',
    numero_comprobante VARCHAR(20),
    url_comprobante TEXT,
    eliminado BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (id_reserva) REFERENCES reserva(id_reserva) ON UPDATE CASCADE ON DELETE RESTRICT,
    FOREIGN KEY (id_metodo_pago) REFERENCES metodo_pago(id_metodo_pago) ON UPDATE CASCADE ON DELETE RESTRICT,
    FOREIGN KEY (id_canal) REFERENCES canal_venta(id_canal) ON UPDATE CASCADE ON DELETE RESTRICT,
    FOREIGN KEY (id_sede) REFERENCES sede(id_sede) ON UPDATE CASCADE ON DELETE RESTRICT
);
CREATE INDEX idx_pago_reserva ON pago(id_reserva);
CREATE INDEX idx_pago_metodo_pago ON pago(id_metodo_pago);
CREATE INDEX idx_pago_canal ON pago(id_canal);
CREATE INDEX idx_pago_sede ON pago(id_sede);
CREATE INDEX idx_pago_fecha ON pago(fecha_pago);
CREATE INDEX idx_pago_estado ON pago(estado);
CREATE INDEX idx_pago_eliminado ON pago(eliminado);

-- Tabla devolucion_pago
CREATE TABLE devolucion_pago (
    id_devolucion SERIAL PRIMARY KEY,
    id_pago INT NOT NULL,
    fecha_devolucion TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    motivo TEXT NOT NULL,
    monto_devolucion DECIMAL(10,2) NOT NULL,
    estado VARCHAR(20) DEFAULT 'PENDIENTE',
    observaciones TEXT,
    FOREIGN KEY (id_pago) REFERENCES pago(id_pago) ON UPDATE CASCADE ON DELETE RESTRICT
);
CREATE INDEX idx_devolucion_pago_pago ON devolucion_pago(id_pago);
CREATE INDEX idx_devolucion_pago_fecha ON devolucion_pago(fecha_devolucion);
CREATE INDEX idx_devolucion_pago_estado ON devolucion_pago(estado);

 