package commands

import "github.com/alexellis/faas-cli/builder"

func pushImage(image string) {
	builder.ExecCommand("./", []string{"docker", "push", image})
}
