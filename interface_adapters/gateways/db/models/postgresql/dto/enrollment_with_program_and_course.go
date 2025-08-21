package dto

import (
	"time"

	"gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/domain/types"
)

type EnrollmentWithProgramAndCourse struct {
	ID                     int        `gorm:"column:ING_ID"`
	IdPg                   int        `gorm:"column:ING_ID_PG"` // Local ID for the enrollment
	EnrollPeriod           *int       `gorm:"column:ING_PERIODO"`
	Nis                    int        `gorm:"column:NIS"`
	CourseId               int        `gorm:"column:FIC_ID" db:"FIC_ID"`
	EnrollStatus           *int       `gorm:"column:ING_ESTADO" db:"ING_ESTADO"`
	EnrollRegDate          *time.Time `gorm:"column:ING_FCH_REGISTRO;autoCreateTime" db:"ING_FCH_REGISTRO"`
	NisRegisFun            *int       `gorm:"column:NIS_FUN_REGISTRO" db:"NIS_FUN_REGISTRO"`
	EnrollAppliesAgreement string     `gorm:"column:ING_APLICO_CONVENIO" db:"ING_APLICO_CONVENIO"`               // CHAR(1)
	EnrollAgreemNumber     *string    `gorm:"column:ING_NUMERO_CONVENIO" db:"ING_NUMERO_CONVENIO"`               // VARCHAR2(50)
	EnrollDecreeReinstated *string    `gorm:"column:ING_DECRETO_REINTEGRADO" db:"ING_DECRETO_REINTEGRADO"`       // VARCHAR2(20)
	EnrollObsReintegrado   *string    `gorm:"column:ING_OBS_REINTEGRADO" db:"ING_OBS_REINTEGRADO"`               // VARCHAR2(1000)
	EnrollPercenWeigh      *float64   `gorm:"column:ING_PORCENTAJE_PONDERACION" db:"ING_PORCENTAJE_PONDERACION"` // NUMBER(13,2)
	EnrollTotalScore       int        `gorm:"column:ING_PUNTAJE_TOTAL" db:"ING_PUNTAJE_TOTAL"`                   // NUMBER(13,2)
	EnrollVirtualPriority  *int       `gorm:"column:ING_PRIORIDAD_VIRTUAL" db:"ING_PRIORIDAD_VIRTUAL"`           // NUMBER(1,0)
	FteIDAccesoPref        *int       `gorm:"column:FTE_ID_ACCESO_PREF" db:"FTE_ID_ACCESO_PREF"`
	State                  string     `gorm:"column:state"`
	RequestData *types.RequestData `gorm:"column:REQUEST_DATA;type:jsonb"`
	AttemptNumbers int 	`gorm:"column:number_attemps"`
	SosId                  *int        `gorm:"column:SOS_ID"`
	ProgramType            string     `gorm:"column:PRF_TIPO_PROGRAMA"`
	NameProgram 		string 	`gorm:"column:PRF_DENOMINACION"`
}
