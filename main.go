package main

import (
	"encoding/hex"
	"encoding/json"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"log"
	"net/http"
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

	r, err := git.PlainOpen("/Users/mhewedy/Work/Code/jirah-api")
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

	http.Handle("/", http.FileServer(http.Dir("./static")))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
