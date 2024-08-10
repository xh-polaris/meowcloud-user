package user

import (
	"context"
	"errors"
	"time"

	"github.com/zeromicro/go-zero/core/stores/monc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/xh-polaris/meowcloud-user/biz/infrastructure/config"
	"github.com/xh-polaris/meowcloud-user/biz/infrastructure/consts"
)

const (
	prefixUserCacheKey = "cache:catalbum:"
	CollectionName     = "user"
)

type (
	// IMongoMapper is an interface to be customized, add more methods here,
	// and implement the added methods in MongoMapper.
	IMongoMapper interface {
		Insert(ctx context.Context, data *CatAlbum) error
		FindOne(ctx context.Context, id string) (*CatAlbum, error)
		Upsert(ctx context.Context, data *CatAlbum) error
		Delete(ctx context.Context, id string) error
		FindOneNoCache(ctx context.Context, id string) (*CatAlbum, error)
		FindMany(ctx context.Context, skip int64, count int64) ([]*CatAlbum, int64, error)
		FindManyByCreatorId(ctx context.Context, creatorId string, skip int64, count int64) ([]*CatAlbum, int64, error)
	}

	MongoMapper struct {
		conn *monc.Model
	}

	CatAlbum struct {
		ID                       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
		Type                     int32              `bson:"type,omitempty" json:"type,omitempty"`
		CreatorId                string             `bson:"creatorId,omitempty" json:"creatorId,omitempty"`
		AlbumName                string             `bson:"albumName,omitempty" json:"albumName,omitempty"`
		Visibility               int32              `bson:"visibility,omitempty" json:"visibility,omitempty"`
		TotalPhotos              int32              `bson:"totalPhotos,omitempty" json:"totalPhotos,omitempty"`
		AvailablePhotos          int32              `bson:"availablePhotos,omitempty" json:"availablePhotos,omitempty"`
		UpdatedAt                time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
		DeletedAt                time.Time          `bson:"deletedAt,omitempty" json:"deletedAt,omitempty"`
		CreatedAt                time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
		CreatedLocation          string             `bson:"createdLocation,omitempty" json:"createdLocation,omitempty"`
		CreatedLocationLongitude float64            `bson:"createdLocationLongitude,omitempty" json:"createdLocationLongitude,omitempty"`
		CreatedLocationLatitude  float64            `bson:"createdLocationLatitude,omitempty" json:"createdLocationLatitude,omitempty"`
		CatInfo                  *CatInfo           `bson:"catInfo,omitempty" json:"catInfo,omitempty"`
	}

	CatInfo struct {
		CoverUrl  string    `bson:"coverUrl,omitempty" json:"coverUrl,omitempty"`
		Color     string    `bson:"color,omitempty" json:"color,omitempty"`
		Gender    string    `bson:"gender,omitempty" json:"gender,omitempty"`
		BirthDate time.Time `bson:"birthDate,omitempty" json:"birthDate,omitempty"`
	}
)

func NewMongoMapper(config *config.Config) IMongoMapper {
	conn := monc.MustNewModel(config.Mongo.URL, config.Mongo.DB, CollectionName, config.CacheConf)
	return &MongoMapper{
		conn: conn,
	}
}

func (m *MongoMapper) Upsert(ctx context.Context, data *CatAlbum) error {
	key := prefixUserCacheKey + data.ID.Hex()

	filter := bson.M{
		consts.ID: data.ID,
	}

	set := bson.M{
		"updatedAt":                time.Now(),
		"type":                     data.Type,
		"creatorId":                data.CreatorId,
		"albumName":                data.AlbumName,
		"visibility":               data.Visibility,
		"createdLocation":          data.CreatedLocation,
		"createdLocationLongitude": data.CreatedLocationLongitude,
		"createdLocationLatitude":  data.CreatedLocationLatitude,
		"catInfo":                  data.CatInfo,
	}

	update := bson.M{
		"$set": set,
		"$setOnInsert": bson.M{
			consts.ID:       data.ID,
			consts.CreateAt: time.Now(),
			consts.UpdateAt: time.Now(),
		},
	}

	option := options.UpdateOptions{}
	option.SetUpsert(true)

	_, err := m.conn.UpdateOne(ctx, key, filter, update, &option)
	return err
}

func (m *MongoMapper) Insert(ctx context.Context, data *CatAlbum) error {
	if data.ID.IsZero() {
		data.ID = primitive.NewObjectID()
		data.CreatedAt = time.Now()
		data.UpdatedAt = time.Now()
	}

	key := prefixUserCacheKey + data.ID.Hex()
	_, err := m.conn.InsertOne(ctx, key, data)
	return err
}

func (m *MongoMapper) FindOne(ctx context.Context, id string) (*CatAlbum, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, consts.ErrInvalidObjectId
	}

	var data CatAlbum
	key := prefixUserCacheKey + id
	err = m.conn.FindOne(ctx, key, &data, bson.M{consts.ID: oid})
	switch {
	case err == nil:
		return &data, nil
	case errors.Is(err, monc.ErrNotFound):
		return nil, consts.ErrNotFound
	default:
		return nil, err
	}
}

func (m *MongoMapper) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return consts.ErrInvalidObjectId
	}
	key := prefixUserCacheKey + id
	_, err = m.conn.DeleteOne(ctx, key, bson.M{consts.ID: oid})
	return err
}

func (m *MongoMapper) FindOneNoCache(ctx context.Context, id string) (*CatAlbum, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, consts.ErrInvalidObjectId
	}

	var data CatAlbum
	err = m.conn.FindOneNoCache(ctx, &data, bson.M{consts.ID: oid})
	switch {
	case err == nil:
		return &data, nil
	case errors.Is(err, monc.ErrNotFound):
		return nil, consts.ErrNotFound
	default:
		return nil, err
	}
}

func (m *MongoMapper) FindMany(ctx context.Context, skip int64, count int64) ([]*CatAlbum, int64, error) {
	data := make([]*CatAlbum, 0, 20)
	err := m.conn.Find(ctx, &data, bson.M{}, &options.FindOptions{
		Skip:  &skip,
		Limit: &count,
		Sort:  bson.M{consts.ID: -1},
	})
	if err != nil {
		return nil, 0, err
	}
	total, err := m.conn.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}
	return data, total, nil
}

func (m *MongoMapper) FindManyByCreatorId(ctx context.Context, creatorId string, skip int64, count int64) ([]*CatAlbum, int64, error) {
	data := make([]*CatAlbum, 0, 20)
	err := m.conn.Find(ctx, &data, bson.M{"creatorId": creatorId}, &options.FindOptions{
		Skip:  &skip,
		Limit: &count,
		Sort:  bson.M{consts.ID: -1},
	})
	if err != nil {
		return nil, 0, err
	}
	total, err := m.conn.CountDocuments(ctx, bson.M{"creatorId": creatorId})
	if err != nil {
		return nil, 0, err
	}
	return data, total, nil
}
