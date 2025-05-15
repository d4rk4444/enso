package main

import (
	"enso/src"
	"fmt"
	"log"
	"strings"
)

func main() {
	data, err := src.ParseFile("data.txt")
	if err != nil {
		fmt.Printf("Error parse %v", err)
	}

	proxy, err := src.ParseFile("proxy.txt")
	if err != nil {
		fmt.Printf("Error parse %v", err)
		return
	}

	if len(proxy) < len(data) {
		fmt.Print("Not enough proxy")
		return
	}

	words, err := src.ParseFile("words.txt")
	if err != nil {
		fmt.Printf("Error parse %v", err)
	}

	for i := 0; i < len(data); i++ {
		parts := strings.Split(data[i], ":")
		if len(parts) != 2 {
			fmt.Printf("Error format 'адрес:uuid'. String: %v", i+1)
			return
		}

		address := parts[0]
		uuid := parts[1]
		projectSlug := src.GenerateProjectSlug(words, 1)

		result, err := src.TrackProjectWithProxy(proxy[i], projectSlug, "shortcuts-widget", address, uuid)
		if err != nil {
			fmt.Printf("Error TrackProjectWithProxy %s: %v\n", proxy, err)
			continue
		}

		log.Printf("%v: %v", address, result)
	}
}
