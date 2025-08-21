package tools

import (
	"errors"

	"github.com/godror/godror"
	"github.com/jackc/pgx/v5/pgconn"
)

func InPosUniqueViolation(err error) bool {
	var pgError *pgconn.PgError
	if errors.As(err, &pgError) {
		return pgError.Code == "23505"
	}
	return false
}

func InOraUniqueViolation(err error) bool {
	var oerr *godror.OraErr
	if errors.As(err, &oerr) {
		return oerr.Code() == 1
	}
	return false
}
