package main

import (
	"encoding/json"
	"log"

	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/nanovms/ops/lepton"
	"github.com/nanovms/ops/types"
)

type InstanceArguments struct {
	ImageName string `json:"longitude" jsonschema:"required,description=The image name of the instance to create"`
}

func instanceLogs(arguments InstanceArguments) (*mcp_golang.ToolResponse, error) {

	c := &types.Config{}

	log.Println(arguments)

	p, ctx, err := getProviderAndContext(c, "onprem")
	if err != nil {
		log.Println(err)
	}
	ctx.Config().RunConfig.JSON = true

	body, err := p.GetInstanceLogs(ctx, arguments.ImageName)
	if err != nil {
		log.Println(err.Error())
	}

	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(body)), nil
}

func instanceCreate(arguments InstanceArguments) (*mcp_golang.ToolResponse, error) {
	c := lepton.NewConfig()
	log.Println(arguments.ImageName)

	log.Println(arguments)

	c.CloudConfig.ImageName = arguments.ImageName

	p, ctx, err := getProviderAndContext(c, "onprem")
	if err != nil {
		log.Println(err)
	}

	c.RunConfig.Kernel = lepton.GetOpsHome() + "/0.1.53-arm/kernel.img"

	err = p.CreateInstance(ctx)
	if err != nil {
		log.Println(err)
	}

	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("spun that instance up")), nil
}

func listInstances(arguments MyFunctionsArguments) (*mcp_golang.ToolResponse, error) {
	c := &types.Config{}

	p, ctx, err := getProviderAndContext(c, "onprem")
	if err != nil {
		log.Println(err)
	}
	ctx.Config().RunConfig.JSON = true

	instances, err := p.GetInstances(ctx)
	if err != nil {
		log.Println(err)
	}

	imgs, err := json.Marshal(instances)
	if err != nil {
		log.Println(err)
	}

	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(string(imgs))), nil
}
