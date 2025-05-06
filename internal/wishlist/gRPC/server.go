package wishlist

import (
	"context"

	wl "github.com/DenisKokorin/WishListProto/gen/go/wishlist"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type WishList interface {
	CreateWishlist(ctx context.Context, userId int64, title string, isPrivate bool) (*wl.WishListResponse, error)
	GetWishList(ctx context.Context, userId int64, wishListId int64) (*wl.WishListResponse, error)
	UpdateWishList(ctx context.Context, userId int64, title string, isPrivate bool) (*wl.WishListResponse, error)
	DeleteWishList(ctx context.Context, userId int64, title string) (string, error)
	AddItem(ctx context.Context, userId int64, wishListId int64, items []*wl.ItemReq) (*wl.WishListResponse, error)
	UpdateItem(ctx context.Context, userId int64, wishListId int64, items []*wl.ItemReq) (*wl.WishListResponse, error)
	DeleteItem(ctx context.Context, userId int64, wishListId int64, title string) (*wl.WishListResponse, error)
}

type serverAPI struct {
	wishlist WishList
	wl.UnimplementedWishlistServiceServer
}

func Register(gRPC *grpc.Server, wishlist WishList) {
	wl.RegisterWishlistServiceServer(gRPC, &serverAPI{wishlist: wishlist})
}

func (s *serverAPI) CreateWishlist(ctx context.Context, req *wl.CreateWishlistRequest) (*wl.WishListResponse, error) {
	res, err := s.wishlist.CreateWishlist(ctx, req.GetUserId(), req.GetTitle(), req.GetIsPrivate())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return res, nil
}

func (s *serverAPI) GetWishList(ctx context.Context, req *wl.GetWishListRequest) (*wl.WishListResponse, error) {
	res, err := s.wishlist.GetWishList(ctx, req.GetUserId(), req.GetWishlistId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return res, nil
}

func (s *serverAPI) UpdateWishList(ctx context.Context, req *wl.UpdateWishListRequest) (*wl.WishListResponse, error) {
	res, err := s.wishlist.UpdateWishList(ctx, req.GetUserId(), req.GetTitle(), req.GetIsPrivate())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return res, nil
}

func (s *serverAPI) DeleteWishList(ctx context.Context, req *wl.DeleteWishListRequest) (*wl.DeleteWishListResponse, error) {
	res, err := s.wishlist.DeleteWishList(ctx, req.GetUserId(), req.GetTitle())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &wl.DeleteWishListResponse{Msg: res}, nil
}

func (s *serverAPI) AddItem(ctx context.Context, req *wl.AddItemRequest) (*wl.WishListResponse, error) {
	res, err := s.wishlist.AddItem(ctx, req.GetUserId(), req.GetWishlistId(), req.GetItems())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return res, nil
}

func (s *serverAPI) UpdateItem(ctx context.Context, req *wl.UpdateItemRequest) (*wl.WishListResponse, error) {
	res, err := s.wishlist.UpdateItem(ctx, req.GetUserId(), req.GetWishlistId(), req.GetItems())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return res, nil
}

func (s *serverAPI) DeleteItem(ctx context.Context, req *wl.DeleteItemRequest) (*wl.WishListResponse, error) {
	res, err := s.wishlist.DeleteItem(ctx, req.GetUserId(), req.GetWishlistId(), req.GetTitle())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return res, nil
}
