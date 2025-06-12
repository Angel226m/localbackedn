package servicios

import (
	"errors"
	"sistema-toursseft/internal/entidades"
	"sistema-toursseft/internal/repositorios"
	"time"
)

// TourProgramadoService maneja la lógica de negocio de los tours programados
type TourProgramadoService struct {
	repo            *repositorios.TourProgramadoRepository
	tipoTourRepo    *repositorios.TipoTourRepository
	embarcacionRepo *repositorios.EmbarcacionRepository
	horarioTourRepo *repositorios.HorarioTourRepository
	sedeRepo        *repositorios.SedeRepository
	usuarioRepo     *repositorios.UsuarioRepository
}

// NewTourProgramadoService crea una nueva instancia del servicio
func NewTourProgramadoService(
	repo *repositorios.TourProgramadoRepository,
	tipoTourRepo *repositorios.TipoTourRepository,
	embarcacionRepo *repositorios.EmbarcacionRepository,
	horarioTourRepo *repositorios.HorarioTourRepository,
	sedeRepo *repositorios.SedeRepository,
	usuarioRepo *repositorios.UsuarioRepository,
) *TourProgramadoService {
	return &TourProgramadoService{
		repo:            repo,
		tipoTourRepo:    tipoTourRepo,
		embarcacionRepo: embarcacionRepo,
		horarioTourRepo: horarioTourRepo,
		sedeRepo:        sedeRepo,
		usuarioRepo:     usuarioRepo,
	}
}

// GetByID obtiene un tour programado por su ID
func (s *TourProgramadoService) GetByID(id int) (*entidades.TourProgramado, error) {
	return s.repo.GetByID(id)
}

// Create crea un nuevo tour programado
func (s *TourProgramadoService) Create(tourProgramado *entidades.NuevoTourProgramadoRequest) (int, error) {
	// Validar que existan las entidades relacionadas
	_, err := s.tipoTourRepo.GetByID(tourProgramado.IDTipoTour)
	if err != nil {
		return 0, errors.New("tipo de tour no encontrado")
	}

	_, err = s.embarcacionRepo.GetByID(tourProgramado.IDEmbarcacion)
	if err != nil {
		return 0, errors.New("embarcación no encontrada")
	}

	_, err = s.horarioTourRepo.GetByID(tourProgramado.IDHorario)
	if err != nil {
		return 0, errors.New("horario no encontrado")
	}

	_, err = s.sedeRepo.GetByID(tourProgramado.IDSede)
	if err != nil {
		return 0, errors.New("sede no encontrada")
	}

	// Validar el chofer si se proporciona
	if tourProgramado.IDChofer != nil {
		usuario, err := s.usuarioRepo.GetByID(*tourProgramado.IDChofer)
		if err != nil {
			return 0, errors.New("chofer no encontrado")
		}
		if usuario.Rol != "CHOFER" {
			return 0, errors.New("el usuario seleccionado no tiene rol de chofer")
		}
	}

	// Validar formato de fechas
	_, err = time.Parse("2006-01-02", tourProgramado.Fecha)
	if err != nil {
		return 0, errors.New("formato de fecha inválido, debe ser YYYY-MM-DD")
	}

	vigenciaDesde, err := time.Parse("2006-01-02", tourProgramado.VigenciaDesde)
	if err != nil {
		return 0, errors.New("formato de vigencia desde inválido, debe ser YYYY-MM-DD")
	}

	vigenciaHasta, err := time.Parse("2006-01-02", tourProgramado.VigenciaHasta)
	if err != nil {
		return 0, errors.New("formato de vigencia hasta inválido, debe ser YYYY-MM-DD")
	}

	// Verificar que la vigencia sea coherente
	if vigenciaHasta.Before(vigenciaDesde) {
		return 0, errors.New("la fecha de vigencia hasta debe ser posterior a la fecha de vigencia desde")
	}

	// ELIMINADA la validación de fechas pasadas
	// ELIMINADA la validación de que la fecha esté dentro del rango de vigencia

	// Verificar compatibilidad con el horario seleccionado
	horario, err := s.horarioTourRepo.GetByID(tourProgramado.IDHorario)
	if err != nil {
		return 0, errors.New("error al verificar el horario: " + err.Error())
	}

	// Verificar que el día de la semana del tour coincida con los días disponibles en el horario
	fechaTour, _ := time.Parse("2006-01-02", tourProgramado.Fecha)
	diaSemana := fechaTour.Weekday()
	diaDisponible := false

	switch diaSemana {
	case 0: // Domingo
		diaDisponible = horario.DisponibleDomingo
	case 1: // Lunes
		diaDisponible = horario.DisponibleLunes
	case 2: // Martes
		diaDisponible = horario.DisponibleMartes
	case 3: // Miércoles
		diaDisponible = horario.DisponibleMiercoles
	case 4: // Jueves
		diaDisponible = horario.DisponibleJueves
	case 5: // Viernes
		diaDisponible = horario.DisponibleViernes
	case 6: // Sábado
		diaDisponible = horario.DisponibleSabado
	}

	// MANTENEMOS esta validación porque es importante que el día seleccionado esté disponible
	if !diaDisponible {
		return 0, errors.New("el día seleccionado no está disponible en el horario configurado")
	}

	return s.repo.Create(tourProgramado)
}

