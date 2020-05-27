package apidoc

import (
	"encoding/json"
	"github.com/go-openapi/spec"
	"testing"

	"github.com/stretchr/testify/assert"
)

var info = spec.Info{
	InfoProps: spec.InfoProps{
		Version: "1.0.9-abcd",
		Title:   "Swagger Sample API",
		Description: "A sample API that uses a petstore as an example to demonstrate features in " +
			"the swagger-2.0 specification",
		TermsOfService: "http://helloreverb.com/terms/",
		Contact:        &spec.ContactInfo{Name: "wordnik api team", URL: "http://developer.wordnik.com"},
		License: &spec.License{
			Name: "Creative Commons 4.0 International",
			URL:  "http://creativecommons.org/licenses/by/4.0/",
		},
	},
	VendorExtensible: spec.VendorExtensible{Extensions: map[string]interface{}{"x-framework": "go-swagger"}},
}

var paths = spec.Paths{
	VendorExtensible: spec.VendorExtensible{Extensions: map[string]interface{}{"x-framework": "go-swagger"}},
	Paths: map[string]spec.PathItem{
		"/": {
			Refable: spec.Refable{Ref: spec.MustCreateRef("cats")},
		},
	},
}

var specc = spec.Swagger{
	SwaggerProps: spec.SwaggerProps{
		ID:          "http://localhost:3849/api-docs",
		Swagger:     "2.0",
		Consumes:    []string{"application/json", "application/x-yaml"},
		Produces:    []string{"application/json"},
		Schemes:     []string{"http", "https"},
		Info:        &info,
		Host:        "some.api.out.there",
		BasePath:    "/",
		Paths:       &paths,
		Definitions: map[string]spec.Schema{"Category": {SchemaProps: spec.SchemaProps{Type: []string{"string"}}}},
		Parameters: map[string]spec.Parameter{
			"categoryParam": {ParamProps: spec.ParamProps{Name: "category", In: "query"}, SimpleSchema: spec.SimpleSchema{Type: "string"}},
		},
		Responses: map[string]spec.Response{
			"EmptyAnswer": {
				ResponseProps: spec.ResponseProps{
					Description: "no data to return for this operation",
				},
			},
		},
		SecurityDefinitions: map[string]*spec.SecurityScheme{
			"internalApiKey": spec.APIKeyAuth("api_key", "header"),
		},
		Security: []map[string][]string{
			{"internalApiKey": {}},
		},
		Tags:         []spec.Tag{spec.NewTag("pets", "", nil)},
		ExternalDocs: &spec.ExternalDocumentation{Description: "the name", URL: "the url"},
	},
	VendorExtensible: spec.VendorExtensible{Extensions: map[string]interface{}{
		"x-some-extension": "vendor",
		"x-schemes":        []interface{}{"unix", "amqp"},
	}},
}

const speccJSON = `{
	"id": "http://localhost:3849/api-docs",
	"consumes": ["application/json", "application/x-yaml"],
	"produces": ["application/json"],
	"schemes": ["http", "https"],
	"swagger": "2.0",
	"info": {
		"contact": {
			"name": "wordnik api team",
			"url": "http://developer.wordnik.com"
		},
		"description": "A sample API that uses a petstore as an example to demonstrate features in the swagger-2.0` +
	` speccification",
		"license": {
			"name": "Creative Commons 4.0 International",
			"url": "http://creativecommons.org/licenses/by/4.0/"
		},
		"termsOfService": "http://helloreverb.com/terms/",
		"title": "Swagger Sample API",
		"version": "1.0.9-abcd",
		"x-framework": "go-swagger"
	},
	"host": "some.api.out.there",
	"basePath": "/",
	"paths": {"x-framework":"go-swagger","/":{"$ref":"cats"}},
	"definitions": { "Category": { "type": "string"} },
	"parameters": {
		"categoryParam": {
			"name": "category",
			"in": "query",
			"type": "string"
		}
	},
	"responses": { "EmptyAnswer": { "description": "no data to return for this operation" } },
	"securityDefinitions": {
		"internalApiKey": {
			"type": "apiKey",
			"in": "header",
			"name": "api_key"
		}
	},
	"security": [{"internalApiKey":[]}],
	"tags": [{"name":"pets"}],
	"externalDocs": {"description":"the name","url":"the url"},
	"x-some-extension": "vendor",
	"x-schemes": ["unix","amqp"]
}`

//
// func verifySpecSerialize(speccJSON []byte, specc Swagger) {
// 	expected := map[string]interface{}{}
// 	json.Unmarshal(speccJSON, &expected)
// 	b, err := json.MarshalIndent(specc, "", "  ")
// 	So(err, ShouldBeNil)
// 	var actual map[string]interface{}
// 	err = json.Unmarshal(b, &actual)
// 	So(err, ShouldBeNil)
// 	compareSpecMaps(actual, expected)
// }

