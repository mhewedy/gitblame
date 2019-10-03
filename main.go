package main

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	gitHttp "gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"syscall"
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
	ErrAlreadyUpToDate        = "already up-to-date"
)

func main() {

	url, username, password := readParams()
	auth := &gitHttp.BasicAuth{
		Username: username,
		Password: password,
	}

	r, err := git.Clone(memory.NewStorage(), memfs.New(), &git.CloneOptions{
		Auth:     auth,
		URL:      url,
		Progress: os.Stdout,
	})
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

		patch := GetPatch(hashSlice, err, r)

		if err != nil {
			logIfError(err)
			writer.WriteHeader(http.StatusNotFound)
			return
		}

		writer.Write([]byte(patch.String()))
	})

	http.HandleFunc("/api/update", func(writer http.ResponseWriter, request *http.Request) {
		err := Pull(r, auth)

		if err != nil && err.Error() == ErrAuthenticationRequired {
			writer.WriteHeader(http.StatusUnauthorized)
		}
		if err != nil && err.Error() == ErrAlreadyUpToDate {
			writer.WriteHeader(http.StatusNotFound)
		}
	})

	http.HandleFunc("/api/settings", func(writer http.ResponseWriter, request *http.Request) {
		settings := struct {
			Path string `json:"path"`
		}{Path: url}

		err = json.NewEncoder(writer).Encode(&settings)
		logIfError(err)
	})

	http.HandleFunc("/", Index)

	fmt.Println("Server starts at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func readParams() (string, string, string) {
	//https://stackoverflow.com/a/32768479/171950
	credentials := func() (string, string) {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Enter Username: ")
		username, _ := reader.ReadString('\n')

		fmt.Print("Enter Password: ")
		bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
		password := string(bytePassword)

		return strings.TrimSpace(username), strings.TrimSpace(password)
	}

	if len(os.Args) < 2 {
		fmt.Println("Usage:", os.Args[0], "<repo url>")
		os.Exit(-1)
	}
	username, password := credentials()
	return os.Args[1], username, password
}

func logIfError(err error) {
	if err != nil {
		log.Println(err)
	}
}
