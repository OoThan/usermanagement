package ds

import (
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

var DB *gorm.DB

type DataSource struct {
	DB  *gorm.DB
	RDB *redis.Client
	MDB *mongo.Client
}

func NewDataSource() (*DataSource, error) {
	db, err := LoadDB()
	if err != nil {
		return nil, err
	}
	DB = db

	rdb, err := LoadRDB()
	if err != nil {
		return nil, err
	}

	mdb, err := LoadMongo()
	if err != nil {
		return nil, err
	}

	return &DataSource{
		DB:  db,
		RDB: rdb,
		MDB: mdb,
	}, nil
}
