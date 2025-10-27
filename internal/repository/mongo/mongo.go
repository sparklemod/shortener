package mongo

import (
	"context"
	"fmt"
	adaptermgo "shortener/internal/adapter/mongo"
	"shortener/internal/model"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	conn *adaptermgo.Mongo
}

func NewMongo(conn *adaptermgo.Mongo) *Mongo {
	return &Mongo{conn: conn}
}

func (m *Mongo) col() *mongo.Collection {
	return m.conn.Collection("shortenings")
}

func (m *Mongo) FilterLinks(ctx context.Context, f model.FilterLinksInput) ([]model.Link, error) {
	const op = "mongo.FilterLinks"

	filter := bson.M{}

	if f.IsActive != nil {
		filter["is_active"] = *f.IsActive
	}

	if f.ShortenUrl != nil && *f.ShortenUrl != "" {
		filter["_id"] = *f.ShortenUrl
	}

	sortField := "original_url"
	switch f.SortBy {
	case "visits":
		sortField = "visits"
	}

	sortOrder := 1
	if strings.ToUpper(f.SortOrder) == "DESC" {
		sortOrder = -1
	}

	findOpts := options.Find().
		SetSort(bson.D{{Key: sortField, Value: sortOrder}})

	if f.Limit > 0 {
		findOpts.SetLimit(int64(f.Limit))
	}
	if f.Offset > 0 {
		findOpts.SetSkip(int64(f.Offset))
	}

	cur, err := m.col().Find(ctx, filter, findOpts)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var mgoDocs []mgoShortening
	if err := cur.All(ctx, &mgoDocs); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	links := make([]model.Link, len(mgoDocs))
	for i, doc := range mgoDocs {
		link := modelShorteningFromMgo(doc)
		links[i] = *link
	}

	return links, nil
}

func (m *Mongo) Post(ctx context.Context, shortening model.Link) (*model.Link, error) {
	const op = "shortening.mgo.Put"

	shortening.CreatedAt = time.Now().UTC()

	count, err := m.col().CountDocuments(ctx, bson.M{"_id": shortening.ShortenUrl})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if count > 0 {
		return nil, fmt.Errorf("%s: %w", op, model.ErrorNonUniq)
	}

	_, err = m.col().InsertOne(ctx, mgoShorteningFromModel(shortening))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &shortening, nil
}

func (m *Mongo) Get(ctx context.Context, shorteningID string) (*model.Link, error) {
	const op = "shortening.mgo.Get"

	var shortening mgoShortening
	if err := m.col().FindOne(ctx, bson.M{"_id": shorteningID}).Decode(&shortening); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("%s: %w", op, model.ErrorNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return modelShorteningFromMgo(shortening), nil
}

func (m *Mongo) IncrementVisits(ctx context.Context, shorteningID string) error {
	const op = "shortening.mgo.IncrementVisits"

	var (
		filter = bson.M{"_id": shorteningID}
		update = bson.M{
			"$inc": bson.M{"visits": 1},
			"$set": bson.M{"updated_at": time.Now().UTC()},
		}
	)

	_, err := m.col().UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

type mgoShortening struct {
	Identifier  string    `bson:"_id"`
	OriginalURL string    `bson:"original_url"`
	Visits      int       `bson:"visits"`
	IsActive    bool      `bson:"is_active"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
}

func mgoShorteningFromModel(shortening model.Link) mgoShortening {
	return mgoShortening{
		Identifier:  shortening.ShortenUrl,
		OriginalURL: shortening.OriginalUrl,
		Visits:      shortening.Visits,
		IsActive:    shortening.IsActive,
		CreatedAt:   shortening.CreatedAt,
		UpdatedAt:   shortening.UpdatedAt,
	}
}

func modelShorteningFromMgo(shortening mgoShortening) *model.Link {
	return &model.Link{
		ShortenUrl:  shortening.Identifier,
		OriginalUrl: shortening.OriginalURL,
		Visits:      shortening.Visits,
		IsActive:    shortening.IsActive,
		CreatedAt:   shortening.CreatedAt,
		UpdatedAt:   shortening.UpdatedAt,
	}
}
