package services

import (
	"celestial-tenant/app"
	"errors"
	"fmt"
	"time"

	"github.com/jjamieson1/celestial-sdk/models"
	"github.com/revel/revel"
	"github.com/revel/revel/cache"
)

func GetProvidersForTenantByType(tenantId, providerType string) ([]models.TenantProvider, error) {
	var tenantProvider models.TenantProvider
	var tenantProviders []models.TenantProvider
	var err error

	if err := cache.Get("tenant-"+providerType+"-"+tenantId, &tenantProviders); err != nil {
		revel.AppLog.Infof("cache miss, getting new data, error: %v", err.Error())

		query := `SELECT tenant_provider.*, eden_adapter.*, eden_provider_type.*, auth_strategy.*
				FROM tenant_provider
				    JOIN eden_adapter ON
				        tenant_provider.eden_adapter_id = eden_adapter.id
                    JOIN eden_provider_type ON
                        eden_adapter.eden_provider_type_id = eden_provider_type.id
                    JOIN auth_strategy ON
                        eden_adapter.auth_strategy_id = auth_strategy.id
					WHERE 
						tenant_provider.tenant_id = ? 
					AND eden_provider_type.eden_provider_type_name = ?
`

		revel.AppLog.Infof("getting provider information for tenantId: %v for provider_type: %v)", tenantId, providerType)

		stmt, err := app.DB.Prepare(query)
		if err != nil {
			error := fmt.Sprintf("error preparing query: %v, error: %v", query, err.Error())
			revel.AppLog.Errorf(error)
			return tenantProviders, errors.New(error)
		}

		results, err := stmt.Query(tenantId, providerType)
		if err != nil {
			error := fmt.Sprintf("error performing query: %v, error: %v", query, err.Error())
			revel.AppLog.Errorf(error)
			return tenantProviders, errors.New(error)
		}
		for results.Next() {
			err := results.Scan(
				&tenantProvider.Id,
				&tenantProvider.Adapter.Id,
				&tenantProvider.TenantId,
				&tenantProvider.CalloutUrl,
				&tenantProvider.UserName,
				&tenantProvider.Password,
				&tenantProvider.ApiKey,
				&tenantProvider.AppKey,
				&tenantProvider.Token,
				&tenantProvider.RefreshToken,
				&tenantProvider.Adapter.Id,
				&tenantProvider.Adapter.Name,
				&tenantProvider.Adapter.PluginName,
				&tenantProvider.Adapter.ProviderType.Id,
				&tenantProvider.Adapter.AuthStrategy.Id,
				&tenantProvider.Adapter.AdapterUrl,
				&tenantProvider.Adapter.Enabled,
				&tenantProvider.Adapter.ProviderType.Id,
				&tenantProvider.Adapter.ProviderType.Name,
				&tenantProvider.Adapter.AuthStrategy.Id,
				&tenantProvider.Adapter.AuthStrategy.Name,
				&tenantProvider.Adapter.AuthStrategy.Parameters,
			)
			if err != nil {
				error := fmt.Sprintf("error mapping query to model: %v, error: %v", query, err.Error())
				revel.AppLog.Errorf(error)
				return tenantProviders, errors.New(error)
			}
			tenantProviders = append(tenantProviders, tenantProvider)
		}
		go cache.Set("tenant-"+tenantId, tenantProviders, 30*time.Minute)
	} else {
		revel.AppLog.Debugf("cache hit, returning tenantId: %v", tenantId)
	}
	revel.AppLog.Debugf("returning provider from request for tenantId: %v, adapter name: %v with url: %v", tenantId, tenantProviders[0].Adapter.Name, tenantProviders[0].CalloutUrl)
	return tenantProviders, err
}

