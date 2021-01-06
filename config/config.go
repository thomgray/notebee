package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/thomgray/notebee/util"
)

// Conf ...
type Conf struct {
	DefaultRoot *string
}

// Config ...
type Config struct {
	SearchPaths    []string
	currentDocRoot *string
	conf           Conf
	// NotePaths   []string
}

// MakeConfig ...
func MakeConfig() *Config {
	return (&Config{}).Init()
}

// Init ...
func (c *Config) Init() *Config {
	c.SearchPaths = loadSeachPaths()
	confFileP := FilePath()

	if _, err := os.Stat(confFileP); err == nil {
		confBytes, err2 := ioutil.ReadFile(confFileP)
		if err2 == nil {
			var conf Conf
			json.Unmarshal(confBytes, &conf)
			c.conf = conf
			c.currentDocRoot = conf.DefaultRoot
		} else {
			log.Panicln(err2)
		}
	} else {
		// log.Panicln(err)
	}
	return c
}

// SetCurrentDocRoot ...
func (c *Config) SetCurrentDocRoot(p string) {
	c.currentDocRoot = &p
}

// SetDefaultDocRoot ...
func (c *Config) SetDefaultDocRoot(p string) {
	c.conf.DefaultRoot = &p
	c.writeConfig()
}

// DocumentRoot ...
func (c *Config) DocumentRoot() *string {
	return c.currentDocRoot
}

func loadSeachPaths() []string {
	bytes, _ := util.ReadFile(NotePathsPath())
	paths := util.ReadLines(bytes)
	return paths
}

func (c *Config) writeConfig() {
	serlaised, err := json.Marshal(c.conf)
	if err == nil {
		ioutil.WriteFile(FilePath(), serlaised, 0644)
	}
}

// func loadNotePaths(searchPaths []string) []string {
// 	files := make([]string, 0)

// 	for _, sp := range searchPaths {
// 		f, err := ioutil.ReadDir(sp)
// 		if err == nil {
// 			for _, file := range f {
// 				if file.Mode().IsRegular() && filepath.Ext(file.Name()) == ".md" {
// 					files = append(files, fmt.Sprintf("%s/%s", sp, file.Name()))
// 				}
// 				log.Printf("File %s\n", fmt.Sprintf("%s/%s", sp, file.Name()))
// 			}
// 		}
// 	}
// 	return files
// }

var _homedir *string = nil

// Directory ...
func Directory() string {
	return filepath.Join(GetAppConfig().HomeDir, ".notebee")
}

// FilePath ...
func FilePath() string {
	return filepath.Join(Directory(), "config")
}

// NotePathsPath ...
func NotePathsPath() string {
	return filepath.Join(Directory(), "paths")
}

// AddSearchPath ...
func (c *Config) AddSearchPath(sp string) {
	c.SearchPaths = append(c.SearchPaths, sp)
	c.updateSearchPathConfig()
}

func (c *Config) updateSearchPathConfig() {
	serlaised := []byte(strings.Join(c.SearchPaths, "\n"))
	ioutil.WriteFile(NotePathsPath(), serlaised, 0644)
}

// RemoveSearchPath ...
func (c *Config) RemoveSearchPath(sp string) {
	for i, p := range c.SearchPaths {
		if p == sp {
			newSp := append(c.SearchPaths[:i], c.SearchPaths[i+1:]...)
			c.SearchPaths = newSp
		}
	}
	c.updateSearchPathConfig()
}

// ReloadNotes ...
// func (c *Config) ReloadNotes() {
// 	c.NotePaths = loadNotePaths(c.SearchPaths)
// }
