package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/tullur/lets-go-chat/internal/domain/chat/token"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrTokeNotFound = errors.New("token Not Found")
)

type MongoRepository struct {
	db      *mongo.Database
	token   *mongo.Collection
	timeout time.Duration
}

type tokenMongo struct {
	id           string `bson:"id"`
	userId       string `bson:"userId"`
	expiresAfter string `bson:"expiresAfter"`
}

func New(ctx context.Context, connectionURI string, timeout int) (*MongoRepository, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionURI))
	if err != nil {
		return nil, err
	}

	db := client.Database("lets-go-chat")
	tokens := db.Collection("tokens")

	return &MongoRepository{
		db:      db,
		token:   tokens,
		timeout: time.Duration(timeout) * time.Second,
	}, nil
}

func (r *MongoRepository) Get(id string) (*token.Token, error) {
	var t tokenMongo

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	result := r.token.FindOne(ctx, bson.M{"id": id})

	err := result.Decode(&t)
	if err != nil {
		return nil, err
	}

	return t.toAggregate(), nil
}

func (repository *MongoRepository) Add(token *token.Token) error {
	ctx, cancel := context.WithTimeout(context.Background(), repository.timeout)
	defer cancel()

	internal := NewFromToken(*token)
	_, err := repository.token.InsertOne(ctx, internal)
	if err != nil {
		return err
	}

	return nil
}

func (repository *MongoRepository) Revoke(id string) error {
	_, err := repository.token.DeleteOne(context.Background(), bson.M{"id": id})
	if err != nil {
		return ErrTokeNotFound
	}

	return nil
}

func NewFromToken(t token.Token) tokenMongo {
	return tokenMongo{
		id:           t.Id(),
		userId:       t.User(),
		expiresAfter: t.ExpiresAfter(),
	}
}

func (tm tokenMongo) toAggregate() *token.Token {
	t := &token.Token{}

	t.SetID(tm.id)
	t.SetUser(tm.userId)
	t.SetExpiresAfter(tm.expiresAfter)

	return t
}
