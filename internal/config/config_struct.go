package config

type App struct {
	Env  string `env:"APP_ENV" env-default:"dev"`
	Port string `env:"APP_PORT" env-default:"3000"`
}

type Postgres struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	Name         string `yaml:"name"`
	SSLMode      string `yaml:"sslMode"`
	MaxOpenConns int    `yaml:"maxOpenConns"`
	MaxIdleConns int    `yaml:"maxIdleConns"`
}

type JWT struct {
	AccessSecret     []byte `env:"JWT_ACCESS_SECRET" env-required:"true"`
	RefreshSecret    []byte `env:"JWT_REFRESH_SECRET" env-required:"true"`
	AccessExpireMin  int    `env:"JWT_ACCESS_EXPIRE_MIN" env-default:"30"`
	RefreshExpireDay int    `env:"JWT_REFRESH_EXPIRE_DAY" env-default:"14"`
}

type Log struct {
	Level string `yaml:"level"`
}

type Cors struct {
	AllowOrigins     []string `yaml:"allowOrigins"`
	AllowMethods     []string `yaml:"allowMethods"`
	AllowHeaders     []string `yaml:"allowHeaders"`
	AllowCredentials bool     `yaml:"allowCredentials"`
}

type Cookie struct {
	Name     string `yaml:"name"`
	Path     string `yaml:"path"`
	HttpOnly bool   `yaml:"httpOnly"`
	Secure   bool   `yaml:"secure"`
	SameSite string `yaml:"sameSite"`
	MaxAge   int    `yaml:"maxAge"`
}
