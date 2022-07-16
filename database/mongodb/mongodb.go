package mongodb

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hirpc/hrpc/database"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func (o Options) URI() string {
	return fmt.Sprintf("mongodb://%s:%s@%s", o.Username, o.Password, o.Address)
}

type MongoDB struct {
	conn    *mongo.Client
	options Options
}

var m *MongoDB

func (m *MongoDB) Load(src []byte) error {
	// If the value of customized is true (enabled),
	// which means DOES NOT use the configurations from the configuration center.
	if m.options.customized {
		return nil
	}
	if err := json.Unmarshal(src, &m.options); err != nil {
		return err
	}
	return nil
}

func (m *MongoDB) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(m.options.URI()))
	if err != nil {
		return err
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	m.conn = client
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

func New(opts ...Option) *MongoDB {
	var options = Options{
		customized: false,
	}
	for _, o := range opts {
		o(&options)
	}

	if m != nil {
		m.Destory()
	}
	m = &MongoDB{
		options: options,
	}
	return m
}

var _ database.Database = (*MongoDB)(nil)
