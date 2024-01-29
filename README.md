# Celestial Tenant Service

An API to store configurations for other services

## To Run

### Create the database

From a console of a Mysql server run:

```
create database celestial_tenant
```

Add your database configuration to app.conf

Example of my local configuration:

```
schema = celestial_tenant
dbHost = localhost
dbUser = root
dbPort = 3306
dbPass = root
```

### Start the web server:

```
revel run .
```

## Code Layout

The directory structure of a generated Revel application:

    conf/             Configuration directory
        app.conf      Main app configuration file
        routes        Routes definition file

    app/              App sources
        init.go       Interceptor registration
        controllers/  App controllers go here
        views/        Templates directory

    messages/         Message files

    public/           Public static assets
        css/          CSS files
        js/           Javascript files
        images/       Image files

    tests/            Test suites

## Help

- The [Getting Started with Revel](http://revel.github.io/tutorial/gettingstarted.html).
- The [Revel guides](http://revel.github.io/manual/index.html).
- The [Revel sample apps](http://revel.github.io/examples/index.html).
- The [API documentation](https://godoc.org/github.com/revel/revel).
