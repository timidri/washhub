package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"golang.org/x/oauth2"

	"github.com/google/go-github/github"
	"github.com/spf13/viper"
)

var githubClient *github.Client

// GithubClient returns a github client as a singleton
func GithubClient() *github.Client {
	if githubClient == nil {
		if token() != "" {
			ts := oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: token()},
			)
			tc := oauth2.NewClient(context.Background(), ts)
			githubClient = github.NewClient(tc)
		} else if UserName() == "" || password() == "" {
			fmt.Fprintf(os.Stderr, "Please specify github_token or github_user and github_password in ~/.washhub.yaml")
			os.Exit(1)
		} else {
			tp := github.BasicAuthTransport{
				Username: UserName(),
				Password: password(),
			}
			githubClient = github.NewClient(tp.Client())
		}
	}

	return githubClient
}

// UserName returns the configured user name
func UserName() string {
	user := viper.Get("github_user")
	if user == nil {
		return ""
	}
	return user.(string)
}

func password() string {
	password := viper.Get("github_password")
	if password == nil {
		return ""
	}
	return password.(string)
}

func token() string {
	token := viper.Get("github_token")
	if token == nil {
		return ""
	}
	return token.(string)
}

// FetchRepositoryContent returns the content at path (a dir listing or file content)
func FetchRepositoryContent(username string, repo string, path string) (*github.RepositoryContent, []*github.RepositoryContent, error) {
	fileContent, dirContent, _, err := GithubClient().Repositories.GetContents(context.Background(), username, repo, path, nil)
	return fileContent, dirContent, err
}

// SplitPath splits a path in a head part and a tail part
// examples:
//   SplitPath("a") -> ("a","")
//   SplitPath("a/b") -> ("a","b")
//   SplitPath("a/b/c") => ("a", "b/c")
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
