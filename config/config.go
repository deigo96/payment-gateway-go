package config

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/elgs/gojq"
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
)

type AppConfig struct {
	Server_Url   string
	JWTKey       string
	Midtrans_url string
	Server_key   string
	Client_key   string
	Midtrans_env int8
	DRIVER       string
	DB_HOST      string
	DB_PORT      string
	DB_USER      string
	DB_PASS      string
	DB_NAME      string
}

var lock = &sync.Mutex{}
var appConfig *AppConfig

func GetConfig() *AppConfig {
	lock.Lock()
	defer lock.Unlock()

	if appConfig == nil {
		appConfig = initConfig()
	}
	return appConfig
}

func initConfig() *AppConfig {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	App_env := os.Getenv("APP_ENV")
	var midtrans_url string
	var Server_key string
	var Client_key string
	var midtrans_env int8
	if App_env == "development" {
		midtrans_url = os.Getenv("SANDBOX_URL")
		Server_key = os.Getenv("SERVER_SANDBOX_KEY")
		Client_key = os.Getenv("CLIENT_SANDBOX_KEY")
		midtrans_env = 1
	} else if App_env == "production" {
		fmt.Println("Production")
		midtrans_url = os.Getenv("SANDBOX_URL")
		Server_key = os.Getenv("SERVER_SANDBOX_KEY")
		Client_key = os.Getenv("CLIENT_SANDBOX_KEY")
		midtrans_env = 2
	} else {
		log.Fatal("Status App Env not allowed")
	}

	return &AppConfig{
		Server_Url:   os.Getenv("SERVER_URL"),
		JWTKey:       os.Getenv("JWTKEY"),
		Midtrans_url: midtrans_url,
		Server_key:   Server_key,
		Client_key:   Client_key,
		Midtrans_env: midtrans_env,
		DRIVER:       os.Getenv("DRIVER"),
		DB_HOST:      os.Getenv("DB_HOST"),
		DB_PORT:      os.Getenv("DB_PORT"),
		DB_USER:      os.Getenv("DB_USER"),
		DB_PASS:      os.Getenv("DB_PASS"),
		DB_NAME:      os.Getenv("DB_NAME"),
	}
}

type ServerConfig struct {
	Host string
	Port string
}

func GetServer() ServerConfig {
	var serverConfig ServerConfig
	config := GetConfig()
	prt, err := http.Get(config.Server_Url)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(prt.Body)
	if err != nil {
		log.Fatal(err)
	}
	sb := string(body)
	p := strings.TrimSuffix(sb, "\n")
	parser, err := gojq.NewStringQuery(p)
	if err != nil {
		fmt.Println(err)
	}

	sName, _ := parser.QueryToString("data.results.[2].serviceName")
	host, _ := parser.QueryToString("data.results.[2].serviceHost")
	sPort, _ := parser.QueryToString("data.results.[2].servicePort")

	port := fmt.Sprintf(":%s", sPort)
	fmt.Printf("server %s started on port %s", sName, sPort)
	serverConfig.Port = port
	serverConfig.Host = host
	return serverConfig
}

func GetJwtServer() ServerConfig {
	var serverConfig ServerConfig
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading dotenv file")
	}
	prt, err := http.Get(os.Getenv("SERVER_URL"))
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(prt.Body)
	if err != nil {
		log.Fatal(err)
	}
	sb := string(body)
	p := strings.TrimSuffix(sb, "\n")
	parser, err := gojq.NewStringQuery(p)
	if err != nil {
		fmt.Println(err)
	}
	host, _ := parser.QueryToString("data.results.[1].serviceHost")
	port, _ := parser.QueryToString("data.results.[1].servicePort")

	serverConfig.Port = port
	serverConfig.Host = host
	return serverConfig
}
