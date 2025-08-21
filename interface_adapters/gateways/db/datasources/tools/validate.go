package tools

import (
	"errors"

	cockroachdbErrors "github.com/cockroachdb/errors"
	internalErrors "gitlab.com/sofia-plus/pg_oracle_etl_sync/internal/tools/errors"
	"golang.org/x/exp/slices"
)

func Validate(toValidate []string, allowedFields []string) error {
	var allErrors []error
	for _, item := range toValidate {
		if !slices.Contains(allowedFields, item) {
			allErrors = append(allErrors, cockroachdbErrors.Wrapf(internalErrors.ErrBadRequest, "infra: el campo %s no esta permitido", item))
		}
	}
	return errors.Join(allErrors...)
}
