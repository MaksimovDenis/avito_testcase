package main

import (
	avito "avito_testcase"
	logger "avito_testcase/logs"
	"avito_testcase/package/handler"
	"avito_testcase/package/repository"
	"avito_testcase/package/service"
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Port string `yaml:"port"`
	DB   struct {
		Username string `yaml:"username"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		DBname   string `yaml:"dbname"`
		SSLmode  string `yaml:"sslmode"`
	}
	Redis struct {
		Addr string `yaml:"addr"`
		DB   int    `yaml:"db"`
	}
}

func initConfig() (*Config, error) {
	var config Config

	file, err := os.Open("configs/config.yaml")
	if err != nil {
		logrus.Errorf("failed to open config file: %v", err)
		return nil, fmt.Errorf("failed to open config file: %v", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		logrus.Errorf("failed to read config file: %v", err)
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		logrus.Errorf("failed to unmarshal config data: %v", err)
		return nil, fmt.Errorf("failed to unmarshal config data: %s", err.Error())
	}

	return &config, nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	config, err := initConfig()
	if err != nil {
		logrus.Fatal("error initializing config:", err)
	}

	//Initializing our PostgresDB
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     config.DB.Host,
		Port:     config.DB.Port,
		Username: config.DB.Username,
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   config.DB.DBname,
		SSLMode:  config.DB.SSLmode,
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	//Initializing our RedisDB
	repository.InitRedisClient(config.Redis.Addr, os.Getenv("REDIS_PASSWORD"), config.Redis.DB)

	ping, err := repository.ClientRedis.Ping(context.Background()).Result()
	if err != nil {
		logrus.Fatalf("failed to initialize redis: %s", err.Error())
	}
	logger.Log.Info(ping)

	//Creating our dependencies
	repositories := repository.NewRepository(db)
	services := service.NewService(repositories)
	handlers := handler.NewHandler(services)

	//Running server
	srv := new(avito.Server)
	go func() {
		if err := srv.Run(config.Port, handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	logger.Log.Info("App started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logger.Log.Info("App is shutting down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logger.Log.Errorf("error occured on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logger.Log.Errorf("error occured on db connection close: %s", err.Error())
	}

}
