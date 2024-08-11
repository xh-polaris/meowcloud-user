package adaptor

import (
	"context"

	"github.com/xh-polaris/service-idl-gen-go/kitex_gen/meowcloud/user"

	"github.com/xh-polaris/meowcloud-user/biz/application/service"
	"github.com/xh-polaris/meowcloud-user/biz/infrastructure/config"
)

type UserServerImpl struct {
	*config.Config
	UserService service.UserService
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
