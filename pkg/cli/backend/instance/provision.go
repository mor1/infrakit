package instance

import (
	"fmt"

	"github.com/docker/infrakit/pkg/cli/backend"
	"github.com/docker/infrakit/pkg/run/scope"
	"github.com/docker/infrakit/pkg/spi/instance"
	"github.com/docker/infrakit/pkg/types"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func init() {
	backend.Register("instanceProvision", Provision,
		func(flags *pflag.FlagSet) {
			flags.String("plugin", "", "plugin")
		})
}

// Provision returns an executable function based on that specification to call the named instance plugin's provision
// method. The optional parameter in the playbook script can be overridden by the value of the `--plugin` flag
// in the command line.
func Provision(scope scope.Scope, test bool, opt ...interface{}) (backend.ExecFunc, error) {

	return func(script string, cmd *cobra.Command, args []string) error {

		var name string

		// Optional parameter for plugin name can be overridden by the value of the flag (--plugin):
		if len(opt) > 0 {
			s, is := opt[0].(string)
			if !is {
				return fmt.Errorf("first param (pluginName) must be string")
			}
			name = s
		}
		if n, err := cmd.Flags().GetString("plugin"); err != nil {
			return err
		} else if n != "" {
			name = n
		}

		plugin, err := scope.Instance(name)
		if err != nil {
			return err
		}

		spec := instance.Spec{}
		if err := types.Decode([]byte(script), &spec); err != nil {
			return err
		}

		id, err := plugin.Provision(spec)
		if err != nil {
			return err
		}
		fmt.Println(*id)
		return nil
	}, nil
}
