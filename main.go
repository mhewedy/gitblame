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
)

type Commit struct {
	Hash    string `json:"hash"`
	Message string `json:"message"`
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
				commits = append(commits, Commit{Message: c.Message, Hash: hex.EncodeToString(c.Hash[:])})
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
		hash := strings.TrimPrefix(request.URL.Path, "/api/diff/")
		if err != nil {
			log.Fatal(err)
		}

		cIter, err := r.Log(&git.LogOptions{All: true})
		if err != nil {
			log.Fatal(err)
		}
		defer cIter.Close()

		cIter.ForEach(func(c *object.Commit) error {
			if hex.EncodeToString(c.Hash[:]) == hash {
				patch, err := GetCommitPatch(c)
				if err != nil {
					log.Fatal(err)
				}
				writer.Write([]byte(patch.String()))
			}
			return nil
		})
	})

	http.HandleFunc("/", Index)

	fmt.Println("Server starts at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
