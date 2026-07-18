package toggle

import (
	"github.com/spf13/cobra"
)

type IYamlToggle struct {
	Profile    *YamlProfile    `yaml:"profile,omitempty"`
	Kit        *YamlKit        `yaml:"kit,omitempty"`
	Service    *YamlService    `yaml:"service,omitempty"`
	Dispatcher *YamlDispatcher `yaml:"dispatcher,omitempty"`
}

func YamlToggle() *cobra.Command {
	return &cobra.Command{
		Use:  "yaml-toggle [filename]",
		Long: "",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// data, err := os.ReadFile(args[0])
			// if err != nil {
			// 	panic(err)
			// }
			// var m IYamlToggle
			// if err := yaml.Unmarshal(data, &m); err != nil {
			// 	panic(err)
			// }
			// p := policy.UnPackInitPolicy(cmd)

			// addr, err := cmd.Flags().GetString(p.CMD_ORG.Key)

			// cli, err := client.New(addr, nil)
		},
	}
}
