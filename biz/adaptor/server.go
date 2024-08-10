package adaptor

import (
	"context"

	"github.com/xh-polaris/service-idl-gen-go/kitex_gen/meowcloud/content"
	"github.com/xh-polaris/service-idl-gen-go/kitex_gen/meowcloud/user"

	"github.com/xh-polaris/meowcloud-user/biz/application/service"
	"github.com/xh-polaris/meowcloud-user/biz/infrastructure/config"
)

type UserServerImpl struct {
	*config.Config
	UserService     service.UserService
	CatAlbumService service.CatAlbumService
}

func (s *UserServerImpl) GetUser(ctx context.Context, req *user.GetUserReq) (res *user.GetUserResp, err error) {
	return s.UserService.GetUser(ctx, req)
}

func (s *UserServerImpl) GetUserDetail(ctx context.Context, req *user.GetUserDetailReq) (res *user.GetUserDetailResp, err error) {
	return s.UserService.GetUserDetail(ctx, req)
}

func (s *UserServerImpl) UpdateUser(ctx context.Context, req *user.UpdateUserReq) (res *user.UpdateUserResp, err error) {
	return s.UserService.UpdateUser(ctx, req)
}

func (s *UserServerImpl) CreateCatAlbum(ctx context.Context, req *content.CreateCatAlbumReq) (res *content.CreateCatAlbumResp, err error) {
	return s.CatAlbumService.CreateCatAlbum(ctx, req)
}

func (s *UserServerImpl) RetrieveCatAlbum(ctx context.Context, req *content.RetrieveCatAlbumReq) (res *content.RetrieveCatAlbumResp, err error) {
	return s.CatAlbumService.RetrieveCatAlbum(ctx, req)
}

func (s *UserServerImpl) UpdateCatAlbum(ctx context.Context, req *content.UpdateCatAlbumReq) (res *content.UpdateCatAlbumResp, err error) {
	return s.CatAlbumService.UpdateCatAlbum(ctx, req)
}

func (s *UserServerImpl) DeleteCatAlbum(ctx context.Context, req *content.DeleteCatAlbumReq) (res *content.DeleteCatAlbumResp, err error) {
	return s.CatAlbumService.DeleteCatAlbum(ctx, req)
}

func (s *UserServerImpl) ListCatAlbum(ctx context.Context, req *content.ListCatAlbumReq) (res *content.ListCatAlbumResp, err error) {
	return s.CatAlbumService.ListCatAlbum(ctx, req)
}
