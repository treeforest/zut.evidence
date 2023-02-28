package conf

type Config struct {
	Redis *Redis
	Mysql *Mysql
}

func Default() *Config {
	return &Config{
		Redis: &Redis{
			Addr:     "localhost:6379",
			Password: "",
			Maxidle:  2048,
		},
		Mysql: &Mysql{
			User:     "user",
			Pass:     "123456",
			Addr:     "localhost:3306",
			Database: "zut",
		},
	}
}

type Redis struct {
	Addr     string
	Password string
	Maxidle  int
}

type Mysql struct {
	User, Pass string
	Addr       string
	Database   string
}
