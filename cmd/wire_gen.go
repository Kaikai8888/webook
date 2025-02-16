// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"webook/internal/interface/web"
	"webook/internal/repository"
	"webook/internal/repository/dao"
	"webook/internal/service"
	"webook/ioc"
)

// Injectors from wire.go:

func InitApp() *App {
	db := ioc.InitDB()
	userDAO := dao.NewUserDAO(db)
	userRepository := repository.NewUserRepository(userDAO)
	userService := service.NewUserService(userRepository)
	userHandler := web.NewUserHandler(userService)
	logger := ioc.InitLogger()
	draftArticleDao := dao.NewDraftArticleDao(logger, db)
	draftArticleRepository := repository.NewDraftArticleRepository(logger, draftArticleDao)
	articleService := service.NewArticleService(logger, draftArticleRepository)
	articleHandler := web.NewArticleHandler(logger, articleService)
	engine := ioc.InitWebServer(userHandler, articleHandler)
	app := &App{
		webServer: engine,
	}
	return app
}
