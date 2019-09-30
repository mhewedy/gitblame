package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

type Commit struct {
	Hash    string    `json:"hash"`
	Message string    `json:"message"`
	When    time.Time `json:"when"`
}

type AuthorWithCommits struct {
	Author  `json:"author"`
	Commits []Commit `json:"commits"`
}

const (
	ErrAuthenticationRequired = "authentication required"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage:", os.Args[0], "<path to local git repository")
		fmt.Println("example:\n", os.Args[0], `"c:\work\myGitProject"`)
		os.Exit(-1)
	}

	projectPath := os.Args[1]
	r, err := git.PlainOpen(projectPath)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/api", func(writer http.ResponseWriter, request *http.Request) {
		response := make([]AuthorWithCommits, 0)

		authors, err := GroupCommitsByAuthor(r)
		logIfError(err)

		for k, v := range *authors {
			commits := make([]Commit, 0)
			for _, c := range v {
				commits = append(commits,
					Commit{Message: c.Message, Hash: hex.EncodeToString(c.Hash[:]), When: c.Author.When})
			}
			response = append(response, AuthorWithCommits{Author: k, Commits: commits})
		}

		// sort entries by commit count
		sort.Slice(response, func(i, j int) bool {
			return len(response[i].Commits) > len(response[j].Commits)
		})
		// sort commits by time
		for _, v := range response {
			sort.Slice(v.Commits, func(i, j int) bool {
				return v.Commits[i].When.After(v.Commits[j].When)
			})
		}

		err = json.NewEncoder(writer).Encode(&response)
		logIfError(err)
	})

	http.HandleFunc("/api/diff/", func(writer http.ResponseWriter, request *http.Request) {

		hashSlice, err := hex.DecodeString(strings.TrimPrefix(request.URL.Path, "/api/diff/"))
		logIfError(err)

		var hashArr [20]byte
		copy(hashArr[:], hashSlice)
		c, err := object.GetCommit(r.Storer, hashArr)
		logIfError(err)

		patch, err := GetCommitPatch(c)
		logIfError(err)

		writer.Write([]byte(patch.String()))
	})

	http.HandleFunc("/api/update", func(writer http.ResponseWriter, request *http.Request) {
		wt, err := r.Worktree()
		logIfError(err)

		err = wt.Pull(&git.PullOptions{})
		logIfError(err)
		if err.Error() == ErrAuthenticationRequired {
			writer.WriteHeader(http.StatusUnauthorized)
		}
	})

	http.HandleFunc("/api/settings", func(writer http.ResponseWriter, request *http.Request) {
		settings := struct {
			Path string `json:"path"`
		}{Path: projectPath}

		err = json.NewEncoder(writer).Encode(&settings)
		logIfError(err)
	})

	http.HandleFunc("/", Index)

	fmt.Println("Server starts at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func logIfError(err error) {
	if err != nil {
		log.Println(err)
	}
}
