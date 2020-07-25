package utils

import (
	server "github.com/harshilkumar/cloud-store-server/utils"
)

const ServerVerifyUser = "Storage.VerifyUser"
const ServerSaveObject = "Storage.SaveObject"
const ServerDownloadObject = "Storage.DownloadObject"
const ServerVerifyObject = "Storage.VerifyObject"
const ServerRegisterUser = "Storage.RegisterUser"

type User struct {
	server.User
}

type UserRequestPackage server.UserRequestPackage
