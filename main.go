package main

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"insider.db/m/internal"
	"insider.db/m/models"
	"insider.db/m/routers"
)

func main() {
	config, err := internal.InitConfig()
	if err != nil {
		fmt.Println("Failed to create config: ", err)
		return
	}

	models.RunMigrate(config.Dsn)
	database := &models.Database{}
	database.Init(config.Dsn)

	defer database.Database.Close()

	node, err := snowflake.NewNode(config.NodeId)
	if err != nil {
		fmt.Println(err)
		return
	}

	routers.Init(database, config, node)
}