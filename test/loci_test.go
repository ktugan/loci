package test

import (
	"testing"

	"github.com/ktugan/loci/localci"
)

func TestDockerfileNotExist(t *testing.T) {
	// We test if an error is correctly thrown.
	c := localci.LociConfig{
		BuildFolder:          "",
		ImmutableBuildFolder: false,
		ExtraVolumes:         nil,
		Dockerfile:           "NOT_EXISTING_DOCKERFILE",
		Image:                "",
		EnvironmentVars:      nil,
		Command:              "ls",
	}
	err := localci.PrepConfig(&c)
	if err, ok := err.(*localci.BadConfig); ok {
		if err.Message != "Dockerfile does not exist at given path" {
			t.Fatalf("%s", "Should fail, does not.")
		}

	}

}
