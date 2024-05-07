package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

// Test para saber si se lee bien el csv
func TestResolveReadCSV(t *testing.T) {
	expectedOutput := `
	Assignee$Worklog Description$User$Start Time$Time Spent (s)
Pol Alsina Domènech$GOVR$Pol Alsina Domènech$22/09/2023 08:50 AM$14400`

	params := graphql.ResolveParams{}
	result, err := resolveReadCSV(params)
	resultStr := fmt.Sprintf("%v", result)
	lines := strings.Split(resultStr, "\n")
	firstTwoLines := strings.Join(lines[:2], "\n")

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if strings.TrimSpace(firstTwoLines) != strings.TrimSpace(expectedOutput) {
		t.Errorf("Expected %q but got %q", expectedOutput, result)
	}
}

// Test para saber si hace bien el parse de la información
func TestParseCSVContent(t *testing.T) {
	validCSVContent := `name,age,city
John,30,New York
Alice,25,Los Angeles
Bob,35,Chicago`

	expectedValidCSVData := [][]string{
		{"name", "age", "city"},
		{"John", "30", "New York"},
		{"Alice", "25", "Los Angeles"},
		{"Bob", "35", "Chicago"},
	}

	validCSVData, err := parseCSVContent(validCSVContent)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	t.Logf("Valid CSV Data:\n%v", validCSVData)

	if !reflect.DeepEqual(validCSVData, expectedValidCSVData) {
		t.Errorf("unexpected CSV data. Got: %v, Expected: %v", validCSVData, expectedValidCSVData)
	}

	emptyCSVContent := ``

	expectedEmptyCSVData := [][]string{}

	emptyCSVData, err := parseCSVContent(emptyCSVContent)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	t.Logf("Empty CSV Data:\n%v", emptyCSVData)
	t.Logf("Empty CSV Data:\n%v", expectedEmptyCSVData)

	if len(emptyCSVData) != len(expectedEmptyCSVData) {
		t.Errorf("unexpected CSV data length. Got: %d, Expected: %d", len(emptyCSVData), len(expectedEmptyCSVData))
	}

	for i := range emptyCSVData {
		if len(emptyCSVData[i]) != len(expectedEmptyCSVData[i]) {
			t.Errorf("unexpected CSV data length in row %d. Got: %d, Expected: %d", i, len(emptyCSVData[i]), len(expectedEmptyCSVData[i]))
		}
	}
}

func TestGraphQLWithCORS(t *testing.T) {

	req, err := http.NewRequest("OPTIONS", "/graphql", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	graphQLHandler := handler.New(&handler.Config{
		Schema: &schema,
		Pretty: true,
	})

	corsHandler := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
	http.Handle("/graphql", corsHandler(graphQLHandler))

	http.DefaultServeMux.ServeHTTP(rr, req)

	expectedHeaders := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "POST, GET, OPTIONS, PUT, DELETE",
		"Access-Control-Allow-Headers": "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With",
	}

	for key, value := range expectedHeaders {
		if rr.Header().Get(key) != value {
			t.Errorf("handler returned wrong CORS header %s: got %v want %v",
				key, rr.Header().Get(key), value)
		}
	}

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
