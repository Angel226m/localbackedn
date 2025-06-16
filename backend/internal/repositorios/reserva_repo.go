package repositorios

import (
	"database/sql"
	"errors"
	"sistema-toursseft/internal/entidades"
	"time"
)

// ReservaRepository maneja las operaciones de base de datos para reservas
type ReservaRepository struct {
	db *sql.DB
}

// NewReservaRepository crea una nueva instancia del repositorio
func NewReservaRepository(db *sql.DB) *ReservaRepository {
	return &ReservaRepository{
		db: db,
	}
}

// GetByID obtiene una reserva por su ID
func (r *ReservaRepository) GetByID(id int) (*entidades.Reserva, error) {
	// Inicializar objeto de reserva
	reserva := &entidades.Reserva{}

	// Consulta para obtener datos básicos de la reserva
	query := `SELECT r.id_reserva, r.id_vendedor, r.id_cliente, r.id_instancia, 
              r.id_canal, r.id_sede, r.fecha_reserva, r.total_pagar, r.notas, r.estado, r.eliminado,
              c.nombres || ' ' || c.apellidos as nombre_cliente,
              COALESCE(u.nombres || ' ' || u.apellidos, 'Web') as nombre_vendedor,
              tt.nombre as nombre_tour,
              to_char(it.fecha_especifica, 'DD/MM/YYYY') as fecha_tour,
              to_char(it.hora_inicio, 'HH24:MI') as hora_inicio_tour,
              to_char(it.hora_fin, 'HH24:MI') as hora_fin_tour,
              cv.nombre as nombre_canal,
              s.nombre as nombre_sede
              FROM reserva r
              INNER JOIN cliente c ON r.id_cliente = c.id_cliente
              LEFT JOIN usuario u ON r.id_vendedor = u.id_usuario
              INNER JOIN instancia_tour it ON r.id_instancia = it.id_instancia
              INNER JOIN tour_programado tp ON it.id_tour_programado = tp.id_tour_programado
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              INNER JOIN canal_venta cv ON r.id_canal = cv.id_canal
              INNER JOIN sede s ON r.id_sede = s.id_sede
              WHERE r.id_reserva = $1 AND r.eliminado = FALSE`

	err := r.db.QueryRow(query, id).Scan(
		&reserva.ID, &reserva.IDVendedor, &reserva.IDCliente, &reserva.IDInstancia,
		&reserva.IDCanal, &reserva.IDSede, &reserva.FechaReserva, &reserva.TotalPagar,
		&reserva.Notas, &reserva.Estado, &reserva.Eliminado,
		&reserva.NombreCliente, &reserva.NombreVendedor, &reserva.NombreTour,
		&reserva.FechaTour, &reserva.HoraInicioTour, &reserva.HoraFinTour,
		&reserva.NombreCanal, &reserva.NombreSede,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("reserva no encontrada")
		}
		return nil, err
	}

	// Obtener las cantidades de pasajes individuales
	queryPasajes := `SELECT pc.id_tipo_pasaje, tp.nombre, pc.cantidad
                    FROM pasajes_cantidad pc
                    INNER JOIN tipo_pasaje tp ON pc.id_tipo_pasaje = tp.id_tipo_pasaje
                    WHERE pc.id_reserva = $1 AND pc.eliminado = FALSE`

	rowsPasajes, err := r.db.Query(queryPasajes, id)
	if err != nil {
		return nil, err
	}
	defer rowsPasajes.Close()

	// Inicializar el slice de cantidades de pasajes
	reserva.CantidadPasajes = []entidades.PasajeCantidad{}

	// Iterar por cada registro de pasajes individuales
	for rowsPasajes.Next() {
		var pasajeCantidad entidades.PasajeCantidad
		err := rowsPasajes.Scan(
			&pasajeCantidad.IDTipoPasaje, &pasajeCantidad.NombreTipo, &pasajeCantidad.Cantidad,
		)
		if err != nil {
			return nil, err
		}
		reserva.CantidadPasajes = append(reserva.CantidadPasajes, pasajeCantidad)
	}

	if err = rowsPasajes.Err(); err != nil {
		return nil, err
	}

	// Obtener los paquetes de pasajes
	queryPaquetes := `SELECT ppd.id_paquete, pp.nombre as nombre_paquete, ppd.cantidad, 
                     pp.precio_total as precio_unitario, (ppd.cantidad * pp.precio_total) as subtotal,
                     pp.cantidad_total
                     FROM paquete_pasaje_detalle ppd
                     INNER JOIN paquete_pasajes pp ON ppd.id_paquete = pp.id_paquete
                     WHERE ppd.id_reserva = $1 AND ppd.eliminado = FALSE`

	rowsPaquetes, err := r.db.Query(queryPaquetes, id)
	if err != nil {
		return nil, err
	}
	defer rowsPaquetes.Close()

	// Inicializar el slice de paquetes
	reserva.Paquetes = []entidades.PaquetePasajeDetalle{}

	// Iterar por cada registro de paquetes
	for rowsPaquetes.Next() {
		var paquete entidades.PaquetePasajeDetalle
		err := rowsPaquetes.Scan(
			&paquete.IDPaquete, &paquete.NombrePaquete, &paquete.Cantidad,
			&paquete.PrecioUnitario, &paquete.Subtotal, &paquete.CantidadTotal,
		)
		if err != nil {
			return nil, err
		}
		reserva.Paquetes = append(reserva.Paquetes, paquete)
	}

	if err = rowsPaquetes.Err(); err != nil {
		return nil, err
	}

	return reserva, nil
}

// Create guarda una nueva reserva en la base de datos
func (r *ReservaRepository) Create(reserva *entidades.NuevaReservaRequest) (int, error) {
	// Iniciar transacción
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	// Si hay error, hacer rollback
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Verificar que al menos haya un pasaje o un paquete
	if len(reserva.CantidadPasajes) == 0 && len(reserva.Paquetes) == 0 {
		return 0, errors.New("debe incluir al menos un pasaje o un paquete en la reserva")
	}

	// Primero, obtener información de la instancia del tour para verificar disponibilidad
	var cupoDisponible int
	queryInstancia := `SELECT cupo_disponible FROM instancia_tour 
                      WHERE id_instancia = $1 AND eliminado = FALSE 
                      AND estado = 'PROGRAMADO'`
	err = tx.QueryRow(queryInstancia, reserva.IDInstancia).Scan(&cupoDisponible)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("la instancia del tour no existe, está eliminada o no está programada")
		}
		return 0, err
	}

	// Calcular el total de pasajeros
	totalPasajeros := 0

	// Sumar pasajeros de pasajes individuales
	for _, pasaje := range reserva.CantidadPasajes {
		totalPasajeros += pasaje.Cantidad
	}

	// Sumar pasajeros de paquetes
	for _, paquete := range reserva.Paquetes {
		// Obtener cantidad total de pasajeros por paquete
		var cantidadPorPaquete int
		queryPaquete := `SELECT cantidad_total FROM paquete_pasajes 
                        WHERE id_paquete = $1 AND eliminado = FALSE`
		err := tx.QueryRow(queryPaquete, paquete.IDPaquete).Scan(&cantidadPorPaquete)
		if err != nil {
			return 0, err
		}

		// Multiplicar por la cantidad de paquetes seleccionados
		totalPasajeros += cantidadPorPaquete * paquete.Cantidad
	}

	// Verificar disponibilidad
	if totalPasajeros > cupoDisponible {
		return 0, errors.New("no hay suficiente cupo disponible para la reserva")
	}

	// Consulta SQL para insertar una nueva reserva
	var idReserva int
	query := `INSERT INTO reserva (id_vendedor, id_cliente, id_instancia, id_canal, id_sede, 
             total_pagar, notas, estado, eliminado)
             VALUES ($1, $2, $3, $4, $5, $6, $7, 'RESERVADO', FALSE)
             RETURNING id_reserva`

	// Ejecutar la consulta con los datos de la reserva
	err = tx.QueryRow(
		query,
		reserva.IDVendedor,
		reserva.IDCliente,
		reserva.IDInstancia,
		reserva.IDCanal,
		reserva.IDSede,
		reserva.TotalPagar,
		reserva.Notas,
	).Scan(&idReserva)

	if err != nil {
		return 0, err
	}

	// Insertar las cantidades de pasajes individuales
	for _, pasaje := range reserva.CantidadPasajes {
		// Solo insertar si la cantidad es mayor que cero
		if pasaje.Cantidad > 0 {
			queryPasaje := `INSERT INTO pasajes_cantidad (id_reserva, id_tipo_pasaje, cantidad, eliminado)
                          VALUES ($1, $2, $3, FALSE)`

			_, err = tx.Exec(queryPasaje, idReserva, pasaje.IDTipoPasaje, pasaje.Cantidad)
			if err != nil {
				return 0, err
			}
		}
	}

	// Insertar los paquetes de pasajes
	for _, paquete := range reserva.Paquetes {
		// Solo insertar si la cantidad es mayor que cero
		if paquete.Cantidad > 0 {
			queryPaquete := `INSERT INTO paquete_pasaje_detalle (id_reserva, id_paquete, cantidad, eliminado)
                           VALUES ($1, $2, $3, FALSE)`

			_, err = tx.Exec(queryPaquete, idReserva, paquete.IDPaquete, paquete.Cantidad)
			if err != nil {
				return 0, err
			}
		}
	}

	// Actualizar el cupo disponible en la instancia del tour
	queryUpdateCupo := `UPDATE instancia_tour 
                       SET cupo_disponible = cupo_disponible - $1 
                       WHERE id_instancia = $2`
	_, err = tx.Exec(queryUpdateCupo, totalPasajeros, reserva.IDInstancia)
	if err != nil {
		return 0, err
	}

	// Commit de la transacción
	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return idReserva, nil
}

