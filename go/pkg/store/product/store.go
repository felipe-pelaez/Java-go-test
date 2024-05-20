package product

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

// ? =================== Interfaces =================== ?

type IProductStore interface {
	InitDatabase(collection string) (*mongo.Collection, error)
}

// ? =================== Structs =================== ?

type Store struct{}

// ? =================== Constructors =================== ?

func NewStore() IProductStore {
	return &Store{}
}

// ? =================== Functions =================== ?

func (s *Store) InitDatabase(collection string) (*mongo.Collection, error) {

	var dsn string

	usernameDb := viper.GetString("database.username")
	passwordDb := viper.GetString("database.password")

	hostDb := viper.GetString("database.host")
	if hostDb == "" {
		hostDb = "localhost"
	}

	portDb := viper.GetString("database.port")
	if portDb == "" {
		portDb = "27017"
	}

	nameDb := viper.GetString("database.name")
	if nameDb == "" {
		nameDb = "microservice_go_template"
	}

	if usernameDb == "" || passwordDb == "" {
		dsn = fmt.Sprintf("mongodb://%s:%s", hostDb, portDb)
	} else {
		dsn = fmt.Sprintf("mongodb://%s:%s@%s:%s", usernameDb, passwordDb, hostDb, portDb)
	}

	db, err := mongo.NewClient(options.Client().ApplyURI(dsn))
	if err != nil {
		return nil, err
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = db.Connect(ctx)

	if err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB!")

	return db.Database(nameDb).Collection(collection), nil

}