// Update actualiza un tour programado existente
func (s *TourProgramadoService) Update(id int, tourProgramado *entidades.ActualizarTourProgramadoRequest) error {
	// Obtener el tour actual para validaciones
	tourActual, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	// Validar las entidades relacionadas si se proporcionan
	if tourProgramado.IDTipoTour > 0 {
		_, err := s.tipoTourRepo.GetByID(tourProgramado.IDTipoTour)
		if err != nil {
			return errors.New("tipo de tour no encontrado")
		}
	}

	if tourProgramado.IDEmbarcacion > 0 {
		_, err := s.embarcacionRepo.GetByID(tourProgramado.IDEmbarcacion)
		if err != nil {
			return errors.New("embarcación no encontrada")
		}
	}

	if tourProgramado.IDHorario > 0 {
		_, err := s.horarioTourRepo.GetByID(tourProgramado.IDHorario)
		if err != nil {
			return errors.New("horario no encontrado")
		}
	}

	if tourProgramado.IDSede > 0 {
		_, err := s.sedeRepo.GetByID(tourProgramado.IDSede)
		if err != nil {
			return errors.New("sede no encontrada")
		}
	}

	// Validar el chofer si se proporciona
	if tourProgramado.IDChofer != nil {
		usuario, err := s.usuarioRepo.GetByID(*tourProgramado.IDChofer)
		if err != nil {
			return errors.New("chofer no encontrado")
		}
		if usuario.Rol != "CHOFER" {
			return errors.New("el usuario seleccionado no tiene rol de chofer")
		}
	}

	// Validar que el cupo disponible no sea mayor que el cupo máximo
	if tourProgramado.CupoDisponible > tourProgramado.CupoMaximo && tourProgramado.CupoMaximo > 0 {
		return errors.New("el cupo disponible no puede ser mayor que el cupo máximo")
	}

	// Variables para validación de fechas
	var fechaTour time.Time
	var vigenciaDesde, vigenciaHasta time.Time
	var errFecha, errVigenciaDesde, errVigenciaHasta error

	// Parsear fechas proporcionadas o usar las existentes
	if tourProgramado.Fecha != "" {
		fechaTour, errFecha = time.Parse("2006-01-02", tourProgramado.Fecha)
		if errFecha != nil {
			return errors.New("formato de fecha inválido, debe ser YYYY-MM-DD")
		}
	} else {
		fechaTour = tourActual.Fecha
	}

	if tourProgramado.VigenciaDesde != "" {
		vigenciaDesde, errVigenciaDesde = time.Parse("2006-01-02", tourProgramado.VigenciaDesde)
		if errVigenciaDesde != nil {
			return errors.New("formato de vigencia desde inválido, debe ser YYYY-MM-DD")
		}
	} else {
		vigenciaDesde = tourActual.VigenciaDesde
	}

	if tourProgramado.VigenciaHasta != "" {
		vigenciaHasta, errVigenciaHasta = time.Parse("2006-01-02", tourProgramado.VigenciaHasta)
		if errVigenciaHasta != nil {
			return errors.New("formato de vigencia hasta inválido, debe ser YYYY-MM-DD")
		}
	} else {
		vigenciaHasta = tourActual.VigenciaHasta
	}

	// Validar que las fechas sean coherentes
	if vigenciaHasta.Before(vigenciaDesde) {
		return errors.New("la fecha de vigencia hasta debe ser posterior a la fecha de vigencia desde")
	}

	// ELIMINADA la validación de fechas pasadas
	// ELIMINADA la validación de que la fecha esté dentro del rango de vigencia

	// Verificar compatibilidad con el horario seleccionado si cambia la fecha o el horario
	var horarioID int
	if tourProgramado.IDHorario > 0 {
		horarioID = tourProgramado.IDHorario
	} else {
		horarioID = tourActual.IDHorario
	}

	if tourProgramado.Fecha != "" || tourProgramado.IDHorario > 0 {
		horario, err := s.horarioTourRepo.GetByID(horarioID)
		if err != nil {
			return errors.New("error al verificar el horario: " + err.Error())
		}

		// Verificar que el día de la semana del tour coincida con los días disponibles en el horario
		diaSemana := fechaTour.Weekday()
		diaDisponible := false

		switch diaSemana {
		case 0: // Domingo
			diaDisponible = horario.DisponibleDomingo
		case 1: // Lunes
			diaDisponible = horario.DisponibleLunes
		case 2: // Martes
			diaDisponible = horario.DisponibleMartes
		case 3: // Miércoles
			diaDisponible = horario.DisponibleMiercoles
		case 4: // Jueves
			diaDisponible = horario.DisponibleJueves
		case 5: // Viernes
			diaDisponible = horario.DisponibleViernes
		case 6: // Sábado
			diaDisponible = horario.DisponibleSabado
		}

		// MANTENEMOS esta validación porque es importante que el día seleccionado esté disponible
		if !diaDisponible {
			return errors.New("el día seleccionado no está disponible en el horario configurado")
		}
	}

	// Validar que no se esté cancelando un tour ya en curso
	if tourProgramado.Estado == "CANCELADO" && tourActual.Estado == "EN_CURSO" {
		return errors.New("no se puede cancelar un tour que ya está en curso")
	}

	return s.repo.Update(id, tourProgramado)
}

