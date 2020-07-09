package utils

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// Storage keeps all local objects
type Storage struct {
	Root    string
	Objects map[string]FObject
	Config  ConfigCloudStore
	User    User
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

func (s *Storage) createObject(path string, info os.FileInfo) {
	obj := FObject{}
	relativePath := strings.TrimPrefix(path, s.Config.BasePath)

	obj.Location = path
	obj.Relativepath = relativePath
	obj.Name = info.Name()
	obj.LastWritten = info.ModTime()

	obj.UpdateHashForObject(s.User, s.Config)
	obj.UpdateHashForObjectBlocks(s.User, s.Config)

	s.Objects[obj.Relativepath] = obj
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

	logrus.Println("All Objects Acknowledges and Marked")

	return nil
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
