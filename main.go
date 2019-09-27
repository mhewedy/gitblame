package main

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"log"
)

func main() {

	r, err := git.PlainOpen("/Users/mhewedy/Work/Code/spring-amqp")
	if err != nil {
		log.Fatal(err)
	}

	authors, err := getAuthors(r)
	if err != nil {
		log.Fatal(err)
	}

	for _, author := range authors {
		fmt.Println(author.Name, author.Email)
	}
}

func getAuthors(r *git.Repository) ([]object.Signature, error) {
	authors := make([]object.Signature, 0)

	cIter, err := r.Log(&git.LogOptions{All: true})
	if err != nil {
		return nil, err
	}

	cIter.ForEach(func(c *object.Commit) error {
		if !contains(authors, c.Author) {
			authors = append(authors, c.Author)
		}
		return nil
	})

	return authors, nil
}

func commitsForAuthor(r *git.Repository, s object.Signature) {

}

func contains(s []object.Signature, e object.Signature) bool {
	for _, elem := range s {
		if elem.Email == e.Email {
			return true
		}
	}
	return false
}
