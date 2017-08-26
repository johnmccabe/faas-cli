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
	deleteCmd.Flags().StringVar(&functionName, "name", "", "Name of the deployed function")
	deleteCmd.Flags().StringVar(&gateway, "gateway", "http://localhost:8080", "Gateway URI - defaults to http://localhost:8080")

	faasCmd.AddCommand(deleteCmd)
}

// faasCmd is the FaaS CLI root command and mimics the legacy client behaviour
// Every other command attached to FaasCmd is a child command to it.
var deleteCmd = &cobra.Command{
	Use:     "delete (FUNCTION_NAME | -f YAML_FILE)",
	Aliases: []string{"rm", "remove"},
	Short:   "Delete deployed FaaS function(s)",
	Long: `Deletes deployed OpenFaaS functions either via the supplied 
YAML config using the "--yaml" flag (which may contain multiple
function definitions), or by explicitly specifying a function
name.`,
	Example: `  faas-cli delete -f https://raw.githubusercontent.com/alexellis/faas-cli/master/samples.yml
  faas-cli delete -f ./samples.yml
  faas-cli delete url-ping`,
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
