package utils

import (
	server "github.com/harshilkumar/cloud-store-server/utils"
	"github.com/spf13/viper"
)

type ConfigCloudStore struct {
	server.ConfigCloudStore

	User             string
	Key              string
	ClientInstanceId string
	LocalDbPath      string
	LocalBasePath    string
	RemotePath       string
	RemoteDbPath     string
}

func GetorCreateConfig() ConfigCloudStore {
	c := ConfigCloudStore{
		User:             viper.GetString("client.Storage.User.Username"),
		Key:              viper.GetString("client.Storage.User.Key"),
		ClientInstanceId: viper.GetString("client.Storage.Config.ClientInstanceId"),
		LocalDbPath:      "",
		LocalBasePath:    viper.GetString("client.Storage.Config.BasePath"),
		RemoteDbPath:     "",
		RemotePath:       viper.GetString("client.Storage.Config.RemotePath"),
	}
	c.BasePath = viper.GetString("client.Storage.Config.BasePath")

	// c.ConfigCloudStore.BasePath = c.BasePath
	return c
}
