package agent

// Protocol names are local to the SQL agent package so consumers do not need
// to import a component-specific connection implementation.
const (
	ProtocolMySQL      = "mysql"
	ProtocolMariaDB    = "mariadb"
	ProtocolPostgreSQL = "postgresql"
	ProtocolSQLServer  = "sqlserver"
	ProtocolOracle     = "oracle"
	ProtocolClickHouse = "clickhouse"
	ProtocolDameng     = "dameng"
	ProtocolDB2        = "db2"
)
