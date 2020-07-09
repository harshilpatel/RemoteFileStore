package utils

import (
	"errors"
	"io/ioutil"
	"net/rpc"
	"os"
	"time"

	server "github.com/harshilkumar/cloud-store-server/utils"
	"github.com/sirupsen/logrus"
)

type Client struct {
	Storage         Storage
	Client          *rpc.Client
	ClientAvailable bool
}

func CreateClient(config ConfigCloudStore, storage Storage) (Client, error) {
	c := Client{
		Storage: storage,
	}
	client, err := rpc.DialHTTP("tcp", config.RemotePath)
	if err != nil {
		logrus.Printf("Error creating a client %v", err)
		return c, errors.New("Error connecting to the Server")
	}

	return Client{
		Storage: storage,
		Client:  client,
	}, nil
}

func (c *Client) MakeUploadRequest(f *FObject) {
	if data, err := ioutil.ReadFile(f.GetRealPath(c.Storage.User, c.Storage.Config)); err == nil {
		requestPackage := server.UserRequestPackage{
			ClientUser: c.Storage.User.User,
			Obj:        f.FObject,
			Operation:  "Create",
			Data:       data,
		}

		if e := c.Client.Call(ServerSaveObject, &requestPackage, nil); e != nil {
			logrus.Printf("Received error when saving file on server %v \n", e)
		} else {
			f.LastPushed = time.Now().UTC()
		}
	}
}

func (c *Client) MakeDownloadRequest(f FObject) {
	pack := UserRequestPackage{
		ClientUser: c.Storage.User.User,
		Obj:        f.FObject,
		Operation:  "Download",
		Data:       make([]byte, 1),
	}

	if e := c.Client.Call(ServerDownloadObject, pack, &pack); e == nil {
		f.WriteData(pack.Data)
		f.Version = pack.Obj.Version
		f.LastPulled = time.Now().UTC()

		c.Storage.Objects[f.Relativepath] = f
	} else {
		logrus.Printf("Received error when saving file locally %v \n", e)
	}
}

func (c *Client) VerifyObjects() {
	for _, obj := range c.Storage.Objects {
		realPath := obj.GetRealPath(c.Storage.User, c.Storage.Config)
		if _, err := os.Stat(realPath); err == nil {

			pack := UserRequestPackage{
				ClientUser: c.Storage.User.User,
				Obj:        obj.FObject,
			}
			response := 0
			if e := c.Client.Call(ServerVerifyObject, &pack, &response); e == nil {
				logrus.WithFields(logrus.Fields{
					"File":     realPath,
					"response": response,
				}).Debugf("Response on VerifyObjects")

				switch response {
				case 0:
					// All Good
					continue
				case 1:
					// Older File than Server. Prepare a Pull.
					c.MakeDownloadRequest(obj)
				case 2:
					// Newer File than Server. Prepare a Push.
					c.MakeUploadRequest(&obj)
				}
				// logrus.Printf("Received error when checking for file on server %v \n", e)
			}
		}

	}
}
