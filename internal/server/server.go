package server

import (
	"context"
	"errors"
	"fmt"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"kekaton/back/internal/service"
	"kekaton/back/internal/storage"
)

var (
	ErrSuccess  = fiber.NewError(fiber.StatusOK, "success")
	ErrInternal = fiber.NewError(fiber.StatusInternalServerError, "something went wrong")
	ErrRequest  = fiber.NewError(fiber.StatusBadRequest, "invalid request")
	ErrData     = fiber.NewError(fiber.StatusBadRequest, "invalid data")
	ErrToken    = fiber.NewError(fiber.StatusUnauthorized, "invalid token")
)

type Config struct {
	Secret    []byte
	TokenName string
}

type Server struct {
	config  Config
	app     *fiber.App
	service *service.Service
}

func New(ctx context.Context, app *fiber.App, service *service.Service, config Config) *Server {
	srv := &Server{
		config:  config,
		app:     app,
		service: service,
	}

	srv.registerHandlers(ctx)

	return srv
}

func (s *Server) Listen(addr string, errc chan error) {
	errc <- s.app.Listen(addr)
}

func (s *Server) Shutdown() error {
	return s.app.Shutdown()
}

func (s *Server) registerHandlers(ctx context.Context) {
	handleJWT := jwtware.New(jwtware.Config{
		TokenLookup: fmt.Sprint("cookie:", s.config.TokenName),
		SigningKey: jwtware.SigningKey{
			JWTAlg: jwtware.HS256,
			Key:    s.config.Secret,
		},
		ErrorHandler: func(fcx *fiber.Ctx, err error) error {
			return ErrToken
		},
	})

	s.app.Use(func(fcx *fiber.Ctx) error {
		fcx.SetUserContext(ctx)

		return fcx.Next()
	}, cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowHeaders:     "Access-Control-Allow-Credentials",
		AllowCredentials: true,
	}))

	api := s.app.Group("/api")
	api.Get("/ping", func(fcx *fiber.Ctx) error {
		return fcx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "pong"})
	})

	v1 := api.Group("/v1")

	public := v1.Group("/public")
	public.Post("/sign-up", s.handleUserSignUp)
	public.Post("/sign-in", s.handleUserSignIn)
	public.Get("/user", s.handleGetUser)
	public.Get("/point", s.handleGetPoint)
	public.Get("/point", s.handleGetClosestPoint)
	public.Get("/points", s.handleGetPoints)
	public.Get("/comment", s.handleGetComment)
	public.Get("/comments", s.handleGetCommentsForPoint)
	public.Get("/tag", s.handleGetTag)
	public.Get("/tags", s.handleGetTags)

	private := v1.Group("/private", handleJWT, s.handleAuth)
	private.Post("/sign-out", s.handleUserSignOut)
	private.Patch("/user", s.handleUpdateUser)
	private.Post("/point", s.handleAddPoint)
	private.Patch("/point", s.handleUpdatePoint)
	private.Post("/comment", s.handleAddComment)
	private.Patch("/comment", s.handleUpdateComment)
}

func (s *Server) handleAuth(fcx *fiber.Ctx) error {
	token, err := s.ValidateJWT(fcx)
	if err != nil {
		return ErrToken
	}

	meta, err := s.ExtractJWTMetadata(token)
	if err != nil {
		return ErrToken
	}

	user := storage.User{
		ID: meta.UserID,
	}

	if err = s.service.GetUserByID(fcx.UserContext(), &user); err != nil {
		return ErrData
	}

	fcx.Locals("user", user)

	return fcx.Next()
}

func ErrorHandler(fcx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	var e *fiber.Error

	if errors.As(err, &e) {
		code = e.Code
	}

	return fcx.Status(code).JSON(fiber.Map{"message": err.Error()})
}
