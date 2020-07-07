package utils

import (
	"errors"
	"os"
	"path/filepath"
)

// Storage keeps all local objects
type Storage struct {
	Root    string
	Objects []FObject
}

func (s *Storage) createObject(path string, info os.FileInfo) {
	obj := FObject{}
	obj.Location = path
	obj.Name = info.Name()
	obj.TagSecure()
	// obj.TagRemoveAll()

	s.Objects = append(s.Objects, obj)
}

func (s *Storage) DisabllAllTags() {
	TagWg.Add(len(s.Objects))
	for _, obj := range s.Objects {
		obj.TagRemoveAll()
	}
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

	filepath.Walk(s.Root, s.walkRoot)

	return nil
}
