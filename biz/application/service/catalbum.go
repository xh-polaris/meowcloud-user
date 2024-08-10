package service

import (
	"context"
	"github.com/google/wire"
	"github.com/jinzhu/copier"
	"github.com/xh-polaris/meowcloud-user/biz/infrastructure/config"
	catalbummapper "github.com/xh-polaris/meowcloud-user/biz/infrastructure/mapper/catalbum"
	"github.com/xh-polaris/service-idl-gen-go/kitex_gen/meowcloud/content"
)

type CatAlbumService interface {
	CreateCatAlbum(ctx context.Context, req *content.CreateCatAlbumReq) (res *content.CreateCatAlbumResp, err error)
	RetrieveCatAlbum(ctx context.Context, req *content.RetrieveCatAlbumReq) (res *content.RetrieveCatAlbumResp, err error)
	UpdateCatAlbum(ctx context.Context, req *content.UpdateCatAlbumReq) (res *content.UpdateCatAlbumResp, err error)
	DeleteCatAlbum(ctx context.Context, req *content.DeleteCatAlbumReq) (res *content.DeleteCatAlbumResp, err error)
	ListCatAlbum(ctx context.Context, req *content.ListCatAlbumReq) (res *content.ListCatAlbumResp, err error)
}

type CatAlbumServiceImpl struct {
	Config              *config.Config
	CatAlbumMongoMapper catalbummapper.IMongoMapper
}

var CatAlbumSet = wire.NewSet(
	wire.Struct(new(CatAlbumServiceImpl), "*"),
	wire.Bind(new(CatAlbumService), new(*CatAlbumServiceImpl)),
)

func (s *CatAlbumServiceImpl) CreateCatAlbum(ctx context.Context, req *content.CreateCatAlbumReq) (res *content.CreateCatAlbumResp, err error) {
	catalbum := &catalbummapper.CatAlbum{}
	err = copier.Copy(catalbum, req.CatAlbum)
	if err != nil {
		return nil, err
	}
	err = s.CatAlbumMongoMapper.Insert(ctx, catalbum)
	if err != nil {
		return nil, err
	}
	return &content.CreateCatAlbumResp{CatAlbum: &content.CatAlbum{
		Id: catalbum.ID.Hex(),
	}}, nil
}

func (s *CatAlbumServiceImpl) RetrieveCatAlbum(ctx context.Context, req *content.RetrieveCatAlbumReq) (res *content.RetrieveCatAlbumResp, err error) {
	catalbum, err := s.CatAlbumMongoMapper.FindOne(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	c := &content.CatAlbum{}
	err = copier.Copy(c, catalbum)
	if err != nil {
		return nil, err
	}
	return &content.RetrieveCatAlbumResp{CatAlbum: &content.CatAlbum{}}, nil
}

func (s *CatAlbumServiceImpl) UpdateCatAlbum(ctx context.Context, req *content.UpdateCatAlbumReq) (res *content.UpdateCatAlbumResp, err error) {
	catalbum := &catalbummapper.CatAlbum{}
	err = copier.Copy(catalbum, req.CatAlbum)
	if err != nil {
		return nil, err
	}
	err = s.CatAlbumMongoMapper.Upsert(ctx, catalbum)
	if err != nil {
		return nil, err
	}
	return &content.UpdateCatAlbumResp{CatAlbum: &content.CatAlbum{
		Id: catalbum.ID.Hex(),
	}}, nil
}

func (s *CatAlbumServiceImpl) DeleteCatAlbum(ctx context.Context, req *content.DeleteCatAlbumReq) (res *content.DeleteCatAlbumResp, err error) {
	err = s.CatAlbumMongoMapper.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &content.DeleteCatAlbumResp{}, nil
}

func (s *CatAlbumServiceImpl) ListCatAlbum(ctx context.Context, req *content.ListCatAlbumReq) (res *content.ListCatAlbumResp, err error) {
	catalbums, count, err := s.CatAlbumMongoMapper.FindMany(ctx, req.PaginationOptions.GetOffset(), req.PaginationOptions.GetLimit())
	if err != nil {
		return nil, err
	}

	var catAlbumList []*content.CatAlbum
	for _, album := range catalbums {
		catAlbum := &content.CatAlbum{}
		err = copier.Copy(catAlbum, album)
		if err != nil {
			return nil, err
		}
		catAlbumList = append(catAlbumList, catAlbum)
	}

	return &content.ListCatAlbumResp{CatAlbums: catAlbumList, Total: count}, nil
}