// Update actualiza la información de una reserva existente
func (r *ReservaRepository) Update(id int, reserva *entidades.ActualizarReservaRequest) error {
	// Iniciar transacción
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// Si hay error, hacer rollback
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Primero, obtener la reserva actual para calcular la diferencia de pasajeros
	var idInstanciaActual int
	var estadoActual string
	queryReservaActual := `SELECT id_instancia, estado FROM reserva 
                          WHERE id_reserva = $1 AND eliminado = FALSE`

	err = tx.QueryRow(queryReservaActual, id).Scan(&idInstanciaActual, &estadoActual)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("reserva no encontrada")
		}
		return err
	}

	// Obtener cantidad actual de pasajeros
	var totalPasajerosActual int
	if estadoActual != "CANCELADA" { // Solo contar si la reserva no estaba cancelada
		// Contar pasajeros de pasajes individuales
		queryPasajesActual := `SELECT COALESCE(SUM(cantidad), 0) FROM pasajes_cantidad 
                             WHERE id_reserva = $1 AND eliminado = FALSE`
		err = tx.QueryRow(queryPasajesActual, id).Scan(&totalPasajerosActual)
		if err != nil {
			return err
		}

		// Contar pasajeros de paquetes
		queryPaquetesActual := `SELECT COALESCE(SUM(ppd.cantidad * pp.cantidad_total), 0)
                              FROM paquete_pasaje_detalle ppd
                              INNER JOIN paquete_pasajes pp ON ppd.id_paquete = pp.id_paquete
                              WHERE ppd.id_reserva = $1 AND ppd.eliminado = FALSE`
		var pasajerosPaquetesActual int
		err = tx.QueryRow(queryPaquetesActual, id).Scan(&pasajerosPaquetesActual)
		if err != nil {
			return err
		}

		totalPasajerosActual += pasajerosPaquetesActual
	}

	// Calcular la cantidad nueva de pasajeros
	totalPasajerosNuevo := 0

	// Sumar pasajeros de pasajes individuales
	for _, pasaje := range reserva.CantidadPasajes {
		totalPasajerosNuevo += pasaje.Cantidad
	}

	// Sumar pasajeros de paquetes
	for _, paquete := range reserva.Paquetes {
		// Obtener cantidad total de pasajeros por paquete
		var cantidadPorPaquete int
		queryPaquete := `SELECT cantidad_total FROM paquete_pasajes 
                        WHERE id_paquete = $1 AND eliminado = FALSE`
		err := tx.QueryRow(queryPaquete, paquete.IDPaquete).Scan(&cantidadPorPaquete)
		if err != nil {
			return err
		}

		// Multiplicar por la cantidad de paquetes seleccionados
		totalPasajerosNuevo += cantidadPorPaquete * paquete.Cantidad
	}

	// Verificar disponibilidad si cambia la cantidad o la instancia
	if idInstanciaActual != reserva.IDInstancia || totalPasajerosActual != totalPasajerosNuevo {
		// Si es la misma instancia, verificar solo la diferencia
		if idInstanciaActual == reserva.IDInstancia {
			var cupoDisponible int
			queryInstancia := `SELECT cupo_disponible FROM instancia_tour 
                            WHERE id_instancia = $1 AND eliminado = FALSE AND estado = 'PROGRAMADO'`
			err := tx.QueryRow(queryInstancia, reserva.IDInstancia).Scan(&cupoDisponible)
			if err != nil {
				if err == sql.ErrNoRows {
					return errors.New("la instancia del tour no existe, está eliminada o no está programada")
				}
				return err
			}

			// Calcular diferencia de pasajeros
			diferenciaPasajeros := totalPasajerosNuevo - totalPasajerosActual

			// Verificar si hay suficiente cupo
			if diferenciaPasajeros > cupoDisponible {
				return errors.New("no hay suficiente cupo disponible para la actualización de la reserva")
			}

			// Actualizar cupo si hay diferencia
			if diferenciaPasajeros != 0 {
				queryUpdateCupo := `UPDATE instancia_tour 
                               SET cupo_disponible = cupo_disponible - $1 
                               WHERE id_instancia = $2`
				_, err = tx.Exec(queryUpdateCupo, diferenciaPasajeros, reserva.IDInstancia)
				if err != nil {
					return err
				}
			}
		} else {
			// Si es una instancia diferente, restaurar cupo en la instancia anterior y verificar en la nueva
			if estadoActual != "CANCELADA" { // Solo restaurar si no estaba cancelada
				queryRestauraCupo := `UPDATE instancia_tour 
                                SET cupo_disponible = cupo_disponible + $1 
                                WHERE id_instancia = $2`
				_, err = tx.Exec(queryRestauraCupo, totalPasajerosActual, idInstanciaActual)
				if err != nil {
					return err
				}
			}

			// Verificar cupo en la nueva instancia
			var cupoDisponible int
			queryInstancia := `SELECT cupo_disponible FROM instancia_tour 
                            WHERE id_instancia = $1 AND eliminado = FALSE AND estado = 'PROGRAMADO'`
			err := tx.QueryRow(queryInstancia, reserva.IDInstancia).Scan(&cupoDisponible)
			if err != nil {
				if err == sql.ErrNoRows {
					return errors.New("la nueva instancia del tour no existe, está eliminada o no está programada")
				}
				return err
			}

			if totalPasajerosNuevo > cupoDisponible {
				return errors.New("no hay suficiente cupo disponible en la nueva instancia seleccionada")
			}

			// Actualizar cupo en la nueva instancia
			queryUpdateCupo := `UPDATE instancia_tour 
                           SET cupo_disponible = cupo_disponible - $1 
                           WHERE id_instancia = $2`
			_, err = tx.Exec(queryUpdateCupo, totalPasajerosNuevo, reserva.IDInstancia)
			if err != nil {
				return err
			}
		}
	}

	// Verificar cambio de estado
	if estadoActual != reserva.Estado {
		// Si se cancela la reserva, restaurar cupo
		if reserva.Estado == "CANCELADA" && estadoActual != "CANCELADA" {
			queryRestauraCupo := `UPDATE instancia_tour 
                               SET cupo_disponible = cupo_disponible + $1 
                               WHERE id_instancia = $2`
			_, err = tx.Exec(queryRestauraCupo, totalPasajerosActual, idInstanciaActual)
			if err != nil {
				return err
			}
		}

		// Si se reactiva una reserva cancelada, verificar cupo
		if estadoActual == "CANCELADA" && reserva.Estado != "CANCELADA" {
			var cupoDisponible int
			queryInstancia := `SELECT cupo_disponible FROM instancia_tour 
                            WHERE id_instancia = $1 AND eliminado = FALSE AND estado = 'PROGRAMADO'`
			err := tx.QueryRow(queryInstancia, reserva.IDInstancia).Scan(&cupoDisponible)
			if err != nil {
				if err == sql.ErrNoRows {
					return errors.New("la instancia del tour no existe, está eliminada o no está programada")
				}
				return err
			}

			if totalPasajerosNuevo > cupoDisponible {
				return errors.New("no hay suficiente cupo disponible para reactivar la reserva")
			}

			// Actualizar cupo al reactivar
			queryUpdateCupo := `UPDATE instancia_tour 
                           SET cupo_disponible = cupo_disponible - $1 
                           WHERE id_instancia = $2`
			_, err = tx.Exec(queryUpdateCupo, totalPasajerosNuevo, reserva.IDInstancia)
			if err != nil {
				return err
			}
		}
	}

	// Actualizar la reserva con los nuevos datos
	query := `UPDATE reserva SET
              id_vendedor = $1,
              id_cliente = $2,
              id_instancia = $3,
              id_canal = $4,
              id_sede = $5,
              total_pagar = $6,
              notas = $7,
              estado = $8
              WHERE id_reserva = $9 AND eliminado = FALSE`

	// Ejecutar la actualización
	_, err = tx.Exec(
		query,
		reserva.IDVendedor,
		reserva.IDCliente,
		reserva.IDInstancia,
		reserva.IDCanal,
		reserva.IDSede,
		reserva.TotalPagar,
		reserva.Notas,
		reserva.Estado,
		id,
	)

	if err != nil {
		return err
	}

	// Marcar como eliminados los registros de pasajes actuales
	queryDeletePasajes := `UPDATE pasajes_cantidad SET eliminado = TRUE WHERE id_reserva = $1`
	_, err = tx.Exec(queryDeletePasajes, id)
	if err != nil {
		return err
	}

	// Marcar como eliminados los registros de paquetes actuales
	queryDeletePaquetes := `UPDATE paquete_pasaje_detalle SET eliminado = TRUE WHERE id_reserva = $1`
	_, err = tx.Exec(queryDeletePaquetes, id)
	if err != nil {
		return err
	}

	// Insertar nuevas cantidades de pasajes
	for _, pasaje := range reserva.CantidadPasajes {
		// Solo insertar si la cantidad es mayor que cero
		if pasaje.Cantidad > 0 {
			queryPasaje := `INSERT INTO pasajes_cantidad (id_reserva, id_tipo_pasaje, cantidad, eliminado)
                         VALUES ($1, $2, $3, FALSE)`

			_, err = tx.Exec(queryPasaje, id, pasaje.IDTipoPasaje, pasaje.Cantidad)
			if err != nil {
				return err
			}
		}
	}

	// Insertar nuevos paquetes
	for _, paquete := range reserva.Paquetes {
		// Solo insertar si la cantidad es mayor que cero
		if paquete.Cantidad > 0 {
			queryPaquete := `INSERT INTO paquete_pasaje_detalle (id_reserva, id_paquete, cantidad, eliminado)
                          VALUES ($1, $2, $3, FALSE)`

			_, err = tx.Exec(queryPaquete, id, paquete.IDPaquete, paquete.Cantidad)
			if err != nil {
				return err
			}
		}
	}

	// Commit de la transacción
	return tx.Commit()
}

