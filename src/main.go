package main

import (
	"gigantic_billing/src/grpc"
	"gigantic_billing/src/invoices"
	"gigantic_billing/src/rest"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Println("Starting GIGANTIC Billing")

	grpcClient := grpc.NewGrpcClient()
	restClient := rest.NewGiganticMessagingClient()
	invoicing := invoices.NewEmailInvoicing(grpcClient, restClient)

	invoicing.GenerateRandomEmailInvoices(100000)
	timingRest := invoicing.SendOutInvoicesRest()
	timingGrpc := invoicing.SendOutInvoicesGrpc()
	timingStream := invoicing.SendOutInvoicesGrpcStream()

	log.WithFields(log.Fields{
		"average": timingRest.AverageTime,
		"total":   timingRest.TotalTime,
	}).Println("timing for REST method")
	log.WithFields(log.Fields{
		"average": timingGrpc.AverageTime,
		"total":   timingGrpc.TotalTime,
	}).Println("timing for gRPC method")
	log.WithFields(log.Fields{
		"averageResponse": timingStream.AverageResponseTime,
		"totalResponse":   timingStream.TotalResponseTime,
		"averageRequest":  timingStream.AverageRequestTime,
		"totalRequest":    timingStream.TotalRequestTime,
	}).Println("timing for gRPC stream method")

	grpcClient.CloseConn()
	log.Println("GIGANTIC Billing has ended his life-long job...")
}
