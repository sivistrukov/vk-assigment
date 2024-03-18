package app

import (
	"database/sql"
	"sync"

	"github.com/sivistrukov/vk-assigment/internal/infrastructure/postgresql"
	"github.com/sivistrukov/vk-assigment/internal/services/actors"
	"github.com/sivistrukov/vk-assigment/internal/services/auth"
	"github.com/sivistrukov/vk-assigment/internal/services/films"
	"github.com/sivistrukov/vk-assigment/internal/services/users"
)

var (
	onceContainer sync.Once
	container     *Container
)

type Container struct {
	psqlConn *sql.DB
}

func GetContainer() *Container {
	onceContainer.Do(func() {
		container = &Container{}
	})

	return container
}

func initContainer(conn *sql.DB) *Container {
	onceContainer.Do(func() {
		container = &Container{psqlConn: conn}
	})

	return container
}

func (c *Container) UserRepo() *postgresql.UserRepo {
	return postgresql.NewUserRepo(c.psqlConn)
}

func (c *Container) ActorRepo() *postgresql.ActorRepo {
	return postgresql.NewActorRepo(c.psqlConn)
}

func (c *Container) FilmRepo() *postgresql.FilmRepo {
	return postgresql.NewFilmRepo(c.psqlConn)
}

func (c *Container) AuthService() *auth.Service {
	return auth.NewService(c.UserRepo())
}

func (c *Container) ActorService() *actors.Service {
	return actors.NewService(c.ActorRepo())
}

func (c *Container) FilmService() *films.Service {
	return films.NewService(c.FilmRepo())
}

func (c *Container) UserService() *users.Service {
	return users.NewService(c.UserRepo())
}
