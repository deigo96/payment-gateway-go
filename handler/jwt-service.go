package handler

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"topup-service/config"

	"github.com/elgs/gojq"
)

type JwtService interface {
	ValidateToken(token string) bool
}

type jwtService struct{}

func NewJWTService() JwtService {
	return &jwtService{}
}

func (c *jwtService) ValidateToken(token string) bool {
	server := config.GetJwtServer()
	url := fmt.Sprintf("http://%s:%s/api/verify-token", server.Host, server.Port)

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", token)
	if err != nil {
		fmt.Println("Request error", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response. \n[ERROR] -", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}

	str := string(body)
	parser, _ := gojq.NewStringQuery(str)
	status, _ := parser.QueryToBool("status")

	return status
}
