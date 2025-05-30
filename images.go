package main

import (
	"encoding/json"
	"log"

	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/nanovms/ops/types"
)

func listImages(arguments MyFunctionsArguments) (*mcp_golang.ToolResponse, error) {
	c := &types.Config{}

	p, ctx, err := getProviderAndContext(c, "onprem")
	if err != nil {
		log.Println(err)
	}
	ctx.Config().RunConfig.JSON = true

	images, err := p.GetImages(ctx, "")
	if err != nil {
		log.Println(err)
	}

	imgs, err := json.Marshal(images)
	if err != nil {
		log.Println(err)
	}

	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(string(imgs))), nil
}