// Delete elimina lógicamente un tour programado
func (s *TourProgramadoService) Delete(id int) error {
	// Verificar que el tour no esté en curso o completado
	tour, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if tour.Estado == "EN_CURSO" || tour.Estado == "COMPLETADO" {
		return errors.New("no se puede eliminar un tour que está en curso o completado")
	}

	return s.repo.SoftDelete(id)
}

// AsignarChofer asigna un chofer a un tour programado
func (s *TourProgramadoService) AsignarChofer(idTour int, idChofer int) error {
	// Validar que el chofer exista y tenga rol de chofer
	usuario, err := s.usuarioRepo.GetByID(idChofer)
	if err != nil {
		return errors.New("chofer no encontrado")
	}

	if usuario.Rol != "CHOFER" {
		return errors.New("el usuario seleccionado no tiene rol de chofer")
	}

	return s.repo.AsignarChofer(idTour, idChofer)
}

// CambiarEstado cambia el estado de un tour programado
func (s *TourProgramadoService) CambiarEstado(id int, estado string) error {
	// Validar la transición de estado
	tour, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	// Reglas de transición de estado
	switch tour.Estado {
	case "PROGRAMADO":
		if estado != "EN_CURSO" && estado != "CANCELADO" {
			return errors.New("desde PROGRAMADO solo se puede cambiar a EN_CURSO o CANCELADO")
		}
	case "EN_CURSO":
		if estado != "COMPLETADO" {
			return errors.New("desde EN_CURSO solo se puede cambiar a COMPLETADO")
		}
	case "COMPLETADO", "CANCELADO":
		return errors.New("no se puede cambiar el estado de un tour COMPLETADO o CANCELADO")
	}

	return s.repo.CambiarEstado(id, estado)
}

