package main

import (
	"fmt"

	"github.com/ppalone/ytsearch"
)

func main() {
	// client
	c := ytsearch.Client{}

	// search for "nocopyrightsounds"
	res, err := c.Search("nocopyrightsounds")
	if err != nil {
		panic(err)
	}

	fmt.Println("Search Results:")

	// res.Results contains an array of Search results
	for _, res := range res.Results {
		fmt.Println("Video Title:", res.Title, "Channel:", res.Channel)
	}

	fmt.Println("Next Search Results:")

	// res.Continuation contains continuation token
	// can be used to fetch next results
	// since youtube limits the number of search results
	res, err = c.Next(res.Continuation)
	if err != nil {
		panic(err)
	}

	for _, res := range res.Results {
		fmt.Println("Video Title:", res.Title, "Channel:", res.Channel)
	}
}