// UpdateEstado actualiza solo el estado de una reserva
func (r *ReservaRepository) UpdateEstado(id int, estado string) error {
	// Iniciar transacción
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// Si hay error, hacer rollback
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Obtener la reserva actual
	var idInstancia int
	var estadoActual string
	queryReservaActual := `SELECT id_instancia, estado FROM reserva 
                          WHERE id_reserva = $1 AND eliminado = FALSE`

	err = tx.QueryRow(queryReservaActual, id).Scan(&idInstancia, &estadoActual)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("reserva no encontrada")
		}
		return err
	}

	// Verificar cambio de estado para manejo de cupos
	if estadoActual != estado {
		// Si se cancela una reserva activa, restaurar cupo
		if estado == "CANCELADA" && estadoActual != "CANCELADA" {
			// Obtener total de pasajeros
			totalPasajeros, err := r.GetCantidadPasajerosByReservaTx(tx, id)
			if err != nil {
				return err
			}

			// Restaurar cupo
			queryRestauraCupo := `UPDATE instancia_tour 
                               SET cupo_disponible = cupo_disponible + $1 
                               WHERE id_instancia = $2`
			_, err = tx.Exec(queryRestauraCupo, totalPasajeros, idInstancia)
			if err != nil {
				return err
			}
		}

		// Si se reactiva una reserva cancelada, verificar y reducir cupo
		if estadoActual == "CANCELADA" && estado != "CANCELADA" {
			// Obtener total de pasajeros
			totalPasajeros, err := r.GetCantidadPasajerosByReservaTx(tx, id)
			if err != nil {
				return err
			}

			// Verificar cupo disponible
			var cupoDisponible int
			queryInstancia := `SELECT cupo_disponible FROM instancia_tour 
                            WHERE id_instancia = $1 AND eliminado = FALSE AND estado = 'PROGRAMADO'`
			err = tx.QueryRow(queryInstancia, idInstancia).Scan(&cupoDisponible)
			if err != nil {
				if err == sql.ErrNoRows {
					return errors.New("la instancia del tour no existe, está eliminada o no está programada")
				}
				return err
			}

			if totalPasajeros > cupoDisponible {
				return errors.New("no hay suficiente cupo disponible para reactivar la reserva")
			}

			// Reducir cupo
			queryUpdateCupo := `UPDATE instancia_tour 
                           SET cupo_disponible = cupo_disponible - $1 
                           WHERE id_instancia = $2`
			_, err = tx.Exec(queryUpdateCupo, totalPasajeros, idInstancia)
			if err != nil {
				return err
			}
		}
	}

	// Actualizar estado
	query := `UPDATE reserva SET estado = $1 WHERE id_reserva = $2 AND eliminado = FALSE`
	_, err = tx.Exec(query, estado, id)
	if err != nil {
		return err
	}

	// Commit de la transacción
	return tx.Commit()
}

// GetCantidadPasajerosByReservaTx obtiene la cantidad total de pasajeros dentro de una transacción
func (r *ReservaRepository) GetCantidadPasajerosByReservaTx(tx *sql.Tx, id int) (int, error) {
	var totalPasajerosIndividuales int
	queryPasajes := `SELECT COALESCE(SUM(cantidad), 0) FROM pasajes_cantidad 
                   WHERE id_reserva = $1 AND eliminado = FALSE`
	err := tx.QueryRow(queryPasajes, id).Scan(&totalPasajerosIndividuales)
	if err != nil {
		return 0, err
	}

	var totalPasajerosPaquetes int
	queryPaquetes := `SELECT COALESCE(SUM(ppd.cantidad * pp.cantidad_total), 0)
                    FROM paquete_pasaje_detalle ppd
                    INNER JOIN paquete_pasajes pp ON ppd.id_paquete = pp.id_paquete
                    WHERE ppd.id_reserva = $1 AND ppd.eliminado = FALSE`
	err = tx.QueryRow(queryPaquetes, id).Scan(&totalPasajerosPaquetes)
	if err != nil {
		return 0, err
	}

	return totalPasajerosIndividuales + totalPasajerosPaquetes, nil
}

// GetCantidadPasajerosByReserva obtiene la cantidad total de pasajeros en una reserva
func (r *ReservaRepository) GetCantidadPasajerosByReserva(id int) (int, error) {
	var totalPasajerosIndividuales int
	queryPasajes := `SELECT COALESCE(SUM(cantidad), 0) FROM pasajes_cantidad 
                   WHERE id_reserva = $1 AND eliminado = FALSE`
	err := r.db.QueryRow(queryPasajes, id).Scan(&totalPasajerosIndividuales)
	if err != nil {
		return 0, err
	}

	var totalPasajerosPaquetes int
	queryPaquetes := `SELECT COALESCE(SUM(ppd.cantidad * pp.cantidad_total), 0)
                    FROM paquete_pasaje_detalle ppd
                    INNER JOIN paquete_pasajes pp ON ppd.id_paquete = pp.id_paquete
                    WHERE ppd.id_reserva = $1 AND ppd.eliminado = FALSE`
	err = r.db.QueryRow(queryPaquetes, id).Scan(&totalPasajerosPaquetes)
	if err != nil {
		return 0, err
	}

	return totalPasajerosIndividuales + totalPasajerosPaquetes, nil
}

