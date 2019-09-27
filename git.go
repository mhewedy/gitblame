package main

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"io"
	"log"
)

type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type AuthorCommits map[Author][]object.Commit

func test() {
	r, err := git.PlainOpen("/Users/mhewedy/Work/gobase/code/src/github.com/mhewedy/gitblame")
	if err != nil {
		log.Fatal(err)
	}

	authors, err := GroupCommitsByAuthor(r)
	if err != nil {
		log.Fatal(err)
	}

	for author, commits := range *authors {
		fmt.Println(author.Email, author.Name, len(commits))
		fmt.Println("****************************")

		for i, commit := range commits {
			fmt.Println("******i=", i)
			patch, err := GetCommitPatch(&commit)
			if err != nil {
				log.Fatal(err, ">>>>")
			}
			fmt.Println(patch, "\n\n\n\n")
		}

		fmt.Println("****************************")
	}
}

func GetCommitPatch(c *object.Commit) (*object.Patch, error) {

	tree, err := c.Tree()
	if err != nil {
		return nil, err
	}

	parents := c.Parents()
	defer parents.Close()

	parent, err := parents.Next()
	if err != nil && err != io.EOF {
		return nil, err
	}

	var prevTree *object.Tree
	if parent != nil {
		prevTree, err = parent.Tree()
		if err != nil {
			return nil, err
		}
	}

	changes, err := prevTree.Diff(tree)
	if err != nil {
		return nil, err
	}

	patch, err := changes.Patch()
	if err != nil {
		return nil, err
	}

	return patch, nil
}

func GroupCommitsByAuthor(r *git.Repository) (*AuthorCommits, error) {
	authorCommits := make(AuthorCommits)

	cIter, err := r.Log(&git.LogOptions{All: true})
	if err != nil {
		return nil, err
	}
	defer cIter.Close()

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
