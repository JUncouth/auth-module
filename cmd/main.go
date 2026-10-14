package main

import (
	authpb "auth/src/proto"
	"context"
	"fmt"
	"net"
	"time"

	auth "auth/src"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewPostgresPool(
	cfg auth.Config,
) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	config, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse database config: %w", err)
	}
	pool, err := pgxpool.NewWithConfig(
		ctx,
		config,
	)
	if err != nil {
		return nil, fmt.Errorf("create database pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()

		return nil, fmt.Errorf("ping database: %w", err)
	}

	return pool, nil
}

func newGRPCServer(
	service auth.AuthService,
	cfg auth.Config,
) (*grpc.Server, error) {
	server := grpc.NewServer()

	grpcService := auth.NewGRPCService(service)

	authpb.RegisterAuthorizationServiceServer(
		server,
		grpcService,
	)

	reflection.Register(server)

	return server, nil
}

func runGRPCServer(
	lifecycle fx.Lifecycle,
	server *grpc.Server,
	cfg auth.Config,
) error {
	listener, err := net.Listen(
		"tcp",
		cfg.GRPCAddress,
	)
	if err != nil {
		return err
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				_ = server.Serve(listener)
			}()

			return nil
		},

		OnStop: func(ctx context.Context) error {
			done := make(chan struct{})

			go func() {
				server.GracefulStop()
				close(done)
			}()

			select {
			case <-done:
				return nil
			case <-ctx.Done():
				server.Stop()
				return ctx.Err()
			}
		},
	})

	return nil
}

func main() {
	fx.New(
		fx.Provide(
			auth.LoadConfig,
			NewPostgresPool,
			auth.NewUserRepository,
			auth.NewPasswordHasher,
			auth.NewTokenManager,
			auth.NewAuthService,
			newGRPCServer,
		),

		fx.Invoke(
			runGRPCServer,
		),

		fx.Invoke(func(lc fx.Lifecycle, pool *pgxpool.Pool) {
			lc.Append(fx.Hook{
				OnStop: func(ctx context.Context) error {
					pool.Close()

					return nil
				},
			})
		}),
	).Run()
}
