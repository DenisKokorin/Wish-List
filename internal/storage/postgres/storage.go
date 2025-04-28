package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/DenisKokorin/Wish-List/internal/domain/models"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error) {
	op := "Storage.pq.new"

	db, err := sql.Open("postgres", path)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveWishList(ctx context.Context, userId int64, title string, isPrivate bool) (models.WishList, error) {
	op := "Storage.pq.saveWL"

	var id int64
	err := s.db.QueryRowContext(ctx, "INSERT INTO wishlist(user_id, title, is_private) VALUES($1, $2, $3) RETURNING id", userId, title, isPrivate).Scan(&id)
	if err != nil {
		return models.WishList{}, fmt.Errorf("%s: %w", op, err)
	}

	return models.WishList{Id: id, Title: title, IsPrivate: isPrivate}, nil
}

func (s *Storage) UpdateWishList(ctx context.Context, userId int64, title string, isPrivate bool) (models.WishList, error) {
	op := "Storage.pq.saveWL"

	var id int64
	err := s.db.QueryRowContext(ctx, "UPDATE wishlist SET title = $1, is_private = $2 WHERE user_id = $3 RETURNING id", title, isPrivate, userId).Scan(&id)
	if err != nil {
		return models.WishList{}, fmt.Errorf("%s: %w", op, err)
	}

	return models.WishList{Id: id, Title: title, IsPrivate: isPrivate}, nil
}

func (s *Storage) DeleteWishList(ctx context.Context, userId int64, title string) error {
	op := "Storage.pq.deleteWL"

	_, err := s.db.ExecContext(ctx, "DELETE FROM wishlist WHERE user_id = $1 AND title = $2", userId, title)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetWishList(ctx context.Context, userId int64, wishListId int64) (models.WishList, error) {
	op := "Storage.pq.getWL"

	var id int64
	var title string
	var isPrivate bool
	err := s.db.QueryRowContext(ctx, "SELECT id, title, is_private FROM wishlist WHERE user_id = $1 AND id = $2", userId, wishListId).Scan(&id, &title, &isPrivate)
	if err != nil {
		return models.WishList{}, fmt.Errorf("%s: %w", op, err)
	}

	return models.WishList{Id: id, Title: title, IsPrivate: isPrivate}, nil
}

func (s *Storage) SaveItem(ctx context.Context, wishListId int64, items []models.Item) error {
	op := "Storage.pq.saveItem"

	for _, item := range items {
		_, err := s.db.ExecContext(ctx, "INSERT INTO item(wishlist_id, title, description) VALUES($1, $2, $3)", wishListId, item.Title, item.Description)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}

func (s *Storage) UpdateItem(ctx context.Context, wishListId int64, items []models.Item) error {
	op := "Storage.pq.updateItem"

	for _, item := range items {
		_, err := s.db.ExecContext(ctx, "UPDATE item SET title = $1, description = $2 WHERE wishlist_id = $3", item.Title, item.Description, wishListId)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}

func (s *Storage) DeleteItem(ctx context.Context, wishListId int64, title string) error {
	op := "Storage.pq.deleteItem"

	_, err := s.db.ExecContext(ctx, "DELETE FROM item WHERE wishlist_id = $1 AND title = $2", wishListId, title)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetItem(ctx context.Context, wishListId int64) ([]models.ItemResponse, error) {
	op := "Storage.pq.getItem"

	items := make([]models.ItemResponse, 0, 0)
	res, err := s.db.QueryContext(ctx, "SELECT id, title, description FROM item WHERE wishlist_id = $1", wishListId)
	if err != nil {
		return []models.ItemResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	for res.Next() {
		var id int64
		var title string
		var description string

		err := res.Scan(&id, &title, &description)
		if err != nil {
			return []models.ItemResponse{}, fmt.Errorf("%s: %w", op, err)
		}

		items = append(items, models.ItemResponse{Id: id, Title: title, Description: description})
	}

	if err := res.Err(); err != nil {
		return []models.ItemResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	return items, nil
}
