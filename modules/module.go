package modules

import (
	"request-redeem/config"
	"request-redeem/controller"
	"request-redeem/handler"
	"request-redeem/helper"
	"request-redeem/router"
)

func RegisterModules(dbCon *config.DatabaseConnection, c *config.AppConfig) router.Controller {
	reward := config.RepositoryFactory(dbCon)
	jwtService := handler.NewJWTService()
	rewardService := helper.NewRewardService(reward)
	controller := router.Controller{
		Auth: controller.NewAuthController(rewardService, jwtService),
	}

	return controller
}
