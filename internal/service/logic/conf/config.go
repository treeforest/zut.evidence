package conf

type Config struct {
	Mysql *Mysql
}

func Default() *Config {
	return &Config{
		Mysql: &Mysql{
			User:     "user",
			Pass:     "123456",
			Addr:     "localhost:3306",
			Database: "zut",
		},
	}
}

type Mysql struct {
	User, Pass string
	Addr       string
	Database   string
}
