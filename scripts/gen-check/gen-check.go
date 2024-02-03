// sparrow
// (C) 2024, Deutsche Telekom IT GmbH
//
// Deutsche Telekom IT GmbH and all other contributors /
// copyright owners license this file to you under the Apache
// License, Version 2.0 (the "License"); you may not use this
// file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/spf13/cobra"
)

// CheckConfig holds the configuration for a new check
type CheckConfig struct {
	PackageName     string
	CheckStructName string
	CheckName       string
}

func main() {
	execute()
}

func execute() {
	rootCmd := &cobra.Command{
		Use:   "gen-check",
		Short: "Generate a new check for the sparrow",
	}
	rootCmd.AddCommand(NewCmdGenCheck())

	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func NewCmdGenCheck() *cobra.Command {
	var checkConfig CheckConfig

	cmd := &cobra.Command{
		Use:   "gen-check",
		Short: "Generate a new check",
		RunE:  runGenCheck(&checkConfig),
	}

	cmd.Flags().StringVarP(&checkConfig.PackageName, "package", "p", "", "Package name for the new check")
	cmd.Flags().StringVarP(&checkConfig.CheckStructName, "struct", "s", "", "Struct name for the new check")
	cmd.Flags().StringVarP(&checkConfig.CheckName, "name", "n", "", "Name of the new check")

	err := cmd.MarkFlagRequired("package")
	if err != nil {
		panic(err)
	}
	err = cmd.MarkFlagRequired("struct")
	if err != nil {
		panic(err)
	}
	err = cmd.MarkFlagRequired("name")
	if err != nil {
		panic(err)
	}

	return cmd
}

func runGenCheck(config *CheckConfig) func(cmd *cobra.Command, args []string) (err error) {
	checkPath := fmt.Sprintf("pkg/checks/%s/%s.go", config.PackageName, config.PackageName)
	return func(cmd *cobra.Command, args []string) (err error) {
		if _, serr := os.Stat(checkPath); serr == nil {
			return fmt.Errorf("check '%s' already exists, aborting to prevent overwriting", config.PackageName)
		}

		err = os.MkdirAll(filepath.Dir(checkPath), 0o755) //nolint:gomnd
		if err != nil {
			return err
		}

		tpl, err := template.ParseFiles("./scripts/gen-check/check_template.go.tpl")
		if err != nil {
			return err
		}

		file, err := os.Create(checkPath)
		if err != nil {
			return err
		}
		defer func() {
			cerr := file.Close()
			err = errors.Join(cerr, err)
		}()

		err = tpl.Execute(file, config)
		if err != nil {
			return err
		}

		fmt.Println("Successfully generated new check:", config.PackageName)
		return nil
	}
}
