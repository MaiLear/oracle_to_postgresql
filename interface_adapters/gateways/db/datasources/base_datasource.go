package datasources

import (
	"context"
	"database/sql"
)

type BaseDatasource struct{
	connection *sql.DB
	placeholder string
}

func (b BaseDatasource) GetAll(ctx context.Context,model any)(data []any,err error){
	b.connection.QueryContext(ctx,`SELECT * FROM `)
}