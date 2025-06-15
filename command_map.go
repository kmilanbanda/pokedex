package main

import (
	"fmt"
	"bytes"
	"net/http"
	"encoding/json"
	"io"
)

func commandMap(c *config, args []string) error {
	if c.next == nil {
		fmt.Println("you're on the last page")
		return nil
	}
	
	var body []byte
	var isCached bool
	if body, isCached = Cache.Get(*c.next); !isCached { 
		res, err := http.Get(*c.next)
		if err != nil {
			fmt.Errorf("Error getting response: %v", err)
			return err
		}
		defer res.Body.Close()

		body, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		Cache.Add(*c.next, body)
	}

	var page Page 
	decoder := json.NewDecoder(bytes.NewReader(body))
	if err := decoder.Decode(&page); err != nil {
		return err
	}

	areas := page.Results
	for i := 0; i < len(areas); i++ {
		fmt.Printf("%s\n", areas[i].Name)
	}

	*c = config{
		next: 		page.Next,
		previous: 	page.Previous,
	}	

	return nil
}

func commandMapb(c *config, args []string) error {
	if c.previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}	

	var body []byte
	var isCached bool
	if body, isCached = Cache.Get(*c.previous); !isCached { 
		res, err := http.Get(*c.previous)
		if err != nil {
			fmt.Errorf("Error getting response: %v", err)
			return err
		}
		defer res.Body.Close()

		body, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		Cache.Add(*c.next, body)
	}

	var page Page 
	decoder := json.NewDecoder(bytes.NewReader(body))
	if err := decoder.Decode(&page); err != nil {
		return err
	}
	
	areas := page.Results
	for i := 0; i < len(areas); i++ {
		fmt.Printf("%s\n", areas[i].Name)
	}

	*c = config{
		next: 		page.Next,
		previous: 	page.Previous,
	}

	return nil
}