func GetAllProvidersForTenant(tenantId string) ([]models.TenantProvider, error) {
	var tenantProvider models.TenantProvider
	var tenantProviders []models.TenantProvider

	query := `SELECT tenant_provider.*, eden_adapter.*, eden_provider_type.*, auth_strategy.*
				FROM tenant_provider
				    JOIN eden_adapter ON
				        tenant_provider.eden_adapter_id = eden_adapter.id
                    JOIN eden_provider_type ON
                        eden_adapter.eden_provider_type_id = eden_provider_type.id
                    JOIN auth_strategy ON
                        eden_adapter.auth_strategy_id = auth_strategy.id
					WHERE 
						tenant_provider.tenant_id = ?
`

	revel.AppLog.Infof("getting provider information for tenantId: %v for provider_type: %v)", tenantId)

	stmt, err := app.DB.Prepare(query)
	if err != nil {
		error := fmt.Sprintf("error preparing query: %v, error: %v", query, err.Error())
		revel.AppLog.Errorf(error)
		return tenantProviders, errors.New(error)
	}

	results, err := stmt.Query(tenantId)
	if err != nil {
		error := fmt.Sprintf("error performing query: %v, error: %v", query, err.Error())
		revel.AppLog.Errorf(error)
		return tenantProviders, errors.New(error)
	}
	for results.Next() {
		err := results.Scan(
			&tenantProvider.Id,
			&tenantProvider.Adapter.Id,
			&tenantProvider.TenantId,
			&tenantProvider.CalloutUrl,
			&tenantProvider.UserName,
			&tenantProvider.Password,
			&tenantProvider.ApiKey,
			&tenantProvider.AppKey,
			&tenantProvider.Token,
			&tenantProvider.RefreshToken,
			&tenantProvider.Adapter.Id,
			&tenantProvider.Adapter.Name,
			&tenantProvider.Adapter.PluginName,
			&tenantProvider.Adapter.ProviderType.Id,
			&tenantProvider.Adapter.AuthStrategy.Id,
			&tenantProvider.Adapter.AdapterUrl,
			&tenantProvider.Adapter.Enabled,
			&tenantProvider.Adapter.ProviderType.Id,
			&tenantProvider.Adapter.ProviderType.Name,
			&tenantProvider.Adapter.AuthStrategy.Id,
			&tenantProvider.Adapter.AuthStrategy.Name,
			&tenantProvider.Adapter.AuthStrategy.Parameters,
		)
		if err != nil {
			error := fmt.Sprintf("error mapping query to model: %v, error: %v", query, err.Error())
			revel.AppLog.Errorf(error)
			return tenantProviders, errors.New(error)
		}
		tenantProviders = append(tenantProviders, tenantProvider)
	}
	return tenantProviders, err

}

func GetProviderByType(providerType string) ([]models.Adapter, error) {
	var adapter models.Adapter
	var adapters []models.Adapter
	var err error

	if err := cache.Get(providerType+"-adapter", &adapters); err != nil {
		revel.AppLog.Infof("cache miss, getting new data, error: %v", err.Error())

		query := `SELECT eden_adapter.*, eden_provider_type.*, auth_strategy.*
				FROM eden_adapter
                    JOIN eden_provider_type ON
                        eden_adapter.eden_provider_type_id = eden_provider_type.id
                    JOIN auth_strategy ON
                        eden_adapter.auth_strategy_id = auth_strategy.id
					WHERE
						eden_provider_type.eden_provider_type_name = ?
					
`

		revel.AppLog.Infof("getting provider information for tenantId: %v for provider_type: %v)", providerType)

		stmt, err := app.DB.Prepare(query)
		if err != nil {
			error := fmt.Sprintf("error preparing query: %v, error: %v", query, err.Error())
			revel.AppLog.Errorf(error)
			return adapters, errors.New(error)
		}

		results, err := stmt.Query(providerType)
		if err != nil {
			error := fmt.Sprintf("error performing query: %v, error: %v", query, err.Error())
			revel.AppLog.Errorf(error)
			return adapters, errors.New(error)
		}
		for results.Next() {
			err := results.Scan(
				&adapter.Id,
				&adapter.Name,
				&adapter.PluginName,
				&adapter.ProviderType.Id,
				&adapter.AuthStrategy.Id,
				&adapter.AdapterUrl,
				&adapter.Enabled,
				&adapter.ProviderType.Id,
				&adapter.ProviderType.Name,
				&adapter.AuthStrategy.Id,
				&adapter.AuthStrategy.Name,
				&adapter.AuthStrategy.Parameters,
			)
			if err != nil {
				error := fmt.Sprintf("error mapping query to model: %v, error: %v", query, err.Error())
				revel.AppLog.Errorf(error)
				return adapters, errors.New(error)
			}
			adapters = append(adapters, adapter)
		}
		go cache.Set(providerType+"-adapter", adapters, 30*time.Minute)
	} else {
		revel.AppLog.Debugf("cache hit, returning adapter list for provider type: %v", providerType)
	}
	return adapters, err
}
