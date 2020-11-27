package grpc

import (
	"context"
	gm "gigantic_billing/proto"
	"gigantic_billing/src/timings"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"io"
	"time"
)

type Client interface {
	SendEmail(request *gm.EmailRequest) *gm.BaseResponse
	SendEmailStream(requests []*gm.EmailRequest) timings.StreamTimingResponse
	CloseConn()
}

type grpcClient struct {
	client     gm.GiganticMessagingClient
	connection *grpc.ClientConn
}

func NewGrpcClient() Client {
	conn, err := grpc.Dial("localhost:6565", grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		log.WithError(err).Fatalln("Couldn't connect to gRPC server")
	}

	return &grpcClient{gm.NewGiganticMessagingClient(conn), conn}
}

func (grpc *grpcClient) SendEmail(request *gm.EmailRequest) *gm.BaseResponse {
	response, err := grpc.client.SendEmail(context.Background(), request)
	if err != nil {
		log.WithError(err).Fatalln("Communication failed with grpc server")
	}
	return response
}

func (grpc *grpcClient) SendEmailStream(requests []*gm.EmailRequest) timings.StreamTimingResponse {
	stream, err := grpc.client.SendEmailStream(context.Background())

	if err != nil {
		log.WithError(err).Fatalln("Communication failed with grpc server")
	}
	wait := make(chan struct{})
	var timing timings.StreamTimingResponse
	go func() {
		for {
			start := time.Now()
			response, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(wait)
				timing.TotalResponseTime += time.Since(start)
				return
			}
			if err != nil {
				log.WithError(err).Fatalln("Failed to receive response status")
			}
			if response.Status != gm.BaseResponse_SUCCESS_OK {
				log.WithField("messageID", response.MessageId).Println("Couldn't send out message, setting email to failed")
			} else {
				log.WithField("messageID", response.MessageId).Println("Message sent")
			}
			timing.TotalResponseTime += time.Since(start)
		}
	}()
	for _, req := range requests {
		start := time.Now()
		if err := stream.Send(req); err != nil {
			log.WithError(err).Fatalln("Failed to send email")
		}
		timing.TotalRequestTime += time.Since(start)
	}
	_ = stream.CloseSend()
	<-wait

	timing.AverageRequestTime = time.Duration(int(timing.TotalRequestTime.Nanoseconds()) / len(requests))
	timing.AverageResponseTime = time.Duration(int(timing.TotalResponseTime.Nanoseconds()) / len(requests))

	return timing
}

func (grpc *grpcClient) CloseConn() {
	grpc.connection.Close()
}
