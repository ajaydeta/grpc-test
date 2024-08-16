package main

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"tablelink/config"
	"tablelink/internal/repository"
	"tablelink/internal/service"
	pb "tablelink/proto/pb/proto"
)

func main() {
	db, err := config.MustPostgres()
	if err != nil {
		panic(err)
	}

	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo)

	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	pb.RegisterUserServiceServer(s, svc)

	fmt.Println("server listening at port 8080")
	err = s.Serve(listen)
	if err != nil {
		panic(err)
	}
}
