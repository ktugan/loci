package main

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"os"
	"path/filepath"
)

//need commands
//build: to check if the docker build is outdated
//invoke command

// todo: feature - immutable volumes
// todo: feature - gitlab ci parser?
// todo: env variables

func loci(config LociConfig) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	//Check if rebuild is set or if we dont have the image locally. On either of these conditions pull or build.
	if config.Rebuild || checkImageExists(cli, config.Image) {
		if config.Dockerfile != "" {
			BuildImage(cli, config)
		} else {
			PullImage(cli, config)
		}
	}

	// Get current cwd and mount it to /build
	cwd, err := filepath.Abs(".")
	if err != nil {
		panic(err)
	}
	resp, err := cli.ContainerCreate(context.Background(), &container.Config{
		Image:      config.Image,
		Cmd:        []string{"bash", "-c", "ls"}, //todo replace with config
		WorkingDir: "/build",
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			mount.Mount{ //todo replace by config
				Type:   mount.TypeBind,
				Source: "/home/kadir/.aws",
				Target: "/root/.ssh",
			},
			mount.Mount{
				Type:   mount.TypeBind,
				Source: cwd,
				Target: config.BuildFolder,
			},
		},
	}, nil, "")
	if err != nil {
		panic(err)
	}

	//Create Container now
	err = cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{})
	if err != nil {
		panic(err)
	}

	//Get logs and attach to stdout
	out, err := cli.ContainerLogs(context.Background(), resp.ID, types.ContainerLogsOptions{
		Follow: true, ShowStdout: true, ShowStderr: true, Tail: "all",
	})
	if err != nil {
		panic(err)
	}

	_, err = stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	if err != nil {
		panic(err)
	}
}

func main() {
	config := LociConfig{
		BuildFolder:          "",
		ImmutableBuildFolder: false,
		ExtraVolumes:         nil,
		Dockerfile:           "/home/kadir/workspace/kubernetes/docker-automation/Dockerfile",
		Image:                "",
		EnvironmentVars:      nil,
		Command:              "ls",
	}

	err := prepConfig(&config)
	if err != nil {
		panic(err)
	}

	loci(config)
}
