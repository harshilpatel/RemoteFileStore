package utils

import (
	"log"
	"net/rpc"
)

type Client struct {
	Storage Storage
	Client  *rpc.Client
	User    User
}

func CreateClient(config ConfigCloudStore, storage Storage) Client {
	client, err := rpc.DialHTTP("tcp", config.RemotePath)
	if err != nil {
		log.Printf("Error creating a client %v", err)
	}

	return Client{
		Storage: storage,
		Client:  client,
		User: User{
			Username: config.User,
			Key:      config.Key,
		},
	}
}
