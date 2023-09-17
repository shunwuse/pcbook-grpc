package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/warnshun/pcbook/pb"
	"github.com/warnshun/pcbook/service"
	"google.golang.org/grpc"
)

func main() {
	port := flag.Int("port", 0, "The port to run the server on")
	flag.Parse()
	log.Printf("Starting server on port %d", *port)

	laptopStore := service.NewInMemoryLaptopStore()
	imageStore := service.NewDiskImageStore("img")

	laptopServer := service.NewLaptopServer(laptopStore, imageStore)

	grpcServer := grpc.NewServer()
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)

	address := fmt.Sprintf("0.0.0.0:%d", *port)
	fmt.Printf("Starting server on %s", address)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
