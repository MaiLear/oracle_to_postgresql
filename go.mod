module gitlab.com/sofia-plus/oracle_to_postgresql

go 1.24.3

require (
	github.com/cockroachdb/errors v1.12.0
	github.com/joho/godotenv v1.5.1
	gitlab.com/sofia-plus/go_db_connectors v1.5.1
	gorm.io/gorm v1.30.1
)

replace gitlab.com/sofia-plus/go_db_connectors v1.5.1 => ../go_db_connectors

require (
	github.com/cockroachdb/logtags v0.0.0-20230118201751-21c54148d20b // indirect
	github.com/cockroachdb/redact v1.1.5 // indirect
	github.com/getsentry/sentry-go v0.27.0 // indirect
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/godror/godror v0.48.2 // indirect
	github.com/godror/knownpb v0.1.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.5.5 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	golang.org/x/crypto v0.36.0 // indirect
	golang.org/x/exp v0.0.0-20250531010427-b6e5de432a8b // indirect
	golang.org/x/sync v0.16.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.28.0 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gorm.io/driver/postgres v1.5.11 // indirect
)
