package tests

import (
	"testing"

	"github.com/taufik-rama/maiden/internal/config"
)

func TestServices(t *testing.T) {

	HTTP(t)

	GRPC(t)
}

func HTTP(t *testing.T) {

	{
		c, err := config.New("definitions/services/http/imports.yaml")
		if err != nil {
			t.Fatal(err)
		}
		if err := c.ResolveServices(); err != nil {
			t.Fatal(err)
		}
		if c.Services == nil || c.Fixtures != nil {
			t.Error()
		}

		if c.Services.HTTP == nil {
			t.Error()
		}
		if c.Services.HTTP["http-service-1"].Imports != nil || c.Services.HTTP["http-service-2"].Imports != nil {
			t.Error()
		}
		if c.Services.HTTP["http-service-1"].Port != 100 || c.Services.HTTP["http-service-2"].Port != 100 {
			t.Error()
		}
		if c.Services.HTTP["http-service-1"].Endpoints == nil || c.Services.HTTP["http-service-2"].Endpoints == nil {
			t.Error()
		}
		{
			api := c.Services.HTTP["http-service-1"].Endpoints["/"][0]
			if api.Method != "GET" || api.Request != "request" || api.Response != "response" {
				t.Error()
			}
		}
		{
			api := c.Services.HTTP["http-service-1"].Endpoints["/"][1]
			if api.Method != "POST" || api.Request == nil || api.Response == nil {
				t.Error()
			}
		}
		{
			api := c.Services.HTTP["http-service-1"].Endpoints["/endpoint-1"][0]
			methods := api.Method.([]interface{})
			if methods[0] != "GET" || methods[1] != "POST" || api.Request != "request" || api.Response != "response" {
				t.Error()
			}
		}
		{
			api := c.Services.HTTP["http-service-2"].Endpoints["/"][0]
			if api.Method != "GET" || api.Request != "request" || api.Response != "response" {
				t.Error()
			}
		}
		{
			api := c.Services.HTTP["http-service-2"].Endpoints["/"][1]
			if api.Method != "POST" || api.Request == nil || api.Response == nil {
				t.Error()
			}
		}
		{
			api := c.Services.HTTP["http-service-2"].Endpoints["/endpoint-1"][0]
			methods := api.Method.([]interface{})
			if methods[0] != "GET" || methods[1] != "POST" || api.Request != "request" || api.Response != "response" {
				t.Error()
			}
		}
	}

	{
		c, err := config.New("definitions/services/http/imports-root.yaml")
		if err != nil {
			t.Fatal(err)
		}
		if err := c.ResolveServices(); err != nil {
			t.Fatal(err)
		}
		if c.Services == nil || c.Fixtures != nil {
			t.Error()
		}

		if c.Services.HTTP == nil {
			t.Error()
		}
		if c.Services.HTTP["http-service-1"].Imports != nil || c.Services.HTTP["http-service-2"].Imports != nil {
			t.Error()
		}
		if c.Services.HTTP["http-service-1"].Port != 90 || c.Services.HTTP["http-service-2"].Port != 90 {
			t.Error()
		}
		if c.Services.HTTP["http-service-1"].Endpoints == nil || c.Services.HTTP["http-service-2"].Endpoints == nil {
			t.Error()
		}
		{
			api := c.Services.HTTP["http-service-1"].Endpoints["/"][0]
			if api.Method != "GET" || api.Request != "root-request" || api.Response != "root-response" {
				t.Error()
			}
		}
		{
			api := c.Services.HTTP["http-service-1"].Endpoints["/"][1]
			if api.Method != "POST" || api.Request == nil || api.Response == nil {
				t.Error()
			}
		}
		{
			api := c.Services.HTTP["http-service-1"].Endpoints["/endpoint-1"][0]
			methods := api.Method.([]interface{})
			if methods[0] != "GET" || methods[1] != "POST" || api.Request != "root-request" || api.Response != "root-response" {
				t.Error()
			}
		}
		{
			api := c.Services.HTTP["http-service-2"].Endpoints["/"][0]
			if api.Method != "GET" || api.Request != "root-request" || api.Response != "root-response" {
				t.Error()
			}
		}
		{
			api := c.Services.HTTP["http-service-2"].Endpoints["/"][1]
			if api.Method != "POST" || api.Request == nil || api.Response == nil {
				t.Error()
			}
		}
		{
			api := c.Services.HTTP["http-service-2"].Endpoints["/endpoint-1"][0]
			methods := api.Method.([]interface{})
			if methods[0] != "GET" || methods[1] != "POST" || api.Request != "root-request" || api.Response != "root-response" {
				t.Error()
			}
		}
	}

	{
		c, err := config.New("definitions/services/http/services.yaml")
		if err != nil {
			t.Fatal(err)
		}
		if err := c.ResolveServices(); err != nil {
			t.Fatal(err)
		}
		if c.Services == nil || c.Fixtures != nil {
			t.Error()
		}

		if c.Services.HTTP == nil {
			t.Error()
		}
		if c.Services.HTTP["http-service-1"].Imports != nil || c.Services.HTTP["http-service-2"].Imports != nil {
			t.Error()
		}
		if c.Services.HTTP["http-service-1"].Port != 100 || c.Services.HTTP["http-service-2"].Port != 100 {
			t.Error()
		}
		if c.Services.HTTP["http-service-1"].Endpoints == nil || c.Services.HTTP["http-service-2"].Endpoints == nil {
			t.Error()
		}
		{
			api := c.Services.HTTP["http-service-1"].Endpoints["/"][0]
			if api.Method != "GET" || api.Request != "request" || api.Response != "response" {
				t.Error()
			}
		}
		{
			api := c.Services.HTTP["http-service-1"].Endpoints["/"][1]
			if api.Method != "POST" || api.Request == nil || api.Response == nil {
				t.Error()
			}
		}
		{
			api := c.Services.HTTP["http-service-1"].Endpoints["/endpoint-1"][0]
			methods := api.Method.([]interface{})
			if methods[0] != "GET" || methods[1] != "POST" || api.Request != "request" || api.Response != "response" {
				t.Error()
			}
		}
		{
			api := c.Services.HTTP["http-service-2"].Endpoints["/"][0]
			if api.Method != "GET" || api.Request != "request" || api.Response != "response" {
				t.Error()
			}
		}
		{
			api := c.Services.HTTP["http-service-2"].Endpoints["/"][1]
			if api.Method != "POST" || api.Request == nil || api.Response == nil {
				t.Error()
			}
		}
		{
			api := c.Services.HTTP["http-service-2"].Endpoints["/endpoint-1"][0]
			methods := api.Method.([]interface{})
			if methods[0] != "GET" || methods[1] != "POST" || api.Request != "request" || api.Response != "response" {
				t.Error()
			}
		}
	}
}

