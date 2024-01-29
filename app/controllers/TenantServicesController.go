package controllers

import (
	"celestial-tenant/services"

	"github.com/revel/revel"
)

func (c Api) GetTenantServiceProvidersByType(tenantId, providerType string) revel.Result {
	tenantProviders, err := services.GetProvidersForTenantByType(tenantId, providerType)
	if err != nil {
		c.Response.SetStatus(400)
		return c.Result
	}
	return c.RenderJSON(tenantProviders)
}

func (c Api) GetAllAvailableProvidersByType(providerType string) revel.Result {
	tenantProviders, err := services.GetProviderByType(providerType)
	if err != nil {
		c.Response.SetStatus(400)
		return c.Result
	}
	return c.RenderJSON(tenantProviders)
}

func (c Api) GetAllProvidersForTenant(tenantId string) revel.Result {
	tenantProviders, err := services.GetAllProvidersForTenant(tenantId)
	if err != nil {
		c.Response.SetStatus(400)
		return c.Result
	}
	return c.RenderJSON(tenantProviders)
}
