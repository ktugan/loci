package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type Volume struct {
	Source    string
	Target    string
	Immutable bool
}

type LociConfig struct {
	BuildFolder          string //default: /build
	ImmutableBuildFolder bool   //default: false
	ExtraVolumes         []Volume
	Dockerfile           string //either set Dockerfile or Image, not both... I think
	Rebuild              bool
	Image                string //either set Dockerfile or Image, not both... I think
	EnvironmentVars      map[string]string
	Command              string
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func prepConfig(config *LociConfig) error {
	if config.BuildFolder == "" {
		config.BuildFolder = "/build"
	}
	if config.Dockerfile != "" && config.Image != "" {
		return &BadConfig{
			Message:   "Cannot set parameter Dockerfile and Image at the same time",
			ExtraInfo: "",
		}
	}
	if config.Dockerfile == "" && config.Image == "" { //both unset
		return &BadConfig{
			Message:   "Need to define one parameter: Dockerfile or Image",
			ExtraInfo: "",
		}
	}

	if config.Dockerfile != "" { // Dockerfile set
		absPath, err := filepath.Abs(config.Dockerfile)
		if err != nil {
			panic(err)
		}
		config.Dockerfile = absPath
		if !fileExists(config.Dockerfile) {
			return &BadConfig{
				Message:   "Dockerfile does not exist at given path",
				ExtraInfo: config.Dockerfile,
			}
		}
		config.Image = tagFromDockerfile(config.Dockerfile)
	}
	return nil
}

type BadConfig struct {
	Message   string
	ExtraInfo string
}

func (e *BadConfig) Error() string {
	return fmt.Sprintf("%s: %s", e.Message, e.ExtraInfo)
}
