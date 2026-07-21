package internal

import (
	"cmp"
	"flag"
	"os"
	"strconv"
)

type Config struct {
	Port  int
	Debug bool
}

func ReadConfig() (Config, error) {
	var cfg Config
	var err error
	flag.IntVar(&cfg.Port, "port", 8080, "Указание порта работы сервиса")
	flag.BoolVar(&cfg.Debug, "debug", false, "Режим отладки")
	flag.Parse()

	port := cmp.Or(
		os.Getenv("SHORTENER_PORT"),
		"8080",
	)
	if port != "8080" && cfg.Port == 8080 {
		cfg.Port, err = strconv.Atoi(port)
		if err != nil {
			return Config{}, err
		}
	}

	return cfg, nil
}
