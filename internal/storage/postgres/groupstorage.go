package postgres

import (
	"context"
	"fmt"

	"github.com/DenisKokorin/Wish-List/internal/domain/models"
)

func (s *Storage) SaveGroup(ctx context.Context, userId int64, title string) (*models.GroupWithMembers, error) {
	op := "Storage.pq.saveGroup"

	var id int64
	err := s.db.QueryRowContext(ctx, "INSERT INTO groups(user_id, title) VALUES($1, $2) RETURNING id", userId, title).Scan(&id)
	if err != nil {
		return &models.GroupWithMembers{}, fmt.Errorf("%s: %w", op, err)
	}

	return &models.GroupWithMembers{Id: id, Title: title, Members: []int64{userId}}, nil
}

func (s *Storage) GetGroup(ctx context.Context, groupId int64) (*models.GroupWithMembers, error) {
	op := "Sorage.pq.getGroup"

	var id int64
	var title string
	err := s.db.QueryRowContext(ctx, "SELECT id, title FROM groups WHERE group_id = $1", groupId).Scan(&id, &title)
	if err != nil {
		return &models.GroupWithMembers{}, fmt.Errorf("%s: %w", op, err)
	}

	members := make([]int64, 0, 0)
	res, err := s.db.QueryContext(ctx, "SELECT user_id FROM group_members WHERE group_id = $1", groupId)
	if err != nil {
		return &models.GroupWithMembers{}, fmt.Errorf("%s: %w", op, err)
	}

	for res.Next() {
		var id int

		err := res.Scan(&id)
		if err != nil {
			return &models.GroupWithMembers{}, fmt.Errorf("%s: %w", op, err)
		}

		members = append(members, int64(id))
	}

	if err := res.Err(); err != nil {
		return &models.GroupWithMembers{}, fmt.Errorf("%s: %w", op, err)
	}

	return &models.GroupWithMembers{Id: id, Title: title, Members: members}, nil
}

func (s *Storage) GetAllGroup(ctx context.Context, userId int64) ([]*models.Group, error) {
	op := "Storage.pq.getAllGroups"

	groups := make([]*models.Group, 0, 0)
	res, err := s.db.QueryContext(ctx, "SELECT id, title FROM groups WHERE user_id = $1", userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for res.Next() {
		var id int64
		var title string

		err := res.Scan(&id, &title)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		groups = append(groups, &models.Group{Id: id, Title: title})
	}

	if err := res.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return groups, nil
}

func (s *Storage) DeleteGroup(ctx context.Context, groupId int64) error {
	op := "Storage.pq.deleteGroup"

	_, err := s.db.ExecContext(ctx, "DELETE FROM groups WHERE group_id = $1", groupId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) AddMember(ctx context.Context, groupId int64, userId int64) error {
	op := "Storage.pq.addMember"

	_, err := s.db.ExecContext(ctx, "INSERT INTO group_member(group_id, user_id) VALUES($1, $2) RETURNING id", groupId, userId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteMember(ctx context.Context, userId int64, groupId int64) error {
	op := "Storage.pq.deleteMember"

	_, err := s.db.ExecContext(ctx, "DELETE FROM group_member WHERE group_id = $1 AND user_id = $2", groupId, userId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
