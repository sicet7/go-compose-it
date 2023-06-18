package dbconfig

type DatabaseConfig struct {
	typeName string
	dsn      string
}

func (c *DatabaseConfig) Type() string {
	return c.typeName
}

func (c *DatabaseConfig) DSN() string {
	return c.dsn
}
