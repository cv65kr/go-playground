package main

import (
	"encoding/json"
	"fmt"

	"github.com/go-pg/pg"
)

type Event struct {
	Data      TestTable `json:"data"`
	Operation string    `json:"operation"`
}

type TestTable struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	db := pg.Connect(&pg.Options{
		Addr:     "localhost:5432",
		User:     "postgres",
		Password: "postgres",
		Database: "test",
	})

	ln := db.Listen("test_table_update")
	defer ln.Close()

	for event := range ln.Channel() {
		go processMessage(event.Payload)
	}
}

func processMessage(event string) {
	fmt.Printf("Plain message: %v\n", event)

	var dataEvent Event
	if err := json.Unmarshal([]byte(event), &dataEvent); err != nil {
		panic(err)
	}

	fmt.Printf("ID: %v\n", dataEvent.Data.Id)
}
