package main

import (
	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
)

func main() {
	done := make(chan struct{})

	server := mcp_golang.NewServer(stdio.NewStdioServerTransport())

	// actually this is foregrounded.
	err := server.RegisterTool("pkg_load", "Load package", loadPackage)
	if err != nil {
		panic(err)
	}

	err = server.RegisterTool("instance_logs", "Instance logs", instanceLogs)
	if err != nil {
		panic(err)
	}

	err = server.RegisterTool("instance_create", "Instance create", instanceCreate)
	if err != nil {
		panic(err)
	}

	err = server.RegisterTool("list_instances", "List instances", listInstances)
	if err != nil {
		panic(err)
	}

	err = server.RegisterTool("list_images", "List images", listImages)
	if err != nil {
		panic(err)
	}

	err = server.Serve()
	if err != nil {
		panic(err)
	}

	<-done
}
