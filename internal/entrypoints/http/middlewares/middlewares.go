package middlewares

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/sivistrukov/vk-assigment/internal/entrypoints/http/schemas"
	"github.com/sivistrukov/vk-assigment/internal/models"
)

type authService interface {
	Authenticate(string, string) (models.User, error)
}

type Key string

func PanicRecover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(r.Method, r.URL.Path, fmt.Sprintf("got error: %v", err))

				response, _ := json.Marshal(schemas.ErrorResponse{
					Error: "internal server error",
				})

				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write(response)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)

		next.ServeHTTP(w, r)
	})
}

func BasicAuth(next http.Handler, auth authService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			unauthorized(w)
			return
		}

		credentials := strings.SplitN(authHeader, " ", 2)
		if len(credentials) != 2 || credentials[0] != "Basic" {
			unauthorized(w)
			return
		}

		payload, err := base64.StdEncoding.DecodeString(credentials[1])
		if err != nil {
			unauthorized(w)
			return
		}

		pair := strings.SplitN(string(payload), ":", 2)

		user, err := auth.Authenticate(pair[0], pair[1])
		if err != nil {
			unauthorized(w)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, Key("user"), user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AdminRoutes(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user := ctx.Value(Key("user"))
		if user == nil {
			forbidden(w)
			return
		}

		if !user.(models.User).IsAdmin {
			forbidden(w)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func unauthorized(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
	w.WriteHeader(http.StatusUnauthorized)
	response, _ := json.Marshal(schemas.ErrorResponse{
		Error: "401 unauthorized",
	})
	_, _ = w.Write(response)
}

func forbidden(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
	response, _ := json.Marshal(schemas.ErrorResponse{
		Error: "403 forbidden",
	})
	_, _ = w.Write(response)
}
