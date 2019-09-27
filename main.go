package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"log"
	"net/http"
	"strings"
)

type Commit struct {
	Hash    string `json:"hash"`
	Message string `json:"message"`
}

func main() {

	r, err := git.PlainOpen("/Users/mhewedy/Work/gobase/code/src/github.com/mhewedy/gitblame")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		response := make(map[string][]Commit)

		authors, err := GroupCommitsByAuthor(r)
		if err != nil {
			log.Fatal(err)
		}

		for k, v := range *authors {

			commits := make([]Commit, 0)
			for _, c := range v {
				commits = append(commits, Commit{Message: c.Message, Hash: hex.EncodeToString(c.Hash[:])})
			}
			response[k.Name+" ("+k.Email+")"] = commits
		}

		err = json.NewEncoder(writer).Encode(&response)
		if err != nil {
			log.Fatal(err)
		}
	})

	http.HandleFunc("/diff/", func(writer http.ResponseWriter, request *http.Request) {
		hash, err := hex.DecodeString(strings.TrimPrefix(request.URL.Path, "/diff/"))
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(hash)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
