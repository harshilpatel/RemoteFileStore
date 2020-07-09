package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"

	"github.com/google/uuid"
	"github.com/harshilkumar/cloud-store-client/utils"
	"github.com/spf13/viper"
)

func CreateConfig() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.TraceLevel)

	configLocation := flag.String("config", ".", "")
	configUser := flag.String("user", "", "")
	configKey := flag.String("key", uuid.New().String(), "")
	configBasePath := flag.String("base_path", ".", "")
	configRemotePath := flag.String("remote_path", "localhost:4533", "")
	flag.Parse()

	if *configUser == "" {
		os.Exit(1)
	}

	viper.SetDefault("client.Storage.Config.ClientInstanceId", uuid.New().String())
	viper.SetDefault("client.Storage.Config.Key", *configKey)
	viper.SetDefault("client.Storage.Config.User", *configUser)
	viper.SetDefault("client.Storage.Config.BasePath", *configBasePath)
	viper.SetDefault("client.Storage.Config.RemotePath", *configRemotePath)

	viper.SetConfigName("client_config.json")
	viper.SetConfigType("json")
	viper.AddConfigPath(filepath.Dir(*configLocation))

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Could not open the config file %s", err))
	}

	config := utils.GetorCreateConfig()
	storage := utils.CreateStorage(config)
	client, err := utils.CreateClient(config, storage)

	if err != nil {
		logrus.Debugln(err)
	} else {
		client.Client.Call(utils.ServerRegisterUser, client.Storage.User.Username, &client.Storage.User)
	}

	viper.Set("client", client)
	viper.WriteConfig()

}
