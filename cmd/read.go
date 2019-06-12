// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	b64 "encoding/base64"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "read <path> ''",
	Short: "Read content at <path>",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		path := strings.TrimPrefix(args[0], "/"+PluginName+"/")
		org, tail := SplitPath(path)
		repo, path := SplitPath(tail)
		printFile(org, repo, path)
	},
}

func printFile(org string, repo string, path string) {
	fileContent, _, err := FetchRepositoryContent(org, repo, path)
	HandleError(err)
	decodedContent, _ := b64.StdEncoding.DecodeString(*fileContent.Content)
	fmt.Println(string(decodedContent))
}

func init() {
	rootCmd.AddCommand(readCmd)
}
