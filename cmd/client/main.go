package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/warnshun/pcbook/pb"
	"github.com/warnshun/pcbook/sample"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func main() {
	serverAddress := flag.String("address", "", "The server address")
	flag.Parse()
	log.Printf("Dial server on %s", *serverAddress)

	conn, err := grpc.Dial(*serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("cannot dial server:", err)
	}
	defer conn.Close()

	laptopClient := pb.NewLaptopServiceClient(conn)

	laptop := sample.NewLaptop()
	// laptop.Id = ""
	// laptop.Id = "b92041c7-04fd-417d-a49c-fed49e5e72b6"
	// laptop.Id = "invalid-uuid"

	req := &pb.CreateLaptopRequest{
		Laptop: laptop,
	}

	// set timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := laptopClient.CreateLaptop(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.AlreadyExists {
			log.Print("laptop already exists")
		} else {
			log.Fatal("cannot create laptop:", err)
		}
		return
	}

	log.Printf("Created laptop with id: %s", res.Id)
}
