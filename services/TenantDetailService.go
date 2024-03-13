package services

import (
	"celestial-tenant/app"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jjamieson1/celestial-sdk/models"
	"github.com/revel/revel"
)

func AddUpdateTenantDetails(tenantId string, tenant models.Tenant) (models.Tenant, error) {
	// Assuming userId is static for this function
	userID := "jamie"
	var query string
	var operation string

	if tenantId == "" {
		tenant.TenantId = uuid.New().String()
		operation = "adding new"
		query = `INSERT INTO tenant (
					tenant_id,
					parent_tenant_id,
 					url, 
 					common_name,
 					logo_primary_url,
 					is_available,
					created_by)
				VALUES (?,?,?,?,?,?,?)`
	} else {
		operation = "updating"
		query = `UPDATE tenant_details SET 
						parent_tenant_id=?,
						url=?, 
						common_name=?, 
						primary_logo_url=?,
						is_available=?
						WHERE tenant_id=?`
	}

	revel.AppLog.Infof("%s tenant (tenantId: %v)", operation, tenant.TenantId)

	stmt, err := app.DB.Prepare(query)
	if err != nil {
		errorMsg := fmt.Sprintf("error preparing query: %v, error: %v", query, err.Error())
		revel.AppLog.Errorf(errorMsg)
		return tenant, errors.New(errorMsg)
	}
	defer stmt.Close()

	var execResult sql.Result
	if tenantId == "" {
		execResult, err = stmt.Exec(
			tenant.TenantId,
			tenant.ParentTenantId,
			tenant.Url,
			tenant.CommonName,
			tenant.LogoUrl,
			tenant.IsAvailable,
			userID,
		)
	} else {
		execResult, err = stmt.Exec(
			tenant.ParentTenantId,
			tenant.Url,
			tenant.CommonName,
			tenant.LogoUrl,
			tenant.IsAvailable,
			tenantId,
		)
	}
	if err != nil {
		errorMsg := fmt.Sprintf("error performing query: %v, error: %v", query, err.Error())
		revel.AppLog.Errorf(errorMsg)
		return tenant, errors.New(errorMsg)
	}

	rowsAffected, err := execResult.RowsAffected()
	if err != nil {
		return tenant, fmt.Errorf("error getting rows affected: %v", err)
	}

	if rowsAffected == 0 {
		revel.AppLog.Infof("No rows affected for %s operation", operation)
	}

	if tenantId == "" {
		for _, tenantType := range tenant.TenantTypes {
			err = AddTenantTypeToTenant(tenantType.Id, tenant.TenantId)
			if err != nil {
				revel.AppLog.Errorf(err.Error())
			}
		}
		t, err := GetTenantTypesByTenantId(tenant.TenantId)
		if err != nil {
			revel.AppLog.Errorf(err.Error())
		}
		apiKey, err := CreateSecretKeys(tenant.TenantId)
		if err != nil {
			revel.AppLog.Errorf(err.Error())
		}
		tenant.SecretKeys.AppKey = tenant.TenantId
		tenant.SecretKeys.ApiKey = apiKey
		tenant.TenantTypes = t
	}
	return tenant, nil
}
func GetTenantDetailsByTenantId(tenantId string) (tenant models.Tenant, err error) {

	query := `SELECT tenant_id, parent_tenant_id, url, common_name, logo_primary_url, is_available  
				FROM  tenant WHERE tenant_id = ?`

	err = app.DB.QueryRow(query, tenantId).Scan(
		&tenant.TenantId,
		&tenant.ParentTenantId,
		&tenant.Url,
		&tenant.CommonName,
		&tenant.LogoUrl,
		&tenant.IsAvailable,
	)

	if err != nil && err != sql.ErrNoRows {
		error := fmt.Sprint(err.Error())
		revel.AppLog.Errorf(error)
		return tenant, errors.New(error)
	} else if err == sql.ErrNoRows {
		e := fmt.Errorf("unable to find tenant with those credentials (tenantId: %s)", tenantId)
		revel.AppLog.Info(e.Error())
		return tenant, e
	} else {
		revel.AppLog.Debugf("retrieved tenant_details for tenantId: %v ", tenant.TenantId)
	}

	tenant.TenantTypes, err = GetTenantTypesByTenantId(tenant.TenantId)
	tenant.ServiceProviders, _ = getServiceProviderData(tenant.TenantId)
	return tenant, err
}
func GetTenantDetails(apiKey, appKey string) (models.Tenant, error) {
	var tenant models.Tenant
	revel.AppLog.Infof("retrieving tenant  (appKey: %v)", appKey)

	query := `SELECT t.tenant_id, t.parent_tenant_id, t.url, t.common_name, t.logo_primary_url, t.is_available  FROM  tenant as t
				JOIN secret_keys as sk on t.tenant_id = sk.tenant_id 
				WHERE sk.tenant_id = ? and sk.api_key = ?`

	err := app.DB.QueryRow(query, appKey, apiKey).Scan(
		&tenant.TenantId,
		&tenant.ParentTenantId,
		&tenant.Url,
		&tenant.CommonName,
		&tenant.LogoUrl,
		&tenant.IsAvailable,
	)

	if err != nil && err != sql.ErrNoRows {
		error := fmt.Sprint(err.Error())
		revel.AppLog.Errorf(error)
		return tenant, errors.New(error)
	} else if err == sql.ErrNoRows {
		e := fmt.Errorf("unable to find tenant with those credentials (AppKey: %s, ApiKey: %s)", appKey, apiKey)
		revel.AppLog.Info(e.Error())
		return tenant, e
	} else {
		revel.AppLog.Debugf("retrieved tenant_details for tenantId: %v ", tenant.TenantId)
	}

	tenant.TenantTypes, err = GetTenantTypesByTenantId(tenant.TenantId)
	tenant.ServiceProviders, _ = getServiceProviderData(tenant.TenantId)
	return tenant, err

}

