package client

import (
	"context"
	"log/slog"

	authProto "github.com/AmadoMuerte/BirthdayWish/API/proto/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	client authProto.AuthServiceClient
	conn   *grpc.ClientConn
	log    *slog.Logger
}

func NewAuthClient(addr string, log *slog.Logger) (*AuthClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &AuthClient{
		client: authProto.NewAuthServiceClient(conn),
		conn:   conn,
		log:    log,
	}, nil
}

func (c *AuthClient) Close() error {
	return c.conn.Close()
}

func (c *AuthClient) SignUp(ctx context.Context, username, email, password string) (*authProto.SignUpResponse, error) {
	c.log.Debug("Calling auth service SignUp", "username", username, "email", email)

	resp, err := c.client.SignUp(ctx, &authProto.SignUpRequest{
		Username: username,
		Email:    email,
		Password: password,
	})

	if err != nil {
		c.log.Error("Auth service SignUp failed", "error", err)
		return nil, err
	}

	return resp, nil
}

func (c *AuthClient) SignIn(ctx context.Context, username, password string) (*authProto.SignInResponse, error) {
	c.log.Debug("Calling auth service SignIn", "username", username)

	resp, err := c.client.SignIn(ctx, &authProto.SignInRequest{
		Username: username,
		Password: password,
	})

	if err != nil {
		c.log.Error("Auth service SignIn failed", "error", err)
		return nil, err
	}

	return resp, nil
}
