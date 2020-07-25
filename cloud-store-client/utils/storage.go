package utils

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"

	"github.com/sirupsen/logrus"
)

// Storage keeps all local objects
type Storage struct {
	Root    string
	Objects map[string]FObject
	Config  ConfigCloudStore
	User    User

	mux sync.Mutex
}

func CreateStorage(c ConfigCloudStore) Storage {
	logrus.Printf("Creating or reading local storage from %v", c.BasePath)
	user := User{}
	user.Username = c.User
	user.Key = c.Key

	return Storage{
		Root:    c.LocalBasePath,
		Config:  c,
		Objects: make(map[string]FObject),
		User:    user,
	}
}

func (s *Storage) UpdateObjectForUser(obj FObject) {
	s.mux.Lock()

	s.Objects[obj.Relativepath] = obj

	s.mux.Unlock()
}

func (s *Storage) GetObjects() map[string]FObject {
	s.mux.Lock()
	defer s.mux.Unlock()

	return s.Objects
}

func (s *Storage) GetObject(relativepath string) FObject {
	s.mux.Lock()
	defer s.mux.Unlock()

	return s.Objects[relativepath]
}

func (s *Storage) createObject(path string, info os.FileInfo) {
	relativePath := strings.TrimPrefix(path, s.Config.BasePath)

	obj := FObject{}
	obj.Name = info.Name()
	obj.Location = path
	obj.Relativepath = relativePath
	obj.LastWritten = info.ModTime().UTC()
	obj.LastPushed = info.ModTime().UTC()
	obj.LastPulled = info.ModTime().UTC()
	obj.Size = info.Size()
	obj.Version = 0

	obj.UpdateHashForObject(s.User, s.Config)
	obj.UpdateHashForObjectBlocks(s.User, s.Config)

	s.UpdateObjectForUser(obj)
}

func (s *Storage) SecureAllObjects() {
	ch := make(chan int, len(s.Objects))
	for _, obj := range s.Objects {
		obj.TagSecure(ch)
	}
	for range make([]int, len(s.Objects)) {
		<-ch
	}

	logrus.Printf("Finished tagging")
}

func (s *Storage) DisabllAllTags() {
	ch := make(chan int, len(s.Objects))
	for _, obj := range s.Objects {
		obj.TagRemoveAll(ch)
	}

	for range make([]int, len(s.Objects)) {
		<-ch
	}

}

func (s *Storage) walkRoot(path string, info os.FileInfo, err error) error {
	// logrus.Printf("Found Object %v\n", path)
	if !info.IsDir() {
		s.createObject(path, info)
	}

	return nil
}

func (s *Storage) CreateObjects() error {

	if len(s.Root) == 0 {
		return errors.New("No root string found")
	}

	if _, e := os.Stat(s.Root); os.IsNotExist(e) {
		return errors.New("Root folder does not exists")
	}

	filepath.Walk(s.Config.LocalBasePath, s.walkRoot)
	s.parseConfigForObjects()

	logrus.Println("All Objects Acknowledges and Marked")

	return nil
}

func (s *Storage) parseConfigForObjects() {
	configUserMap := viper.GetStringMap("client.storage")
	if configUserMap["objects"] != nil {
		configUserObjectsMap := configUserMap["objects"].(map[string]interface{})

		for configFileRelativePath, configObj := range configUserObjectsMap {
			configObjMap := configObj.(map[string]interface{})

			lastWritten := configObjMap["lastwritten"].(string)
			lastPushed := configObjMap["lastpushed"].(string)
			lastPulled := configObjMap["lastpulled"].(string)
			version := configObjMap["version"].(float64)

			if userObject, ok := s.Objects[configFileRelativePath]; ok {
				if LastPulled, err := time.Parse(time.RFC3339Nano, lastPulled); err == nil {
					userObject.LastPulled = LastPulled
				}
				if LastWritten, err := time.Parse(time.RFC3339Nano, lastWritten); err == nil {
					userObject.LastWritten = LastWritten
				}
				if LastPushed, err := time.Parse(time.RFC3339Nano, lastPushed); err == nil {
					userObject.LastPushed = LastPushed
				}

				userObject.Version = int64(version)
				s.UpdateObjectForUser(userObject)
			}

		}

	}
}

func (s *Storage) keepWatchingForChanges() {
	started := time.Now().UTC()
	for {
		time.Sleep(60 * time.Second)

		for _, v := range s.Objects {
			lastWritten := v.GetOrSetLastWritten()
			if started.After(lastWritten) {
				v.RequiresPush = true
			}
			// TODO
		}

	}
}
