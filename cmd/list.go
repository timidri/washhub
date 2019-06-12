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
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/github"
	"github.com/spf13/cobra"
	"github.com/tidwall/pretty"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list <path> ''",
	Short: "List content at <path>",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		path := strings.TrimSuffix(args[0], "/")
		user, _, err := GithubClient().Users.Get(context.Background(), "")
		HandleError(err)

		if path == "/"+PluginName { // we are at top level
			listOrganisations(user)
		} else {
			path = strings.TrimPrefix(path, "/"+PluginName+"/")
			if strings.Contains(path, "/") { // we have at least org and repo in the path
				listContent(path)
			} else { // we have only an org or a user in the path
				if user.GetLogin() == path {
					// list user's repos
					listUserRepos(path)
				} else {
					// list org's repos
					listOrgRepos(path)
				}
			}
		}
		os.Exit(0)
	},
}

type attributes struct {
	Size int `json:"size"`
}

type entry struct {
	Name       string     `json:"name"`
	Methods    []string   `json:"methods"`
	Attributes attributes `json:"attributes"`
}

func listContent(path string) {
	org, tail := SplitPath(path)
	repo, filePath := SplitPath(tail)
	_, dirContent, err := FetchRepositoryContent(org, repo, filePath)
	HandleError(err)
	entries := directoryToEntries(dirContent)
	PrintEntries(entries)
}

func listOrgRepos(org string) {
	opt := &github.RepositoryListByOrgOptions{
		// setting PerPage to 99 instead of 100 to work around
		// https://github.com/google/go-github/issues/999
		// but work-around doesn't seem to work...
		ListOptions: github.ListOptions{PerPage: 99},
	}
	// get all pages of results
	var allRepos []*github.Repository
	for {
		repos, resp, err := GithubClient().Repositories.ListByOrg(context.Background(), org, opt)
		HandleError(err)

		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	entries := reposToEntries(allRepos)
	PrintEntries(entries)
}

func listUserRepos(userName string) {
	opt := &github.RepositoryListOptions{
		Affiliation: "owner", // we want only owner's repos
		// setting PerPage to 99 instead of 100 to work around
		// https://github.com/google/go-github/issues/999
		// but work-around doesn't seem to work...
		ListOptions: github.ListOptions{PerPage: 99},
	}
	// get all pages of results
	var allRepos []*github.Repository
	for {
		// passing empty string for user name to get repos for authenticated user
		repos, resp, err := GithubClient().Repositories.List(context.Background(), "", opt)
		HandleError(err)

		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	entries := reposToEntries(allRepos)
	PrintEntries(entries)
}

func listOrganisations(user *github.User) {
	orgs, _, err := GithubClient().Organizations.ListOrgMemberships(context.Background(), nil)
	HandleError(err)

	var entries []*entry

	userEntry := &entry{
		Name:    user.GetLogin(),
		Methods: []string{"list"},
	}

	entries = append(entries, userEntry)

	for _, org := range orgs {
		orgEntry := &entry{
			Name:    *org.GetOrganization().Login,
			Methods: []string{"list"},
		}
		entries = append(entries, orgEntry)
	}
	PrintEntries(entries)
}

// PrintEntries prints the entries
func PrintEntries(entries []*entry) {
	json, _ := json.Marshal(entries)
	fmt.Println(string(pretty.Pretty(json)))
}

// PrintEntry prints one entry
func PrintEntry(entry *entry) {
	json, _ := json.Marshal(entry)
	fmt.Println(string(pretty.Pretty(json)))
}

// directoryToEntries converts github dirContent to entries
func directoryToEntries(dirEntries []*github.RepositoryContent) []*entry {
	entries := make([]*entry, len(dirEntries))
	for i, dirEntry := range dirEntries {
		myEntry := &entry{
			Name: *dirEntry.Name,
			Attributes: attributes{
				Size: *dirEntry.Size,
			},
		}
		if *dirEntry.Type == "file" {
			myEntry.Methods = []string{"read"}
		} else {
			myEntry.Methods = []string{"list"}
		}
		entries[i] = myEntry
	}
	return entries

}

// reposToEntries converts repos to entries
func reposToEntries(repos []*github.Repository) []*entry {
	entries := make([]*entry, len(repos))
	for i, repo := range repos {
		myEntry := &entry{
			Name:    *repo.Name,
			Methods: []string{"list"},
		}
		entries[i] = myEntry
	}
	return entries
}

func init() {
	rootCmd.AddCommand(listCmd)
}
