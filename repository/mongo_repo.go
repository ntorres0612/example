package repository

import (
	"context"
	"errors"
	"time"
	"user-backend/graph/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoRepository gestiona la conexión a MongoDB
type MongoRepository struct {
	client     *mongo.Client
	dbName     string
	collection *mongo.Collection
}

func NewMongoRepository(uri string) (*MongoRepository, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	dbName := "test-db"
	collection := client.Database(dbName).Collection("users")

	return &MongoRepository{
		client:     client,
		dbName:     dbName,
		collection: collection,
	}, nil
}

func (r *MongoRepository) CreateUser(user *model.User) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	session, err := r.client.StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)

	var createdUser *model.User

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		result, err := r.collection.InsertOne(sessCtx, user)
		if err != nil {
			return nil, err
		}

		user.ID = result.InsertedID.(primitive.ObjectID).Hex()
		createdUser = user

		return nil, nil
	}

	_, err = session.WithTransaction(ctx, callback)
	if err != nil {
		return nil, errors.New("transacción fallida: " + err.Error())
	}

	return createdUser, nil
}

func (r *MongoRepository) GetUsers() ([]*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var Users []*model.User
	for cursor.Next(ctx) {
		var User model.User
		if err := cursor.Decode(&User); err != nil {
			return nil, err
		}
		Users = append(Users, &User)
	}

	return Users, nil
}

func (r *MongoRepository) GetUser(id string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var User model.User
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&User)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &User, nil
}

func (r *MongoRepository) UpdateUser(User *model.User) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(User.ID)
	if err != nil {
		return nil, err
	}

	update := bson.M{"$set": User}
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return nil, err
	}

	var updatedCustomer model.User
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&updatedCustomer)
	if err != nil {
		return nil, err
	}

	return &updatedCustomer, nil
}

func (r *MongoRepository) DeleteUser(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{"$set": bson.M{"status": false}}
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return err
	}

	return nil
}
