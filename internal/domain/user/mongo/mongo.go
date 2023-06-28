package mongo

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/tullur/lets-go-chat/internal/domain/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	db   *mongo.Database
	user *mongo.Collection
}

type mongoUser struct {
	ID       string `bson:"id"`
	Name     string `bson:"name"`
	Password string `bson:"password"`
}

func NewFromUser(u user.User) mongoUser {
	return mongoUser{
		ID:       u.Id(),
		Name:     u.Name(),
		Password: u.Password(),
	}
}

func New(ctx context.Context, connectionURI string) (*MongoRepository, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionURI))
	if err != nil {
		return nil, err
	}

	db := client.Database("lets-go-chat")
	users := db.Collection("users")

	return &MongoRepository{
		db:   db,
		user: users,
	}, nil
}

func (u mongoUser) toAggregate() (*user.User, error) {
	user := &user.User{}

	id, err := uuid.Parse(u.ID)
	if err != nil {
		return nil, err
	}

	user.SetID(id)
	user.SetName(u.Name)
	user.SetPassword(u.Password)

	return user, nil
}

func (repository *MongoRepository) List() []user.User {
	var users []user.User

	cursor, err := repository.user.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	for cursor.Next(context.Background()) {
		var u user.User

		err := cursor.Decode(&u)
		if err != nil {
			log.Fatal(err)
		}

		log.Println(cursor)

		users = append(users, u)
	}

	return users
}

func (repository *MongoRepository) Create(user *user.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	internal := NewFromUser(*user)
	_, err := repository.user.InsertOne(ctx, internal)
	if err != nil {
		return err
	}

	return nil
}

func (repository *MongoRepository) FindByName(name string) (*user.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result := repository.user.FindOne(ctx, bson.M{"name": name})

	var u mongoUser
	err := result.Decode(&u)
	if err != nil {
		return nil, err
	}

	return u.toAggregate()
}
