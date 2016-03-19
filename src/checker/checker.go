package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"io/ioutil"
)

type DB struct {
	UserList []User `json:"userList"`
}

type User struct {
	ID string `json:"id"`
}

type Config struct {
	AccessToken string `json:"AccessToken"`
	UserID      string `json:"UserID"`
	DBFile      string `json:"DBFile"`
}

func main() {
	// Load Config
	config := Config{}
	configData, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(configData, &config)
	if err != nil {
		fmt.Println(err)
	}

	// Get OAuth session
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: config.AccessToken})
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	// Load DB
	ExistingDB := DB{}
	data, err := ioutil.ReadFile(config.DBFile)
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(data, &ExistingDB)
	if err != nil {
		fmt.Println(err)
	}

	// Get followers
	users, _, err := client.Users.ListFollowers(config.UserID, &github.ListOptions{PerPage: 200})
	if err != nil {
		fmt.Println(err)
	}

	NewDB := DB{UserList: []User{}}
	for _, user := range users {
		NewDB.UserList = append(NewDB.UserList, User{ID: *user.Login})
	}

	// Display diff
	NewFollowers := compare(NewDB.UserList, ExistingDB.UserList)
	fmt.Print("NewFollowers: ")
	fmt.Println(NewFollowers)

	UnFollowers := compare(ExistingDB.UserList, NewDB.UserList)
	fmt.Print("UnFollowers: ")
	fmt.Println(UnFollowers)

	// Save DB
	dumpData, err := json.Marshal(NewDB)
	if err != nil {
		fmt.Println(err)
	}
	ioutil.WriteFile(config.DBFile, dumpData, 0644)
}

func compare(X, Y []User) []User {
	m := make(map[User]int)

	for _, y := range Y {
		m[y]++
	}

	var ret []User
	for _, x := range X {
		if m[x] > 0 {
			m[x]--
			continue
		}
		ret = append(ret, x)
	}

	return ret
}