// Delete realiza una eliminación lógica de una reserva
func (r *ReservaRepository) Delete(id int) error {
	// Iniciar transacción
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// Si hay error, hacer rollback
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Verificar si hay pagos asociados a esta reserva (que no estén eliminados)
	var countPagos int
	queryCheckPagos := `SELECT COUNT(*) FROM pago WHERE id_reserva = $1 AND eliminado = FALSE`
	err = tx.QueryRow(queryCheckPagos, id).Scan(&countPagos)
	if err != nil {
		return err
	}

	if countPagos > 0 {
		return errors.New("no se puede eliminar esta reserva porque tiene pagos asociados")
	}

	// Verificar si hay comprobantes asociados a esta reserva (que no estén eliminados)
	var countComprobantes int
	queryCheckComprobantes := `SELECT COUNT(*) FROM comprobante_pago WHERE id_reserva = $1 AND eliminado = FALSE`
	err = tx.QueryRow(queryCheckComprobantes, id).Scan(&countComprobantes)
	if err != nil {
		return err
	}

	if countComprobantes > 0 {
		return errors.New("no se puede eliminar esta reserva porque tiene comprobantes asociados")
	}

	// Obtener información de la reserva para restaurar cupo
	var idInstancia int
	var estado string
	queryReserva := `SELECT id_instancia, estado FROM reserva WHERE id_reserva = $1 AND eliminado = FALSE`
	err = tx.QueryRow(queryReserva, id).Scan(&idInstancia, &estado)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("reserva no encontrada")
		}
		return err
	}

	// Si la reserva no está cancelada, restaurar cupo
	if estado != "CANCELADA" {
		// Obtener total de pasajeros
		totalPasajeros, err := r.GetCantidadPasajerosByReservaTx(tx, id)
		if err != nil {
			return err
		}

		// Restaurar cupo
		queryRestauraCupo := `UPDATE instancia_tour 
                          SET cupo_disponible = cupo_disponible + $1 
                          WHERE id_instancia = $2`
		_, err = tx.Exec(queryRestauraCupo, totalPasajeros, idInstancia)
		if err != nil {
			return err
		}
	}

	// Marcar los registros de pasajes como eliminados (eliminación lógica)
	queryDeletePasajes := `UPDATE pasajes_cantidad SET eliminado = TRUE WHERE id_reserva = $1`
	_, err = tx.Exec(queryDeletePasajes, id)
	if err != nil {
		return err
	}

	// Marcar los registros de paquetes como eliminados (eliminación lógica)
	queryDeletePaquetes := `UPDATE paquete_pasaje_detalle SET eliminado = TRUE WHERE id_reserva = $1`
	_, err = tx.Exec(queryDeletePaquetes, id)
	if err != nil {
		return err
	}

	// Marcar la reserva como eliminada (eliminación lógica)
	queryDeleteReserva := `UPDATE reserva SET eliminado = TRUE WHERE id_reserva = $1`
	_, err = tx.Exec(queryDeleteReserva, id)
	if err != nil {
		return err
	}

	// Commit de la transacción
	return tx.Commit()
}

