package commands

import (
	"fmt"
	"log"

	"github.com/alexellis/faas-cli/builder"
	"github.com/alexellis/faas-cli/stack"
	"github.com/spf13/cobra"
)

func init() {
	faasCmd.AddCommand(pushCmd)
}

// pushCmd handles pushing function container images to a remote repo
var pushCmd = &cobra.Command{
	Use:   "push -f YAML_FILE",
	Short: "Push OpenFaaS function(s) defined in YAML (currently to Docker Hub)",
	Long: `Pushes the OpenFaaS function container image(s) defined in the
supplied YAML config to a remote repository.

These container images must already be present in your local image
cache.

NOTE - this command currently supports pushing to docker hub only,
	   support for additional container repos is planned.`,

	Example: `  faas-cli push -f https://raw.githubusercontent.com/alexellis/faas-cli/master/samples.yml
  faas-cli push -f ./samples.yml`,
	Run: runPush,
}

func runPush(cmd *cobra.Command, args []string) {

	var services stack.Services
	if len(yamlFile) > 0 {
		parsedServices, err := stack.ParseYAML(yamlFile)
		if err != nil {
			log.Fatalln(err.Error())
			return
		}

		if parsedServices != nil {
			services = *parsedServices
		}
	}

	if len(services.Functions) > 0 {
		for k, function := range services.Functions {
			function.Name = k
			fmt.Printf("Pushing: %s to remote repository.\n", function.Name)
			pushImage(function.Image)
		}
	} else {
		fmt.Println("You must supply a valid YAML file.")
		return
	}

}

func pushImage(image string) {
	builder.ExecCommand("./", []string{"docker", "push", image})
}
