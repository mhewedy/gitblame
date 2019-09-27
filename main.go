package main

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"log"
)

type Author struct {
	Name  string
	Email string
}

type AuthorCommits map[Author][]object.Commit

func main() {

	r, err := git.PlainOpen("/Users/mhewedy/Work/Code/spring-amqp")
	if err != nil {
		log.Fatal(err)
	}

	authors, err := groupCommitsByAuthor(r)
	if err != nil {
		log.Fatal(err)
	}

	for author, commits := range *authors {
		fmt.Println(author.Email, author.Name, len(commits))
	}
}

func groupCommitsByAuthor(r *git.Repository) (*AuthorCommits, error) {
	authorCommits := make(AuthorCommits)

	cIter, err := r.Log(&git.LogOptions{All: true})
	if err != nil {
		return nil, err
	}

	cIter.ForEach(func(c *object.Commit) error {

		author := Author{Name: c.Author.Name, Email: c.Author.Email}
		commits, found := authorCommits[author]
		if !found {
			commits = make([]object.Commit, 0, 10)
		}
		commits = append(commits, *c)
		authorCommits[author] = commits

		return nil
	})

	return &authorCommits, nil
}
