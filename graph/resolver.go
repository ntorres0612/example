package graph

import "user-backend/graph/service"

//go:generate go run github.com/99designs/gqlgen
type Resolver struct {
	Service *service.UserService
}
