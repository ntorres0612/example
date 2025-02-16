package repository

import "example/graph/model"

// Database define los métodos que cualquier base de datos debe implementar
type Database interface {
	CreateTodo(todo *model.Todo) error
	GetTodos() ([]*model.Todo, error)
}
