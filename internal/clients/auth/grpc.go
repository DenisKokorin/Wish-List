package clients

import (
	"context"
	"fmt"

	authv1 "github.com/Gergenus/Protobuf/gen/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CLient struct {
	api authv1.AuthClient
}

func New(ctx context.Context, addr string) (*CLient, error) {
	op := "clients.Auth.New"

	client, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &CLient{api: authv1.NewAuthClient(client)}, nil
}

func (c *CLient) GetUserId(ctx context.Context, username string) (int64, error) {
	op := "clients.Auth.GetUserId"

	res, err := c.api.GetUserId(ctx, &authv1.GetUserIdRequest{Username: username})
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return res.GetUserId(), nil
}