func GetTenants() ([]models.Tenant, error) {
	var tenant models.Tenant
	var tenants []models.Tenant
	revel.AppLog.Info("retrieving tenants")

	query := `SELECT * FROM  tenant_details WHERE parent_tenant_id = "none"`

	stmt, err := app.DB.Prepare(query)
	if err != nil {
		error := fmt.Sprintf("error preparing query: %v, error: %v", query, err.Error())
		revel.AppLog.Errorf(error)
		return tenants, errors.New(error)
	}

	results, err := stmt.Query()
	if err != nil {
		error := fmt.Sprintf("error performing query: %v, error: %v", query, err.Error())
		revel.AppLog.Errorf(error)
		return tenants, errors.New(error)
	}
	defer stmt.Close()

	for results.Next() {
		err := results.Scan(
			&tenant.TenantId,
			&tenant.ParentTenantId,
			&tenant.Url,
			&tenant.CommonName,
			&tenant.LogoUrl,
			&tenant.IsAvailable,
		)
		if err != nil {
			error := fmt.Sprintf("error mapping query to model: %v, error: %v", query, err.Error())
			revel.AppLog.Errorf(error)
			return tenants, errors.New(error)
		}
		tenants = append(tenants, tenant)
	}

	return tenants, err

}

func GetAllTenantChildrenDetails(tenantId string) ([]models.Tenant, error) {
	var tenant models.Tenant
	var tenants []models.Tenant

	query := `SELECT * FROM  tenant_details WHERE parent_tenant_id = ?`

	stmt, err := app.DB.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("error preparing query: %v", err)
	}
	defer stmt.Close()

	results, err := stmt.Query(tenantId)
	if err != nil {
		return nil, fmt.Errorf("error performing query: %v", err)
	}
	defer results.Close()

	for results.Next() {
		err := results.Scan(
			&tenant.TenantId,
			&tenant.ParentTenantId,
			&tenant.Url,
			&tenant.CommonName,
			&tenant.LogoUrl,
			&tenant.IsAvailable,
		)
		if err != nil {
			return tenants, fmt.Errorf("error mapping query to model: %v", err)
		}
		tenants = append(tenants, tenant)
	}
	return tenants, nil
}

func GetTenantIdByUrl(url string) (string, error) {
	revel.AppLog.Debugf("looking up tenant by url: %s", url)

	query := `SELECT t.tenant_id FROM  tenant as t
				WHERE t.url = ?`

	revel.AppLog.Infof("getting tenant  by url: %v", url)
	var tenantId string
	err := app.DB.QueryRow(query, url).Scan(&tenantId)

	if err != nil && err != sql.ErrNoRows {
		error := fmt.Sprint(err.Error())
		revel.AppLog.Errorf(error)
		return tenantId, errors.New(error)
	} else if err == sql.ErrNoRows {
		error := fmt.Sprint(err.Error())
		revel.AppLog.Info(error)
		return tenantId, sql.ErrNoRows
	} else {
		revel.AppLog.Debugf("retrieved tenant_details for url: %v ", url)
	}

	return tenantId, err

}

func DeleteTenant(tenantId string) error {

	query := `DELETE FROM tenant_details 
						WHERE tenant_id=?`

	revel.AppLog.Infof("updating CMS article cmsId: %v for tenant: %v", tenantId)

	stmt, err := app.DB.Prepare(query)
	if err != nil {
		error := fmt.Sprintf("error performing query: %v, error: %v", query, err.Error())
		revel.AppLog.Errorf(error)
		return errors.New(error)
	}
	defer stmt.Close()

	_, err = stmt.Exec(tenantId)

	if err != nil && err != sql.ErrNoRows {
		error := fmt.Sprint(err.Error())
		revel.AppLog.Errorf(error)
		return errors.New(error)
	} else if err == sql.ErrNoRows {
		error := fmt.Sprint(err.Error())
		revel.AppLog.Info(error)
		return sql.ErrNoRows
	} else {
		revel.AppLog.Debugf("deleted tenantId: %v ", tenantId)
	}

	return err
}

// Tenant Type services

