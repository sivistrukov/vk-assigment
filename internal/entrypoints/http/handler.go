package http

import (
	"net/http"

	"github.com/go-playground/validator"
	_ "github.com/sivistrukov/vk-assigment/docs"
	mw "github.com/sivistrukov/vk-assigment/internal/entrypoints/http/middlewares"
	v1 "github.com/sivistrukov/vk-assigment/internal/entrypoints/http/v1"
	httpSwag "github.com/swaggo/http-swagger/v2"
)

func NewHandler(
	validator *validator.Validate,
	authService authService,
	userService v1.UserService,
	actorsService v1.ActorService,
	filmsService v1.FilmService,
) http.Handler {
	mux := http.NewServeMux()

	// users
	usersHandlers := v1.NewUsersHandler(userService, validator)
	mux.Handle("POST /api/v1/users", usersHandlers.Create())

	// actors
	actorsHandler := v1.NewActorHandler(actorsService, validator)
	readerActorsMux := http.NewServeMux()
	readerActorsMux.Handle("/", actorsHandler.GetList())

	adminActorsMux := http.NewServeMux()
	adminActorsMux.Handle("POST /", actorsHandler.Add())
	adminActorsMux.Handle("PUT /{id}", actorsHandler.Update())
	adminActorsMux.Handle("PATCH /{id}", actorsHandler.PartialUpdate())
	adminActorsMux.Handle("DELETE /{id}", actorsHandler.Remove())

	adminActorRouter := mw.BasicAuth(mw.AdminRoutes(adminActorsMux), authService)
	readerActorsRouter := mw.BasicAuth(readerActorsMux, authService)
	mux.Handle("GET /api/v1/actors", readerActorsRouter)
	mux.Handle("/api/v1/actors", adminActorRouter)

	//films
	filmsHandler := v1.NewFilmsHandler(filmsService, validator)
	readerFilmsMux := http.NewServeMux()
	readerFilmsMux.Handle("/", filmsHandler.GetList())

	adminFilmsMux := http.NewServeMux()
	adminFilmsMux.Handle("POST /", filmsHandler.Add())
	adminFilmsMux.Handle("PUT /{id}", filmsHandler.Update())
	adminFilmsMux.Handle("PATCH /{id}", filmsHandler.PartialUpdate())
	adminFilmsMux.Handle("DELETE /{id}", filmsHandler.Remove())

	readerFilmsRouter := mw.BasicAuth(readerFilmsMux, authService)
	adminFilmsRouter := mw.BasicAuth(mw.AdminRoutes(adminFilmsMux), authService)
	mux.Handle("GET /api/v1/films", readerFilmsRouter)
	mux.Handle("/api/v1/films", adminFilmsRouter)

	mux.Handle("/swagger/", httpSwag.Handler(
		httpSwag.URL("http://localhost:8080/swagger/doc.json"),
		httpSwag.DeepLinking(true),
		httpSwag.DocExpansion("list"),
		httpSwag.DomID("swagger-ui"),
	))

	handler := mw.Logging(mw.PanicRecover(mux))

	return handler
}
