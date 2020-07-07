package utils

import (
	"path/filepath"
	"time"
)

type FObject struct {
	Name     string
	Location string

	IsDir    bool
	IsBinary bool

	Lastwritten time.Time
	Lastpulled  time.Time

	Version int16
}

type Range struct {
	Start int64
	End   int64
}

func (f *FObject) readData() error {
	return nil
}

func (f *FObject) writeDate(from int64, data []int) error {
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

func (f *FObject) TagSecure() {
	TagWg.Add(1)
	t := Tag{}
	t.UpdateTagToFile(f.Location, "Green")
}

func (f *FObject) TagRemoveAll() {
	TagWg.Add(1)
	t := Tag{}
	t.RemoveAll(f.Location)
}

func (f *FObject) TagNotSecure() {
	TagWg.Add(1)
	t := Tag{}
	t.UpdateTagToFile(f.Location, "Red")
}

func (f *FObject) TagWorking() {
	t := Tag{}
	t.UpdateTagToFile(f.Location, "Orange")
}
