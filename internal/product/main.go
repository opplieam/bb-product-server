package product

import (
	"context"
	"errors"

	pb "github.com/opplieam/bb-grpc/protogen/go/product"
	"github.com/opplieam/bb-product-server/internal/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Storer interface {
	GetAllProducts(userID uint32) ([]store.ResultGetAllProducts, error)
}

type Server struct {
	Storer Storer
	pb.UnimplementedProductServiceServer
}

func NewServer(s Storer) *Server {
	return &Server{Storer: s}
}

func (s *Server) GetProductsByUser(ctx context.Context, req *pb.GetProductsByUserReq) (*pb.GetProductsByUserRes, error) {
	result, err := s.Storer.GetAllProducts(req.UserId)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrRecordNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	var products []*pb.Products
	for _, v := range result {
		product := &pb.Products{
			Id:   uint32(v.ID),
			Name: v.Name,
			Url:  v.URL,
		}
		var priceNow []*pb.Price
		for _, pn := range v.PriceNow {
			price := &pb.Price{
				Id:        uint32(pn.ID),
				Price:     pn.Price,
				CreatedAt: timestamppb.New(pn.CreatedAt),
			}
			priceNow = append(priceNow, price)
		}
		product.Prices = priceNow

		var images []*pb.Image
		for _, im := range v.ImageProduct {
			image := &pb.Image{
				Id:  uint32(im.ID),
				Url: im.ImageURL,
			}
			images = append(images, image)
		}
		product.Images = images

		product.Seller = &pb.Seller{
			Id:   uint32(v.Sellers.ID),
			Name: v.Sellers.Name,
			Url:  *v.Sellers.URL,
		}

		products = append(products, product)
	}
	return &pb.GetProductsByUserRes{Products: products}, nil
}
