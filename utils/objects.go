package utils

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	server "github.com/harshilkumar/cloud-store-server/utils"
)

type FObject struct {
	server.FObject
}

type Range struct {
	Start int64
	End   int64
}

func (f *FObject) readData() error {
	return nil
}

func (f *FObject) WriteData(data []byte) error {
	os.Truncate(f.Location, 0)

	file, err := os.OpenFile(f.Location, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	defer file.Close()

	if err != nil {
		return err
	}
	file.Write(data)

	return nil
}

func (f *FObject) EditData(from int64, data []byte) error {
	file, err := os.OpenFile(f.Location, os.O_CREATE|os.O_RDWR, os.ModePerm)
	defer file.Close()

	if err != nil {
		return err
	}
	file.WriteAt(data, from)

	return nil
}

func (f *FObject) appendDate(data []int) error {
	return nil
}

func (f *FObject) replaceData(old Range, to Range, data []int) error {
	return nil
}

func (f *FObject) getParent() string {
	return filepath.Dir(f.Location)
}

func (f *FObject) getChildren() error {
	return nil
}

func (f *FObject) TagSecure(c chan int) {
	t := Tag{}
	t.AddTagToFile(f.Location, "Green")

	c <- 1
}

func (f *FObject) TagRemoveAll(c chan int) {
	t := Tag{}
	t.RemoveAll(f.Location)
	c <- 1
}

func (f *FObject) TagNotSecure(c chan int) {
	t := Tag{}
	t.AddTagToFile(f.Location, "Red")
	c <- 1
}

func (f *FObject) TagWorking(c chan int) {
	t := Tag{}
	t.AddTagToFile(f.Location, "Orange")
	c <- 1
}

func (f *FObject) GetOrSetLastWritten() time.Time {
	if f.LastWritten.IsZero() {
		f.LastWritten = time.Now().UTC()
	}

	if f.LastPulled.IsZero() {
		f.LastPulled = time.Now().UTC()
	}

	return f.LastWritten
}

func (f *FObject) GetOrSetLastPulled() time.Time {

	if f.LastPulled.IsZero() {
		f.LastPulled = time.Now().UTC()
	}

	return f.LastPulled
}

func (f *FObject) GetRealPath(u User, c ConfigCloudStore) string {
	return filepath.Join(c.BasePath, f.Relativepath)
}

func (f *FObject) CreateHashForObject(u User, c ConfigCloudStore) ([]byte, error) {
	realPath := f.GetRealPath(u, c)
	if data, err := ioutil.ReadFile(realPath); err == nil {
		h := sha256.New()
		h.Write(data)
		return h.Sum(nil), nil
	}

	return nil, errors.New("Could not find the file " + realPath)
}

func (f *FObject) CreateHashForObjectBlocks(u User, c ConfigCloudStore) ([][]byte, error) {
	realPath := f.GetRealPath(u, c)

	hash := make([][]byte, 0)
	buf := make([]byte, 500)
	if file, err := os.Open(realPath); err == nil {
		defer file.Close()
		for {
			if n, e := file.Read(buf); e == nil {
				if n > 0 {
					h := sha256.New()
					h.Write(buf)
					res := h.Sum(nil)
					hash = append(hash, res)
				}
			} else if e == io.EOF {
				return hash, nil
			}
		}

	}

	return nil, errors.New("Something went wrong")
}

func (f *FObject) UpdateHashForObject(u User, c ConfigCloudStore) {
	if h, err := f.CreateHashForObject(u, c); err == nil {
		f.HashOfFile = h
	} else {
		fmt.Println(err)
	}
}

func (f *FObject) UpdateHashForObjectBlocks(u User, c ConfigCloudStore) {
	if h, err := f.CreateHashForObjectBlocks(u, c); err == nil {
		f.Hash = h
	} else {
		fmt.Println(err)
	}
}
