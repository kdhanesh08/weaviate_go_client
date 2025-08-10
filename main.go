package main

import (
	"context"
	"fmt"
	"log"

	"github.com/weaviate/weaviate-go-client/v5/weaviate"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/graphql"
	"github.com/weaviate/weaviate/entities/models"
)

func main() {
	client := weaviate.New(weaviate.Config{
		Host:   "localhost:8080",
		Scheme: "http",
	})
	className := "Article"
	schema := &models.Class{
		Class: className,
		Properties: []*models.Property{
			{
				DataType: []string{"text"},
				Name:     "title",
			},
			{
				DataType: []string{"text"},
				Name:     "content",
			},
		},
	}

	ctx := context.Background()
	classes, err := client.Schema().Getter().Do(ctx)
	if err != nil {
		log.Fatalf("Error fetching schema: %v", err)
	}
	exists := false
	for _, c := range classes.Classes {
		if c.Class == className {
			exists = true
			break
		}
	}
	if !exists {
		err = client.Schema().ClassCreator().WithClass(schema).Do(ctx)
		if err != nil {
			log.Fatalf("Error creating schema: %v", err)
		}
		fmt.Println("Class 'Article' created.")
	} else {
		fmt.Println("Class 'Article' already exists.")
	}
	_, err = client.Data().Creator().
		WithClassName(className).
		WithProperties(map[string]interface{}{
			"title":   "Getting Started with Weaviate",
			"content": "This is a sample article stored using GraphQL!",
		}).
		Do(ctx)

	if err != nil {
		log.Fatalf("Error creating object: %v", err)
	}
	fmt.Println("Object inserted successfully.")
	query := client.GraphQL().Get().
		WithClassName(className).
		WithFields(
			graphql.Field{Name: "title"},
			graphql.Field{Name: "content"},
		)

	response, err := query.Do(ctx)
	if err != nil {
		log.Fatalf("Error in GraphQL query: %v", err)
	}

	fmt.Println("GraphQL Response:")
	for _, obj := range response.Data["Get"].(map[string]interface{})[className].([]interface{}) {
		article := obj.(map[string]interface{})
		fmt.Printf("• Title: %s\n", article["title"])
		fmt.Printf("  Content: %s\n", article["content"])
	}

}
