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
	Handler().ServeHTTP(w, r)
}
