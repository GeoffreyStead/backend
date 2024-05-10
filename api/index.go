package handler

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

// Define the schema
var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	},
)

// URL to the CSV file
const csvFileURL = "https://backend-lac-seven.vercel.app/csvtest.csv"

// Define a resolver function for reading the CSV file
func resolveReadCSV(p graphql.ResolveParams) (interface{}, error) {
	// Fetch CSV content from the URL
	resp, err := http.Get(csvFileURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the content of the response
	csvData, err := ioutil.ReadAll(resp.Body)
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
	file, ok := p.Args["file"].(string)
	if !ok {
		return nil, fmt.Errorf("file content must be a string")
	}

	// Parse the CSV content
	csvData, err := parseCSVContent(file)
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
					"file": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: resolveUploadCSV,
			},
		},
	},
)

// Exported function to handle GraphQL requests
func Handler() http.Handler {
	// Create a new GraphQL handler
	return handler.New(&handler.Config{
		Schema: &schema,
		Pretty: true,
	})
}

// Exported entry point for the serverless function
func Index(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers to allow cross-origin requests
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight OPTIONS request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Serve GraphQL requests
	Handler().ServeHTTP(w, r)
}
