package main

import (
	"testing"

	"gotest.tools/assert"
)

func TestDockerfileNotExist(t *testing.T) {
	//Should throw error if dockerfile does not exist.
	c := LociConfig{
		BuildFolder:          "",
		ImmutableBuildFolder: false,
		ExtraVolumes:         nil,
		Dockerfile:           "NOT_EXISTING_DOCKERFILE",
		Image:                "",
		EnvironmentVars:      nil,
		Command:              "ls",
	}
	err := prepConfig(&c)
	if err, ok := err.(*BadConfig); ok {
		assert.Equal(t, err.Message, "Dockerfile does not exist at given path")
	}
	print("hello")

}