// List obtiene todas las reservas activas del sistema
func (r *ReservaRepository) List() ([]*entidades.Reserva, error) {
	query := `SELECT r.id_reserva, r.id_vendedor, r.id_cliente, r.id_instancia, 
              r.id_canal, r.id_sede, r.fecha_reserva, r.total_pagar, r.notas, r.estado, r.eliminado,
              c.nombres || ' ' || c.apellidos as nombre_cliente,
              COALESCE(u.nombres || ' ' || u.apellidos, 'Web') as nombre_vendedor,
              tt.nombre as nombre_tour,
              to_char(it.fecha_especifica, 'DD/MM/YYYY') as fecha_tour,
              to_char(it.hora_inicio, 'HH24:MI') as hora_inicio_tour,
              to_char(it.hora_fin, 'HH24:MI') as hora_fin_tour,
              cv.nombre as nombre_canal,
              s.nombre as nombre_sede
              FROM reserva r
              INNER JOIN cliente c ON r.id_cliente = c.id_cliente
              LEFT JOIN usuario u ON r.id_vendedor = u.id_usuario
              INNER JOIN instancia_tour it ON r.id_instancia = it.id_instancia
              INNER JOIN tour_programado tp ON it.id_tour_programado = tp.id_tour_programado
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              INNER JOIN canal_venta cv ON r.id_canal = cv.id_canal
              INNER JOIN sede s ON r.id_sede = s.id_sede
              WHERE r.eliminado = FALSE
              ORDER BY r.fecha_reserva DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reservas := []*entidades.Reserva{}

	// Iterar por cada reserva encontrada
	for rows.Next() {
		reserva := &entidades.Reserva{}
		err := rows.Scan(
			&reserva.ID, &reserva.IDVendedor, &reserva.IDCliente, &reserva.IDInstancia,
			&reserva.IDCanal, &reserva.IDSede, &reserva.FechaReserva, &reserva.TotalPagar,
			&reserva.Notas, &reserva.Estado, &reserva.Eliminado,
			&reserva.NombreCliente, &reserva.NombreVendedor, &reserva.NombreTour,
			&reserva.FechaTour, &reserva.HoraInicioTour, &reserva.HoraFinTour,
			&reserva.NombreCanal, &reserva.NombreSede,
		)
		if err != nil {
			return nil, err
		}

		// Obtener las cantidades de pasajes para cada reserva
		queryPasajes := `SELECT pc.id_tipo_pasaje, tp.nombre, pc.cantidad
                       FROM pasajes_cantidad pc
                       INNER JOIN tipo_pasaje tp ON pc.id_tipo_pasaje = tp.id_tipo_pasaje
                       WHERE pc.id_reserva = $1 AND pc.eliminado = FALSE`

		rowsPasajes, err := r.db.Query(queryPasajes, reserva.ID)
		if err != nil {
			return nil, err
		}

		reserva.CantidadPasajes = []entidades.PasajeCantidad{}

		// Iterar por cada registro de pasaje
		for rowsPasajes.Next() {
			var pasajeCantidad entidades.PasajeCantidad
			err := rowsPasajes.Scan(
				&pasajeCantidad.IDTipoPasaje, &pasajeCantidad.NombreTipo, &pasajeCantidad.Cantidad,
			)
			if err != nil {
				rowsPasajes.Close()
				return nil, err
			}
			reserva.CantidadPasajes = append(reserva.CantidadPasajes, pasajeCantidad)
		}

		rowsPasajes.Close()
		if err = rowsPasajes.Err(); err != nil {
			return nil, err
		}

		// Obtener los paquetes de pasajes
		queryPaquetes := `SELECT ppd.id_paquete, pp.nombre as nombre_paquete, ppd.cantidad, 
                        pp.precio_total as precio_unitario, (ppd.cantidad * pp.precio_total) as subtotal,
                        pp.cantidad_total
                        FROM paquete_pasaje_detalle ppd
                        INNER JOIN paquete_pasajes pp ON ppd.id_paquete = pp.id_paquete
                        WHERE ppd.id_reserva = $1 AND ppd.eliminado = FALSE`

		rowsPaquetes, err := r.db.Query(queryPaquetes, reserva.ID)
		if err != nil {
			return nil, err
		}

		reserva.Paquetes = []entidades.PaquetePasajeDetalle{}

		// Iterar por cada registro de paquete
		for rowsPaquetes.Next() {
			var paquete entidades.PaquetePasajeDetalle
			err := rowsPaquetes.Scan(
				&paquete.IDPaquete, &paquete.NombrePaquete, &paquete.Cantidad,
				&paquete.PrecioUnitario, &paquete.Subtotal, &paquete.CantidadTotal,
			)
			if err != nil {
				rowsPaquetes.Close()
				return nil, err
			}
			reserva.Paquetes = append(reserva.Paquetes, paquete)
		}

		rowsPaquetes.Close()
		if err = rowsPaquetes.Err(); err != nil {
			return nil, err
		}

		reservas = append(reservas, reserva)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reservas, nil
}

// ListByCliente lista todas las reservas activas de un cliente específico
func (r *ReservaRepository) ListByCliente(idCliente int) ([]*entidades.Reserva, error) {
	query := `SELECT r.id_reserva, r.id_vendedor, r.id_cliente, r.id_instancia, 
              r.id_canal, r.id_sede, r.fecha_reserva, r.total_pagar, r.notas, r.estado, r.eliminado,
              c.nombres || ' ' || c.apellidos as nombre_cliente,
              COALESCE(u.nombres || ' ' || u.apellidos, 'Web') as nombre_vendedor,
              tt.nombre as nombre_tour,
              to_char(it.fecha_especifica, 'DD/MM/YYYY') as fecha_tour,
              to_char(it.hora_inicio, 'HH24:MI') as hora_inicio_tour,
              to_char(it.hora_fin, 'HH24:MI') as hora_fin_tour,
              cv.nombre as nombre_canal,
              s.nombre as nombre_sede
              FROM reserva r
              INNER JOIN cliente c ON r.id_cliente = c.id_cliente
              LEFT JOIN usuario u ON r.id_vendedor = u.id_usuario
              INNER JOIN instancia_tour it ON r.id_instancia = it.id_instancia
              INNER JOIN tour_programado tp ON it.id_tour_programado = tp.id_tour_programado
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              INNER JOIN canal_venta cv ON r.id_canal = cv.id_canal
              INNER JOIN sede s ON r.id_sede = s.id_sede
              WHERE r.id_cliente = $1 AND r.eliminado = FALSE
              ORDER BY r.fecha_reserva DESC`

	rows, err := r.db.Query(query, idCliente)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reservas := []*entidades.Reserva{}

	// Iterar por cada reserva encontrada
	for rows.Next() {
		reserva := &entidades.Reserva{}
		err := rows.Scan(
			&reserva.ID, &reserva.IDVendedor, &reserva.IDCliente, &reserva.IDInstancia,
			&reserva.IDCanal, &reserva.IDSede, &reserva.FechaReserva, &reserva.TotalPagar,
			&reserva.Notas, &reserva.Estado, &reserva.Eliminado,
			&reserva.NombreCliente, &reserva.NombreVendedor, &reserva.NombreTour,
			&reserva.FechaTour, &reserva.HoraInicioTour, &reserva.HoraFinTour,
			&reserva.NombreCanal, &reserva.NombreSede,
		)
		if err != nil {
			return nil, err
		}

		// Obtener las cantidades de pasajes para cada reserva
		queryPasajes := `SELECT pc.id_tipo_pasaje, tp.nombre, pc.cantidad
                       FROM pasajes_cantidad pc
                       INNER JOIN tipo_pasaje tp ON pc.id_tipo_pasaje = tp.id_tipo_pasaje
                       WHERE pc.id_reserva = $1 AND pc.eliminado = FALSE`

		rowsPasajes, err := r.db.Query(queryPasajes, reserva.ID)
		if err != nil {
			return nil, err
		}

		reserva.CantidadPasajes = []entidades.PasajeCantidad{}

		// Iterar por cada registro de pasaje
		for rowsPasajes.Next() {
			var pasajeCantidad entidades.PasajeCantidad
			err := rowsPasajes.Scan(
				&pasajeCantidad.IDTipoPasaje, &pasajeCantidad.NombreTipo, &pasajeCantidad.Cantidad,
			)
			if err != nil {
				rowsPasajes.Close()
				return nil, err
			}
			reserva.CantidadPasajes = append(reserva.CantidadPasajes, pasajeCantidad)
		}

		rowsPasajes.Close()
		if err = rowsPasajes.Err(); err != nil {
			return nil, err
		}

		// Obtener los paquetes de pasajes
		queryPaquetes := `SELECT ppd.id_paquete, pp.nombre as nombre_paquete, ppd.cantidad, 
                        pp.precio_total as precio_unitario, (ppd.cantidad * pp.precio_total) as subtotal,
                        pp.cantidad_total
                        FROM paquete_pasaje_detalle ppd
                        INNER JOIN paquete_pasajes pp ON ppd.id_paquete = pp.id_paquete
                        WHERE ppd.id_reserva = $1 AND ppd.eliminado = FALSE`

		rowsPaquetes, err := r.db.Query(queryPaquetes, reserva.ID)
		if err != nil {
			return nil, err
		}

		reserva.Paquetes = []entidades.PaquetePasajeDetalle{}

		// Iterar por cada registro de paquete
		for rowsPaquetes.Next() {
			var paquete entidades.PaquetePasajeDetalle
			err := rowsPaquetes.Scan(
				&paquete.IDPaquete, &paquete.NombrePaquete, &paquete.Cantidad,
				&paquete.PrecioUnitario, &paquete.Subtotal, &paquete.CantidadTotal,
			)
			if err != nil {
				rowsPaquetes.Close()
				return nil, err
			}
			reserva.Paquetes = append(reserva.Paquetes, paquete)
		}

		rowsPaquetes.Close()
		if err = rowsPaquetes.Err(); err != nil {
			return nil, err
		}

		reservas = append(reservas, reserva)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reservas, nil
}

// ListByInstancia lista todas las reservas asociadas a una instancia específica de tour
func (r *ReservaRepository) ListByInstancia(idInstancia int) ([]*entidades.Reserva, error) {
	query := `SELECT r.id_reserva, r.id_vendedor, r.id_cliente, r.id_instancia, 
              r.id_canal, r.id_sede, r.fecha_reserva, r.total_pagar, r.notas, r.estado, r.eliminado,
              c.nombres || ' ' || c.apellidos as nombre_cliente,
              COALESCE(u.nombres || ' ' || u.apellidos, 'Web') as nombre_vendedor,
              tt.nombre as nombre_tour,
              to_char(it.fecha_especifica, 'DD/MM/YYYY') as fecha_tour,
              to_char(it.hora_inicio, 'HH24:MI') as hora_inicio_tour,
              to_char(it.hora_fin, 'HH24:MI') as hora_fin_tour,
              cv.nombre as nombre_canal,
              s.nombre as nombre_sede
              FROM reserva r
              INNER JOIN cliente c ON r.id_cliente = c.id_cliente
              LEFT JOIN usuario u ON r.id_vendedor = u.id_usuario
              INNER JOIN instancia_tour it ON r.id_instancia = it.id_instancia
              INNER JOIN tour_programado tp ON it.id_tour_programado = tp.id_tour_programado
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              INNER JOIN canal_venta cv ON r.id_canal = cv.id_canal
              INNER JOIN sede s ON r.id_sede = s.id_sede
              WHERE r.id_instancia = $1 AND r.eliminado = FALSE
              ORDER BY r.fecha_reserva DESC`

	rows, err := r.db.Query(query, idInstancia)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reservas := []*entidades.Reserva{}

	// Iterar por cada reserva encontrada
	for rows.Next() {
		reserva := &entidades.Reserva{}
		err := rows.Scan(
			&reserva.ID, &reserva.IDVendedor, &reserva.IDCliente, &reserva.IDInstancia,
			&reserva.IDCanal, &reserva.IDSede, &reserva.FechaReserva, &reserva.TotalPagar,
			&reserva.Notas, &reserva.Estado, &reserva.Eliminado,
			&reserva.NombreCliente, &reserva.NombreVendedor, &reserva.NombreTour,
			&reserva.FechaTour, &reserva.HoraInicioTour, &reserva.HoraFinTour,
			&reserva.NombreCanal, &reserva.NombreSede,
		)
		if err != nil {
			return nil, err
		}

		// Obtener las cantidades de pasajes para cada reserva
		queryPasajes := `SELECT pc.id_tipo_pasaje, tp.nombre, pc.cantidad
                       FROM pasajes_cantidad pc
                       INNER JOIN tipo_pasaje tp ON pc.id_tipo_pasaje = tp.id_tipo_pasaje
                       WHERE pc.id_reserva = $1 AND pc.eliminado = FALSE`

		rowsPasajes, err := r.db.Query(queryPasajes, reserva.ID)
		if err != nil {
			return nil, err
		}

		reserva.CantidadPasajes = []entidades.PasajeCantidad{}

		// Iterar por cada registro de pasaje
		for rowsPasajes.Next() {
			var pasajeCantidad entidades.PasajeCantidad
			err := rowsPasajes.Scan(
				&pasajeCantidad.IDTipoPasaje, &pasajeCantidad.NombreTipo, &pasajeCantidad.Cantidad,
			)
			if err != nil {
				rowsPasajes.Close()
				return nil, err
			}
			reserva.CantidadPasajes = append(reserva.CantidadPasajes, pasajeCantidad)
		}

		rowsPasajes.Close()
		if err = rowsPasajes.Err(); err != nil {
			return nil, err
		}

		// Obtener los paquetes de pasajes
		queryPaquetes := `SELECT ppd.id_paquete, pp.nombre as nombre_paquete, ppd.cantidad, 
                        pp.precio_total as precio_unitario, (ppd.cantidad * pp.precio_total) as subtotal,
                        pp.cantidad_total
                        FROM paquete_pasaje_detalle ppd
                        INNER JOIN paquete_pasajes pp ON ppd.id_paquete = pp.id_paquete
                        WHERE ppd.id_reserva = $1 AND ppd.eliminado = FALSE`

		rowsPaquetes, err := r.db.Query(queryPaquetes, reserva.ID)
		if err != nil {
			return nil, err
		}

		reserva.Paquetes = []entidades.PaquetePasajeDetalle{}

		// Iterar por cada registro de paquete
		for rowsPaquetes.Next() {
			var paquete entidades.PaquetePasajeDetalle
			err := rowsPaquetes.Scan(
				&paquete.IDPaquete, &paquete.NombrePaquete, &paquete.Cantidad,
				&paquete.PrecioUnitario, &paquete.Subtotal, &paquete.CantidadTotal,
			)
			if err != nil {
				rowsPaquetes.Close()
				return nil, err
			}
			reserva.Paquetes = append(reserva.Paquetes, paquete)
		}

		rowsPaquetes.Close()
		if err = rowsPaquetes.Err(); err != nil {
			return nil, err
		}

		reservas = append(reservas, reserva)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reservas, nil
}

// ListByFecha lista todas las reservas para una fecha específica de instancia
func (r *ReservaRepository) ListByFecha(fecha time.Time) ([]*entidades.Reserva, error) {
	query := `SELECT r.id_reserva, r.id_vendedor, r.id_cliente, r.id_instancia, 
              r.id_canal, r.id_sede, r.fecha_reserva, r.total_pagar, r.notas, r.estado, r.eliminado,
              c.nombres || ' ' || c.apellidos as nombre_cliente,
              COALESCE(u.nombres || ' ' || u.apellidos, 'Web') as nombre_vendedor,
              tt.nombre as nombre_tour,
              to_char(it.fecha_especifica, 'DD/MM/YYYY') as fecha_tour,
              to_char(it.hora_inicio, 'HH24:MI') as hora_inicio_tour,
              to_char(it.hora_fin, 'HH24:MI') as hora_fin_tour,
              cv.nombre as nombre_canal,
              s.nombre as nombre_sede
              FROM reserva r
              INNER JOIN cliente c ON r.id_cliente = c.id_cliente
              LEFT JOIN usuario u ON r.id_vendedor = u.id_usuario
              INNER JOIN instancia_tour it ON r.id_instancia = it.id_instancia
              INNER JOIN tour_programado tp ON it.id_tour_programado = tp.id_tour_programado
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              INNER JOIN canal_venta cv ON r.id_canal = cv.id_canal
              INNER JOIN sede s ON r.id_sede = s.id_sede
              WHERE it.fecha_especifica = $1 AND r.eliminado = FALSE
              ORDER BY it.hora_inicio ASC, r.fecha_reserva DESC`

	rows, err := r.db.Query(query, fecha)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reservas := []*entidades.Reserva{}

	// Iterar por cada reserva encontrada
	for rows.Next() {
		reserva := &entidades.Reserva{}
		err := rows.Scan(
			&reserva.ID, &reserva.IDVendedor, &reserva.IDCliente, &reserva.IDInstancia,
			&reserva.IDCanal, &reserva.IDSede, &reserva.FechaReserva, &reserva.TotalPagar,
			&reserva.Notas, &reserva.Estado, &reserva.Eliminado,
			&reserva.NombreCliente, &reserva.NombreVendedor, &reserva.NombreTour,
			&reserva.FechaTour, &reserva.HoraInicioTour, &reserva.HoraFinTour,
			&reserva.NombreCanal, &reserva.NombreSede,
		)
		if err != nil {
			return nil, err
		}

		// Obtener las cantidades de pasajes para cada reserva
		queryPasajes := `SELECT pc.id_tipo_pasaje, tp.nombre, pc.cantidad
                       FROM pasajes_cantidad pc
                       INNER JOIN tipo_pasaje tp ON pc.id_tipo_pasaje = tp.id_tipo_pasaje
                       WHERE pc.id_reserva = $1 AND pc.eliminado = FALSE`

		rowsPasajes, err := r.db.Query(queryPasajes, reserva.ID)
		if err != nil {
			return nil, err
		}

		reserva.CantidadPasajes = []entidades.PasajeCantidad{}

		// Iterar por cada registro de pasaje
		for rowsPasajes.Next() {
			var pasajeCantidad entidades.PasajeCantidad
			err := rowsPasajes.Scan(
				&pasajeCantidad.IDTipoPasaje, &pasajeCantidad.NombreTipo, &pasajeCantidad.Cantidad,
			)
			if err != nil {
				rowsPasajes.Close()
				return nil, err
			}
			reserva.CantidadPasajes = append(reserva.CantidadPasajes, pasajeCantidad)
		}

		rowsPasajes.Close()
		if err = rowsPasajes.Err(); err != nil {
			return nil, err
		}

		// Obtener los paquetes de pasajes
		queryPaquetes := `SELECT ppd.id_paquete, pp.nombre as nombre_paquete, ppd.cantidad, 
                        pp.precio_total as precio_unitario, (ppd.cantidad * pp.precio_total) as subtotal,
                        pp.cantidad_total
                        FROM paquete_pasaje_detalle ppd
                        INNER JOIN paquete_pasajes pp ON ppd.id_paquete = pp.id_paquete
                        WHERE ppd.id_reserva = $1 AND ppd.eliminado = FALSE`

		rowsPaquetes, err := r.db.Query(queryPaquetes, reserva.ID)
		if err != nil {
			return nil, err
		}

		reserva.Paquetes = []entidades.PaquetePasajeDetalle{}

		// Iterar por cada registro de paquete
		for rowsPaquetes.Next() {
			var paquete entidades.PaquetePasajeDetalle
			err := rowsPaquetes.Scan(
				&paquete.IDPaquete, &paquete.NombrePaquete, &paquete.Cantidad,
				&paquete.PrecioUnitario, &paquete.Subtotal, &paquete.CantidadTotal,
			)
			if err != nil {
				rowsPaquetes.Close()
				return nil, err
			}
			reserva.Paquetes = append(reserva.Paquetes, paquete)
		}

		rowsPaquetes.Close()
		if err = rowsPaquetes.Err(); err != nil {
			return nil, err
		}

		reservas = append(reservas, reserva)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reservas, nil
}

// ListByEstado lista todas las reservas por estado específico (RESERVADO, CANCELADA, CONFIRMADA)
func (r *ReservaRepository) ListByEstado(estado string) ([]*entidades.Reserva, error) {
	query := `SELECT r.id_reserva, r.id_vendedor, r.id_cliente, r.id_instancia, 
              r.id_canal, r.id_sede, r.fecha_reserva, r.total_pagar, r.notas, r.estado, r.eliminado,
              c.nombres || ' ' || c.apellidos as nombre_cliente,
              COALESCE(u.nombres || ' ' || u.apellidos, 'Web') as nombre_vendedor,
              tt.nombre as nombre_tour,
              to_char(it.fecha_especifica, 'DD/MM/YYYY') as fecha_tour,
              to_char(it.hora_inicio, 'HH24:MI') as hora_inicio_tour,
              to_char(it.hora_fin, 'HH24:MI') as hora_fin_tour,
              cv.nombre as nombre_canal,
              s.nombre as nombre_sede
              FROM reserva r
              INNER JOIN cliente c ON r.id_cliente = c.id_cliente
              LEFT JOIN usuario u ON r.id_vendedor = u.id_usuario
              INNER JOIN instancia_tour it ON r.id_instancia = it.id_instancia
              INNER JOIN tour_programado tp ON it.id_tour_programado = tp.id_tour_programado
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              INNER JOIN canal_venta cv ON r.id_canal = cv.id_canal
              INNER JOIN sede s ON r.id_sede = s.id_sede
              WHERE r.estado = $1 AND r.eliminado = FALSE
              ORDER BY r.fecha_reserva DESC`

	rows, err := r.db.Query(query, estado)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reservas := []*entidades.Reserva{}

	// Iterar por cada reserva encontrada
	for rows.Next() {
		reserva := &entidades.Reserva{}
		err := rows.Scan(
			&reserva.ID, &reserva.IDVendedor, &reserva.IDCliente, &reserva.IDInstancia,
			&reserva.IDCanal, &reserva.IDSede, &reserva.FechaReserva, &reserva.TotalPagar,
			&reserva.Notas, &reserva.Estado, &reserva.Eliminado,
			&reserva.NombreCliente, &reserva.NombreVendedor, &reserva.NombreTour,
			&reserva.FechaTour, &reserva.HoraInicioTour, &reserva.HoraFinTour,
			&reserva.NombreCanal, &reserva.NombreSede,
		)
		if err != nil {
			return nil, err
		}

		// Obtener las cantidades de pasajes para cada reserva
		queryPasajes := `SELECT pc.id_tipo_pasaje, tp.nombre, pc.cantidad
                       FROM pasajes_cantidad pc
                       INNER JOIN tipo_pasaje tp ON pc.id_tipo_pasaje = tp.id_tipo_pasaje
                       WHERE pc.id_reserva = $1 AND pc.eliminado = FALSE`

		rowsPasajes, err := r.db.Query(queryPasajes, reserva.ID)
		if err != nil {
			return nil, err
		}

		reserva.CantidadPasajes = []entidades.PasajeCantidad{}

		// Iterar por cada registro de pasaje
		for rowsPasajes.Next() {
			var pasajeCantidad entidades.PasajeCantidad
			err := rowsPasajes.Scan(
				&pasajeCantidad.IDTipoPasaje, &pasajeCantidad.NombreTipo, &pasajeCantidad.Cantidad,
			)
			if err != nil {
				rowsPasajes.Close()
				return nil, err
			}
			reserva.CantidadPasajes = append(reserva.CantidadPasajes, pasajeCantidad)
		}

		rowsPasajes.Close()
		if err = rowsPasajes.Err(); err != nil {
			return nil, err
		}

		// Obtener los paquetes de pasajes
		queryPaquetes := `SELECT ppd.id_paquete, pp.nombre as nombre_paquete, ppd.cantidad, 
                        pp.precio_total as precio_unitario, (ppd.cantidad * pp.precio_total) as subtotal,
                        pp.cantidad_total
                        FROM paquete_pasaje_detalle ppd
                        INNER JOIN paquete_pasajes pp ON ppd.id_paquete = pp.id_paquete
                        WHERE ppd.id_reserva = $1 AND ppd.eliminado = FALSE`

		rowsPaquetes, err := r.db.Query(queryPaquetes, reserva.ID)
		if err != nil {
			return nil, err
		}

		reserva.Paquetes = []entidades.PaquetePasajeDetalle{}

		// Iterar por cada registro de paquete
		for rowsPaquetes.Next() {
			var paquete entidades.PaquetePasajeDetalle
			err := rowsPaquetes.Scan(
				&paquete.IDPaquete, &paquete.NombrePaquete, &paquete.Cantidad,
				&paquete.PrecioUnitario, &paquete.Subtotal, &paquete.CantidadTotal,
			)
			if err != nil {
				rowsPaquetes.Close()
				return nil, err
			}
			reserva.Paquetes = append(reserva.Paquetes, paquete)
		}

		rowsPaquetes.Close()
		if err = rowsPaquetes.Err(); err != nil {
			return nil, err
		}

		reservas = append(reservas, reserva)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reservas, nil
}

// ListBySede lista todas las reservas de una sede específica o todas las reservas si es ADMIN
func (r *ReservaRepository) ListBySede(idSede *int) ([]*entidades.Reserva, error) {
	query := `SELECT r.id_reserva, r.id_vendedor, r.id_cliente, r.id_instancia, 
              r.id_canal, r.id_sede, r.fecha_reserva, r.total_pagar, r.notas, r.estado, r.eliminado,
              c.nombres || ' ' || c.apellidos as nombre_cliente,
              COALESCE(u.nombres || ' ' || u.apellidos, 'Web') as nombre_vendedor,
              tt.nombre as nombre_tour,
              to_char(it.fecha_especifica, 'DD/MM/YYYY') as fecha_tour,
              to_char(it.hora_inicio, 'HH24:MI') as hora_inicio_tour,
              to_char(it.hora_fin, 'HH24:MI') as hora_fin_tour,
              cv.nombre as nombre_canal,
              s.nombre as nombre_sede
              FROM reserva r
              INNER JOIN cliente c ON r.id_cliente = c.id_cliente
              LEFT JOIN usuario u ON r.id_vendedor = u.id_usuario
              INNER JOIN instancia_tour it ON r.id_instancia = it.id_instancia
              INNER JOIN tour_programado tp ON it.id_tour_programado = tp.id_tour_programado
              INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
              INNER JOIN canal_venta cv ON r.id_canal = cv.id_canal
              INNER JOIN sede s ON r.id_sede = s.id_sede
              WHERE r.eliminado = FALSE`

	// Si se proporciona un ID de sede, filtrar por ella
	if idSede != nil {
		query += " AND r.id_sede = $1"
	}

	query += " ORDER BY r.fecha_reserva DESC"

	var rows *sql.Rows
	var err error

	if idSede != nil {
		rows, err = r.db.Query(query, *idSede)
	} else {
		rows, err = r.db.Query(query)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reservas := []*entidades.Reserva{}

	// Iterar por cada reserva encontrada
	for rows.Next() {
		reserva := &entidades.Reserva{}
		err := rows.Scan(
			&reserva.ID, &reserva.IDVendedor, &reserva.IDCliente, &reserva.IDInstancia,
			&reserva.IDCanal, &reserva.IDSede, &reserva.FechaReserva, &reserva.TotalPagar,
			&reserva.Notas, &reserva.Estado, &reserva.Eliminado,
			&reserva.NombreCliente, &reserva.NombreVendedor, &reserva.NombreTour,
			&reserva.FechaTour, &reserva.HoraInicioTour, &reserva.HoraFinTour,
			&reserva.NombreCanal, &reserva.NombreSede,
		)
		if err != nil {
			return nil, err
		}

		// Obtener las cantidades de pasajes para cada reserva
		queryPasajes := `SELECT pc.id_tipo_pasaje, tp.nombre, pc.cantidad
                       FROM pasajes_cantidad pc
                       INNER JOIN tipo_pasaje tp ON pc.id_tipo_pasaje = tp.id_tipo_pasaje
                       WHERE pc.id_reserva = $1 AND pc.eliminado = FALSE`

		rowsPasajes, err := r.db.Query(queryPasajes, reserva.ID)
		if err != nil {
			return nil, err
		}

		reserva.CantidadPasajes = []entidades.PasajeCantidad{}

		// Iterar por cada registro de pasaje
		for rowsPasajes.Next() {
			var pasajeCantidad entidades.PasajeCantidad
			err := rowsPasajes.Scan(
				&pasajeCantidad.IDTipoPasaje, &pasajeCantidad.NombreTipo, &pasajeCantidad.Cantidad,
			)
			if err != nil {
				rowsPasajes.Close()
				return nil, err
			}
			reserva.CantidadPasajes = append(reserva.CantidadPasajes, pasajeCantidad)
		}

		rowsPasajes.Close()
		if err = rowsPasajes.Err(); err != nil {
			return nil, err
		}

		// Obtener los paquetes de pasajes
		queryPaquetes := `SELECT ppd.id_paquete, pp.nombre as nombre_paquete, ppd.cantidad, 
                        pp.precio_total as precio_unitario, (ppd.cantidad * pp.precio_total) as subtotal,
                        pp.cantidad_total
                        FROM paquete_pasaje_detalle ppd
                        INNER JOIN paquete_pasajes pp ON ppd.id_paquete = pp.id_paquete
                        WHERE ppd.id_reserva = $1 AND ppd.eliminado = FALSE`

		rowsPaquetes, err := r.db.Query(queryPaquetes, reserva.ID)
		if err != nil {
			return nil, err
		}

		reserva.Paquetes = []entidades.PaquetePasajeDetalle{}

		// Iterar por cada registro de paquete
		for rowsPaquetes.Next() {
			var paquete entidades.PaquetePasajeDetalle
			err := rowsPaquetes.Scan(
				&paquete.IDPaquete, &paquete.NombrePaquete, &paquete.Cantidad,
				&paquete.PrecioUnitario, &paquete.Subtotal, &paquete.CantidadTotal,
			)
			if err != nil {
				rowsPaquetes.Close()
				return nil, err
			}
			reserva.Paquetes = append(reserva.Paquetes, paquete)
		}

		rowsPaquetes.Close()
		if err = rowsPaquetes.Err(); err != nil {
			return nil, err
		}

		reservas = append(reservas, reserva)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reservas, nil
}

// GetTotalReservasByInstancia obtiene el número total de reservas para una instancia específica
func (r *ReservaRepository) GetTotalReservasByInstancia(idInstancia int) (int, error) {
	var total int
	query := `SELECT COUNT(*) FROM reserva 
             WHERE id_instancia = $1 AND estado != 'CANCELADA' AND eliminado = FALSE`

	err := r.db.QueryRow(query, idInstancia).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

// GetTotalPasajerosByInstancia obtiene el total de pasajeros reservados para una instancia específica
func (r *ReservaRepository) GetTotalPasajerosByInstancia(idInstancia int) (int, error) {
	// Primero, obtener todas las reservas no canceladas para esta instancia
	query := `SELECT id_reserva FROM reserva 
             WHERE id_instancia = $1 AND estado != 'CANCELADA' AND eliminado = FALSE`

	rows, err := r.db.Query(query, idInstancia)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	// Calcular el total de pasajeros sumando los de cada reserva
	totalPasajeros := 0

	for rows.Next() {
		var idReserva int
		err := rows.Scan(&idReserva)
		if err != nil {
			return 0, err
		}

		// Obtener pasajeros de esta reserva
		pasajerosReserva, err := r.GetCantidadPasajerosByReserva(idReserva)
		if err != nil {
			return 0, err
		}

		totalPasajeros += pasajerosReserva
	}

	if err = rows.Err(); err != nil {
		return 0, err
	}

	return totalPasajeros, nil
}

// VerificarDisponibilidadInstancia verifica si hay suficiente cupo en una instancia para un número de pasajeros
func (r *ReservaRepository) VerificarDisponibilidadInstancia(idInstancia int, cantidadPasajeros int) (bool, error) {
	var cupoDisponible int
	query := `SELECT cupo_disponible FROM instancia_tour 
             WHERE id_instancia = $1 AND eliminado = FALSE AND estado = 'PROGRAMADO'`

	err := r.db.QueryRow(query, idInstancia).Scan(&cupoDisponible)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, errors.New("la instancia del tour no existe, está eliminada o no está programada")
		}
		return false, err
	}

	return cantidadPasajeros <= cupoDisponible, nil
}

// ReservarInstanciaMercadoPago crea una reserva a través de Mercado Pago
func (r *ReservaRepository) ReservarInstanciaMercadoPago(reserva *entidades.NuevaReservaRequest) (int, string, error) {
	// Iniciar transacción
	tx, err := r.db.Begin()
	if err != nil {
		return 0, "", err
	}

	// Si hay error, hacer rollback
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Verificar que al menos haya un pasaje o un paquete
	if len(reserva.CantidadPasajes) == 0 && len(reserva.Paquetes) == 0 {
		return 0, "", errors.New("debe incluir al menos un pasaje o un paquete en la reserva")
	}

	// Primero, obtener información de la instancia del tour para verificar disponibilidad
	var cupoDisponible int
	var nombreTour string
	queryInstancia := `SELECT it.cupo_disponible, tt.nombre 
                      FROM instancia_tour it
                      INNER JOIN tour_programado tp ON it.id_tour_programado = tp.id_tour_programado
                      INNER JOIN tipo_tour tt ON tp.id_tipo_tour = tt.id_tipo_tour
                      WHERE it.id_instancia = $1 AND it.eliminado = FALSE 
                      AND it.estado = 'PROGRAMADO'`
	err = tx.QueryRow(queryInstancia, reserva.IDInstancia).Scan(&cupoDisponible, &nombreTour)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, "", errors.New("la instancia del tour no existe, está eliminada o no está programada")
		}
		return 0, "", err
	}

	// Calcular el total de pasajeros
	totalPasajeros := 0

	// Sumar pasajeros de pasajes individuales
	for _, pasaje := range reserva.CantidadPasajes {
		totalPasajeros += pasaje.Cantidad
	}

	// Sumar pasajeros de paquetes
	for _, paquete := range reserva.Paquetes {
		// Obtener cantidad total de pasajeros por paquete
		var cantidadPorPaquete int
		queryPaquete := `SELECT cantidad_total FROM paquete_pasajes 
                        WHERE id_paquete = $1 AND eliminado = FALSE`
		err := tx.QueryRow(queryPaquete, paquete.IDPaquete).Scan(&cantidadPorPaquete)
		if err != nil {
			return 0, "", err
		}

		// Multiplicar por la cantidad de paquetes seleccionados
		totalPasajeros += cantidadPorPaquete * paquete.Cantidad
	}

	// Verificar disponibilidad
	if totalPasajeros > cupoDisponible {
		return 0, "", errors.New("no hay suficiente cupo disponible para la reserva")
	}

	// Consulta SQL para insertar una nueva reserva
	var idReserva int
	query := `INSERT INTO reserva (id_vendedor, id_cliente, id_instancia, id_canal, id_sede, 
             total_pagar, notas, estado, eliminado)
             VALUES ($1, $2, $3, $4, $5, $6, $7, 'RESERVADO', FALSE)
             RETURNING id_reserva`

	// Ejecutar la consulta con los datos de la reserva
	err = tx.QueryRow(
		query,
		reserva.IDVendedor,
		reserva.IDCliente,
		reserva.IDInstancia,
		reserva.IDCanal,
		reserva.IDSede,
		reserva.TotalPagar,
		reserva.Notas,
	).Scan(&idReserva)

	if err != nil {
		return 0, "", err
	}

	// Insertar las cantidades de pasajes individuales
	for _, pasaje := range reserva.CantidadPasajes {
		// Solo insertar si la cantidad es mayor que cero
		if pasaje.Cantidad > 0 {
			queryPasaje := `INSERT INTO pasajes_cantidad (id_reserva, id_tipo_pasaje, cantidad, eliminado)
                          VALUES ($1, $2, $3, FALSE)`

			_, err = tx.Exec(queryPasaje, idReserva, pasaje.IDTipoPasaje, pasaje.Cantidad)
			if err != nil {
				return 0, "", err
			}
		}
	}

	// Insertar los paquetes de pasajes
	for _, paquete := range reserva.Paquetes {
		// Solo insertar si la cantidad es mayor que cero
		if paquete.Cantidad > 0 {
			queryPaquete := `INSERT INTO paquete_pasaje_detalle (id_reserva, id_paquete, cantidad, eliminado)
                           VALUES ($1, $2, $3, FALSE)`

			_, err = tx.Exec(queryPaquete, idReserva, paquete.IDPaquete, paquete.Cantidad)
			if err != nil {
				return 0, "", err
			}
		}
	}

	// Actualizar el cupo disponible en la instancia del tour
	queryUpdateCupo := `UPDATE instancia_tour 
                       SET cupo_disponible = cupo_disponible - $1 
                       WHERE id_instancia = $2`
	_, err = tx.Exec(queryUpdateCupo, totalPasajeros, reserva.IDInstancia)
	if err != nil {
		return 0, "", err
	}

	// Commit de la transacción
	err = tx.Commit()
	if err != nil {
		return 0, "", err
	}

	return idReserva, nombreTour, nil
}
