package grpcClients

import (
	"fmt"
	"google.golang.org/grpc/credentials/insecure"

	"mail/config"

	"google.golang.org/grpc"
	"mail/gen/go/auth"
)

type Clients struct {
	AuthConn *proto.AuthServiceClient
}

func Init(cfg *config.Config) (*Clients, error) {
	//auth microservice
	authConn, err := grpc.NewClient(fmt.Sprintf("%s:%s", cfg.AuthServer.IP, cfg.AuthServer.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf(fmt.Sprintf("the microservice 'authorization' is not available: %v", err))
		return nil, err
	}
	authClient := proto.NewAuthServiceClient(authConn)
	return &Clients{
		AuthConn: &authClient,
	}, nil
}
