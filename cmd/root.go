package cmd

import (
	"errors"
	"os"

	"github.com/brownhounds/openapi-tsgen/schema"
	"github.com/brownhounds/openapi-tsgen/version"
	"github.com/spf13/cobra"
)

var errOutputPathRequired = errors.New("output path is required")

var rootCmd = &cobra.Command{
	Use:           "openapi-tsgen [schema.yml]",
	Short:         "Generate TypeScript types from an OpenAPI schema",
	Version:       version.Ver,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		schema.CLIVersion = cmd.Version
		in, err := cmd.Flags().GetString("schema")
		if err != nil {
			return err
		}
		if in == "" && len(args) > 0 {
			in = args[0]
		}
		if in == "" {
			_ = cmd.Help()
			return nil
		}

		out, err := cmd.Flags().GetString("output")
		if err != nil {
			return err
		}
		if out == "" {
			return errOutputPathRequired
		}

		inputJSON, err := cmd.Flags().GetBool("input-json")
		if err != nil {
			return err
		}

		format := schema.InputYAML
		if inputJSON {
			format = schema.InputJSON
		}

		if err := schema.WriteSchema(in, out, format); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.Flags().StringP("schema", "s", "", "Path to OpenAPI schema (YAML)")
	rootCmd.Flags().StringP("output", "o", "type.ts", "Output file path")
	rootCmd.Flags().Bool("input-json", false, "Treat schema input as JSON")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
