package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/harshilkumar/cloud-store-client/utils"
	"github.com/spf13/viper"
)

type Tag string

var wg sync.WaitGroup
var counter time.Duration

// var configuration Config

func deltag(path string, color string, wg *sync.WaitGroup, c time.Duration) {
	wg.Add(1)
	defer wg.Done()

	cmd := exec.Command("tag", "-r", color, path)

	time.Sleep(c)
	fmt.Printf("%v \n", path)

	cmd.Run()
}

func addtag(path string, color string, wg *sync.WaitGroup, c time.Duration) {
	wg.Add(1)
	defer wg.Done()

	cmd := exec.Command("tag", "-a", color, path)

	time.Sleep(c)
	fmt.Printf("%v \n", path)

	cmd.Run()
}

func walkingDelTag(path string, info os.FileInfo, err error) error {
	// fmt.Printf("%v %v \n", path, info.IsDir())
	counter = counter + 1*time.Second
	go deltag(path, "Green", &wg, counter)
	// go deltag(path, "Red", &wg)
	return nil
}

func walkingAddTag(path string, info os.FileInfo, err error) error {

	counter = counter + 1*time.Second
	go addtag(path, "Green", &wg, counter)
	return nil
}

// func main() {
// 	counter = 1 * time.Millisecond
// 	redTag := Tag("Red")
// 	dir := filepath.Dir("/Users/harshilpatel/Projects/cloud-store-test/files/0/")

// 	fmt.Println("%v", redTag)
// 	filepath.Walk(dir, walkingAddTag)

// 	fmt.Println("Now waiting for all routines to finish")

// 	wg.Wait()

// 	counter = 1 * time.Millisecond
// 	time.Sleep(2 * time.Second)
// 	filepath.Walk(dir, walkingDelTag)

// 	wg.Wait()
// }

type Args struct {
	A, B int
}

func main() {

	client, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	// args := Args{3, 4}
	var reply int
	log.Println("reply: %v", reply)
	err = client.Call("Server.VerifyUser", "1234", "1234")
	if err != nil {
		log.Printf("received error %v \n", err)
	} else {
		log.Println("reply: %v", reply)
	}

	viper.SetConfigName("client_config.json")
	viper.SetConfigType("json")
	viper.AddConfigPath("/Users/harshilpatel/Projects/cloud-store-test/files/")

	config := utils.GetorCreateConfig()
	storage := utils.Storage{config.LocalBasePath, make([]utils.FObject, 1)}

	err := storage.CreateObjects()
	if err == nil {
		fmt.Printf("Created fobjects %v \n", len(storage.Objects))
	}

	for _, obj := range storage.Objects {
		fmt.Printf("%v %v \n", obj.Location, obj.Name)
	}

	sigs := make(chan os.Signal, 1)
	// done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs

	viper.WriteConfig()
	storage.DisabllAllTags()
	utils.TagWg.Wait()
	fmt.Println("%v EXIT", sig)
}
