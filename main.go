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

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage:", os.Args[0], "<path to local git repository")
		fmt.Println("example:\n", os.Args[0], `"c:\work\myGitProject"`)
		os.Exit(-1)
	}

	r, err := git.PlainOpen(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/api", func(writer http.ResponseWriter, request *http.Request) {
		response := make([]AuthorWithCommits, 0)

		authors, err := GroupCommitsByAuthor(r)
		if err != nil {
			log.Fatal(err)
		}

		for k, v := range *authors {

			commits := make([]Commit, 0)
			for _, c := range v {
				commits = append(commits,
					Commit{Message: c.Message, Hash: hex.EncodeToString(c.Hash[:]), When: c.Author.When})
			}
			response = append(response, AuthorWithCommits{Author: k, Commits: commits})
		}

		sort.Slice(response, func(i, j int) bool {
			return len(response[i].Commits) > len(response[j].Commits)
		})

		err = json.NewEncoder(writer).Encode(&response)
		if err != nil {
			log.Fatal(err)
		}
	})

	http.HandleFunc("/api/diff/", func(writer http.ResponseWriter, request *http.Request) {

		hashSlice, err := hex.DecodeString(strings.TrimPrefix(request.URL.Path, "/api/diff/"))
		if err != nil {
			log.Fatal(err)
		}

		var hashArr [20]byte
		copy(hashArr[:], hashSlice)
		c, err := object.GetCommit(r.Storer, hashArr)
		if err != nil {
			log.Fatal(err)
		}

		patch, err := GetCommitPatch(c)
		if err != nil {
			log.Fatal(err)
		}
		writer.Write([]byte(patch.String()))
	})

	http.HandleFunc("/", Index)

	fmt.Println("Server starts at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
