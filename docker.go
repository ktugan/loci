package loci

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/archive"

	"log"
	"reflect"

	"github.com/docker/docker/client"
)

func checkImageExists(cli *client.Client, image string) bool {
	inspect, _, err := cli.ImageInspectWithRaw(context.Background(), image)

	if err != nil {

		// Not sure if this is the right way, I tried to check for the type directly without reflect
		// but couldn't access the internal error struct type.
		if reflect.TypeOf(err).String() == "client.objectNotFoundError" {
			// Docker image not found.
			return false
		} else {
			// Some other error.
			panic(err)
		}
	}

	// Image ID found
	if inspect.ID != "" {
		return true
	}

	log.Fatal("Reaching end of function but shouldn't. Please verify.")
	return false
}

type dockerBuildJson struct {
	Content string `json:"stream"`
}

func BuildImage(cli *client.Client, config LociConfig) string {

	if config.Dockerfile == "" {
		log.Fatal("Path to Dockerfile cannot be empty")
	}
	if config.Image == "" {
		log.Fatal(
			"Image is not set but should be by convention. " +
				"Shouldn't reach here except there is a gap in the code flow.")
	}

	base, filename := filepath.Split(config.Dockerfile)

	ctx, _ := archive.TarWithOptions(base, &archive.TarOptions{})
	resp, err := cli.ImageBuild(context.Background(), ctx, types.ImageBuildOptions{
		Dockerfile: filename,
		Tags:       []string{config.Image},
	})
	if err != nil {
		panic(err)
	}

	buffer := new(bytes.Buffer)
	_, err = buffer.ReadFrom(resp.Body)
	if err != nil {
		panic(err)
	}

	resultLines := strings.Split(buffer.String(), "\n")

	//var js map[string]interface{}

	var parsed dockerBuildJson
	for _, line := range resultLines {
		_ = json.Unmarshal([]byte(line), &parsed)
		if err == nil {
			fmt.Printf("%s", parsed.Content)
		}
	}

	return config.Image
}

func tagFromDockerfile(dockerFilePath string) string {
	// Tags only support 128 chars. Need some other def. here.
	dockerFilePath, err := filepath.Abs(dockerFilePath)
	if err != nil {
		panic(err)
	}

	tag := strings.ReplaceAll(dockerFilePath, "/", "_")
	tag = strings.ToLower(tag)
	tag = strings.Trim(tag, "_")
	return tag
}

func PullImage(cli *client.Client, config LociConfig) string {
	reader, err := cli.ImagePull(context.Background(), config.Image, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(os.Stdout, reader)
	if err != nil {
		panic(err)
	}

	return config.Image
}
