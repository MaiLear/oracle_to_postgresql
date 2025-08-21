package oracle

import (
	"time"
)

type CourseModel struct {
	ID                   int64      `db:"FIC_ID"`
	Quota                int64      `db:"FIC_CUPO"`
	RegionID             int64      `db:"RGN_ID"`
	BranchID             int64      `db:"SED_ID"`
	SubsiteID            *int64     `db:"SSD_ID"`
	ProgramId            int        `db:"PRF_ID"`
	RegisteredByNIS      *int64     `db:"NIS_FUN_REGISTRO"`
	MunicipalityID       int64      `db:"MPO_ID"`
	StartDate            time.Time  `db:"FIC_FCH_INICIALIZACION"`
	EndDate              time.Time  `db:"FIC_FCH_FINALIZACION"`
	Responsible          *string    `db:"FIC_RESPONSABLE"`
	LearningCenterID     int64      `db:"LTC_ID"`
	TrainingRouteID      int64      `db:"TRD_ID"`
	OfferedServiceID     int64      `db:"NFS_ID_OFRECIDO"`
	TrainingModality     string     `db:"FIC_MOD_FORMACION"`
	ShiftID              *int64     `db:"JOR_ID"`
	EmployerNIS          *int64     `db:"NIS_EMP"`
	Status               int        `db:"FIC_ESTADO"`
	RegistrationDate     time.Time  `db:"FIC_FCH_REGISTRO"`
	CancellationDate     *time.Time `db:"FIC_FCH_CANCELACION"`
	CancellationReason   *string    `db:"FIC_MOTIVO_CANCELACION"`
	LearningMethodID     *int64     `db:"LMS_ID"`
	SemesterID           *int64     `db:"SEM_ID"`
	MinimumQuota         *int64     `db:"FIC_CUPO_MINIMO"`
	TrialQuotaCount      *int64     `db:"FIC_VECES_CUPO_PRUEBA"`
	ManagerNIS           *int64     `db:"NIS_FUN_GESTOR"`
	ApprovedScheduling   string     `db:"FIC_PROGRAMACION_APROBADA"`
	ApprovedScheduleDate *time.Time `db:"FIC_FCH_PROG_APROBADA"`
	OfferProfileID       *int64     `db:"POF_ID"`
	PreconditionID       *int64     `db:"PRE_ID"`
	PreconditionName     *string    `db:"PRE_NOMBRE"`
}

func (c CourseModel) TableName() string {
	return "INTEGRACION.V_FICHA_CARACTERIZACION_B FC"
}
