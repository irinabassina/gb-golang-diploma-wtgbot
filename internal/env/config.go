package env

import (
	"net/url"
	"strconv"
	"time"
)

type Config struct {
	PostgresConfig PostgresConfig `env:",prefix=DB_"`
	TokenConfig    TokenConfig    `env:",prefix=TOKEN_"`
	Logger         LoggerConfig   `env:",prefix=LOGGER_"`
}

type LoggerConfig struct {
	Level string `env:"LEVEL,default=debug"`
	Debug bool   `env:"DEBUG,default=true"`
}

type TokenConfig struct {
	TgToken     string `env:"TG,default="`
	OpenAiToken string `env:"OPEN_AI,default="`
}

type PostgresConfig struct {
	Name         string        `env:"NAME,default=wtgbot" `
	User         string        `env:"USER,default=postgres" `
	Host         string        `env:"HOST,default=localhost" `
	Port         int           `env:"PORT,default=5432" `
	SSLMode      string        `env:"SSLMODE,default=disable" `
	ConnTimeout  int           `env:"CONN_TIMEOUT,default=5" `
	Password     string        `env:"PASSWORD,default=postgres" `
	PoolMinConns int           `env:"POOL_MIN_CONNS,default=10" `
	PoolMaxConns int           `env:"POOL_MAX_CONNS,default=50" `
	DBTimeout    time.Duration `env:"TIMEOUT,default=5s"`
}

func (c PostgresConfig) ConnectionURL() string {
	host := c.Host
	if v := c.Port; v != 0 {
		host = host + ":" + strconv.Itoa(c.Port)
	}

	u := &url.URL{
		Scheme: "postgres",
		Host:   host,
		Path:   c.Name,
	}

	if c.User != "" || c.Password != "" {
		u.User = url.UserPassword(c.User, c.Password)
	}

	q := u.Query()
	if v := c.ConnTimeout; v > 0 {
		q.Add("connect_timeout", strconv.Itoa(v))
	}
	if v := c.SSLMode; v != "" {
		q.Add("sslmode", v)
	}

	u.RawQuery = q.Encode()

	return u.String()
}
