package cmd

import (
	"encoding/json"
	"okp4/template-go/internal/version"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

const flagLong = "long"
const flagOutput = "output"

// NewVersionCommand returns a CLI command to interactively print the application binary version information.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the application binary version information",
	RunE: func(cmd *cobra.Command, args []string) error {
		verInfo := version.NewInfo()

		if long, _ := cmd.Flags().GetBool(flagLong); !long {
			cmd.Println(verInfo.Version)
			return nil
		}

		var (
			bz  []byte
			err error
		)

		output, _ := cmd.Flags().GetString(flagOutput)
		switch strings.ToLower(output) {
		case "json":
			bz, err = json.Marshal(verInfo)

		default:
			bz, err = yaml.Marshal(&verInfo)
		}

		if err != nil {
			return err
		}

		cmd.Println(string(bz))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	versionCmd.Flags().Bool(flagLong, false, "Print long version information")
	versionCmd.Flags().StringP(flagOutput, "o", "text", "Output format (text|json)")
}
