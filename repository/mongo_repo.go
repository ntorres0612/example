package repository

import (
	"context"
	"errors"
	"example/graph/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoRepository gestiona la conexión a MongoDB
type MongoRepository struct {
	client *mongo.Client
}

// NewMongoRepository inicializa la conexión a MongoDB
func NewMongoRepository(uri string) (*MongoRepository, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	return &MongoRepository{client: client}, nil
}

func (r *MongoRepository) CreateTodo(todo *model.Todo) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbName := "gqlgen-todo"
	db := r.client.Database(dbName)
	collection := db.Collection("todos")

	session, err := r.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		_, err := collection.InsertOne(sessCtx, todo)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	_, err = session.WithTransaction(ctx, callback)
	if err != nil {
		return errors.New("transacción fallida: " + err.Error())
	}

	return nil
}

func (r *MongoRepository) GetTodos() ([]*model.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbName := "gqlgen-todo"
	collection := r.client.Database(dbName).Collection("todos")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var todos []*model.Todo
	for cursor.Next(ctx) {
		var todo model.Todo
		if err := cursor.Decode(&todo); err != nil {
			return nil, err
		}
		todos = append(todos, &todo)
	}

	return todos, nil
}
