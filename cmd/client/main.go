package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/Matt-Kelly-/go-memory-cache/api"
	"google.golang.org/grpc"
	"log"
	"os"
)

const (
	address = "localhost:50051"
)

type commandFunc func(ctx context.Context, client api.CacheClient) error

func main() {
	var (
		ok  bool
		err error
	)

	args := os.Args[1:]

	// Read the command argument
	command, ok := readArgument(args, 0)
	if !ok {
		log.Fatal("No command specified")
	}

	// Parse the command handler
	commandHandler, err := parseCommandHandler(command, args)
	if err != nil {
		log.Fatalf("Failed to parse command: %v", err)
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create client
	client := api.NewCacheClient(conn)

	// Execute command handler
	err = commandHandler(context.Background(), client)
	if err != nil {
		log.Fatalf("Failed to execute command: %v", err)
	}
}

func readArgument(args []string, index int) (string, bool) {
	if index >= len(args) {
		return "", false
	}
	return args[index], true
}

func parseCommandHandler(command string, args []string) (commandFunc, error) {
	switch command {
	case "has":
		return parseHasHandler(args)
	case "get":
		return parseGetHandler(args)
	case "put":
		return parsePutHandler(args)
	case "delete":
		return parseDeleteHandler(args)
	default:
		return nil, fmt.Errorf("Invalid command: %v", command)
	}
}

func parseHasHandler(args []string) (commandFunc, error) {
	key, ok := readArgument(args, 1)
	if !ok {
		return nil, errors.New("No key specified")
	}

	log.Printf("Request: Has key:\"%v\"", key)

	return func(ctx context.Context, client api.CacheClient) error {
		response, err := client.Has(ctx, &api.HasRequest{
			Key: key,
		})
		if err != nil {
			return err
		}
		log.Printf("Response: %v", response)
		return nil
	}, nil
}

func parseGetHandler(args []string) (commandFunc, error) {
	key, ok := readArgument(args, 1)
	if !ok {
		return nil, errors.New("No key specified")
	}

	log.Printf("Command: Get key:\"%v\"", key)

	return func(ctx context.Context, client api.CacheClient) error {
		response, err := client.Get(ctx, &api.GetRequest{
			Key: key,
		})
		if err != nil {
			return err
		}
		log.Printf("Response: %v", response)
		return nil
	}, nil
}

func parsePutHandler(args []string) (commandFunc, error) {
	key, ok := readArgument(args, 1)
	if !ok {
		return nil, errors.New("No key specified")
	}
	value, ok := readArgument(args, 2)
	if !ok {
		return nil, errors.New("No value specified")
	}

	log.Printf("Request: Put key:\"%v\" value:\"%v\"", key, value)

	return func(ctx context.Context, client api.CacheClient) error {
		response, err := client.Put(ctx, &api.PutRequest{
			Key:   key,
			Value: value,
		})
		if err != nil {
			return err
		}
		log.Printf("Response: %v", response)
		return nil
	}, nil
}

func parseDeleteHandler(args []string) (commandFunc, error) {
	key, ok := readArgument(args, 1)
	if !ok {
		return nil, errors.New("No key specified")
	}

	log.Printf("Request: Delete key:\"%v\"", key)

	return func(ctx context.Context, client api.CacheClient) error {
		response, err := client.Delete(ctx, &api.DeleteRequest{
			Key: key,
		})
		if err != nil {
			return err
		}
		log.Printf("Response: %v", response)
		return nil
	}, nil
}
