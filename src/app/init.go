package app

import (
	healthcheckrepo "github.com/idilhaq/go-template/src/app/domain/repository/healthcheck"
	masterprojectrepo "github.com/idilhaq/go-template/src/app/domain/repository/master/project"

	healthcheckusecase "github.com/idilhaq/go-template/src/app/usecase/healthcheck"
	projecthomeusecase "github.com/idilhaq/go-template/src/app/usecase/home"
	masterprojectusecase "github.com/idilhaq/go-template/src/app/usecase/master/project"

	"github.com/idilhaq/go-template/src/pkg/lib/apicalls"
	"github.com/idilhaq/go-template/src/pkg/lib/apicalls/request"
	"github.com/idilhaq/go-template/src/pkg/lib/config"
	"github.com/idilhaq/go-template/src/pkg/lib/storage/database"
	"github.com/idilhaq/go-template/src/pkg/lib/storage/redis"
)

var (
	healthCheckRepo    healthcheckrepo.RepositoryItf
	healthCheckUsecase healthcheckusecase.UsecaseItf

	masterProjectRepo    masterprojectrepo.RepositoryItf
	masterProjectUsecase masterprojectusecase.UsecaseItf

	projectHomeUsecase projecthomeusecase.UsecaseItf
)

type AppUsecaseDepedency struct {
	HealthCheck     healthcheckusecase.UsecaseItf
	MasterProject masterprojectusecase.UsecaseItf
	ProjectHome   projecthomeusecase.UsecaseItf
}

func Init(cfg *config.Config) (*AppUsecaseDepedency, error) {
	dbCore, err := database.Init(cfg.Storage.Database, database.DatabaseCore)
	if err != nil {
		return nil, err
	}

	redisCore, err := redis.Init(cfg.Storage.Redis, redis.RedisCore)
	if err != nil {
		return nil, err
	}

	requestClient, err := request.Init(nil)
	if err != nil {
		return nil, err
	}

	apicallsClient := apicalls.Init(requestClient, cfg.APICalls)

	initRepositories(
		dbCore,
		redisCore,
		apicallsClient,
	)

	initUsecases(
		cfg.Main,
	)

	return &AppUsecaseDepedency{
		HealthCheck:     healthCheckUsecase,
		MasterProject: masterProjectUsecase,
		ProjectHome:   projectHomeUsecase,
	}, nil
}

func initRepositories(
	dbCore database.SQLDBItf,
	redisCore redis.RedisItf,
	apicallsClient apicalls.Apicalls,
) {
	healthCheckRepo = healthcheckrepo.InitRepository(
		dbCore, redisCore,
	)

	masterProjectRepo = masterprojectrepo.InitRepository(
		dbCore,
	)
}

func initUsecases(
	cfg config.MainConfig,
) {
	healthCheckUsecase = healthcheckusecase.InitUsecase(
		healthCheckRepo,
	)

	masterProjectUsecase = masterprojectusecase.InitUsecase(
		masterProjectRepo,
	)

	projectHomeUsecase = projecthomeusecase.InitUsecase(
		masterProjectRepo,
	)

}
