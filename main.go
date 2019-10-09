package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const url string = "https://api.github.com/users/%s/repos"

// GithubRepo represents a github repository
type GithubRepo struct {
	Name        string `json:"name"`
	BranchesURL string `json:"branches_url"`
}

// GithubBranch represents a branch inside of a repository
type GithubBranch struct {
	Name string `json:"name"`
}

// GithubError represents an error returned from github api
type GithubError struct {
	Message string `json:"message"`
}

// Repository represents the desired data structure for a repository with its branches
type Repository struct {
	Name     string
	Branches []string
	Error    error
}

// Gets all the repositories for a given user and initializes goroutines to fetch the branches for each repository
func main() {
	repositories, err := getRepositories("RicardoLinck")

	if err != nil {
		log.Fatal(err.Error())
	}

	ch := make(chan Repository)
	for _, repository := range repositories {
		go getBranches(repository, ch)
	}

	repos := readRepositoriesFromChannel(ch, len(repositories))
	counts := getBranchCount(repos)

	for _, repo := range repos {
		fmt.Printf("%+v\n", repo)
	}

	fmt.Println()
	fmt.Println("Percentages")
	for key, value := range counts {
		percentage := float64(value) / float64(len(repos)) * 100
		fmt.Printf("Branch [%s] present in %.2f%% of valid repositories\n", key, percentage)
	}
}

func readRepositoriesFromChannel(ch <-chan Repository, length int) []Repository {
	var repositories []Repository
	for i := 0; i < length; i++ {
		repo := <-ch
		if repo.Error == nil {
			repositories = append(repositories, repo)
		} else {
			log.Printf("Error in repository [%s]: %s\n", repo.Name, repo.Error.Error())
		}
	}

	return repositories
}

func getBranchCount(repositories []Repository) map[string]int {
	counts := make(map[string]int)
	for _, repo := range repositories {
		for _, branch := range repo.Branches {
			counts[branch]++
		}
	}
	return counts
}

func getBranches(repository GithubRepo, ch chan<- Repository) {
	repo := Repository{
		Name: repository.Name,
	}
	branches, err := fetchBranchesFromGithub(prepareBranchURL(repository.BranchesURL))

	if err != nil {
		repo.Error = err
	}

	for _, branch := range branches {
		repo.Branches = append(repo.Branches, branch.Name)
	}

	ch <- repo
}

func fetchBranchesFromGithub(url string) ([]GithubBranch, error) {
	var branches []GithubBranch
	httpResponse, err := http.Get(url)

	if err != nil {
		return branches, err
	}

	defer httpResponse.Body.Close()
	err = handleGithubError(httpResponse)
	if err != nil {
		return branches, err
	}

	err = json.NewDecoder(httpResponse.Body).Decode(&branches)
	if err != nil {
		return branches, err
	}

	return branches, nil
}

func prepareBranchURL(url string) string {
	return url[:strings.Index(url, "{")]
}

func getRepositories(user string) ([]GithubRepo, error) {
	var repositories []GithubRepo
	httpResponse, err := http.Get(fmt.Sprintf(url, user))

	if err != nil {
		return repositories, err
	}

	defer httpResponse.Body.Close()

	err = handleGithubError(httpResponse)
	if err != nil {
		return repositories, err
	}

	err = json.NewDecoder(httpResponse.Body).Decode(&repositories)
	if err != nil {
		return repositories, err
	}

	return repositories, nil
}

func handleGithubError(httpResponse *http.Response) error {
	if httpResponse.StatusCode != 200 {
		var githubError GithubError
		err := json.NewDecoder(httpResponse.Body).Decode(&githubError)

		if err != nil {
			return err
		}

		return errors.New(githubError.Message)
	}
	return nil
}
