package grpcClients

import (
	"fmt"

	cfg "mail/config"

	"google.golang.org/grpc"
)

type Clients struct {
	AuthConn    *grpc.ClientConn
}

func Init(cfg *config.Config) *Clients {
	//auth microservice
	authConn, err := grpc.NewClient(fmt.Sprintf("%s:%d", cfg.AuthServer.IP, cfg.AuthServer.Port))
	if err != nil {
		fmt.Printf(fmt.Sprintf("the microservice 'authorization' is not available: %v", err))
		return nil
	}
	return &Clients{
		AuthConn:    authConn,
	}
}