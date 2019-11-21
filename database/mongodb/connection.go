package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

const connectTimeout = 10 * time.Second
const pingDBTimeout = 10 * time.Second

type MongoDBConfig struct {
	ConnectionString string
	DBName           string
}

type MongoDBConnection struct {
	config   MongoDBConfig
	client   *mongo.Client
	database *mongo.Database
}

func NewMongoDBConnection(config MongoDBConfig) *MongoDBConnection {
	return &MongoDBConnection{config: config}
}

func (m *MongoDBConnection) Connect() {
	// Connect creates a new Client and then initializes it using the Connect method.
	connectCtx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	client, err := mongo.Connect(connectCtx, options.Client().ApplyURI(m.config.ConnectionString))
	if err != nil {
		log.Fatal("Can not connect creates a new Client and then initializes it using the Connect method. Error: " + err.Error())
	}
	defer cancel()

	// Ping verifies that the client can connect to the topology.
	pingCtx, _ := context.WithTimeout(context.Background(), pingDBTimeout)
	err = client.Ping(pingCtx, readpref.Primary())
	if err != nil {
		log.Fatal("Can not ping verifies that the client can connect to the topology. Error: " + err.Error())
	}

	m.client = client
	m.database = client.Database(m.config.DBName)
}

func (m *MongoDBConnection) Client() *mongo.Client {
	return m.client
}

func (m *MongoDBConnection) Database() *mongo.Database {
	return m.database
}
