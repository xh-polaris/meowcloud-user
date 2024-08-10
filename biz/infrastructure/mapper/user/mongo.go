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
	prefixUserCacheKey = "cache:user:"
	CollectionName     = "user"
)

type (
	// IMongoMapper is an interface to be customized, add more methods here,
	// and implement the added methods in MongoMapper.
	IMongoMapper interface {
		Insert(ctx context.Context, data *User) error
		FindOne(ctx context.Context, id string) (*User, error)
		Update(ctx context.Context, data *User) error
		Delete(ctx context.Context, id string) error
		UpsertUser(ctx context.Context, data *User) error
		FindOneNoCache(ctx context.Context, id string) (*User, error)
	}

	MongoMapper struct {
		conn *monc.Model
	}

	UserPreview struct {
		ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
		Nickname string             `bson:"nickname,omitempty" json:"nickname,omitempty"`
		Avatar   string             `bson:"avatar,omitempty" json:"avatar,omitempty"`
	}

	User struct {
		ID                 primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
		Nickname           string             `bson:"username,omitempty" json:"username,omitempty"`
		Bio                string             `bson:"bio,omitempty" json:"bio,omitempty"`
		Avatar             string             `bson:"avatar,omitempty" json:"avatar,omitempty"`
		UpdatedAt          time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
		CreatedAt          time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
		DeletedAt          time.Time          `bson:"deletedAt,omitempty" json:"deletedAt,omitempty"`
		Membership         *Membership        `bson:"membership,omitempty" json:"membership,omitempty"`
		TeamCount          int32              `bson:"teamCount,omitempty" json:"teamCount,omitempty"`
		TeamIds            []string           `bson:"teamIds,omitempty" json:"teamIds,omitempty"`
		MyAlbumCount       int32              `bson:"myAlbumCount,omitempty" json:"myAlbumCount,omitempty"`
		MyAlbumIds         []int32            `bson:"myAlbumIds,omitempty" json:"myAlbumIds,omitempty"`
		FollowedAlbumCount int32              `bson:"followedAlbumCount,omitempty" json:"followedAlbumCount,omitempty"`
		FollowedAlbumIds   []int32            `bson:"followedAlbumIds,omitempty" json:"followedAlbumIds,omitempty"`
		StorageInfo        *StorageInfo       `bson:"storageInfo,omitempty" json:"storageInfo,omitempty"`
		Points             int32              `bson:"points,omitempty" json:"points,omitempty"`
		Achievements       []string           `bson:"achievements,omitempty" json:"achievements,omitempty"`
	}

	Membership struct {
		MemberId    string `bson:"memberId,omitempty" json:"memberId,omitempty"`
		MemberLevel int32  `bson:"memberLevel,omitempty" json:"memberLevel,omitempty"` // 0: 普通用户, 1: 普通会员, 2: 高级会员
	}

	StorageInfo struct {
		AvailablePhotos int32 `bson:"availablePhotos,omitempty" json:"availablePhotos,omitempty"`
		UsedPhotos      int32 `bson:"usedPhotos,omitempty" json:"usedPhotos,omitempty"`
		AvailableMemory int64 `bson:"availableMemory,omitempty" json:"availableMemory,omitempty"`
		UsedMemory      int64 `bson:"usedMemory,omitempty" json:"usedMemory,omitempty"`
		AvailableAlbums int32 `bson:"availableAlbums,omitempty" json:"availableAlbums,omitempty"`
		UsedAlbums      int32 `bson:"usedAlbums,omitempty" json:"usedAlbums,omitempty"`
	}
)

func NewMongoMapper(config *config.Config) IMongoMapper {
	conn := monc.MustNewModel(config.Mongo.URL, config.Mongo.DB, CollectionName, config.CacheConf)
	return &MongoMapper{
		conn: conn,
	}
}

func (m *MongoMapper) UpsertUser(ctx context.Context, data *User) error {
	key := prefixUserCacheKey + data.ID.Hex()

	filter := bson.M{
		consts.ID: data.ID,
	}

	set := bson.M{
		consts.UpdateAt: time.Now(),
	}
	if data.Bio != "" {
		set["bio"] = data.Bio
	}
	if data.Avatar != "" {
		set["avatar"] = data.Avatar
	}
	if data.Nickname != "" {
		set["username"] = data.Nickname
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

func (m *MongoMapper) Insert(ctx context.Context, data *User) error {
	if data.ID.IsZero() {
		data.ID = primitive.NewObjectID()
		data.CreatedAt = time.Now()
		data.UpdatedAt = time.Now()
	}

	key := prefixUserCacheKey + data.ID.Hex()
	_, err := m.conn.InsertOne(ctx, key, data)
	return err
}

func (m *MongoMapper) FindOne(ctx context.Context, id string) (*User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, consts.ErrInvalidObjectId
	}

	var data User
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

func (m *MongoMapper) Update(ctx context.Context, data *User) error {
	data.UpdatedAt = time.Now()
	key := prefixUserCacheKey + data.ID.Hex()
	_, err := m.conn.UpdateOne(ctx, key, bson.M{consts.ID: data.ID}, bson.M{"$set": data})
	return err
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

func (m *MongoMapper) FindOneNoCache(ctx context.Context, id string) (*User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, consts.ErrInvalidObjectId
	}

	var data User
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
