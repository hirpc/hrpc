package mongodb

import (
	"context"
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Option struct {
	Address  string `json:"address"`
	DB       int    `json:"db"`
	Username string `json:"username"`
	Password string `json:"password"`
	Port     int    `json:"port"`
	Network  string `json:"network"`
}

type MongoDB struct {
	conn   *mongo.Client
	option Option
}

var m *MongoDB

func (m *MongoDB) Load(src []byte) error {
	if err := json.Unmarshal(src, &m.option); err != nil {
		return err
	}
	return nil
}

func (m *MongoDB) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return err
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	return nil
}

func (m *MongoDB) Destory() {
	if m.conn != nil {
		m.conn.Disconnect(context.Background())
	}
}

func Client() *mongo.Client {
	return m.conn
}

func (m MongoDB) Name() string {
	return "mongodb"
}

// Valid returns a bool valud to determine whether the connection is ready to use
func Valid() bool {
	if m == nil {
		return false
	}
	if m.conn == nil {
		return false
	}
	if err := m.conn.Ping(context.Background(), readpref.Primary()); err != nil {
		return false
	}
	return true
}

func New() *MongoDB {
	if m != nil {
		m.Destory()
	}
	m = &MongoDB{
		option: Option{},
	}
	return m
}
