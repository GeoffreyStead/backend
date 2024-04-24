package main

import (
    "net/http"

    "github.com/graphql-go/handler"
    "back/internal"
)

func main() {
    schema := internal.DefineSchema() // Importing the schema from internal package

    h := handler.New(&handler.Config{
        Schema: &schema,
        Pretty: true,
    })

    http.Handle("/graphql", h)

    http.ListenAndServe(":8080", nil)
}
