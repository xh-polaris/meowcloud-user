package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/wire"
	genuser "github.com/xh-polaris/service-idl-gen-go/kitex_gen/meowcloud/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/xh-polaris/meowcloud-user/biz/infrastructure/config"
	"github.com/xh-polaris/meowcloud-user/biz/infrastructure/consts"
	usermapper "github.com/xh-polaris/meowcloud-user/biz/infrastructure/mapper/user"
)

type UserService interface {
	GetUser(ctx context.Context, req *genuser.GetUserReq) (res *genuser.GetUserResp, err error)
	GetUserDetail(ctx context.Context, req *genuser.GetUserDetailReq) (res *genuser.GetUserDetailResp, err error)
	UpdateUser(ctx context.Context, req *genuser.UpdateUserReq) (res *genuser.UpdateUserResp, err error)
}

type UserServiceImpl struct {
	Config          *config.Config
	UserMongoMapper usermapper.IMongoMapper
}

var UserSet = wire.NewSet(
	wire.Struct(new(UserServiceImpl), "*"),
	wire.Bind(new(UserService), new(*UserServiceImpl)),
)

func (s *UserServiceImpl) GetUser(ctx context.Context, req *genuser.GetUserReq) (res *genuser.GetUserResp, err error) {
	user1, err := s.UserMongoMapper.FindOne(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return &genuser.GetUserResp{
		User: &genuser.User{
			Id:       user1.ID.Hex(),
			Nickname: user1.Nickname,
			Avatar:   user1.Avatar,
		},
	}, nil
}

func (s *UserServiceImpl) GetUserDetail(ctx context.Context, req *genuser.GetUserDetailReq) (res *genuser.GetUserDetailResp, err error) {
	user, err := s.UserMongoMapper.FindOne(ctx, req.UserId)
	if err != nil {
		if !errors.Is(err, consts.ErrNotFound) {
			return nil, err
		}
		user = &usermapper.User{}
		user.ID, err = primitive.ObjectIDFromHex(req.GetUserId())
		if err != nil {
			return nil, err
		}
		user.Avatar = "https://static.xhpolaris.com/cat_world.jpg"
		user.Nickname = "用户_" + req.GetUserId()[:13]
		user.UpdatedAt = time.Now()
		user.CreatedAt = time.Now()
		err = s.UserMongoMapper.Insert(ctx, user)
		// 处理并发冲突
		if mongo.IsDuplicateKeyError(err) {
			user, err = s.UserMongoMapper.FindOneNoCache(ctx, req.UserId)
			if err != nil {
				return nil, err
			}
		} else if err != nil {
			return nil, err
		}
	}

	return &genuser.GetUserDetailResp{
		User: &genuser.User{
			Id:       user.ID.Hex(),
			Avatar:   user.Avatar,
			Nickname: user.Nickname,
			Bio:      user.Bio,
		},
	}, nil
}

func (s *UserServiceImpl) UpdateUser(ctx context.Context, req *genuser.UpdateUserReq) (res *genuser.UpdateUserResp, err error) {
	oid, err := primitive.ObjectIDFromHex(req.User.Id)
	if err != nil {
		return nil, consts.ErrInvalidObjectId
	}

	err = s.UserMongoMapper.UpsertUser(ctx, &usermapper.User{
		ID:       oid,
		Avatar:   req.User.Avatar,
		Nickname: req.User.Nickname,
		Bio:      req.User.Bio,
	})
	if err != nil {
		return nil, err
	}

	return &genuser.UpdateUserResp{}, nil
}