/*
// assertEquivalent is currently unused
func assertEquivalent(t testing.TB, actual, expected interface{}) bool {
	if actual == nil || expected == nil || reflect.DeepEqual(actual, expected) {
		return true
	}

	actualType := reflect.TypeOf(actual)
	expectedType := reflect.TypeOf(expected)
	if reflect.TypeOf(actual).ConvertibleTo(expectedType) {
		expectedValue := reflect.ValueOf(expected)
		if swag.IsZero(expectedValue) && swag.IsZero(reflect.ValueOf(actual)) {
			return true
		}

		// Attempt comparison after type conversion
		if reflect.DeepEqual(actual, expectedValue.Convert(actualType).Interface()) {
			return true
		}
	}

	// Last ditch effort
	if fmt.Sprintf("%#v", expected) == fmt.Sprintf("%#v", actual) {
		return true
	}
	errFmt := "Expected: '%T(%#v)'\nActual:   '%T(%#v)'\n(Should be equivalent)!"
	return assert.Fail(t, errFmt, expected, expected, actual, actual)
}

// ShouldBeEquivalentTo is currently unused
func ShouldBeEquivalentTo(actual interface{}, expecteds ...interface{}) string {
	expected := expecteds[0]
	if actual == nil || expected == nil {
		return ""
	}

	if reflect.DeepEqual(expected, actual) {
		return ""
	}

	actualType := reflect.TypeOf(actual)
	expectedType := reflect.TypeOf(expected)
	if reflect.TypeOf(actual).ConvertibleTo(expectedType) {
		expectedValue := reflect.ValueOf(expected)
		if swag.IsZero(expectedValue) && swag.IsZero(reflect.ValueOf(actual)) {
			return ""
		}

		// Attempt comparison after type conversion
		if reflect.DeepEqual(actual, expectedValue.Convert(actualType).Interface()) {
			return ""
		}
	}

	// Last ditch effort
	if fmt.Sprintf("%#v", expected) == fmt.Sprintf("%#v", actual) {
		return ""
	}
	errFmt := "Expected: '%T(%#v)'\nActual:   '%T(%#v)'\n(Should be equivalent)!"
	return fmt.Sprintf(errFmt, expected, expected, actual, actual)

}

// assertSpecMaps is currently unused
func assertSpecMaps(t testing.TB, actual, expected map[string]interface{}) bool {
	res := true
	if id, ok := expected["id"]; ok {
		res = assert.Equal(t, id, actual["id"])
	}
	res = res && assert.Equal(t, expected["consumes"], actual["consumes"])
	res = res && assert.Equal(t, expected["produces"], actual["produces"])
	res = res && assert.Equal(t, expected["schemes"], actual["schemes"])
	res = res && assert.Equal(t, expected["swagger"], actual["swagger"])
	res = res && assert.Equal(t, expected["info"], actual["info"])
	res = res && assert.Equal(t, expected["host"], actual["host"])
	res = res && assert.Equal(t, expected["basePath"], actual["basePath"])
	res = res && assert.Equal(t, expected["paths"], actual["paths"])
	res = res && assert.Equal(t, expected["definitions"], actual["definitions"])
	res = res && assert.Equal(t, expected["responses"], actual["responses"])
	res = res && assert.Equal(t, expected["securityDefinitions"], actual["securityDefinitions"])
	res = res && assert.Equal(t, expected["tags"], actual["tags"])
	res = res && assert.Equal(t, expected["externalDocs"], actual["externalDocs"])
	res = res && assert.Equal(t, expected["x-some-extension"], actual["x-some-extension"])
	res = res && assert.Equal(t, expected["x-schemes"], actual["x-schemes"])

	return res
}
*/
func assertSpecs(t testing.TB, actual, expected spec.Swagger) bool {
	expected.Swagger = "2.0"
	return assert.Equal(t, actual, expected)
}

/*
// assertSpecJSON is currently unused
func assertSpecJSON(t testing.TB, speccJSON []byte) bool {
	var expected map[string]interface{}
	if !assert.NoError(t, json.Unmarshal(speccJSON, &expected)) {
		return false
	}

	obj := Swagger{}
	if !assert.NoError(t, json.Unmarshal(speccJSON, &obj)) {
		return false
	}

	cb, err := json.MarshalIndent(obj, "", "  ")
	if assert.NoError(t, err) {
		return false
	}
	var actual map[string]interface{}
	if !assert.NoError(t, json.Unmarshal(cb, &actual)) {
		return false
	}
	return assertSpecMaps(t, actual, expected)
}
*/

func TestSwaggerSpec_Serialize(t *testing.T) {
	expected := make(map[string]interface{})
	_ = json.Unmarshal([]byte(speccJSON), &expected)
	b, err := json.MarshalIndent(specc, "", "  ")
	if assert.NoError(t, err) {
		var actual map[string]interface{}
		err := json.Unmarshal(b, &actual)
		if assert.NoError(t, err) {
			assert.EqualValues(t, actual, expected)
		}
	}
}

func TestSwaggerSpec_Deserialize(t *testing.T) {
	var actual spec.Swagger
	err := json.Unmarshal([]byte(speccJSON), &actual)
	if assert.NoError(t, err) {
		assert.EqualValues(t, actual, specc)
	}
}
