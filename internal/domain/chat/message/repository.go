package message

import (
	"context"
	"time"

	"github.com/tullur/lets-go-chat/internal/domain/chat"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	Add(message chat.Message) error
}

type MongoRepository struct {
	db      *mongo.Database
	message *mongo.Collection
}

type messageMongo struct {
	Sender  string `bson:"sender"`
	Content string `bson:"content"`
}

func New(ctx context.Context, connectionURI string) (*MongoRepository, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionURI))
	if err != nil {
		return nil, err
	}

	db := client.Database("lets-go-chat")
	messages := db.Collection("messages")

	return &MongoRepository{
		db:      db,
		message: messages,
	}, nil
}

func (repository *MongoRepository) Add(message chat.Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := repository.message.InsertOne(ctx, fromModel(message))
	if err != nil {
		return err
	}

	return nil
}

func fromModel(in chat.Message) messageMongo {
	return messageMongo{
		Sender:  in.Sender.LocalAddr().String(),
		Content: string(in.Content),
	}
}
