package rest

import (
	"bytes"
	"encoding/json"
	gm "gigantic_billing/proto"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type giganticMessagingClient struct {
	client *http.Client
}

type GiganticMessagingClient interface {
	SendEmail(request *gm.EmailRequest) *gm.BaseResponse
}

func NewGiganticMessagingClient() GiganticMessagingClient {
	return &giganticMessagingClient{&http.Client{}}
}

func (client giganticMessagingClient) SendEmail(request *gm.EmailRequest) *gm.BaseResponse {
	reqBody, err := json.Marshal(request)
	if err != nil {
		log.WithError(err).Fatalln("Couldn't serialize to json")
	}
	res, err := client.client.Post("http://localhost:8080/emails/send", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		print(err)
	}
	defer res.Body.Close()

	var response *gm.BaseResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		log.WithError(err).Fatalln("Couldn't decode response")
	}
	return response
}
