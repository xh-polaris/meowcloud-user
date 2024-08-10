package provider

import (
	"github.com/google/wire"

	"github.com/xh-polaris/meowcloud-user/biz/application/service"
	"github.com/xh-polaris/meowcloud-user/biz/infrastructure/config"
	"github.com/xh-polaris/meowcloud-user/biz/infrastructure/mapper/user"
	"github.com/xh-polaris/meowcloud-user/biz/infrastructure/stores/redis"
)

var AllProvider = wire.NewSet(
	ApplicationSet,
	InfrastructureSet,
)

var ApplicationSet = wire.NewSet(
	service.UserSet,
	service.CatAlbumSet,
)

var InfrastructureSet = wire.NewSet(
	config.NewConfig,
	redis.NewRedis,
	MapperSet,
)

var MapperSet = wire.NewSet(
	user.NewMongoMapper,
)