func GRPC(t *testing.T) {

	{
		c, err := config.New("definitions/services/grpc/imports.yaml")
		if err != nil {
			t.Fatal(err)
		}
		if err := c.ResolveServices(); err != nil {
			t.Fatal(err)
		}
		if c.Services == nil || c.Fixtures != nil {
			t.Error()
		}

		if c.Services.GRPC == nil {
			t.Error()
		}
		if c.Services.GRPC["grpc-service-1"].Imports != nil || c.Services.GRPC["grpc-service-2"].Imports != nil {
			t.Error()
		}
		if c.Services.GRPC["grpc-service-1"].Port != 100 || c.Services.GRPC["grpc-service-2"].Port != 100 {
			t.Error()
		}
		if c.Services.GRPC["grpc-service-1"].Methods == nil || c.Services.GRPC["grpc-service-2"].Methods == nil {
			t.Error()
		}
		if c.Services.GRPC["grpc-service-1"].Conditions == nil || c.Services.GRPC["grpc-service-2"].Conditions == nil {
			t.Error()
		}
		{
			service := c.Services.GRPC["grpc-service-1"]
			methods := service.Methods
			if service.Port != 100 || service.Definition != "definition-1" {
				t.Error()
			}
			{
				if methods["method-1"].Request != "method-1-request" || methods["method-1"].Response != "method-1-response" {
					t.Error()
				}
			}
			{
				conditions := service.Conditions
				if conditions["method-1"][0].Request == nil || conditions["method-1"][0].Response == nil {
					t.Error()
				}
				if conditions["method-2"][0].Request == nil || conditions["method-2"][0].Response == nil {
					t.Error()
				}
			}
		}
		{
			service := c.Services.GRPC["grpc-service-2"]
			methods := service.Methods
			if service.Port != 100 || service.Definition != "definition-2" {
				t.Error()
			}
			{
				if methods["method-1"].Request != "method-1-request" || methods["method-1"].Response != "method-1-response" {
					t.Error()
				}
			}
			{
				conditions := service.Conditions
				if conditions["method-1"][0].Request == nil || conditions["method-1"][0].Response == nil {
					t.Error()
				}
				if conditions["method-2"][0].Request == nil || conditions["method-2"][0].Response == nil {
					t.Error()
				}
			}
		}
	}

	{
		c, err := config.New("definitions/services/grpc/imports-root.yaml")
		if err != nil {
			t.Fatal(err)
		}
		if err := c.ResolveServices(); err != nil {
			t.Fatal(err)
		}
		if c.Services == nil || c.Fixtures != nil {
			t.Error()
		}

		if c.Services.GRPC == nil {
			t.Error()
		}
		if c.Services.GRPC["grpc-service-1"].Imports != nil || c.Services.GRPC["grpc-service-2"].Imports != nil {
			t.Error()
		}
		if c.Services.GRPC["grpc-service-1"].Port != 90 || c.Services.GRPC["grpc-service-2"].Port != 90 {
			t.Error()
		}
		if c.Services.GRPC["grpc-service-1"].Methods == nil || c.Services.GRPC["grpc-service-2"].Methods == nil {
			t.Error()
		}
		if c.Services.GRPC["grpc-service-1"].Conditions == nil || c.Services.GRPC["grpc-service-2"].Conditions == nil {
			t.Error()
		}
		{
			service := c.Services.GRPC["grpc-service-1"]
			methods := service.Methods
			if service.Port != 90 || service.Definition != "root-definition-1" {
				t.Error()
			}
			{
				if methods["root-method-1"].Request != "root-method-1-request" || methods["root-method-1"].Response != "root-method-1-response" {
					t.Error()
				}
			}
			{
				conditions := service.Conditions
				if conditions["root-method-1"][0].Request == nil || conditions["root-method-1"][0].Response == nil {
					t.Error()
				}
				if conditions["root-method-2"][0].Request == nil || conditions["root-method-2"][0].Response == nil {
					t.Error()
				}
			}
		}
		{
			service := c.Services.GRPC["grpc-service-2"]
			methods := service.Methods
			if service.Port != 90 || service.Definition != "root-definition-2" {
				t.Error()
			}
			{
				if methods["method-1"].Request != "method-1-request" || methods["method-1"].Response != "method-1-response" {
					t.Error()
				}
			}
			{
				conditions := service.Conditions
				if conditions["method-1"][0].Request == nil || conditions["method-1"][0].Response == nil {
					t.Error()
				}
				if conditions["method-2"][0].Request == nil || conditions["method-2"][0].Response == nil {
					t.Error()
				}
			}
		}
	}

	{
		c, err := config.New("definitions/services/grpc/services.yaml")
		if err != nil {
			t.Fatal(err)
		}
		if err := c.ResolveServices(); err != nil {
			t.Fatal(err)
		}
		if c.Services == nil || c.Fixtures != nil {
			t.Error()
		}

		if c.Services.GRPC == nil {
			t.Error()
		}
		if c.Services.GRPC["grpc-service-1"].Imports != nil || c.Services.GRPC["grpc-service-2"].Imports != nil {
			t.Error()
		}
		if c.Services.GRPC["grpc-service-1"].Port != 100 || c.Services.GRPC["grpc-service-2"].Port != 100 {
			t.Error()
		}
		if c.Services.GRPC["grpc-service-1"].Methods == nil || c.Services.GRPC["grpc-service-2"].Methods == nil {
			t.Error()
		}
		if c.Services.GRPC["grpc-service-1"].Conditions == nil || c.Services.GRPC["grpc-service-2"].Conditions == nil {
			t.Error()
		}
		{
			service := c.Services.GRPC["grpc-service-1"]
			methods := service.Methods
			if service.Port != 100 || service.Definition != "definition-1" {
				t.Error()
			}
			{
				if methods["method-1"].Request != "method-1-request" || methods["method-1"].Response != "method-1-response" {
					t.Error()
				}
			}
			{
				conditions := service.Conditions
				if conditions["method-1"][0].Request == nil || conditions["method-1"][0].Response == nil {
					t.Error()
				}
				if conditions["method-2"][0].Request == nil || conditions["method-2"][0].Response == nil {
					t.Error()
				}
			}
		}
		{
			service := c.Services.GRPC["grpc-service-2"]
			methods := service.Methods
			if service.Port != 100 || service.Definition != "definition-2" {
				t.Error()
			}
			{
				if methods["method-1"].Request != "method-1-request" || methods["method-1"].Response != "method-1-response" {
					t.Error()
				}
			}
			{
				conditions := service.Conditions
				if conditions["method-1"][0].Request == nil || conditions["method-1"][0].Response == nil {
					t.Error()
				}
				if conditions["method-2"][0].Request == nil || conditions["method-2"][0].Response == nil {
					t.Error()
				}
			}
		}
	}
}
