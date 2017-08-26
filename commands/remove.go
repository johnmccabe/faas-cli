package commands

import (
	"fmt"
	"log"

	"github.com/alexellis/faas-cli/proxy"
	"github.com/alexellis/faas-cli/stack"
	"github.com/spf13/cobra"
)

func init() {
	// Setup flags that are used by multiple commands (variables defined in faas.go)
	removeCmd.Flags().StringVar(&functionName, "name", "", "Name of the deployed function")
	removeCmd.Flags().StringVar(&gateway, "gateway", "http://localhost:8080", "Gateway URI - defaults to http://localhost:8080")

	faasCmd.AddCommand(removeCmd)
}

// removeCmd deletes/removes OpenFaaS function containers
var removeCmd = &cobra.Command{
	Use:     "remove (FUNCTION_NAME | -f YAML_FILE)",
	Aliases: []string{"rm", "delete"},
	Short:   "Remove deployed FaaS function(s)",
	Long: `Removes/deletes deployed OpenFaaS functions either via the supplied 
YAML config using the "--yaml" flag (which may contain multiple
function definitions), or by explicitly specifying a function
name.`,
	Example: `  faas-cli remove -f https://raw.githubusercontent.com/alexellis/faas-cli/master/samples.yml
  faas-cli remove -f ./samples.yml
  faas-cli remove url-ping`,
	Run: runDelete,
}

func runDelete(cmd *cobra.Command, args []string) {
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
			fmt.Printf("Deleting: %s.\n", function.Name)

			proxy.DeleteFunction(services.Provider.GatewayURL, function.Name)
		}
	} else {
		if len(args) < 1 {
			fmt.Println("Please provide the name of a function to delete")
			return
		}
		functionName = args[0]
		fmt.Printf("Deleting: %s.\n", functionName)
		proxy.DeleteFunction(gateway, functionName)
	}
}
