package constants

// Estados de registro académico (RGA_ESTADO) usados en validaciones
const (
	AcademicRecordStatusCancelled         int = 2 // Cancelado
	AcademicRecordStatusVoluntaryWithdraw int = 6 // Retiro Voluntario
)

// Estados de novedad académica (NAP_ESTADO)
const (
	AcademicNoveltyStatusAccepted   int = 1 // Aceptada
	AcademicNoveltyStatusNoResponse int = 3 // Sin respuesta
	AcademicNoveltyStatusProcessed  int = 4 // Tramitada
)

// Flag que indica si la novedad valida en inscripción (SN.SUN_VALIDA_INSCRIPCION)
const AcademicNoveltyValidatesEnrollmentFlag = "1"

// Valores especiales de novedad
const AcademicNoveltyDurationIndefinite = -1

// Parámetros COMUN.PARAMETRO
const (
	ParamIDConfrontationIdentity = 259 // habilitar confrontación identidad
	ParamIDConfirmationDays      = 97  // días confirmación inscripción
)

// Lugares de realización de la ficha
const (
	PlaceOfRealizationEnterprise = "VE"
	PlaceOfRealizationNational   = "VN"
)
