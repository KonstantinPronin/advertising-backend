package repository

import (
	"context"
	"fmt"
	"github.com/KonstantinPronin/advertising-backend/internal/advertisement"
	"github.com/KonstantinPronin/advertising-backend/internal/advertisement/model"
	"github.com/KonstantinPronin/advertising-backend/pkg/constants"
	pkg "github.com/KonstantinPronin/advertising-backend/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math"

	"github.com/KonstantinPronin/advertising-backend/pkg/infrastructure"
	"go.uber.org/zap"
	"time"
)

const (
	collection = "advertisement"
)

type MongoDbClient struct {
	db     *infrastructure.Database
	logger *zap.Logger
}

func (m *MongoDbClient) CreateAd(ad *model.Ad) (string, error) {
	if ad == nil {
		return "", pkg.NewInvalidArgument("Nil argument")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	ad.Created = time.Now().Format(time.RFC3339)

	result, err := m.db.GetConnection().Collection(collection).InsertOne(ctx, ad)
	if err != nil {
		m.logger.Error(fmt.Sprintf("unexpected error: %s", err.Error()))
		return "", err
	}

	id := result.InsertedID.(primitive.ObjectID).Hex()
	m.logger.Info(fmt.Sprintf("Advertisement %s was saved to database", id))

	return id, nil
}

func (m *MongoDbClient) GetAdsOrderByPrice(page uint32, desc bool) (model.AdList, uint32, error) {
	opt := options.Find()

	opt.SetLimit(constants.PageSize)
	opt.SetSkip(int64((page - 1) * constants.PageSize))
	if desc {
		opt.SetSort(bson.D{{"price", -1}, {"_id", 1}})
	} else {
		opt.SetSort(bson.D{{"price", 1}, {"_id", 1}})
	}

	return m.GetAds(page, opt)
}

func (m *MongoDbClient) GetAdsOrderByDate(page uint32, desc bool) (model.AdList, uint32, error) {
	opt := options.Find()

	opt.SetLimit(constants.PageSize)
	opt.SetSkip(int64((page - 1) * constants.PageSize))
	if desc {
		opt.SetSort(bson.D{{"created", -1}, {"_id", 1}})
	} else {
		opt.SetSort(bson.D{{"created", 1}, {"_id", 1}})
	}

	return m.GetAds(page, opt)
}

func (m *MongoDbClient) GetAd(id string) (*model.Ad, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, pkg.NewInvalidArgument("wrong id format")
	}

	ad := new(model.Ad)
	filter := bson.M{"_id": _id}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	result := m.db.GetConnection().
		Collection(collection).FindOne(ctx, filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, pkg.NewNotFoundError(result.Err().Error())
		}
		return nil, result.Err()
	}

	if err := result.Decode(ad); err != nil {
		m.logger.Error(fmt.Sprintf("Decode error: %s", err.Error()))
		return nil, err
	}

	return ad, nil
}

func (m *MongoDbClient) GetAds(page uint32, opt *options.FindOptions) (model.AdList, uint32, error) {
	result := model.AdList{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	total, err := m.db.GetConnection().
		Collection(collection).CountDocuments(ctx, bson.M{})
	if err != nil {
		m.logger.Error(fmt.Sprintf("unexpected error: %s", err.Error()))
		return nil, 0, err
	}

	maxPage := uint32(math.Ceil(float64(total) / float64(constants.PageSize)))

	if page > maxPage {
		return nil, maxPage, pkg.NewInvalidArgument("page over limit")
	}

	cursor, err := m.db.GetConnection().Collection(collection).Find(ctx, bson.M{}, opt)
	if err != nil {
		return nil, maxPage, err
	}
	defer func() {
		if err = cursor.Close(ctx); err != nil {
			m.logger.Error(fmt.Sprintf("Resource release error: %s", err.Error()))
		}
	}()

	for cursor.Next(ctx) {
		ad := new(model.Ad)
		if err = cursor.Decode(ad); err != nil {
			m.logger.Error(fmt.Sprintf("Decode error: %s", err.Error()))
			return nil, maxPage, err
		}

		result = append(result, *ad)
	}

	return result, maxPage, nil
}

func NewMongoDbClient(db *infrastructure.Database, logger *zap.Logger) advertisement.Repository {
	return &MongoDbClient{
		db:     db,
		logger: logger,
	}
}