// List obtiene una lista de tours programados con filtros opcionales
func (s *TourProgramadoService) List(filtros entidades.FiltrosTourProgramado) ([]*entidades.TourProgramado, error) {
	return s.repo.List(filtros)
}

// GetProgramacionSemanal obtiene los tours programados para una semana específica
func (s *TourProgramadoService) GetProgramacionSemanal(fechaInicio string, idSede int) ([]*entidades.TourProgramado, error) {
	return s.repo.GetProgramacionSemanal(fechaInicio, idSede)
}

// GetToursDisponiblesEnFecha obtiene tours disponibles para una fecha específica
func (s *TourProgramadoService) GetToursDisponiblesEnFecha(fecha string, idSede int) ([]*entidades.TourProgramado, error) {
	return s.repo.GetToursDisponiblesEnFecha(fecha, idSede)
}

// GetToursDisponiblesEnRangoFechas obtiene los tours disponibles para reserva en un rango de fechas
func (s *TourProgramadoService) GetToursDisponiblesEnRangoFechas(fechaInicio, fechaFin string, idSede int) ([]*entidades.TourProgramado, error) {
	return s.repo.GetToursDisponiblesEnRangoFechas(fechaInicio, fechaFin, idSede)
}

// VerificarDisponibilidadHorario verifica si un horario está disponible para una fecha específica
func (s *TourProgramadoService) VerificarDisponibilidadHorario(idHorario int, fecha string) (bool, error) {
	// 1. Verificar si el horario existe
	horario, err := s.horarioTourRepo.GetByID(idHorario)
	if err != nil {
		return false, errors.New("horario no encontrado: " + err.Error())
	}

	// 2. Verificar si el día de la semana está disponible en el horario
	fechaObj, err := time.Parse("2006-01-02", fecha)
	if err != nil {
		return false, errors.New("formato de fecha inválido, debe ser YYYY-MM-DD")
	}

	// Verificar que el día de la semana del tour coincida con los días disponibles en el horario
	diaSemana := fechaObj.Weekday()
	diaDisponible := false

	switch diaSemana {
	case 0: // Domingo
		diaDisponible = horario.DisponibleDomingo
	case 1: // Lunes
		diaDisponible = horario.DisponibleLunes
	case 2: // Martes
		diaDisponible = horario.DisponibleMartes
	case 3: // Miércoles
		diaDisponible = horario.DisponibleMiercoles
	case 4: // Jueves
		diaDisponible = horario.DisponibleJueves
	case 5: // Viernes
		diaDisponible = horario.DisponibleViernes
	case 6: // Sábado
		diaDisponible = horario.DisponibleSabado
	}

	if !diaDisponible {
		// El día de la semana no está disponible en este horario
		return false, nil
	}

	// 3. Verificar si ya hay tours programados con este horario en esta fecha
	var filtros entidades.FiltrosTourProgramado
	fechaStr := fecha
	filtros.FechaInicio = &fechaStr
	filtros.FechaFin = &fechaStr

	tours, err := s.repo.List(filtros)
	if err != nil {
		return false, errors.New("error al verificar tours existentes: " + err.Error())
	}

	// Verificar si hay conflictos con otros tours en el mismo horario
	for _, tour := range tours {
		if tour.IDHorario == idHorario && !tour.Eliminado && tour.Estado != "CANCELADO" {
			// Ya existe un tour con este horario en esta fecha
			return false, nil
		}
	}

	// No hay conflictos, el horario está disponible
	return true, nil
}

