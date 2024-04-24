package internal

import "github.com/graphql-go/graphql"

// DefineSchema function defines the GraphQL schema
func DefineSchema() graphql.Schema {
    var queryType = graphql.NewObject(
        graphql.ObjectConfig{
            Name: "Query",
            Fields: graphql.Fields{
                "hello": &graphql.Field{
                    Type: graphql.String,
                    Resolve: func(p graphql.ResolveParams) (interface{}, error) {
                        return "world", nil
                    },
                },
            },
        },
    )

    schema, _ := graphql.NewSchema(
        graphql.SchemaConfig{
            Query: queryType,
        },
    )

    return schema
}
