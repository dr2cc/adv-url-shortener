package config

import (
	"flag"
	"os"
	"time"
)

// переменная FlagRunAddr содержит адрес и порт для запуска сервера
var FlagRunAddr string

// переменная FlagURL отвечает за базовый адрес результирующего сокращённого URL
var FlagURL string

// ParseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func ParseFlags() {
	// регистрируем переменные
	flag.StringVar(&FlagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&FlagURL, "b", "none", "host and port")
	// разбираем переданные серверу аргументы коммандной строки в зарегистрированные переменные
	flag.Parse()

	// Добавляем переменные окружения
	// $env:SERVER_ADDRESS = "localhost:8089"
	// $env:BASE_URL  = "http://localhost:9999"
	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		FlagRunAddr = envRunAddr
	}
	if envURL := os.Getenv("BASE_URL"); envURL != "" {
		FlagURL = envURL
	}
}

// adv
// Структуры для анмаршаллинга
type Config struct {
	Env         string `yaml:"env" env-default:"development"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"0.0.0.0:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {

	// $env:CONFIG_PATH = "C:\__git\URLsShortener\config\local.yaml"       (на drkk)
	// $env:CONFIG_PATH = "C:\Mega\__git\URLsShortener\config\local.yaml"  (на ноуте)

	// // Если будут проблемы с переменной окружения, то писать путь так (\ экранируется \\):
	//configPath := "C:\\__git\\URLsShortener\\config\\local.yaml"

	// ya #start# - заполняем нашу структуру сразу
	cfg := Config{
		Env:         "local",
		StoragePath: "./storage.db",
		HTTPServer: HTTPServer{
			Address: FlagRunAddr, //"localhost:8080",
			//Timeout:     5,
			//IdleTimeout: 60,
		},
	}
	// ya #end#

	// // adv #start#
	// // Получаем путь до конфиг-файла из env-переменной CONFIG_PATH
	// configPath := os.Getenv("CONFIG_PATH")
	// if configPath == "" {
	// 	log.Fatal("CONFIG_PATH environment variable is not set")
	// }

	// // Проверяем существование конфиг-файла
	// if _, err := os.Stat(configPath); err != nil {
	// 	log.Fatalf("error opening config file: %s", err)
	// }

	// // Читаем конфиг-файл и заполняем нашу структуру

	// var cfg Config

	// err := cleanenv.ReadConfig(configPath, &cfg)
	// fmt.Println(cfg)
	// if err != nil {
	// 	log.Fatalf("error reading config file: %s", err)
	// }
	// // adv #end#

	return &cfg
}
