package wishlist

import (
	"context"

	wl "github.com/DenisKokorin/WishListProto/gen/go"
	"google.golang.org/grpc"
)

type serverAPI struct {
	wl.UnimplementedWishlistServiceServer
}

func Register(gRPC *grpc.Server) {
	wl.RegisterWishlistServiceServer(gRPC, &serverAPI{})
}

func (s *serverAPI) CreateWishlist(ctx context.Context, req *wl.CreateWishlistRequest) (*wl.WishListResponse, error) {
	return nil, nil
}

func (s *serverAPI) GetWishList(ctx context.Context, req *wl.GetWishListRequest) (*wl.WishListResponse, error) {
	return nil, nil
}

func (s *serverAPI) UpdateWishList(ctx context.Context, req *wl.UpdateWishListRequest) (*wl.WishListResponse, error) {
	return nil, nil
}

func (s *serverAPI) DeleteWishList(ctx context.Context, req *wl.DeleteWishListRequest) (*wl.DeleteWishListResponse, error) {
	return nil, nil
}

func (s *serverAPI) AddItem(ctx context.Context, req *wl.AddItemRequest) (*wl.ItemResponse, error) {
	return nil, nil
}

func (s *serverAPI) UpdateItem(ctx context.Context, req *wl.UpdateItemRequest) (*wl.ItemResponse, error) {
	return nil, nil
}

func (s *serverAPI) DeleteItem(ctx context.Context, req *wl.DeleteItemRequest) (*wl.DeleteItemResponse, error) {
	return nil, nil
}
