package message

import (
	"context"
	"log"
	"time"

	"github.com/tullur/lets-go-chat/internal/domain/chat"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	GetMessages(limit int) ([]*MessageMongo, error)
	Add(message chat.Message) error
}

type MongoRepository struct {
	db      *mongo.Database
	message *mongo.Collection
}

type MessageMongo struct {
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

func (repository *MongoRepository) GetMessages(limit int) ([]*MessageMongo, error) {
	opts := options.FindOptions{}
	opts.SetLimit(int64(limit))

	query := bson.M{}

	cursor, err := repository.message.Find(context.Background(), query, &opts)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.Background())

	var messages []*MessageMongo

	for cursor.Next(context.Background()) {
		msg := &MessageMongo{}
		log.Println(cursor.Decode(msg))
		err := cursor.Decode(msg)
		if err != nil {
			return nil, err
		}

		messages = append(messages, msg)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	if len(messages) == 0 {
		return []*MessageMongo{}, nil
	}

	return messages, nil
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

func fromModel(in chat.Message) MessageMongo {
	return MessageMongo{
		Sender:  in.Sender.LocalAddr().String(),
		Content: string(in.Content),
	}
}
