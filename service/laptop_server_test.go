package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/warnshun/pcbook/pb"
	"github.com/warnshun/pcbook/sample"
	"github.com/warnshun/pcbook/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestServerCreateLaptop(t *testing.T) {
	t.Parallel()

	laptopNoId := sample.NewLaptop()
	laptopNoId.Id = ""

	laptopIncalidId := sample.NewLaptop()
	laptopIncalidId.Id = "invalid-uuid"

	laptopDuplicateId := sample.NewLaptop()
	storeDuplicateId := service.NewInMemoryLaptopStore()
	err := storeDuplicateId.Save(laptopDuplicateId)

	require.Nil(t, err)

	testCases := []struct {
		name        string
		laptop      *pb.Laptop
		laptopStore service.LaptopStore
		code        codes.Code
	}{
		{
			name:        "success_with_id",
			laptop:      sample.NewLaptop(),
			laptopStore: service.NewInMemoryLaptopStore(),
			code:        codes.OK,
		},
		{
			name:        "success_without_id",
			laptop:      laptopNoId,
			laptopStore: service.NewInMemoryLaptopStore(),
			code:        codes.OK,
		},
		{
			name:        "failure_invalid_id",
			laptop:      laptopIncalidId,
			laptopStore: service.NewInMemoryLaptopStore(),
			code:        codes.InvalidArgument,
		},
		{
			name:        "failure_duplicate_id",
			laptop:      laptopDuplicateId,
			laptopStore: storeDuplicateId,
			code:        codes.AlreadyExists,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := &pb.CreateLaptopRequest{
				Laptop: tc.laptop,
			}

			server := service.NewLaptopServer(tc.laptopStore, nil)
			res, err := server.CreateLaptop(context.Background(), req)
			if tc.code == codes.OK {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.NotEmpty(t, res.Id)
				if len(tc.laptop.Id) > 0 {
					require.Equal(t, tc.laptop.Id, res.Id)
				}
			} else {
				require.Error(t, err)
				require.Nil(t, res)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tc.code, st.Code())
			}
		})
	}
}
