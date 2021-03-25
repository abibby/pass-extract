package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var rootCmd = &cobra.Command{
	Use:   "pass-extract <pass-name> <key>",
	Short: "Pass extract is a simple utility to extract extra data from pass entries.",
	Long:  ``,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		passName := args[0]
		key := args[1]

		b, err := exec.Command("pass", passName).Output()
		if err != nil {
			if err, ok := err.(*exec.ExitError); ok {
				fmt.Fprintf(os.Stderr, "%s\n", err.Stderr)
			}
			return err
		}
		re := regexp.MustCompile(`(^[^\n]*)\n---*\n([^$]*)$`)
		matches := re.FindSubmatch(b)
		rawYaml := matches[2]
		out := map[string]interface{}{}
		err = yaml.Unmarshal(rawYaml, &out)

		if err != nil {
			return err
		}
		value, ok := out[key]
		if !ok {
			// spew.Dump(value)
			return fmt.Errorf("no key '%s' in data", key)
		}

		fmt.Printf("%v\n", value)
		return nil
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
