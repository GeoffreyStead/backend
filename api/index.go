package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil" // Import ioutil for file reading
	"net/http"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hello from Go!</h1>")
}

// CSV file path
const csvFilePath = "csvtest.csv"

// Define the schema
var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	},
)

// Define a resolver function for reading the CSV file
func resolveReadCSV(p graphql.ResolveParams) (interface{}, error) {
	// Read the content of the CSV file
	csvData, err := ioutil.ReadFile(csvFilePath)
	if err != nil {
		return nil, err
	}

	// Parse CSV content
	reader := csv.NewReader(bytes.NewReader(csvData))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Convert CSV data to structured text response
	var text bytes.Buffer
	for i, record := range records {
		if i > 0 {
			text.WriteString("\n")
		}
		for j, column := range record {
			if j > 0 {
				text.WriteString("$")
			}
			if column == "" {
				column = " " // Replace empty string with space
			}
			text.WriteString(column)
		}
	}

	return text.String(), nil
}

// Define a resolver function for uploading the CSV file
func resolveUploadCSV(p graphql.ResolveParams) (interface{}, error) {
	// Extract file content from the resolver arguments
	fileContent, ok := p.Args["fileContent"].(string)
	if !ok {
		return nil, fmt.Errorf("file content must be a string")
	}

	// Parse the CSV content
	csvData, err := parseCSVContent(fileContent)
	if err != nil {
		return nil, err
	}

	// Convert CSV data to structured text response
	var text bytes.Buffer
	for i, record := range csvData {
		if i > 0 {
			text.WriteString("\n")
		}
		for j, column := range record {
			if j > 0 {
				text.WriteString("$")
			}
			if column == "" {
				column = " " // Replace empty string with space
			}
			text.WriteString(column)
		}
	}

	return text.String(), nil
}

// Function to parse CSV content into [][]string slice
func parseCSVContent(content string) ([][]string, error) {
	reader := csv.NewReader(strings.NewReader(content))
	var csvData [][]string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		csvData = append(csvData, record)
	}
	return csvData, nil
}

// Define the query type
var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"read": &graphql.Field{
				Type:    graphql.String,
				Resolve: resolveReadCSV,
			},
		},
	},
)

// Define the mutation type
var mutationType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"uploadCSV": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					"fileContent": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: resolveUploadCSV,
			},
		},
	},
)

func main() {
	// Create a new GraphQL handler
	graphQLHandler := handler.New(&handler.Config{
		Schema: &schema,
		Pretty: true,
	})

	// Define a handler function for CORS
	corsHandler := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set CORS headers
			w.Header().Set("Access-Control-Allow-Origin", "*") // Allow requests from your frontend URL
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")

			// If it's a preflight request, respond with 200 OK
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			// Call the actual handler
			h.ServeHTTP(w, r)
		})
	}

	// Serve GraphQL requests with CORS support
	http.Handle("/graphql", corsHandler(graphQLHandler))

	// Start the HTTP server
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