func AddTenantType(tenantType models.TenantType) (string, error) {
	i := uuid.New().String()

	query := `INSERT INTO tenant_type (tenant_type_id, tenant_type_name) VALUES (?,?)`

	revel.AppLog.Infof("adding new tenant type: %v, for tenantId: %v", tenantType.Name, tenantType.Id)

	stmt, err := app.DB.Prepare(query)
	if err != nil {
		error := fmt.Sprintf("error performing query: %v, error: %v", query, err.Error())
		revel.AppLog.Errorf(error)
		return i, errors.New(error)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		i,
		tenantType.Name,
	)
	if err != nil {
		error := fmt.Sprintf("error performing query: %v, error: %v", query, err.Error())
		revel.AppLog.Errorf(error)
		return i, errors.New(error)
	}

	return i, err
}

func AddTenantTypeToTenant(tenantTypeId, tenantId string) error {
	query := `INSERT INTO tenant_type_to_tenant (tenant_type_id, tenant_id) 
		VALUES (?,?)`
	stmt, err := app.DB.Prepare(query)
	if err != nil {
		error := fmt.Sprintf("error performing query: %v, error: %v", query, err.Error())
		revel.AppLog.Errorf(error)
		return errors.New(error)
	}
	defer stmt.Close()

	_, err = stmt.Exec(tenantTypeId, tenantId)
	if err != nil {
		revel.AppLog.Errorf(err.Error())
	}
	return err
}

func GetTenantTypesByTenantId(tenantId string) ([]models.TenantType, error) {
	var tenantType models.TenantType
	var tenantTypes []models.TenantType
	var err error

	query := `SELECT tt.tenant_type_id, tenant_type_name FROM tenant_type as tt
			JOIN tenant_type_to_tenant as t2t on tt.tenant_type_id = t2t.tenant_type_id 
			WHERE t2t.tenant_id = ?
	`
	stmt, err := app.DB.Prepare(query)
	if err != nil {
		error := fmt.Sprintf("error preparing query: %v, error: %v", query, err.Error())
		revel.AppLog.Errorf(error)
		return tenantTypes, errors.New(error)
	}

	results, err := stmt.Query(tenantId)
	if err != nil {
		error := fmt.Sprintf("error performing query: %v, error: %v", query, err.Error())
		revel.AppLog.Errorf(error)
		return tenantTypes, errors.New(error)
	}
	defer stmt.Close()
	for results.Next() {
		err := results.Scan(
			&tenantType.Id,
			&tenantType.Name,
		)
		if err != nil {
			error := fmt.Sprintf("error mapping query to model: %v, error: %v", query, err.Error())
			revel.AppLog.Errorf(error)
			return tenantTypes, errors.New(error)
		}
		tenantTypes = append(tenantTypes, tenantType)
	}
	return tenantTypes, err
}

func getServiceProviderData(tenantId string) (serviceProviders []models.ServiceProviders, err error) {
	q := `SELECT sptd.id, sptd.URL, sp.id, sp.name FROM service_provider_tenant_data AS sptd 
			JOIN service_provider AS sp ON sptd.service_provider_id = sp.id 
			WHERE sptd.tenant_id = ?`
	stmt, err := app.DB.Prepare(q)
	if err != nil {
		error := fmt.Sprintf("error performing query: %v, error: %v", q, err.Error())
		revel.AppLog.Errorf(error)
		return serviceProviders, errors.New(error)
	}
	defer stmt.Close()

	results, err := stmt.Query(tenantId)
	if err != nil {
		error := fmt.Sprintf("error performing query: %v, error: %v", q, err.Error())
		revel.AppLog.Errorf(error)
		return serviceProviders, errors.New(error)
	}

	var sp models.ServiceProviders

	for results.Next() {
		err := results.Scan(
			&sp.Id,
			&sp.BaseURL,
			&sp.ServiceProvider.Id,
			&sp.ServiceProvider.Name,
		)
		if err != nil {
			error := fmt.Sprintf("error mapping query to model: %v, error: %v", q, err.Error())
			revel.AppLog.Errorf(error)
			return serviceProviders, errors.New(error)
		}
		serviceProviders = append(serviceProviders, sp)
	}

	return serviceProviders, err
}

func CreateSecretKeys(tenantId string) (string, error) {
	apiKey := uuid.New().String()
	revel.AppLog.Debugf("creating a new key set for appKey: %s apiKey %s", tenantId, apiKey)
	query := `INSERT INTO secret_keys (tenant_id, api_key) 
		VALUES (?,?)`
	stmt, err := app.DB.Prepare(query)
	if err != nil {
		error := fmt.Sprintf("error performing query: %v, error: %v", query, err.Error())
		revel.AppLog.Errorf(error)
		return "", errors.New(error)
	}
	defer stmt.Close()

	_, err = stmt.Exec(tenantId, apiKey)
	if err != nil {
		revel.AppLog.Errorf(err.Error())
	}
	return apiKey, err
}

func DoesTenantNameExist(name string) bool {
	var count int
	q := `select count(*) from tenant where common_name = ?`
	app.DB.QueryRow(q, name).Scan(&count)
	return count != 0
}
