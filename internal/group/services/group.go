package groupservice

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	clients "github.com/DenisKokorin/Wish-List/internal/clients/auth"
	"github.com/DenisKokorin/Wish-List/internal/domain/models"
	gr "github.com/DenisKokorin/WishListProto/gen/go/group"
)

type Group struct {
	log          *slog.Logger
	groupStorage GroupStorage
	authClient   *clients.CLient
}

type GroupStorage interface {
	SaveGroup(ctx context.Context, userId int64, title string) (*models.GroupWithMembers, error)
	GetGroup(ctx context.Context, groupId int64) (*models.GroupWithMembers, error)
	GetAllGroup(ctx context.Context, userId int64) ([]*models.Group, error)
	DeleteGroup(ctx context.Context, groupId int64) error
	AddMember(ctx context.Context, groupId int64, userId int64) error
	DeleteMember(ctx context.Context, userId int64, groupId int64) error
}

func New(log *slog.Logger, groupStorage GroupStorage, authClient *clients.CLient) *Group {
	return &Group{
		log:          log,
		groupStorage: groupStorage,
		authClient:   authClient,
	}
}

func (g *Group) CreateGroup(ctx context.Context, userId int64, title string) (*gr.GroupResponse, error) {
	const op = "GR.CreateGroup"
	log := g.log.With(slog.String("op", op), slog.Int64("uid", userId))

	log.Info("creating group")

	res, err := g.groupStorage.SaveGroup(ctx, userId, title)
	if err != nil {
		return &gr.GroupResponse{}, fmt.Errorf("%s: %w", op, err)

	}

	return toGroupResponse(*res), nil
}

func (g *Group) GetAllGroups(ctx context.Context, userId int64) (*gr.AllGroupsResponse, error) {
	const op = "GR.GetAllgroups"
	log := g.log.With(slog.String("op", op), slog.Int64("uid", userId))

	log.Info("getting groups")

	res, err := g.groupStorage.GetAllGroup(ctx, userId)
	if err != nil {
		return &gr.AllGroupsResponse{}, fmt.Errorf("%s: %w", op, err)

	}

	groups := make([]*gr.Group, 0, len(res))
	for _, group := range res {
		groups = append(groups, toGroup(*group))
	}

	return &gr.AllGroupsResponse{Groups: groups}, nil
}

func (g *Group) GetGroup(ctx context.Context, groupId int64) (*gr.GroupResponse, error) {
	const op = "GR.GetGroup"
	log := g.log.With(slog.String("op", op), slog.Int64("groupId", groupId))

	log.Info("getting group by id")

	res, err := g.groupStorage.GetGroup(ctx, groupId)
	if err != nil {
		return &gr.GroupResponse{}, fmt.Errorf("%s: %w", op, err)

	}

	return toGroupResponse(*res), nil
}

func (g *Group) DeleteGroup(ctx context.Context, groupId int64) (string, error) {
	const op = "GR.DeleteGroup"
	log := g.log.With(slog.String("op", op), slog.Int64("groupId", groupId))

	log.Info("deleting group")

	err := g.groupStorage.DeleteGroup(ctx, groupId)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)

	}

	return "deleted", nil
}

func (g *Group) InviteUser(ctx context.Context, groupId int64, username string) (string, error) {
	const op = "GR.InviteUser"
	log := g.log.With(slog.String("op", op), slog.Int64("groupId", groupId))

	log.Info("inviting user")

	id, err := g.authClient.GetUserId(ctx, username)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	group, err := g.groupStorage.GetGroup(ctx, groupId)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	for _, member := range group.Members {
		if id == member {
			return "", errors.New("user already in group")
		}
	}

	secret := os.Getenv("SECRET")
	token, err := GenerateInviteToken(id, groupId, secret)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (g *Group) AcceptInvite(ctx context.Context, token string) (string, error) {
	const op = "GR.AcceptInvite"

	secret := os.Getenv("SECRET")
	invite, err := ValidateInviteToken(token, secret)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	err = g.groupStorage.AddMember(ctx, invite.GroupID, invite.UserID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return "invite accepted", nil
}

func (g *Group) LeaveGroup(ctx context.Context, userId int64, groupId int64) (string, error) {
	const op = "GR.LeavaGroup"
	log := g.log.With(slog.String("op", op), slog.Int64("groupId", groupId))

	log.Info("user leaving group")

	err := g.groupStorage.DeleteMember(ctx, userId, groupId)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)

	}

	return "left the group", nil
}

func toGroup(res models.Group) *gr.Group {
	return &gr.Group{
		GroupId: res.Id,
		Title:   res.Title,
	}
}

func toGroupResponse(res models.GroupWithMembers) *gr.GroupResponse {
	title := models.Group{Id: res.Id, Title: res.Title}
	return &gr.GroupResponse{
		Title:   toGroup(title),
		Members: res.Members,
	}
}

func GenerateInviteToken(userId int64, groupId int64, secret string) (string, error) {
	expiresAt := time.Now().Add(7 * 24 * time.Hour).Unix()

	payload := models.InvitationToken{UserID: userId, GroupID: groupId, ExpiresAt: expiresAt}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	payloadBase64 := base64.URLEncoding.EncodeToString(payloadBytes)

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(payloadBase64))
	signature := base64.URLEncoding.EncodeToString(h.Sum(nil))

	return payloadBase64 + "." + signature, nil
}

func ValidateInviteToken(token string, secret string) (*models.InvitationToken, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return nil, errors.New("invalid token format")
	}

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(parts[0]))
	expectedSig := base64.URLEncoding.EncodeToString(h.Sum(nil))

	if parts[1] != expectedSig {
		return nil, errors.New("invalid signature")
	}

	payloadBytes, err := base64.URLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, err
	}

	var payload models.InvitationToken
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return nil, err
	}

	if time.Now().Unix() > payload.ExpiresAt {
		return nil, errors.New("token expired")
	}

	return &payload, nil
}
