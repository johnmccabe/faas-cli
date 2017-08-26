package commands

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/alexellis/faas-cli/proxy"
	"github.com/alexellis/faas-cli/stack"
	"github.com/spf13/cobra"
)

// Flags that are to be added to commands.

var (
	envvarOpts []string
	replace    bool
)

func init() {
	// Setup flags that are used by multiple commands (variables defined in faas.go)
	deployCmd.Flags().StringVar(&fprocess, "fprocess", "", "Fprocess to be run by the watchdog")
	deployCmd.Flags().StringVar(&gateway, "gateway", "http://localhost:8080", "Gateway URI (defaults to http://localhost:8080)")
	deployCmd.Flags().StringVar(&handler, "handler", "", "Handler for function, i.e. handler.js")
	deployCmd.Flags().StringVar(&image, "image", "", "Docker image name to build")
	deployCmd.Flags().StringVar(&language, "lang", "node", "Programming language template (defaults  node)")
	deployCmd.Flags().StringVar(&functionName, "name", "", "Name of the deployed function")

	// Setup flags that are used only by this command (variables defined above)
	deployCmd.Flags().StringArrayVarP(&envvarOpts, "env", "e", []string{}, "Set environment variable ENVVAR=VALUE. Can be repeated.")
	deployCmd.Flags().BoolVar(&replace, "replace", false, "Replace any existing function (defaults to false)")

	// Set bash-completion.
	_ = deployCmd.Flags().SetAnnotation("handler", cobra.BashCompSubdirsInDir, []string{})

	faasCmd.AddCommand(deployCmd)
}

// deployCmd handles deploying OpenFaaS function containers
var deployCmd = &cobra.Command{
	Use: "deploy (-f YAML_FILE | --image IMAGE_NAME --name FUNCTION_NAME [--lang <ruby|python|node|csharp>] [--handler DIR] [--env ENVVAR=VALUE ...]) [--replace]",

	Short: "Deploy OpenFaaS functions",
	Long: `Deploys OpenFaaS function containers either via the supplied 
YAML config using the "--yaml" flag (which may contain multiple
function definitions), or directly via flags.

Pass the --replace flag to overwrite existing functions.`,
	Example: `  faas-cli deploy -f https://raw.githubusercontent.com/alexellis/faas-cli/master/samples.yml
  faas-cli deploy -f ./samples.yml
  faas-cli deploy -f ./samples.yml --replace
  faas-cli deploy --image=alexellis/faas-url-ping --name=url-ping
  faas-cli deploy --image=alexellis/faas-url-ping --name=url-ping --lang=python --hander=./url-ping/ --env=MYVAR=myval --env=MYOTHERVAR=myotherval --replace`,
	Run: runDeploy,
}

func runDeploy(cmd *cobra.Command, args []string) {
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
		if len(services.Provider.Network) == 0 {
			services.Provider.Network = defaultNetwork
		}

		for k, function := range services.Functions {
			function.Name = k
			fmt.Printf("Deploying: %s.\n", function.Name)

			proxy.DeployFunction(function.FProcess, services.Provider.GatewayURL, function.Name, function.Image, function.Language, replace, function.Environment, services.Provider.Network)
		}
	} else {
		if len(image) == 0 {
			fmt.Println("Please provide a --image to be deployed.")
			return
		}
		if len(functionName) == 0 {
			fmt.Println("Please provide a --name for your function as it will be deployed on FaaS")
			return
		}

		envvars, err := parseEnvvars(envvarOpts)
		if err != nil {
			fmt.Printf("Error parsing envvars: %v\n", err)
			os.Exit(1)
		}

		proxy.DeployFunction(fprocess, gateway, functionName, image, language, replace, envvars, defaultNetwork)
	}
}

func parseEnvvars(envvars []string) (map[string]string, error) {
	result := map[string]string{}
	for _, envvar := range envvars {
		s := strings.SplitN(strings.TrimSpace(envvar), "=", 2)
		envvarName := s[0]
		envvarValue := s[1]

		if !(len(envvarName) > 0) {
			return nil, fmt.Errorf("Empty envvar name: [%s]", envvar)
		}
		if !(len(envvarValue) > 0) {
			return nil, fmt.Errorf("Empty envvar value: [%s]", envvar)
		}

		result[envvarName] = envvarValue
	}
	return result, nil
}
