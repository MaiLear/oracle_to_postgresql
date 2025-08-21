package constants

// Estados de inscripción (ING_ESTADO)
const (
	EnrollmentStatusPending                int = 1  // Pendiente
	EnrollmentStatusEnrolled               int = 2  // Inscrito
	EnrollmentStatusNotAdmitted            int = 3  // No Admitido
	EnrollmentStatusCancelledEnrollment    int = 4  // Anulado Matricula
	EnrollmentStatusCalledEnrollment       int = 5  // Convocado Inscripcion
	EnrollmentStatusRegistered             int = 6  // Matriculado
	EnrollmentStatusPreEnrolled            int = 7  // Preinscrito
	EnrollmentStatusSelected               int = 8  // Seleccionado
	EnrollmentStatusNotSelected            int = 9  // No Seleccionado
	EnrollmentStatusAbsentEnrollment       int = 10 // Ausente Inscripcion
	EnrollmentStatusCancelledInscription   int = 11 // Anulado Inscripcion
	EnrollmentStatusCalledSelection        int = 12 // Convocado Selección
	EnrollmentStatusAbsentSelection        int = 13 // Ausente Selección
	EnrollmentStatusCalledRegistration     int = 14 // Convocado Matricula
	EnrollmentStatusPendingOfflineTest     int = 15 // Pendiente prueba offline
	EnrollmentStatusCertified              int = 16 // Certificado
	EnrollmentStatusTransferred            int = 17 // Trasladado
	EnrollmentStatusCancelled              int = 18 // Cancelado
	EnrollmentStatusFailed                 int = 19 // Reprobado
	EnrollmentStatusPostponed              int = 20 // Aplazado
	EnrollmentStatusVoluntaryWithdrawal    int = 21 // Retiro Voluntario
	EnrollmentStatusAcademicCancelled      int = 22 // Cancelado Academico
	EnrollmentStatusDisciplinaryCancelled  int = 23 // Cancelado Disciplinario
	EnrollmentStatusAtypicalEnrollment     int = 24 // Atipico Inscripcion
	EnrollmentStatusAssigned               int = 25 // Asignado
	EnrollmentStatusPendingVirtualPreEnrollment int = 26 // Pendiente Preinscripción Virtual
	EnrollmentStatusPendingLevelTest       int = 27 // Pendiente_prueba_Nivel
	EnrollmentStatusPendingSecondOption    int = 28 // Pendiente Segunda Opcion
	EnrollmentStatusCancelledSecondOption  int = 29 // Cancelado Segunda Opcion
	EnrollmentStatusTempCancelled          int = 30 // Cancelado Temp
	EnrollmentStatusUnauthorizedEnrollment int = 31 // Inscripción No autorizada
	EnrollmentStatusPendingConfirmation    int = 32 // Pendiente Confirmación
)

// Mapeo de códigos de estado a descripciones
var EnrollmentStatusMap = map[int]string{
	EnrollmentStatusPending:                     "PENDIENTE",
	EnrollmentStatusEnrolled:                    "INSCRITO",
	EnrollmentStatusNotAdmitted:                 "NO ADMITIDO",
	EnrollmentStatusCancelledEnrollment:         "ANULADO MATRICULA",
	EnrollmentStatusCalledEnrollment:            "CONVOCADO INSCRIPCION",
	EnrollmentStatusRegistered:                  "MATRICULADO",
	EnrollmentStatusPreEnrolled:                 "PREINSCRITO",
	EnrollmentStatusSelected:                    "SELECCIONADO",
	EnrollmentStatusNotSelected:                 "NO SELECCIONADO",
	EnrollmentStatusAbsentEnrollment:            "AUSENTE INSCRIPCION",
	EnrollmentStatusCancelledInscription:        "ANULADO INSCRIPCION",
	EnrollmentStatusCalledSelection:             "CONVOCADO SELECCIÓN",
	EnrollmentStatusAbsentSelection:             "AUSENTE SELECCIÓN",
	EnrollmentStatusCalledRegistration:          "CONVOCADO MATRICULA",
	EnrollmentStatusPendingOfflineTest:          "PENDIENTE PRUEBA OFFLINE",
	EnrollmentStatusCertified:                   "CERTIFICADO",
	EnrollmentStatusTransferred:                 "TRASLADADO",
	EnrollmentStatusCancelled:                   "CANCELADO",
	EnrollmentStatusFailed:                      "REPROBADO",
	EnrollmentStatusPostponed:                   "APLAZADO",
	EnrollmentStatusVoluntaryWithdrawal:         "RETIRO VOLUNTARIO",
	EnrollmentStatusAcademicCancelled:           "CANCELADO ACADEMICO",
	EnrollmentStatusDisciplinaryCancelled:       "CANCELADO DISCIPLINARIO",
	EnrollmentStatusAtypicalEnrollment:          "ATIPICO INSCRIPCION",
	EnrollmentStatusAssigned:                    "ASIGNADO",
	EnrollmentStatusPendingVirtualPreEnrollment: "PENDIENTE PREINSCRIPCIÓN VIRTUAL",
	EnrollmentStatusPendingLevelTest:            "PENDIENTE PRUEBA NIVEL",
	EnrollmentStatusPendingSecondOption:         "PENDIENTE SEGUNDA OPCION",
	EnrollmentStatusCancelledSecondOption:       "CANCELADO SEGUNDA OPCION",
	EnrollmentStatusTempCancelled:               "CANCELADO TEMP",
	EnrollmentStatusUnauthorizedEnrollment:      "INSCRIPCIÓN NO AUTORIZADA",
	EnrollmentStatusPendingConfirmation:         "PENDIENTE CONFIRMACIÓN",
}

// GetEnrollmentStatusDescription retorna la descripción del estado de inscripción
func GetEnrollmentStatusDescription(status int) string {
	if description, exists := EnrollmentStatusMap[status]; exists {
		return description
	}
	return "Estado Desconocido"
}