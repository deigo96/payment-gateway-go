package api

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"topup-service/helper"

	"github.com/elgs/gojq"
)

func GetUserSeen(token string) (e helper.ExtractToken, err error) {
	seenUrl := os.Getenv("USER_URL")
	url := fmt.Sprintf(seenUrl + "/profile")
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}

	str := string(body)

	parser, _ := gojq.NewStringQuery(str)
	status, _ := parser.QueryToBool("status")
	message, _ := parser.QueryToString("message")
	data, _ := parser.QueryToMap("data")

	res := helper.DataToken{
		Id:       int(data["id"].(float64)),
		Email:    data["email"].(string),
		Username: data["username"].(string),
		Phone:    data["phone"].(string),
	}

	e.Status = status
	e.Message = message
	e.Data = res

	return e, err
}

func GetStatusTransaction(orderId string) bool {
	sandboxURL := os.Getenv("SANDBOX_URL")
	url := fmt.Sprintf(sandboxURL + "/" + orderId + "/status")
	fmt.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	// req.Header.Add("Authorization", orderId)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}

	str := string(body)

	// parser, _ := gojq.NewStringQuery(str)
	// status, _ := parser.QueryToBool("status")
	// message, _ := parser.QueryToString("message")
	// data, _ := parser.QueryToMap("data")

	fmt.Println(str)

	return true
}

func GetUserProfile(username string) (e helper.RequestUser) {
	seenUrl := os.Getenv("USER_URL")
	url := fmt.Sprintf(seenUrl + "/username?username=" + username)
	req, err := http.NewRequest("GET", url, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}

	str := string(body)

	parser, _ := gojq.NewStringQuery(str)
	id, _ := parser.QueryToInt64("data.id")
	name, _ := parser.QueryToString("data.username")
	email, _ := parser.QueryToString("data.email")
	phone, _ := parser.QueryToString("data.phone")

	e.Id = int(id)
	e.Email = email
	e.Phone = phone
	e.Username = name

	return e
}
