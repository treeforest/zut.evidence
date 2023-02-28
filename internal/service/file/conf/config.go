package conf

type Config struct {
	Redis *Redis
	Ipfs  *Ipfs
}

func Default() *Config {
	return &Config{
		Redis: &Redis{
			Addr:     "localhost:6379",
			Password: "",
			Maxidle:  2048,
		},
		Ipfs: &Ipfs{
			Url: "localhost:5001",
		},
	}
}

type Redis struct {
	Addr     string
	Password string
	Maxidle  int
}

type Ipfs struct {
	Url string
}
