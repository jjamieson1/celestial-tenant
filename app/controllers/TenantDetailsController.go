package controllers

import (
	"celestial-tenant/services"
	"errors"
	"fmt"
	"strings"

	userClient "github.com/jjamieson1/celestial-sdk/clients/user"
	"github.com/jjamieson1/celestial-sdk/models"
	"github.com/jjamieson1/celestial-sdk/utilities"
	"github.com/revel/revel"
)

type Api struct {
	*revel.Controller
}

func (c Api) SetTenantUserServiceProvider(tenantId string, provider string) revel.Result {
	c.Response.SetStatus(200)
	return c.Result
}

func (c Api) AddTenant() revel.Result {
	var tenant models.Tenant
	err := c.Params.BindJSON(&tenant)
	if err != nil {
		revel.AppLog.Errorf("error binding the JSON to the model with error: %s", err.Error())
		return c.RenderJSON(err.Error())
	}
	// Validation
	c.Validation.Required(tenant.CommonName).Key("commonName").Message("commonName is required")

	if services.DoesTenantNameExist(tenant.CommonName) {
		err := errors.New("common_name already exists as a tenant")
		return errorReturn(err, c.Controller)
	}
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.Response.Status = 400
		return errorReturn(err, c.Controller)
	}

	tenant, err = services.AddUpdateTenantDetails("", tenant)
	if err != nil {
		return errorReturn(err, c.Controller)
	}
	return c.RenderJSON(tenant)
}

func (c Api) GetTenantById() revel.Result {

	// Use the api/app key in the headers to return the tenantId
	apiKey := c.Request.Header.Get("api-key")
	appKey := c.Request.Header.Get("app-key")

	c.Validation.Required(apiKey).Key("api-key").Message("api-key missing in the header")
	c.Validation.Required(appKey).Key("app-key").Message("app-key missing in the header")
	c.Validation.Length(appKey, 36).Key("app-key").Message("app-key is not a valid uuid")
	c.Validation.Length(apiKey, 36).Key("api-key").Message("api-key is not a valid uuid")
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.Response.Status = 400
		return c.RenderJSON(c.Validation.Errors)
	}

	revel.AppLog.Debugf("requesting tenant details with app-key: %v", appKey)

	tenant, err := services.GetTenantDetails(apiKey, appKey)
	if err != nil {
		c.Response.SetStatus(400)
		revel.AppLog.Debugf("unable to get tenant with error: %s", err.Error())
		return c.RenderText(err.Error())
	}
	return c.RenderJSON(tenant)
}

func (c Api) GetTenantDetails(tenantId string) revel.Result {
	jwt := utilities.HandelHeaderJWT(c.Controller)
	if c.Validation.HasErrors() {
		return utilities.HandleValidationError(c.Validation.Errors, 403, c.Controller)
	}
	revel.AppLog.Debugf("jwt found to be: %v", jwt)
	user, status, err := userClient.GetAccount(jwt, tenantId)
	if status != 200 {
		return utilities.HandleInternalError("authentication", fmt.Sprintf("http status %v", status), 500, c.Controller)
	}

	if err != nil {
		return utilities.HandleInternalError("tenantId", err.Error(), 500, c.Controller)
	}

	if !utilities.IsAdmin(user.Roles) {
		return utilities.HandleAuthorizationError("role", "missing required admin role", 401, c.Controller)
	}

	tenant, err := services.GetTenantDetailsByTenantId(tenantId)
	if err != nil {
		return utilities.HandleInternalError("tenant Service", err.Error(), 500, c.Controller)
	}

	return c.RenderJSON(tenant)
}

func (c Api) GetTenants() revel.Result {
	revel.AppLog.Debug("requesting all tenants")

	tenant, err := services.GetTenants()
	if err != nil {
		c.Response.SetStatus(400)
		return c.RenderJSON(err)
	}
	return c.RenderJSON(tenant)
}

func (c Api) GetAllChildrenOfTenant(tenantId string) revel.Result {
	revel.AppLog.Debugf("getting all  tenant details for children of tenantId: %v", tenantId)
	tenants, _ := services.GetAllTenantChildrenDetails(tenantId)

	return c.RenderJSON(tenants)
}

func (c Api) GetTenantByUrl(requestUrl string) revel.Result {

	revel.AppLog.Debugf("requesting tenant details for url: %v", requestUrl)
	// In case a port number is appended remove
	u := strings.Split(requestUrl, ":")
	t := make(map[string]string)
	var err error
	t["tenantId"], err = services.GetTenantIdByUrl(u[0])
	if err != nil {
		c.Response.SetStatus(400)
		return c.RenderJSON(err)
	}
	return c.RenderJSON(t)
}

func (c Api) UpdateTenant(tenantId string) revel.Result {
	var tenant models.Tenant
	err := c.Params.BindJSON(&tenant)
	if err != nil {
		c.Response.SetStatus(400)
		return c.RenderJSON(err)
	}
	tenant, err = services.AddUpdateTenantDetails(tenantId, tenant)
	if err != nil {
		c.Response.SetStatus(400)
		return c.RenderJSON(err)
	}
	return c.RenderJSON(tenant)
}

func (c Api) DeleteTenant(tenantId string) revel.Result {
	err := services.DeleteTenant(tenantId)
	if err != nil {
		c.Response.SetStatus(400)
		return c.RenderJSON(err)
	}
	return c.Result
}

func (c Api) GetTenantType() revel.Result {
	tenantId := c.Request.Header.Get("tenantId")
	tenantTypes, err := services.GetTenantTypesByTenantId(tenantId)
	if err != nil {
		c.Response.SetStatus(400)
		return c.RenderJSON(err)
	}
	return c.RenderJSON(tenantTypes)
}

func (c Api) AddTenantType() revel.Result {
	var tenantType models.TenantType
	err := c.Params.BindJSON(&tenantType)
	if err != nil {
		return c.RenderJSON(err)
	}
	i, err := services.AddTenantType(tenantType)
	if err != nil {
		c.Response.SetStatus(400)
		return c.RenderJSON(err)
	}
	tenantType.Id = i
	return c.RenderJSON(tenantType)
}

func (c Api) DeleteTenantType() revel.Result {
	return c.RenderError(errors.New("not implemented"))
}

func (c Api) UpdateTenantType() revel.Result {
	return c.RenderError(errors.New("not implemented"))

}

func errorReturn(err error, c *revel.Controller) revel.Result {
	c.Response.Status = 400
	c.Response.ContentType = "application/json"
	errorResponse := make(map[string]string)
	errorResponse["message"] = err.Error()
	revel.AppLog.Debugf("%+s", errorResponse)
	return c.RenderJSON(errorResponse)
}
