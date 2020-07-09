package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/harshilkumar/cloud-store-client/utils"
	"github.com/spf13/viper"
)

var globalWatch sync.WaitGroup

func watchForSignals() {

	// logrus.Println("Received Exit Signal")
}

func watchForConfigChanges() {

}

func main() {

	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.TraceLevel)

	configLocation := flag.String("config", ".", "use for saving and reading config")
	flag.Parse()

	logrus.Printf("Looking for config in %v", *configLocation)

	viper.SetDefault("client.Storage.Config.ClientInstanceId", uuid.New().String())
	viper.SetConfigName("client_config.json")
	viper.SetConfigType("json")
	viper.AddConfigPath(filepath.Dir(*configLocation))
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Could not open the config file %s", err))
	}

	config := utils.GetorCreateConfig()
	storage := utils.CreateStorage(config)
	client, _ := utils.CreateClient(config, storage)

	if callError := client.Client.Call(utils.ServerVerifyUser, &storage.User.Username, &storage.User.Key); callError != nil {
		logrus.Printf("Could not auth user %v", callError)
		logrus.Printf("Registering User")

		key := ""
		if err := client.Client.Call(utils.ServerRegisterUser, storage.User.Username, &key); err != nil {
			storage.User.Key = key
		} else {
			logrus.Error("Sevrer rejected client")
			os.Exit(1)
		}

	}
	logrus.Printf("Verify Local User from Config %v", storage.User.Username)

	err := storage.CreateObjects()
	if err != nil {
		logrus.Printf("Error %v", err)
	}
	// storage.SecureAllObjects()
	fmt.Printf("Created fobjects %v \n", len(storage.Objects))
	client.VerifyObjects()
	client.HeartBeat()
	client.InitiateWatchers()

	watchForConfigChanges()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	logrus.Printf("Safe for shutdown signals")
	<-sigs

	client.Watcher.Close()
	client.Watcher = nil

	viper.Set("client", client)
	viper.WriteConfig()
	// storage.DisabllAllTags()

	logrus.Println("")
	logrus.Printf("Shutting down client")
}
