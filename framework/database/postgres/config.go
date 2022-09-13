package postgres

type Config struct {
	host     string
	database string
	port     int
	user     string
	password string
	sslmode  string
}

func NewConfig() *Config {
	return &Config{
		host:     "0.0.0.0",
		database: "postgres",
		port:     5432,
		user:     "postgres",
		password: "postgres",
		sslmode:  "disable",
	}
}

func (c *Config) WithHost(host string) *Config {
	c.host = host
	return c
}

func (c *Config) WithDatabase(database string) *Config {
	c.database = database
	return c
}

func (c *Config) WithPort(port int) *Config {
	c.port = port
	return c
}

func (c *Config) WithUser(user string) *Config {
	c.user = user
	return c
}

func (c *Config) WithPassword(password string) *Config {
	c.password = password
	return c
}

func (c *Config) WithSSLMode(sslmode string) *Config {
	c.sslmode = sslmode
	return c
}
