package app

import (
	"context"
	"time"

	"github.com/sivistrukov/vk-assigment/internal/entrypoints/http"
	"github.com/sivistrukov/vk-assigment/internal/infrastructure/postgresql"
	"github.com/sivistrukov/vk-assigment/internal/services/validator"
)

const ShutdownTimeout = time.Second * 10

func Run() error {
	cfg := NewConfig()

	closer := GetCloser()

	db, err := postgresql.NewConnection(cfg.Database)
	if err != nil {
		return err
	}
	initContainer(db)

	validate := validator.New()

	httpHandler := http.NewHandler(
		validate,
		container.AuthService(),
		container.UserService(),
		container.ActorService(),
		container.FilmService(),
	)

	srv := http.NewServer(cfg.Http, httpHandler)
	closer.Add(srv.Shutdown)

	closer.Add(func(_ context.Context) error { return db.Close() })

	exit := make(chan error)
	go func() {
		err := srv.ListenAndServe()

		exit <- err
	}()

	return <-exit
}

func Shutdown() error {
	ctx, cancel := context.WithTimeout(
		context.Background(), ShutdownTimeout,
	)
	defer cancel()

	closer := GetCloser()
	return closer.Close(ctx)
}
