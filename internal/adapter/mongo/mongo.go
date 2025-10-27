package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	urlConnect string
	dbName     string
	client     *mongo.Client
}

func New(urlConnect, dbName string) *Mongo {
	return &Mongo{urlConnect: urlConnect, dbName: dbName}
}
func (m *Mongo) Collection(name string) *mongo.Collection {
	return m.client.Database(m.dbName).Collection(name)
}

func (mgo *Mongo) Connect(ctx context.Context) error {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mgo.urlConnect))
	if err != nil {
		return err
	}

	if err := client.Ping(ctx, nil); err != nil {
		client.Disconnect(ctx)
		return err
	}

	mgo.client = client
	return nil
}

func (m *Mongo) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}
