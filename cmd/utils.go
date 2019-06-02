package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/github"
)

var githubClient *github.Client
var userName string

// GithubClient returns a github client as a singleton
func GithubClient() *github.Client {
	if githubClient == nil {
		if UserName() == "" || password() == "" {
			fmt.Printf("Please specify GITHUB_USER and GITHUB_PASSWORD in ~/.washhubrc")
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
	return os.Getenv("GITHUB_USER")
}

func password() string {
	return os.Getenv("GITHUB_PASSWORD")
}

// FetchRepositoryContent returns the content at path (a dir listing or file content)
func FetchRepositoryContent(username string, repo string, path string) (*github.RepositoryContent, []*github.RepositoryContent, error) {
	fileContent, dirContent, _, err := GithubClient().Repositories.GetContents(context.Background(), username, repo, path, nil)
	return fileContent, dirContent, err
}

// SplitPath splits a path in a head part and a tail part
func SplitPath(path string) (string, string) {
	// fmt.Println("SplitPath called with:", path)
	parts := strings.SplitN(path, string(os.PathSeparator), 2)
	head := parts[0]
	tail := ""
	if len(parts) > 1 {
		tail = parts[1]
	}
	return head, tail
}
