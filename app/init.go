package app

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pressly/goose"
	_ "github.com/revel/modules"
	"github.com/revel/revel"
)

var (
	// AppVersion revel app version (ldflags)
	AppVersion string
	// Database connection
	DB *sql.DB
	// BuildTime revel app build-time (ldflags)
	BuildTime string
)

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.BeforeAfterFilter,       // Call the before and after filter functions
		revel.ActionInvoker,           // Invoke the action.
	}

	revel.OnAppStart(InitDB)
	revel.OnAppStart(TestDataBase)
}

// HeaderFilter adds common security headers
// There is a full implementation of a CSRF filter in
// https://github.com/revel/modules/tree/master/csrf
var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")
	c.Response.Out.Header().Add("Referrer-Policy", "strict-origin-when-cross-origin")

	fc[0](c, fc[1:]) // Execute the next filter stage.
}

func InitDB() {
	mysqlDsnSchema := ""
	mysqlDsn := ""
	schema := revel.Config.StringDefault("schema", "celesital_cms")
	dbHost := revel.Config.StringDefault("dbHost", "localhost")
	dbUser := revel.Config.StringDefault("dbUser", "root")
	dbPort := revel.Config.StringDefault("dbPort", "3306")
	dbPass := revel.Config.StringDefault("dbPass", "root")

	fmt.Printf("using database schema: %s \n\n", schema)

	if dbUser != "" && dbPass != "" && dbHost != "" && dbPort != "" {
		mysqlDsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/", dbUser, dbPass, dbHost, dbPort)
		mysqlDsnSchema = fmt.Sprintf("%s%s", mysqlDsn, schema)
		flags := os.Getenv("MYSQL_FLAGS")
		if flags != "" {
			mysqlDsn = fmt.Sprintf("%s?%s", mysqlDsn, flags)
			mysqlDsnSchema = fmt.Sprintf("%s?%s", mysqlDsnSchema, flags)
		}
	} else {
		revel.AppLog.Error("no database configuration was provided")
	}
	var err error
	DB, err = sql.Open("mysql", mysqlDsnSchema)
	if err != nil {
		revel.AppLog.Errorf("Failed to connect to database with error: %s", err.Error())
	}
	if err := DB.Ping(); err != nil {
		revel.AppLog.Error("Failed to connect to schema, attempting to create it")
		DB, err = sql.Open("mysql", mysqlDsn)
		if err != nil {
			revel.AppLog.Error(err.Error())
		}

		DB.Exec(fmt.Sprintf("create database if not exists %s", schema))
		DB.Close()

		DB, err = sql.Open("mysql", mysqlDsnSchema)
		if err != nil {
			revel.AppLog.Error(err.Error())
		}
		revel.AppLog.Info("DB Connected")
	}

	if err := DB.Ping(); err != nil {
		revel.AppLog.Errorf("unable to ping database with error: %s", err.Error())
	}

	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(2 * time.Minute)

	revel.AppLog.Infof("database connections in use: %v", DB.Stats().InUse)
	revel.AppLog.Infof("database open connections: %v", DB.Stats().OpenConnections)

	err = doMigrations("./migrations", DB)
	if err != nil {
		revel.AppLog.Errorf("Failed to run migrations: %s", err.Error())
	}
}

func doMigrations(dir string, db *sql.DB) error {
	goose.SetDialect("mysql")
	return goose.Run("up", db, dir)
}

func TestDataBase() {
	q := "select count(*) from tenant_type"
	var num int
	err := DB.QueryRow(q).Scan(&num)
	if err != nil {
		revel.AppLog.Errorf("error running query with error: %s", err.Error())
	}
	revel.AppLog.Infof("Database test returned %v rows", num)
}
