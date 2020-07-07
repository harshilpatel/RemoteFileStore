package utils

import (
	"os/exec"
	"sync"
	"time"
)

var TagWg sync.WaitGroup

type Tag struct {
	Expired    string `"Red`
	Consistent string `Green`
	Processing string `Orange`
}

func (t *Tag) AddTagToFile(path string, color string) bool {

	time.Sleep(5 * time.Millisecond)
	cmd := exec.Command("tag", "-a", color, path)
	cmd.Run()

	// TagWg.Done()
	return true
}

func (t *Tag) DelTagToFile(path string, color string) bool {
	time.Sleep(5 * time.Millisecond)
	cmd := exec.Command("tag", "-r", color, path)
	cmd.Run()

	// TagWg.Done()
	return true
}

func (t *Tag) RemoveAll(path string) bool {
	time.Sleep(5 * time.Millisecond)
	cmd := exec.Command("tag", "-r", t.Expired, path)
	cmd.Run()
	cmd = exec.Command("tag", "-r", "Green", path)
	cmd.Run()
	cmd = exec.Command("tag", "-r", "Orange", path)
	cmd.Run()

	// TagWg.Done()
	return true
}

func (t *Tag) UpdateTagToFile(path string, color string) bool {
	time.Sleep(5 * time.Millisecond)
	t.RemoveAll(path)

	cmd := exec.Command("tag", "-a", color, path)
	cmd.Run()

	// TagWg.Done()
	return true
}
