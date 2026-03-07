package main

import (
	"fmt"

	c "github.com/ostafen/clover/v2"
	"github.com/ostafen/clover/v2/query"
)

func main() {
	db, _ := c.Open("dummy")
	doc, _ := db.FindById("parties", "123")
	docs, _ := db.FindAll(query.NewQuery("parties").Where(query.Field("code").Eq("123")))
	exists, _ := db.Exists(query.NewQuery("parties").Where(query.Field("code").Eq("123")))
	fmt.Println(doc, docs, exists)
}
