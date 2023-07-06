package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	once     sync.Once
	instance *Config
)

type Config struct {
	Server   *serverConfig
	Database *dbConfig
}

type serverConfig struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type dbConfig struct {
	Host    string
	Timeout string
}

func New() *Config {
	once.Do(func() {
		host := os.Getenv("HOST")

		port, err := strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			log.Fatal(err)
		}

		readTimeout, err := strconv.Atoi(os.Getenv("READ_TIMEOUT"))
		if err != nil {
			log.Fatal(err)
		}

		writeTimeout, err := strconv.Atoi(os.Getenv("WRITE_TIMEOUT"))
		if err != nil {
			log.Fatal(err)
		}

		idleTimeout, err := strconv.Atoi(os.Getenv("IDLE_TIMEOUT"))
		if err != nil {
			log.Fatal(err)
		}

		dbHost := os.Getenv("MONGODB_HOST")

		instance = &Config{
			Server: &serverConfig{
				Addr:         fmt.Sprintf("%s:%d", host, port),
				ReadTimeout:  time.Duration(readTimeout) * time.Second,
				WriteTimeout: time.Duration(writeTimeout) * time.Second,
				IdleTimeout:  time.Duration(idleTimeout) * time.Second,
			},
			Database: &dbConfig{
				Host: dbHost,
			},
		}
	})

	return instance
}
