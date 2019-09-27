package main

import (
	"encoding/hex"
	"encoding/json"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"log"
	"net/http"
	"strings"
)

type Commit struct {
	Hash    string `json:"hash"`
	Message string `json:"message"`
}

func main() {

	r, err := git.PlainOpen("/Users/mhewedy/Work/Code/spring-amqp")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/api", func(writer http.ResponseWriter, request *http.Request) {
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

	log.Fatal(http.ListenAndServe(":8080", nil))
}
