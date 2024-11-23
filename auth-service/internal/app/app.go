package app

import (
	"mail/config"
	"mail/internal/delivery/httpserver"
	"mail/pkg/logger"
	"mail/service/postgres"
	"mail/service/redis"
)

func Run(cfg *config.Config, l logger.Logger) error {
	dbPostgres, err := postgres.Init(cfg)
	if err != nil {
		return err
	}

	redisSessionClient, err := redis.Init(cfg, 0)
	if err != nil {
		return err
	}

	redisCSRFClient, err := redis.Init(cfg, 1)
	if err != nil {
		return err
	}
	
	
	port := fmt.Sprintf(":%d", cfg.AuthServer.Port)
	conn, err := net.Listen("tcp", port)
	fmt.Println("auth started")
	
	server := gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
        recovery.UnaryServerInterceptor(),
    ))
    
	auth.RegisterAuthManagerServer(server, NewAuthServer(dbPostgres, redisSessionClient, redisCSRFClient))
	
	err = server.Serve(conn)
	if err != nil {
		fmt.Println("oh nooooo")
	}

}
