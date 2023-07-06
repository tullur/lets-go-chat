package message

import (
	"context"
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
	timeout time.Duration
}

type MessageMongo struct {
	Sender  string `bson:"sender"`
	Content string `bson:"content"`
}

func New(ctx context.Context, connectionURI string, timeout int) (*MongoRepository, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionURI))
	if err != nil {
		return nil, err
	}

	db := client.Database("lets-go-chat")
	messages := db.Collection("messages")

	return &MongoRepository{
		db:      db,
		message: messages,
		timeout: time.Duration(timeout) * time.Second,
	}, nil
}

func (r *MongoRepository) GetMessages(limit int) ([]*MessageMongo, error) {
	opts := options.FindOptions{}

	opts.SetMaxTime(r.timeout)
	opts.SetLimit(int64(limit))

	query := bson.M{}

	cursor, err := r.message.Find(context.Background(), query, &opts)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.Background())

	var messages []*MessageMongo

	for cursor.Next(context.Background()) {
		msg := &MessageMongo{}

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

func (r *MongoRepository) Add(message chat.Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	_, err := r.message.InsertOne(ctx, fromModel(message))
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
