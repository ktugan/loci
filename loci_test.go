package main

import (
	"testing"
)

func TestDockerfileNotExist(t *testing.T) {
	// We test if an error is correctly thrown.
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
		if err.Message != "Dockerfile does not exist at given path" {
			t.Fatalf("%s", "Should fail, does not.")
		}

	}

}
