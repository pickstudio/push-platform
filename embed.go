package pushplatform

import "embed"

// StaticSwaggerUI Static is a collection of pre-built static files for swagger web ui.
//go:embed static/swagger-ui
var StaticSwaggerUI embed.FS

// StaticAdmin Static is a collection of pre-built static files for admin web ui.
//go:embed static/admin
var StaticAdmin embed.FS

// OAPISpecYAML is the Open API Specifications Manifest document that defines golive HTTP API.
//go:embed api/oapi
var OAPISpecYAML embed.FS
