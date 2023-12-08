package cortex_test

import (
	"context"
	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testResourceDefinitionResponse = &cortex.ResourceDefinition{
	Type:        "squid-proxy",
	Name:        "squid-proxy",
	Description: "Cortex's customized squid proxy that is used to make requests to firewalled self-managed resources with a static IP",
	Source:      "CUSTOM",
	Schema: map[string]interface{}{
		"properties": map[string]interface{}{
			"ip": map[string]interface{}{
				"type": "string",
			},
			"resources": map[string]interface{}{
				"type": "string",
			},
			"vpc": map[string]interface{}{
				"type": "string",
			},
		},
		"required": map[string]interface{}{
			"0": "ip",
			"1": "vpc",
		},
	},
}

func TestGetResourceDefinition(t *testing.T) {
	typeName := "squid-proxy"
	c, teardown, err := setupClient(cortex.Route("resource_definitions", typeName), testResourceDefinitionResponse, AssertRequestMethod(t, "GET"))
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	res, err := c.ResourceDefinitions().Get(context.Background(), typeName)
	assert.Nil(t, err, "error retrieving a resource definition")
	assert.Equal(t, testResourceDefinitionResponse.Type, res.Type)
	assert.Equal(t, testResourceDefinitionResponse.Name, res.Name)
	assert.Equal(t, testResourceDefinitionResponse.Description, res.Description)
	assert.Equal(t, testResourceDefinitionResponse.Source, res.Source)
	assert.Equal(t, testResourceDefinitionResponse.Schema, res.Schema)
}

var testListResourceDefinitionsResponse = &cortex.ResourceDefinitionsResponse{
	ResourceDefinitions: []cortex.ResourceDefinition{
		*testResourceDefinitionResponse,
		*testResourceDefinitionResponse,
	},
}

func TestListResourceDefinitions(t *testing.T) {
	c, teardown, err := setupClient(cortex.Route("resource_definitions", ""), testListResourceDefinitionsResponse, AssertRequestMethod(t, "GET"))
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	req := &cortex.ResourceDefinitionListParams{}
	res, err := c.ResourceDefinitions().List(context.Background(), req)
	assert.Nil(t, err, "error retrieving resource definitions")

	assert.Len(t, res.ResourceDefinitions, 2)
}

func TestCreateResourceDefinition(t *testing.T) {
	req := cortex.CreateResourceDefinitionRequest{
		Type:        testResourceDefinitionResponse.Type,
		Name:        testResourceDefinitionResponse.Name,
		Description: testResourceDefinitionResponse.Description,
		Source:      testResourceDefinitionResponse.Source,
		Schema:      testResourceDefinitionResponse.Schema,
	}
	c, teardown, err := setupClient(
		cortex.Route("resource_definitions", ""),
		testResourceDefinitionResponse,
		AssertRequestMethod(t, "POST"),
		AssertRequestBody(t, req),
	)
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	res, err := c.ResourceDefinitions().Create(context.Background(), req)
	assert.Nil(t, err, "error creating a resource definition")
	assert.Equal(t, res.Type, req.Type)
	assert.Equal(t, res.Name, req.Name)
	assert.Equal(t, res.Description, req.Description)
	assert.Equal(t, res.Source, req.Source)
	assert.Equal(t, res.Schema, req.Schema)
}

func TestUpdateResourceDefinition(t *testing.T) {
	req := cortex.UpdateResourceDefinitionRequest{
		Name:        testResourceDefinitionResponse.Name + "-updated",
		Description: testResourceDefinitionResponse.Description + "-updated",
		Schema:      testResourceDefinitionResponse.Schema,
	}
	typeName := testResourceDefinitionResponse.Type

	resp := &cortex.ResourceDefinition{
		Type:        typeName,
		Name:        testResourceDefinitionResponse.Name + "-updated",
		Description: testResourceDefinitionResponse.Description + "-updated",
		Schema:      testResourceDefinitionResponse.Schema,
	}

	c, teardown, err := setupClient(
		cortex.Route("resource_definitions", typeName),
		resp,
		AssertRequestMethod(t, "PUT"),
		AssertRequestBody(t, req),
	)
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	res, err := c.ResourceDefinitions().Update(context.Background(), typeName, req)

	assert.Nil(t, err, "error updating a resource definition")
	assert.Equal(t, res.Type, typeName)
	assert.Equal(t, res.Name, req.Name)
	assert.Equal(t, res.Description, req.Description)
	assert.Equal(t, res.Schema, req.Schema)
}

func TestDeleteResourceDefinition(t *testing.T) {
	typeName := testResourceDefinitionResponse.Type
	c, teardown, err := setupClient(
		cortex.Route("resource_definitions", typeName),
		cortex.DeleteResourceDefinitionResponse{},
		AssertRequestMethod(t, "DELETE"),
	)
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	err = c.ResourceDefinitions().Delete(context.Background(), typeName)
	assert.Nil(t, err, "error deleting a resource definition")
}
