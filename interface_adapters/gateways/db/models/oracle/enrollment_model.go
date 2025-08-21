package oracle

import "time"

type EnrollmentModel struct {
	ID                     int        `db:"ING_ID"`
	EnrollPeriod           *int       `db:"ING_PERIODO"`
	Nis                    int        `db:"NIS"`
	CourseId               int        `db:"FIC_ID"`
	EnrollStatus           *int       `db:"ING_ESTADO"`
	EnrollRegDate          *time.Time `db:"ING_FCH_REGISTRO"`
	NisRegisFun            *int       `db:"NIS_FUN_REGISTRO"`
	EnrollAppliesAgreement string     `db:"ING_APLICO_CONVENIO"`
	EnrollAgreemNumber     *string    `db:"ING_NUMERO_CONVENIO"`
	EnrollDecreeReinstated *string    `db:"ING_DECRETO_REINTEGRADO"`
	EnrollObsReintegrado   *string    `db:"ING_OBS_REINTEGRADO"`
	EnrollPercenWeigh      *float64   `db:"ING_PORCENTAJE_PONDERACION"`
	EnrollTotalScore       int        `db:"ING_PUNTAJE_TOTAL"`
	EnrollVirtualPriority  *int       `db:"ING_PRIORIDAD_VIRTUAL"`
	FteIDAccesoPref        *int       `db:"FTE_ID_ACCESO_PREF"`
}

func (e EnrollmentModel) TableName() string {
	return "INSCRIPCION.INGRESO_ASPIRANTE"
}
