package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	// "go-prj-skeleton/app/interface/restful"
	"go-prj-skeleton/app/interface/gin"
	"go-prj-skeleton/app/pgutil"
	"go-prj-skeleton/app/registry"
	"go-prj-skeleton/app/setting"
)

func main() {
	initEnvSettings()

	pgutil.StartUp(pgutil.Configuration{
		URL:             os.Getenv("DATABASE_URL"), //make work with heroku
		Host:            setting.ProjectEnvSettings.PostgreHost,
		Port:            setting.ProjectEnvSettings.PostgrePort,
		Database:        setting.ProjectEnvSettings.PostgreDatabaseName,
		User:            setting.ProjectEnvSettings.PostgreUser,
		Password:        setting.ProjectEnvSettings.PostgrePassword,
		ApplicationName: "HRS",
	})

	// make it work on heroku
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	ctn, err := registry.NewContainer()
	if err != nil {
		log.Fatalf("failed to build container: %v", err)
	}
	_ = ctn
	// server := http.Server{
	// 	Addr:    ":" + port,
	// 	Handler: restful.Handlers(ctn),
	// }

	// fmt.Printf("starting server on port: %s\n", port)
	// if err := server.ListenAndServe(); err != nil {
	// 	fmt.Printf("start sever fail: %s", err.Error())
	// }

	handler := gin.Handler()
	fmt.Printf("starting server on port: %s\n", port)
	if err := handler.Run(":" + port); err != nil {
		fmt.Printf("start sever fail: %s", err.Error())
	}
}

func initEnvSettings() {
	//initialize env settings and read from env
	setting.EnvSettingsInit([]string{
		"SETTING_POSTGRE_HOST",
		"SETTING_POSTGRE_PORT",
		"SETTING_POSTGRE_DATABASE_NAME",
		"SETTING_POSTGRE_USER",
	})
}
