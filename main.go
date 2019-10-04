package main

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gobuffalo/packr/v2"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	gitHttp "gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"syscall"
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

	http.HandleFunc("/api/stats", func(writer http.ResponseWriter, request *http.Request) {
		stats, err := GetCommitStats(r)
		logIfError(err)

		err = json.NewEncoder(writer).Encode(stats)
	})

	http.HandleFunc("/api", func(writer http.ResponseWriter, request *http.Request) {

		response, err := GroupCommitsByAuthor(r)
		logIfError(err)

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

		patch, err := GetPatch(hashSlice, err, r)

		if err != nil {
			logIfError(err)
			writer.WriteHeader(http.StatusNotFound)
			return
		}

		writer.Write([]byte(patch))
	})

	http.HandleFunc("/api/update", func(writer http.ResponseWriter, request *http.Request) {
		err := Pull(r, auth)

		if err != nil {
			if err == transport.ErrAuthenticationRequired {
				writer.WriteHeader(http.StatusUnauthorized)
			}
			if err == git.NoErrAlreadyUpToDate {
				writer.WriteHeader(http.StatusNotFound)
			}
		}
	})

	http.HandleFunc("/api/settings", func(writer http.ResponseWriter, request *http.Request) {
		settings := struct {
			Path string `json:"path"`
		}{Path: url}

		err = json.NewEncoder(writer).Encode(&settings)
		logIfError(err)
	})

	BuildHttpHandlers(packr.New("myBox", "./templates"))

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
