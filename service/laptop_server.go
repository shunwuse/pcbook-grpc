package service

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"

	"github.com/google/uuid"
	"github.com/warnshun/pcbook/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	MAX_IMAGE_SIZE = 1 << 20 // 1 megabyte
	// MAX_IMAGE_SIZE = 1 << 10 // 1 kilobyte
)

// LaptopServer is the server API for Laptop service.
type LaptopServer struct {
	pb.UnimplementedLaptopServiceServer
	laptopStore LaptopStore
	imageStore  ImageStore
	ratingStore RatingStore
}

// NewLaptopServer returns a new LaptopServer.
func NewLaptopServer(laptopStore LaptopStore, imageStore ImageStore, ratingStore RatingStore) *LaptopServer {
	return &LaptopServer{
		laptopStore: laptopStore,
		imageStore:  imageStore,
		ratingStore: ratingStore,
	}
}

// CreateLaptop is a unary RPC to create a new laptop
func (s *LaptopServer) CreateLaptop(ctx context.Context, req *pb.CreateLaptopRequest) (*pb.CreateLaptopResponse, error) {
	laptop := req.GetLaptop()
	log.Printf("Received a CreateLaptopRequest with id: %s", laptop.Id)

	if len(laptop.Id) > 0 {
		// check if it's a valid UUID
		_, err := uuid.Parse(laptop.Id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "Laptop ID is not a valid UUID: %v", err)
		}
	} else {
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "cannot generate a new laptop ID: %v", err)
		}
		laptop.Id = id.String()
	}

	// some heavy processing
	// time.Sleep(6 * time.Second)

	// check context error
	if err := contextError(ctx); err != nil {
		return nil, err
	}

	// save the laptop in memory store
	err := s.laptopStore.Save(laptop)
	if err != nil {
		code := codes.Internal
		if errors.Is(err, ErrAlreadyExists) {
			code = codes.AlreadyExists
		}

		return nil, status.Errorf(code, "cannot save laptop to the store: %v", err)
	}

	log.Printf("Saved laptop with id: %s", laptop.Id)

	res := &pb.CreateLaptopResponse{
		Id: laptop.Id,
	}

	return res, nil
}

// SearchLaptop is a server-streaming RPC to search for laptops
func (s *LaptopServer) SearchLaptop(req *pb.SearchLaptopRequest, stream pb.LaptopService_SearchLaptopServer) error {
	filter := req.GetFilter()
	log.Printf("receive a search-laptop request with filter: %v", filter)

	err := s.laptopStore.Search(
		stream.Context(),
		filter,
		func(laptop *pb.Laptop) error {
			res := &pb.SearchLaptopResponse{Laptop: laptop}
			err := stream.Send(res)
			if err != nil {
				return err
			}

			log.Printf("sent laptop with id: %s", laptop.GetId())
			return nil
		},
	)

	if err != nil {
		return status.Errorf(codes.Internal, "unexpected error: %v", err)
	}

	return nil
}

// UploadImage is a client-streaming RPC to upload the laptop image
func (s *LaptopServer) UploadImage(stream pb.LaptopService_UploadImageServer) error {
	req, err := stream.Recv()
	if err != nil {
		return logErrors(status.Error(codes.Unknown, "cannot receive image info"))
	}

	laptopId := req.GetInfo().GetLaptopId()
	imageType := req.GetInfo().GetImageType()
	log.Printf("receive an image upload request for laptop %s with image type %s", laptopId, imageType)

	laptop, err := s.laptopStore.Find(laptopId)
	if err != nil {
		return logErrors(status.Errorf(codes.Internal, "cannot find laptop: %v", err))
	}
	if laptop == nil {
		return logErrors(status.Errorf(codes.InvalidArgument, "laptop %s not found", laptopId))
	}

	imageData := bytes.Buffer{}
	imageSize := 0

	for {
		// check context error
		if err := contextError(stream.Context()); err != nil {
			return err
		}

		log.Print("waiting to receive more data")

		req, err := stream.Recv()
		if err == io.EOF {
			log.Print("no more data")
			break
		}
		if err != nil {
			return logErrors(status.Errorf(codes.Unknown, "cannot receive chunk data: %v", err))
		}

		chunk := req.GetChunkData()
		size := len(chunk)

		log.Printf("received a chunk with size: %d", size)

		imageSize += size
		if imageSize > MAX_IMAGE_SIZE {
			return logErrors(status.Errorf(codes.InvalidArgument, "image is too large: %d > %d", imageSize, MAX_IMAGE_SIZE))
		}

		// write slowly
		// time.Sleep(1 * time.Second)

		_, err = imageData.Write(chunk)
		if err != nil {
			return logErrors(status.Errorf(codes.Internal, "cannot write image data: %v", err))
		}
	}

	imageId, err := s.imageStore.Save(laptopId, imageType, imageData)
	if err != nil {
		return logErrors(status.Errorf(codes.Internal, "cannot save image to the store: %v", err))
	}

	res := &pb.UploadImageResponse{
		Id:   imageId,
		Size: uint32(imageSize),
	}

	err = stream.SendAndClose(res)
	if err != nil {
		return logErrors(status.Errorf(codes.Unknown, "cannot send response: %v", err))
	}

	log.Printf("saved image with id: %s, size: %d", imageId, imageSize)
	return nil
}

// RateLaptop is a bidirectional-streaming RPC to rate a laptop
func (s *LaptopServer) RateLaptop(stream pb.LaptopService_RateLaptopServer) error {
	for {
		err := contextError(stream.Context())
		if err != nil {
			return err
		}

		req, err := stream.Recv()
		if err == io.EOF {
			log.Print("no more data")
			break
		}
		if err != nil {
			return logErrors(status.Errorf(codes.Unknown, "cannot receive stream request: %v", err))
		}

		laptopId := req.GetLaptopId()
		score := req.GetScore()

		log.Printf("received a rate-laptop request: id = %s, score = %.2f", laptopId, score)

		found, err := s.laptopStore.Find(laptopId)
		if err != nil {
			return logErrors(status.Errorf(codes.Internal, "cannot find laptop: %v", err))
		}
		if found == nil {
			return logErrors(status.Errorf(codes.InvalidArgument, "laptopID %s not found", laptopId))
		}

		rating, err := s.ratingStore.Add(laptopId, score)
		if err != nil {
			return logErrors(status.Errorf(codes.Internal, "cannot add rating to the store: %v", err))
		}

		res := &pb.RateLaptopResponse{
			LaptopId:     laptopId,
			RatedCount:   rating.Count,
			AverageScore: rating.Sum / float64(rating.Count),
		}

		err = stream.Send(res)
		if err != nil {
			return logErrors(status.Errorf(codes.Unknown, "cannot send stream response: %v", err))
		}

		// log.Printf("sent a rate-laptop response: %v", res)
	}
	return nil
}

func contextError(ctx context.Context) error {
	switch ctx.Err() {
	// check context canceled
	case context.Canceled:
		return logErrors(status.Error(codes.Canceled, "the client canceled the request"))
		// check context deadline
	case context.DeadlineExceeded:
		return logErrors(status.Error(codes.DeadlineExceeded, "deadline exceeded"))
	}

	return nil
}

func logErrors(err error) error {
	if err != nil {
		log.Print(err)
	}

	return err
}
