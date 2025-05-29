package main

import (
	"encoding/json"
	"fmt"

	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"

	api "github.com/nanovms/ops/lepton"
	"github.com/nanovms/ops/provider"
	"github.com/nanovms/ops/types"
)

type Content struct {
	Title       string  `json:"title" jsonschema:"required,description=The title to submit"`
	Description *string `json:"description" jsonschema:"description=The description to submit"`
}

type MyFunctionsArguments struct {
	Submitter string  `json:"submitter" jsonschema:"required,description=The name of the thing calling this tool (openai, google, claude, etc)"`
	Content   Content `json:"content" jsonschema:"required,description=The content of the message"`
}

func getProviderAndContext(c *types.Config, providerName string) (api.Provider, *api.Context, error) {
	p, err := provider.CloudProvider(providerName, &c.CloudConfig)
	if err != nil {
		return nil, nil, err
	}

	ctx := api.NewContext(c)

	return p, ctx, nil
}

func main() {
	done := make(chan struct{})

	server := mcp_golang.NewServer(stdio.NewStdioServerTransport())
	err := server.RegisterTool("hello", "Say hello to a person", func(arguments MyFunctionsArguments) (*mcp_golang.ToolResponse, error) {
		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(fmt.Sprintf("Hello, %server!", arguments.Submitter))), nil
	})
	if err != nil {
		panic(err)
	}

	err = server.RegisterTool("list_images", "List images", func(arguments MyFunctionsArguments) (*mcp_golang.ToolResponse, error) {

		c := &types.Config{}

		p, ctx, err := getProviderAndContext(c, "onprem")
		if err != nil {
			//			fmt.Println(err)
		}
		ctx.Config().RunConfig.JSON = true

		images, err := p.GetImages(ctx, "")
		if err != nil {
			//			fmt.Println(err)
		}

		imgs, err := json.Marshal(images)
		if err != nil {
			//			fmt.Println(err)
		}

		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(string(imgs))), nil
	})
	if err != nil {
		panic(err)
	}

	err = server.RegisterPrompt("promt_test", "This is a test prompt", func(arguments Content) (*mcp_golang.PromptResponse, error) {
		return mcp_golang.NewPromptResponse("description", mcp_golang.NewPromptMessage(mcp_golang.NewTextContent(fmt.Sprintf("Hello, %server!", arguments.Title)), mcp_golang.RoleUser)), nil
	})
	if err != nil {
		panic(err)
	}

	err = server.RegisterResource("test://resource", "resource_test", "This is a test resource", "application/json", func() (*mcp_golang.ResourceResponse, error) {
		return mcp_golang.NewResourceResponse(mcp_golang.NewTextEmbeddedResource("test://resource", "This is a test resource", "application/json")), nil
	})

	err = server.Serve()
	if err != nil {
		panic(err)
	}

	<-done
}
