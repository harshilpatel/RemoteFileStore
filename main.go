package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/harshilkumar/cloud-store-client/utils"
	"github.com/spf13/viper"
)

var globalWatch sync.WaitGroup

func watchForSignals() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	log.Println("Received Exit Signal")
}

func watchForConfigChanges() {

}

func main() {

	configLocation := flag.String("config", "", "use for saving and reading config")
	flag.Parse()

	viper.SetConfigName("client_config.json")
	viper.SetConfigType("json")
	viper.AddConfigPath(filepath.Dir(*configLocation))

	config := utils.GetorCreateConfig()
	storage := utils.CreateStorage(config)
	client := utils.CreateClient(config, storage)

	err := storage.CreateObjects()
	if err != nil {
		log.Printf("Error %v", err)
	}
	fmt.Printf("Created fobjects %v \n", len(storage.Objects))

	for _, obj := range storage.Objects {
		fmt.Printf("%v %v \n", obj.Location, obj.Name)
	}

	watchForConfigChanges()

	viper.WriteConfig()
	storage.DisabllAllTags()
	utils.TagWg.Wait()
}
