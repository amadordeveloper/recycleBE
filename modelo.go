package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

var fsjson = getEnv("FS_JSON", "")
var projectID = getEnv("FS_PROJECT_ID", "")

func findAll(collection string) []interface{} {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(fsjson))
	if err != nil {
		return nil
	}
	defer client.Close()

	iter := client.Collection(collection).Documents(ctx)
	var docs []interface{}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil
		}
		docs = append(docs, doc.Data())
	}
	return docs
}

func save(data interface{}, collection string) bool {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(fsjson))
	if err != nil {
		panic(err)
	}

	defer client.Close()

	_, _, err = client.Collection(collection).Add(ctx, data)

	if err != nil {
		panic(err)
	}
	return true
}

func saveAll(data []dataInterface, collection string) bool {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(fsjson))
	if err != nil {
		panic(err)
	}

	defer client.Close()

	batch := client.Batch()

	// Add all documents in a batch if field nombre doesnt exist in collection
	collectionRef := client.Collection(collection)
	operationCount := 0
	for _, d := range data {
		_, err := client.Collection(collection).Doc(d.NombreResiduo()).Get(ctx)
		if err != nil {
			batch.Set(collectionRef.Doc(d.NombreResiduo()), d)
			operationCount++
		}
	}

	if operationCount > 0 {
		_, err = batch.Commit(ctx)
		if err != nil {
			panic(err)
		}
	}

	return true

}

func findBy(collection string, field string, keyword string, isArrayCompare bool) []interface{} {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(fsjson))
	if err != nil {
		panic(err)
	}

	defer client.Close()

	query := client.Collection(collection)

	if isArrayCompare {
		query.Where(field, "array-contains-any", keyword)
		fmt.Println("keyword: ", keyword)
	} else {
		query.Where(field, "==", keyword)
	}

	iter := query.Documents(ctx)
	var results []interface{}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			panic(err)
		}

		results = append(results, doc.Data())
	}

	return results

}
