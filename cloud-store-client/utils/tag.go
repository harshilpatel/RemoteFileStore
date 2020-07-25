package utils

import (
	"os/exec"
	"sync"

	"github.com/sirupsen/logrus"
)

var TagWg sync.WaitGroup

type Tag struct {
	Expired    string `"Red`
	Consistent string `Green`
	Processing string `Orange`
}

func (t *Tag) AddTagToFile(path string, color string) bool {

	cmd := exec.Command("tag", "-a", color, path)
	if err := cmd.Run(); err != nil {
		logrus.Printf("Could not tag %v", path)
		return false
	}
	logrus.Printf("tagged: %v %v\n", path, color)
	return true
}

func (t *Tag) DelTagToFile(path string, color string) bool {
	logrus.Printf("tagging: %v %v\n", path, color)
	cmd := exec.Command("tag", "-r", color, path)
	if err := cmd.Run(); err != nil {
		logrus.Printf("Could not tag %v", path)
		return false
	}

	return true
}

func (t *Tag) RemoveAll(path string) bool {
	logrus.Printf("tagging: %v %v\n", path, "\\*")
	// cmd := exec.Command("tag", "-r", "Red", path)
	cmd := exec.Command("tag", "-r", "Green", path)
	if err := cmd.Run(); err != nil {
		logrus.Printf("Could not tag %v", path)
		return false
	}
	// cmd = exec.Command("tag", "-r", "Orange", path)

	logrus.Printf("untagged: %v\n", path)
	return true
}

func (t *Tag) UpdateTagToFile(path string, color string) bool {
	// time.Sleep(5 * time.Millisecond)
	logrus.Printf("tagging: %v %v\n", path, color)
	t.RemoveAll(path)

	cmd := exec.Command("tag", "-a", color, path)
	if err := cmd.Run(); err != nil {
		logrus.Printf("Could not tag %v", path)
		return false
	}

	return true
}
