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

	// res.Results contains an array of Search results
	for _, res := range res.Results {
		fmt.Println("Video Title:", res.Title, "Channel:", res.Channel)
	}
}
