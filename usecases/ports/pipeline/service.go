package pipeline

import "context"

type Service interface{
	SynchronizeData(ctx context.Context)error
}