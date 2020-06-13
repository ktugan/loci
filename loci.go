package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/stdcopy"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//need commands
//build: to check if the docker build is outdated
//invoke command

// todo: feature - immutable volumes
// todo: feature - gitlab ci parser?
// todo: env variables

func tagFromDockerfile(dockerFilePath string) string {
	dockerFilePath, err := filepath.Abs(dockerFilePath)
	if err != nil {
		panic(err)
	}

	tag := strings.ReplaceAll(dockerFilePath, "/", "_")
	tag = strings.ToLower(tag)
	tag = strings.Trim(tag, "_")
	return tag
}

func BuildImage(cli *client.Client, dockerFilePath string) string {
	if dockerFilePath == "" {
		log.Fatal("Path to Dockerfile cannot be empty")
	}

	dockerFilePath, err := filepath.Abs(dockerFilePath)
	tag := tagFromDockerfile(dockerFilePath)

	base, filename := filepath.Split(dockerFilePath)

	ctx, _ := archive.TarWithOptions(base, &archive.TarOptions{})
	resp, err := cli.ImageBuild(context.Background(), ctx, types.ImageBuildOptions{
		Dockerfile: filename,
		Tags:       []string{tag},
	})
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	newStr := buf.String()

	fmt.Printf(newStr)
	return tag
}

type Volume struct {
	Source    string
	Target    string
	Immutable bool
}

type Config struct {
	BuildFolder          string //default: /build
	ImmutableBuildFolder bool   //default: false
	ExtraVolumes         []Volume
	Dockerfile           string //either set Dockerfile or Image, not both... I think
	Rebuild              bool
	Image                string //either set Dockerfile or Image, not both... I think
	EnvironmentVars      map[string]string
	Command              string
}

func checkImageExists(cli *client.Client, image string) bool {
	//todo check if a specific image exists in the `docker image` command
	return false
}

func prepConfig(config *Config) {
	if config.BuildFolder == "" {
		config.BuildFolder = "/build"
	}
	if config.Dockerfile != "" && config.Image != "" {
		log.Fatal("Cannot both define parameter Dockerfile and Image.")
	}
	//todo check if tag already exists and if yes exists EXCEPT the rebuild flag is set

	//todo docker pull if image

}

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	//todo need to build or build image first
	//-------------- PARAMETER PART --------------------
	config := Config{
		BuildFolder:          "",
		ImmutableBuildFolder: false,
		ExtraVolumes:         nil,
		Dockerfile:           "/home/kadir/workspace/kubernetes/docker-automation/Dockerfile",
		Image:                "",
		EnvironmentVars:      nil,
		Command:              "ls",
	}
	prepConfig(&config)
	image := ""
	buildFlag := false
	if buildFlag {
		image = BuildImage(cli, config.Dockerfile)
	} else {
		image = tagFromDockerfile(config.Dockerfile)
	}

	// ----------- END PARAMETER PART --------
	if image == "" {
		log.Fatal("No tag given") // todo probably just temporary and needs to be put into the parser.
	}

	cwd, err := filepath.Abs(".")
	if err != nil {
		panic(err)
	}

	resp, err := cli.ContainerCreate(context.Background(), &container.Config{
		Image:      image,
		Cmd:        []string{"bash", "-c", "ls"},
		WorkingDir: "/build",
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			mount.Mount{
				Type:   mount.TypeBind,
				Source: "/home/kadir/.aws",
				Target: "/root/.ssh",
			},
			mount.Mount{
				Type:   mount.TypeBind,
				Source: cwd,
				Target: "/build",
			},
		},
	}, nil, "")
	if err != nil {
		panic(err)
	}

	err = cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{})
	if err != nil {
		panic(err)
	}

	out, err := cli.ContainerLogs(context.Background(), resp.ID, types.ContainerLogsOptions{
		Follow: true, ShowStdout: true, ShowStderr: true, Tail: "all",
	})
	if err != nil {
		panic(err)
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	//fmt.Printf(resp.ID, resp.Warnings)

}
