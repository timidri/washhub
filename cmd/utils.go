package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/github"
	"github.com/spf13/viper"
)

var githubClient *github.Client
var userName string

// GithubClient returns a github client as a singleton
func GithubClient() *github.Client {
	if githubClient == nil {
		if UserName() == "" || password() == "" {
			fmt.Printf("Please specify github_user and github_password in ~/.washhub.yaml")
			os.Exit(1)
		}

		tp := github.BasicAuthTransport{
			Username: UserName(),
			Password: password(),
		}
		githubClient = github.NewClient(tp.Client())
	}
	return githubClient
}

// UserName returns the configured user name
func UserName() string {
	return viper.Get("github_user").(string)
}

func password() string {
	return viper.Get("github_password").(string)
}

// FetchRepositoryContent returns the content at path (a dir listing or file content)
func FetchRepositoryContent(username string, repo string, path string) (*github.RepositoryContent, []*github.RepositoryContent, error) {
	fileContent, dirContent, _, err := GithubClient().Repositories.GetContents(context.Background(), username, repo, path, nil)
	return fileContent, dirContent, err
}

// SplitPath splits a path in a head part and a tail part
func SplitPath(path string) (string, string) {
	parts := strings.SplitN(path, string(os.PathSeparator), 2)
	head := parts[0]
	tail := ""
	if len(parts) > 1 {
		tail = parts[1]
	}
	return head, tail
}

// HandleError provides trivial error handling
func HandleError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err.Error())
		os.Exit(1)
	}
}
