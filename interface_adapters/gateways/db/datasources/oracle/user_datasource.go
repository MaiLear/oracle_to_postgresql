package oracle

import (
	"context"
	"database/sql"

	cockroachdbErrors "github.com/cockroachdb/errors"
	oraclePort "gitlab.com/sofia-plus/go_db_connectors/ports/repositories/oracle"
	oracleModels "gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/infrastructure/datasources/models/oracle"
)

type UserDataSource struct {
	connection *sql.DB
	db         oraclePort.OracleRepository
}

func NewUserDataSource(db oraclePort.OracleRepository, connection *sql.DB) UserDataSource {
	return UserDataSource{
		connection: connection,
		db:         db,
	}
}

func (u UserDataSource) SaveUserCompletely(ctx context.Context, user oracleModels.UserModel, basicDataUser oracleModels.BasicDataUserModel, userLocation oracleModels.UserLocationModel, applicant oracleModels.ApplicantModel) (nis int, err error) {
	_, err = u.connection.ExecContext(ctx, `
BEGIN
	COMUN.USP_INSERTAR_USUARIO(
		:p_tipo_documento,
		:p_num_doc_identidad,
		:p_usr_fch_registro,
		:p_usr_clave,
		:p_asp_nombre,
		:p_asp_primer_apellido,
		:p_asp_segundo_apellido,
		:p_asp_correo_e,
		:p_nis_fun_registro,
		:p_dbu_fch_exd_doc_identidad,
		:p_dbu_genero,
		:p_dbu_fch_nacimiento,
		:p_mpo_id_nacimiento,
		:p_mpo_nombre_nacimiento,
		:p_dbu_lib_militar,
		:p_dbu_estado_civil,
		:p_dbu_estrato,
		:p_dbu_tipo_sangre,
		:p_dbu_afiliado_eps,
		:p_dbu_eps,
		:p_dbu_puntaje_icfes,
		:p_dbu_es_media_tecnica,
		:p_mpo_id_exp_doc_identidad,
		:p_ubu_dir_residencia,
		:p_pai_id_residencia,
		:p_pai_nombre_residencia,
		:p_dpt_id_residencia,
		:p_dpt_nombre_residencia,
		:p_mpo_id_residencia,
		:p_mpo_nombre_residencia,
		:p_zon_id_residencia,
		:p_zon_nombre_residencia,
		:p_bar_id_residencia,
		:p_bar_nombre_residencia,
		:p_ubu_tel_principal,
		:p_ubu_tel_alternativo,
		:p_ubu_tel_movil,
		:p_dbu_nombre_contacto,
		:p_dbu_tel_fijo_contacto,
		:p_dbu_parentesco_contacto,
		:p_dbu_fecha_vencimiento,
		:p_nis
	);
END;
`,
		sql.Named("p_tipo_documento", user.DocumentType),
		sql.Named("p_num_doc_identidad", user.DocumentNumber),
		sql.Named("p_usr_fch_registro", user.RegistrationDate),
		sql.Named("p_usr_clave", user.Password),
		sql.Named("p_asp_nombre", applicant.Name),
		sql.Named("p_asp_primer_apellido", applicant.FirstSurname),
		sql.Named("p_asp_segundo_apellido", applicant.SecondSurname),
		sql.Named("p_asp_correo_e", applicant.Email),
		sql.Named("p_nis_fun_registro", applicant.IdRegOfficer),
		sql.Named("p_dbu_fch_exd_doc_identidad", basicDataUser.DocumentIssueDate),
		sql.Named("p_dbu_genero", basicDataUser.Gender),
		sql.Named("p_dbu_fch_nacimiento", basicDataUser.BirthDate),
		sql.Named("p_mpo_id_nacimiento", basicDataUser.BirthMunicipalityID),
		sql.Named("p_mpo_nombre_nacimiento", basicDataUser.BirthMunicipalityName),
		sql.Named("p_dbu_lib_militar", basicDataUser.MilitaryCard),
		sql.Named("p_dbu_estado_civil", basicDataUser.MaritalStatus),
		sql.Named("p_dbu_estrato", basicDataUser.Stratum),
		sql.Named("p_dbu_tipo_sangre", basicDataUser.BloodType),
		sql.Named("p_dbu_afiliado_eps", basicDataUser.EpsAffiliated),
		sql.Named("p_dbu_eps", basicDataUser.EpsName),
		sql.Named("p_dbu_puntaje_icfes", basicDataUser.IcfesScore),
		sql.Named("p_dbu_es_media_tecnica", basicDataUser.IsTechnicalHighSchool),
		sql.Named("p_mpo_id_exp_doc_identidad", basicDataUser.DocExpMunicipalityID),
		sql.Named("p_ubu_dir_residencia", userLocation.Address),
		sql.Named("p_pai_id_residencia", userLocation.CountryID),
		sql.Named("p_pai_nombre_residencia", userLocation.CountryName),
		sql.Named("p_dpt_id_residencia", userLocation.DepartmentID),
		sql.Named("p_dpt_nombre_residencia", userLocation.DepartmentName),
		sql.Named("p_mpo_id_residencia", userLocation.MunicipalityID),
		sql.Named("p_mpo_nombre_residencia", userLocation.MunicipalityName),
		sql.Named("p_zon_id_residencia", userLocation.ZoneID),
		sql.Named("p_zon_nombre_residencia", userLocation.ZoneName),
		sql.Named("p_bar_id_residencia", userLocation.NeighborhoodID),
		sql.Named("p_bar_nombre_residencia", userLocation.NeighborhoodName),
		sql.Named("p_ubu_tel_principal", userLocation.MainPhone),
		sql.Named("p_ubu_tel_alternativo", userLocation.AlternativePhone),
		sql.Named("p_ubu_tel_movil", userLocation.MobilePhone),
		sql.Named("p_dbu_nombre_contacto", basicDataUser.EmergencyContactName),
		sql.Named("p_dbu_tel_fijo_contacto", basicDataUser.EmergencyPhoneLandline),
		sql.Named("p_dbu_parentesco_contacto", basicDataUser.EmergencyRelationship),
		sql.Named("p_dbu_fecha_vencimiento", basicDataUser.ExpirationDate),
		sql.Named("p_nis", sql.Out{Dest: &nis}),
	)
	if err != nil {
		err = cockroachdbErrors.Wrap(err, "infra: ocurrio un error insertando los datos del usuario")
		return
	}
	return
}
