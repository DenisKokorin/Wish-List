package wishlistservice

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/DenisKokorin/Wish-List/internal/domain/models"
	wl "github.com/DenisKokorin/WishListProto/gen/go"
)

type WishList struct {
	log         *slog.Logger
	wlStorage   WlStorage
	itemStorage ItemStorage
}

type WlStorage interface {
	SaveWishList(ctx context.Context, userId int64, title string, isPrivate bool) (models.WishList, error)
	UpdateWishList(ctx context.Context, userId int64, title string, isPrivate bool) (models.WishList, error)
	GetWishList(ctx context.Context, userId int64, wishListId int64) (models.WishList, error)
	DeleteWishList(ctx context.Context, userId int64, title string) error
}

type ItemStorage interface {
	SaveItem(ctx context.Context, wishListId int64, items []models.Item) error
	GetItem(ctx context.Context, wishListId int64) ([]models.ItemResponse, error)
	UpdateItem(ctx context.Context, wishListId int64, items []models.Item) error
	DeleteItem(ctx context.Context, wishListId int64, title string) error
}

func New(log *slog.Logger, wlStorage WlStorage, itemStorage ItemStorage) *WishList {
	return &WishList{
		log:         log,
		wlStorage:   wlStorage,
		itemStorage: itemStorage,
	}
}

func (w *WishList) CreateWishlist(ctx context.Context, userId int64, title string, isPrivate bool) (*wl.WishListResponse, error) {
	const op = "WL.SaveWl"
	log := w.log.With(slog.String("op", op), slog.Int64("uid", userId))
	log.Info("saving new wishlist")

	res, err := w.wlStorage.SaveWishList(ctx, userId, title, isPrivate)
	if err != nil {
		return &wl.WishListResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("wishlist saved", slog.String("title", title))

	return &wl.WishListResponse{Id: res.Id, Title: res.Title, IsPrivate: res.IsPrivate}, nil
}

func (w *WishList) GetWishList(ctx context.Context, userId int64, wishListId int64) (*wl.WishListResponse, error) {
	const op = "WL.GetWl"
	log := w.log.With(slog.String("op", op), slog.Int64("uid", userId))

	log.Info("getting wishlist")

	wishList, err := w.wlStorage.GetWishList(ctx, userId, wishListId)
	if err != nil {
		return &wl.WishListResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	items, err := w.itemStorage.GetItem(ctx, wishListId)
	if err != nil {
		return &wl.WishListResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	return toWishlistResponse(wishList, items), nil
}

func (w *WishList) UpdateWishList(ctx context.Context, userId int64, title string, isPrivate bool) (*wl.WishListResponse, error) {
	const op = "WL.UpdateWL"
	log := w.log.With(slog.String("op", op), slog.Int64("uid", userId))
	log.Info("updating wishlist")

	wishList, err := w.wlStorage.UpdateWishList(ctx, userId, title, isPrivate)
	if err != nil {
		return &wl.WishListResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	items, err := w.itemStorage.GetItem(ctx, wishList.Id)
	if err != nil {
		return &wl.WishListResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("wishlist updated", slog.String("title", title))

	return toWishlistResponse(wishList, items), nil
}

func (w *WishList) DeleteWishList(ctx context.Context, userId int64, title string) (string, error) {
	const op = "WL.DeleteWL"
	log := w.log.With(slog.String("op", op), slog.Int64("uid", userId))
	log.Info("deleting wishlist")

	err := w.wlStorage.DeleteWishList(ctx, userId, title)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("wishlist deleted", slog.String("title", title))

	return "deleted", nil
}

func (w *WishList) AddItem(ctx context.Context, userId int64, wishListId int64, items []*wl.ItemReq) (*wl.WishListResponse, error) {
	const op = "WL.AddItem"
	log := w.log.With(slog.String("op", op), slog.Int64("uid", userId))

	log.Info("adding new items")

	modelItems := make([]models.Item, 0, len(items))
	for _, item := range items {
		modelItems = append(modelItems, models.Item{Title: item.GetTitle(), Description: item.GetDescription()})
	}

	err := w.itemStorage.SaveItem(ctx, wishListId, modelItems)
	if err != nil {
		return &wl.WishListResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("added new items", slog.Int64("WLID", wishListId))

	wishList, err := w.wlStorage.GetWishList(ctx, userId, wishListId)
	if err != nil {
		return &wl.WishListResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	respItems, err := w.itemStorage.GetItem(ctx, wishList.Id)
	if err != nil {
		return &wl.WishListResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	return toWishlistResponse(wishList, respItems), nil
}

func (w *WishList) UpdateItem(ctx context.Context, userId int64, wishListId int64, items []*wl.ItemReq) (*wl.WishListResponse, error) {
	const op = "WL.UpdateItem"
	log := w.log.With(slog.String("op", op), slog.Int64("uid", userId))

	log.Info("updating items")

	modelItems := make([]models.Item, 0, len(items))
	for _, item := range items {
		modelItems = append(modelItems, models.Item{Title: item.GetTitle(), Description: item.GetDescription()})
	}

	err := w.itemStorage.UpdateItem(ctx, wishListId, modelItems)
	if err != nil {
		return &wl.WishListResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("items updated", slog.Int64("WLID", wishListId))

	wishList, err := w.wlStorage.GetWishList(ctx, userId, wishListId)
	if err != nil {
		return &wl.WishListResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	respItems, err := w.itemStorage.GetItem(ctx, wishList.Id)
	if err != nil {
		return &wl.WishListResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	return toWishlistResponse(wishList, respItems), nil
}

func (w *WishList) DeleteItem(ctx context.Context, userId int64, wishListId int64, title string) (*wl.WishListResponse, error) {
	const op = "WL.Delete Item"
	log := w.log.With(slog.String("op", op), slog.Int64("uid", userId))

	log.Info("deleting items")

	err := w.itemStorage.DeleteItem(ctx, wishListId, title)
	if err != nil {
		return &wl.WishListResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("items deleted", slog.Int64("WLID", wishListId))

	wishList, err := w.wlStorage.GetWishList(ctx, userId, wishListId)
	if err != nil {
		return &wl.WishListResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	respItems, err := w.itemStorage.GetItem(ctx, wishList.Id)
	if err != nil {
		return &wl.WishListResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	return toWishlistResponse(wishList, respItems), nil
}

func toWishlistResponse(w models.WishList, items []models.ItemResponse) *wl.WishListResponse {
	itemResp := make([]*wl.ItemResp, 0, len(items))
	for _, item := range items {
		itemResp = append(itemResp, &wl.ItemResp{Id: item.Id, Title: item.Title, Description: item.Description})
	}

	return &wl.WishListResponse{Id: w.Id, Title: w.Title, IsPrivate: w.IsPrivate, Items: itemResp}
}
