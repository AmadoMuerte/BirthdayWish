package client

import (
	"context"
	"log/slog"

	filerProto "github.com/AmadoMuerte/BirthdayWish/API/proto/filer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type FilerClient struct {
	client filerProto.FilerServiceClient
	conn   *grpc.ClientConn
	log    *slog.Logger
}

func NewFilerClient(addr string, log *slog.Logger) (*FilerClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &FilerClient{
		client: filerProto.NewFilerServiceClient(conn),
		conn:   conn,
		log:    log,
	}, nil
}

func (c *FilerClient) Close() error {
	return c.conn.Close()
}

func (c *FilerClient) LoadImage(ctx context.Context, imageBase64 string) (*filerProto.LoadImageResponse, error) {
	resp, err := c.client.LoadImage(ctx, &filerProto.LoadImageRequest{
		ImageBase64: imageBase64,
	})

	return resp, err
}
