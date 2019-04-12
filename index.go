package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	// "github.com/graph-gophers/graphql-go"
	"github.com/graphql-go/graphql"
)

// type query struct{}

// func (*query) Hello() string {
// 	return "Hello, world!"
// }

// Handler is the entry point for the lambda
func Handler(w http.ResponseWriter, r *http.Request) {
	// log.Printf("Processing Lambda request")
	// kraken_api := "https://api.kraken.com/0/public/Assets"

	// s := `
	// 			schema {
	// 				query: Query
	// 			}
	// 			type Query {
	// 					hello: String!
	// 			}
	// `

	// schema := graphql.MustParseSchema(s, &query{})

	// fmt.Fprintf(w, "Hello from Go on Now 2.0!")
	// Schema
	fields := graphql.Fields{
		"hello": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return "world", nil
			},
		},
	}
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	// Query
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	query := string(body)

	params := graphql.Params{Schema: schema, RequestString: query}
	var res = graphql.Do(params)
	if len(res.Errors) > 0 {
		log.Fatalf("failed to execute graphql operation, errors: %+v", res.Errors)
	}
	rJSON, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.Write(rJSON) // {“data”:{“hello”:”world”}}
}
