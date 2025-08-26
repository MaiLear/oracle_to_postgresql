package loggers

import (
	"context"

	usecasesDto "gitlab.com/sofia-plus/oracle_to_postgresql/usecases/dto"
)

type Logger interface {
	Save(context.Context, usecasesDto.LogError) error
}
