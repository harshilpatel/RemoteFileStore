package utils

import (
	"errors"
	"os"
	"path/filepath"
	"time"
)

// Storage keeps all local objects
type Storage struct {
	Root    string
	Objects map[string]FObject
	Config  ConfigCloudStore
	User    User
}

func CreateStorage(c ConfigCloudStore) Storage {
	return Storage{
		Root:    c.LocalBasePath,
		Objects: make(map[string]FObject),
		Config:  c,
		User: User{
			Username: c.User,
			Key:      c.Key,
		},
	}
}

func (s *Storage) createObject(path string, info os.FileInfo) {
	obj := FObject{}
	obj.Location = path
	obj.Name = info.Name()
	obj.TagSecure()

	s.Objects[obj.Location] = obj
}

func (s *Storage) DisabllAllTags() {
	TagWg.Add(len(s.Objects))
	for _, obj := range s.Objects {
		go obj.TagRemoveAll()
	}

	TagWg.Wait()
}

func (s *Storage) walkRoot(path string, info os.FileInfo, err error) error {
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

	filepath.Walk(s.Root, s.walkRoot)
	return nil
}

func (s *Storage) keepWatchingForChanges() {
	started := time.Now().UTC()
	for {
		time.Sleep(60 * time.Second)

		// now := time.Now().UTC()

		for _, v := range s.Objects {
			lastWritten := v.GetOrSetLastWritten()
			if started.After(lastWritten) {
				v.RequiresPushed = true
			}

			if v.LastPushed.After()
		}

	}
}
