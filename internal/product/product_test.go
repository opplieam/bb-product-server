package product

import (
	"context"
	"errors"
	"net"
	"testing"

	pb "github.com/opplieam/bb-grpc/protogen/go/product"
	"github.com/opplieam/bb-product-server/.gen/buy-better-core/public/model"
	dbStore "github.com/opplieam/bb-product-server/internal/store"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

type ProductUnitTestSuite struct {
	suite.Suite
}

func TestProductHandler(t *testing.T) {
	suite.Run(t, new(ProductUnitTestSuite))
}

func (s *ProductUnitTestSuite) TestGetProductsByUserUnit() {
	sellerUrl := "s1.com"
	testData := []dbStore.ResultGetAllProducts{
		{
			ID:   1,
			Name: "P1",
			URL:  "p1.com",
			Sellers: model.Sellers{
				ID:   1,
				Name: "Seller1",
				URL:  &sellerUrl,
			},
			ImageProduct: []dbStore.ImageProductCustom{
				{ID: 1, ImageURL: "p1.png"},
			},
			PriceNow: []dbStore.PriceNowCustom{
				{ID: 1, Price: 1.0},
			},
		},
	}

	testCases := []struct {
		name       string
		buildStubs func(store *MockStorer)
		req        *pb.GetProductsByUserReq
		expStatus  codes.Code
		expLen     int
	}{
		{
			name: "successful get products",
			buildStubs: func(store *MockStorer) {
				store.EXPECT().GetAllProducts(mock.Anything).Return(testData, nil).Once()
			},
			req:       &pb.GetProductsByUserReq{UserId: 1},
			expStatus: codes.OK,
			expLen:    len(testData),
		},
		{
			name: "products not found",
			buildStubs: func(store *MockStorer) {
				store.EXPECT().GetAllProducts(mock.Anything).Return(nil, dbStore.ErrRecordNotFound).Once()
			},
			req:       &pb.GetProductsByUserReq{UserId: 150},
			expStatus: codes.NotFound,
			expLen:    0,
		},
		{
			name: "failed get products",
			buildStubs: func(store *MockStorer) {
				store.EXPECT().GetAllProducts(mock.Anything).Return(nil, errors.New("another error")).Once()
			},
			req:       &pb.GetProductsByUserReq{UserId: 2},
			expStatus: codes.Internal,
			expLen:    0,
		},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			mockStore := NewMockStorer(s.T())
			tc.buildStubs(mockStore)

			lis := bufconn.Listen(1024 * 1024)
			gServer := grpc.NewServer()
			tp := trace.NewTracerProvider()
			mt := tp.Tracer("test")
			pb.RegisterProductServiceServer(gServer, NewServer(mockStore, mt))

			go func() {
				err := gServer.Serve(lis)
				s.Require().NoError(err)
			}()
			defer gServer.Stop()

			dialOps := []grpc.DialOption{
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
					return lis.Dial()
				}),
			}
			resolver.SetDefaultScheme("passthrough")

			conn, err := grpc.NewClient(lis.Addr().String(), dialOps...)
			s.Require().NoError(err)
			defer conn.Close()

			client := pb.NewProductServiceClient(conn)
			ctx := context.Background()
			resp, err := client.GetProductsByUser(ctx, tc.req)

			s.Assert().Equal(tc.expStatus, status.Code(err))
			if resp != nil {
				s.Assert().Equal(tc.expLen, len(resp.Products))
			}
		})

	}
}
