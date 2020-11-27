package invoices

import (
	gm "gigantic_billing/proto"
	grpc2 "gigantic_billing/src/grpc"
	rest2 "gigantic_billing/src/rest"
	"gigantic_billing/src/timings"
	"github.com/brianvoe/gofakeit/v5"
	log "github.com/sirupsen/logrus"
	"time"
)

type invoices struct {
	emailQueue []*gm.EmailRequest
	grpc       grpc2.Client
	rest       rest2.GiganticMessagingClient
}

type Invoices interface {
	GenerateRandomEmailInvoices(amount int)
	SendOutInvoicesGrpc() timings.TimingResponse
	SendOutInvoicesRest() timings.TimingResponse
	SendOutInvoicesGrpcStream() timings.StreamTimingResponse
}

func NewEmailInvoicing(grpc grpc2.Client, rest rest2.GiganticMessagingClient) Invoices {
	return &invoices{
		emailQueue: []*gm.EmailRequest{},
		grpc:       grpc,
		rest:       rest,
	}
}

func (in *invoices) GenerateRandomEmailInvoices(amount int) {
	gofakeit.Seed(0)
	log.WithField("amount", amount).Println("Generating random email invoices to send out")
	start := time.Now()
	for i := 0; i < amount; i++ {
		in.emailQueue = append(in.emailQueue, &gm.EmailRequest{
			Email:     gofakeit.Email(),
			MessageId: int32(i),
			Body:      gofakeit.LoremIpsumParagraph(5, 5, 50, "\n"),
		})
	}
	log.WithField("executionTime", time.Since(start)).Println("Finished generating invoices")
}

func (in *invoices) SendOutInvoicesGrpc() timings.TimingResponse {
	log.Println("Starting sending of invoices through gRPC")

	var timing timings.TimingResponse
	for _, email := range in.emailQueue {
		start := time.Now()
		in.processEmailGrpc(email)
		executionTime := time.Since(start)
		timing.TotalTime += executionTime
		log.WithField("executionTime", executionTime).Println("Execution time")
	}
	averageDuration := int(timing.TotalTime.Nanoseconds()) / len(in.emailQueue)
	timing.AverageTime = time.Duration(averageDuration)
	return timing
}

func (in *invoices) SendOutInvoicesGrpcStream() timings.StreamTimingResponse {
	log.Println("Starting sending of invoices through gRPC stream")
	start := time.Now()
	res := in.grpc.SendEmailStream(in.emailQueue)
	log.WithField("executionTime", time.Since(start)).Println("Finished sending through gRPC stream")
	return res
}

func (in *invoices) SendOutInvoicesRest() timings.TimingResponse {
	var timing timings.TimingResponse
	for _, email := range in.emailQueue {
		start := time.Now()
		in.processEmailRest(email)
		executionTime := time.Since(start)
		timing.TotalTime += executionTime
		log.WithField("executionTime", executionTime).Println("Execution time")
	}
	averageDuration := int(timing.TotalTime.Nanoseconds()) / len(in.emailQueue)
	timing.AverageTime = time.Duration(averageDuration)
	return timing
}

func (in *invoices) processEmailGrpc(email *gm.EmailRequest) {
	response := in.grpc.SendEmail(email)
	if response.Status != gm.BaseResponse_SUCCESS_OK {
		log.WithField("messageID", response.MessageId).Println("Couldn't send out message, setting email to failed")
	} else {
		log.WithField("messageID", response.MessageId).Println("Message sent")
	}
}

func (in *invoices) processEmailRest(email *gm.EmailRequest) {
	response := in.rest.SendEmail(email)
	if response.Status != gm.BaseResponse_SUCCESS_OK {
		log.WithField("messageID", response.MessageId).Println("Couldn't send out message, setting email to failed")
	} else {
		log.WithField("messageID", response.MessageId).Println("Message sent")
	}
}
