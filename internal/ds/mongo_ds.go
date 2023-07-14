package ds

import (
	"context"
	"time"

	"github.com/OoThan/usermanagement/config"
	"github.com/OoThan/usermanagement/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func LoadMongo() (*mongo.Client, error) {
	uri := config.MongoDSN()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	logger.Sugar.Info("Successfully connect to mongodb.")

	return client, nil
}