// ProgramarToursSemanal crea tours programados para una semana basados en un template
func (s *TourProgramadoService) ProgramarToursSemanal(fechaInicio string, tourBase *entidades.NuevoTourProgramadoRequest, cantidadDias int) ([]int, error) {
	fechaInicioObj, err := time.Parse("2006-01-02", fechaInicio)
	if err != nil {
		return nil, errors.New("formato de fecha inválido, debe ser YYYY-MM-DD")
	}

	// Validar entidades relacionadas una sola vez
	_, err = s.tipoTourRepo.GetByID(tourBase.IDTipoTour)
	if err != nil {
		return nil, errors.New("tipo de tour no encontrado")
	}

	_, err = s.embarcacionRepo.GetByID(tourBase.IDEmbarcacion)
	if err != nil {
		return nil, errors.New("embarcación no encontrada")
	}

	_, err = s.horarioTourRepo.GetByID(tourBase.IDHorario)
	if err != nil {
		return nil, errors.New("horario no encontrado")
	}

	_, err = s.sedeRepo.GetByID(tourBase.IDSede)
	if err != nil {
		return nil, errors.New("sede no encontrada")
	}

	// Validar el chofer si se proporciona
	if tourBase.IDChofer != nil {
		usuario, err := s.usuarioRepo.GetByID(*tourBase.IDChofer)
		if err != nil {
			return nil, errors.New("chofer no encontrado")
		}
		if usuario.Rol != "CHOFER" {
			return nil, errors.New("el usuario seleccionado no tiene rol de chofer")
		}
	}

	// Validar fechas de vigencia
	vigenciaDesde, err := time.Parse("2006-01-02", tourBase.VigenciaDesde)
	if err != nil {
		return nil, errors.New("formato de vigencia desde inválido, debe ser YYYY-MM-DD")
	}

	vigenciaHasta, err := time.Parse("2006-01-02", tourBase.VigenciaHasta)
	if err != nil {
		return nil, errors.New("formato de vigencia hasta inválido, debe ser YYYY-MM-DD")
	}

	// Verificar que la vigencia sea coherente
	if vigenciaHasta.Before(vigenciaDesde) {
		return nil, errors.New("la fecha de vigencia hasta debe ser posterior a la fecha de vigencia desde")
	}

	// Obtener el horario para verificar días disponibles
	horario, err := s.horarioTourRepo.GetByID(tourBase.IDHorario)
	if err != nil {
		return nil, errors.New("error al verificar el horario: " + err.Error())
	}

	// Crear tours para cada día
	tourIDs := []int{}
	for i := 0; i < cantidadDias; i++ {
		currentDate := fechaInicioObj.AddDate(0, 0, i)

		// Verificar si el día de la semana está disponible en el horario
		diaSemana := currentDate.Weekday()
		diaDisponible := false

		switch diaSemana {
		case 0: // Domingo
			diaDisponible = horario.DisponibleDomingo
		case 1: // Lunes
			diaDisponible = horario.DisponibleLunes
		case 2: // Martes
			diaDisponible = horario.DisponibleMartes
		case 3: // Miércoles
			diaDisponible = horario.DisponibleMiercoles
		case 4: // Jueves
			diaDisponible = horario.DisponibleJueves
		case 5: // Viernes
			diaDisponible = horario.DisponibleViernes
		case 6: // Sábado
			diaDisponible = horario.DisponibleSabado
		}

		// Solo crear tour si el día está disponible
		if diaDisponible {
			// Crear una copia del tour base con la fecha actualizada
			nuevoTour := *tourBase
			nuevoTour.Fecha = currentDate.Format("2006-01-02")

			// ELIMINADA la validación de que la fecha esté dentro del rango de vigencia

			// Crear el tour
			id, err := s.repo.Create(&nuevoTour)
			if err != nil {
				// Continuar con el siguiente día si hay error en este
				continue
			}

			tourIDs = append(tourIDs, id)
		}
	}

	if len(tourIDs) == 0 {
		return nil, errors.New("no se pudo crear ningún tour para las fechas seleccionadas")
	}

	return tourIDs, nil
}

// En servicios/TourProgramadoService.go

// GetToursDisponibles obtiene los tours disponibles para reserva
func (s *TourProgramadoService) GetToursDisponibles() ([]*entidades.TourProgramado, error) {
	return s.repo.GetToursDisponibles()
}
