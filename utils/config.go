package utils

import (
	"github.com/spf13/viper"
)

type ConfigCloudStore struct {
	User          string
	Key           string
	LocalDbPath   string
	LocalBasePath string
	RemotePath    string
	RemoteDbPath  string
}

func GetorCreateConfig() ConfigCloudStore {
	return ConfigCloudStore{
		User:          viper.GetString("User"),
		Key:           viper.GetString("Key"),
		LocalBasePath: "/Users/harshilpatel/Projects/cloud-store-test/files/",
		LocalDbPath:   "",
		RemoteDbPath:  "",
		RemotePath:    "localhost:1234",
	}
}
