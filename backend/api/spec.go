package api

import _ "embed"

// Spec это сырой OpenAPI-документ, вшитый в бинарь для раздачи через Swagger UI
//
//go:embed swagger.yaml
var Spec []byte
