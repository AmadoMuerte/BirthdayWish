package client

import (
	"context"
	"log/slog"

	wishProto "github.com/AmadoMuerte/BirthdayWish/API/proto/wish"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type WishlisterClient struct {
	client wishProto.WishServiceClient
	conn   *grpc.ClientConn
	log    *slog.Logger
}

func NewWishlisterClient(addr string, log *slog.Logger) (*WishlisterClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &WishlisterClient{
		client: wishProto.NewWishServiceClient(conn),
		conn:   conn,
		log:    log,
	}, nil
}

func (c *WishlisterClient) Close() error {
	return c.conn.Close()
}

func (c *WishlisterClient) CreateWish(ctx context.Context, req *wishProto.CreateWishRequest) (*wishProto.WishResponse, error) {
	return c.client.CreateWish(ctx, req)
}

func (c *WishlisterClient) GetWish(ctx context.Context, req *wishProto.GetWishRequest) (*wishProto.WishResponse, error) {
	c.log.Debug("Calling wishlister service GetWish", "user_id", req.UserId, "wish_id", req.WishId)

	resp, err := c.client.GetWish(ctx, req)
	if err != nil {
		c.log.Error("Wishlister service GetWish failed", "error", err)
		return nil, err
	}

	return resp, nil
}

func (c *WishlisterClient) UpdateWish(ctx context.Context, req *wishProto.UpdateWishRequest) (*wishProto.WishResponse, error) {
	return c.client.UpdateWish(ctx, req)
}

func (c *WishlisterClient) DeleteWish(ctx context.Context, req *wishProto.DeleteWishRequest) (*wishProto.EmptyResponse, error) {
	return c.client.DeleteWish(ctx, req)
}

func (c *WishlisterClient) ListWishes(ctx context.Context, req *wishProto.ListWishesRequest) (*wishProto.WishListResponse, error) {
	return c.client.ListWishes(ctx, req)
}

func (c *WishlisterClient) CreateShareLink(ctx context.Context, req *wishProto.CreateShareLinkRequest) (*wishProto.ShareLinkResponse, error) {
	return c.client.CreateShareLink(ctx, req)
}

func (c *WishlisterClient) GetSharedWishes(ctx context.Context, req *wishProto.GetSharedRequest) (*wishProto.WishListResponse, error) {
	return c.client.GetSharedWishes(ctx, req)
}
