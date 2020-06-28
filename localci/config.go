package localci

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type LociConfig struct {
	Image                string            `yaml:"image"`
	Dockerfile           string            `yaml:"dockerfile"`
	BuildFolder          string            `yaml:"build_folder"`           //default: /build
	ImmutableBuildFolder bool              `yaml:"immutable_build_folder"` //default: false
	ExtraVolumes         []Volume          `yaml:"extra_volumes"`
	Rebuild              bool              `yaml:"rebuild"`
	EnvironmentVars      map[string]string `yaml:"environment"`
	Command              string            `yaml:"command"`
}

type Volume struct {
	Source    string `yaml:"source"`
	Target    string `yaml:"target"`
	Immutable bool   `yaml:"immutable"`
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func PrepConfig(config *LociConfig) error {
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

func LoadConfig(fname string) LociConfig {
	if !fileExists(fname) {
		log.Fatal(".loci.yml not found in current folder.")
	}

	f, err := os.Open(fname)
	if err != nil {
		panic(err)
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	var config LociConfig
	err = yaml.Unmarshal(b, &config)
	if err != nil {
		panic(err)
	}

	return config
}
