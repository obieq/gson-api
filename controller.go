package gsonapi

import (
	"encoding/json"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/martini-contrib/render"
)

// JSONApiServerInfo => contains necessary info for building an api's route
type JSONApiServerInfo struct {
	BaseURL string
	Prefix  string
}

// GetBaseURL => api routes base url
// EX: https://test.myapi.com/.....
func (jasi JSONApiServerInfo) GetBaseURL() string {
	return jasi.BaseURL
}

// GetPrefix => api route's prefix
// EX: https://xxxxx.com/v1/xxxxx (v1 is the prefix)
func (jasi JSONApiServerInfo) GetPrefix() string {
	return jasi.Prefix
}

func HandleIndexResponse(jasi JSONApiServerInfo, err *JsonApiError, result interface{}, r render.Render) {
	var j []byte
	var jsonError error
	var response interface{}

	if err == nil {
		// r.JSON(200, map[string]interface{}{"links": link, "data": result}) // TODO: return links before data
		if j, jsonError = jsonapi.MarshalToJSONWithURLs(result, jasi); jsonError != nil {
			r.JSON(400, map[string]interface{}{"errors": jsonError})
		}
		if jsonError = json.Unmarshal(j, &response); jsonError != nil {
			r.JSON(400, map[string]interface{}{"errors": jsonError})
		}
		r.JSON(200, response)
	} else {
		r.JSON(404, map[string]interface{}{"errors": err})
	}
}

func HandleGetResponse(jasi JSONApiServerInfo, err *JsonApiError, result interface{}, r render.Render) {
	var j []byte
	var jsonError error
	var response interface{}

	if err == nil {
		if j, jsonError = jsonapi.MarshalToJSONWithURLs(result, jasi); jsonError != nil {
			r.JSON(400, map[string]interface{}{"errors": jsonError})
		}
		if jsonError = json.Unmarshal(j, &response); jsonError != nil {
			r.JSON(400, map[string]interface{}{"errors": jsonError})
		}
		r.JSON(200, response)
	} else {
		r.JSON(404, map[string]interface{}{"errors": err})
	}
}

// HandlePostResponse => formats appropriate JSON response based on success vs. error
func HandlePostResponse(success bool, err *JsonApiError, resource JsonApiResourcer, r render.Render) {
	// TODO: return 404 if resource not found
	if success {
		// TODO: retrieve from the database instead of re-using instance
		r.Header().Set("Location", LinkSelfInstance(resource))
		r.JSON(201, map[string]interface{}{"data": resource})
	} else if err != nil {
		// TODO: how do I parse the status code?
		r.JSON(400, map[string]interface{}{"errors": err})
	} else {
		r.JSON(422, map[string]interface{}{"errors": resource.Errors()})
	}
}

// HandlePatchResponse => formats appropriate JSON response based on success vs. error
func HandlePatchResponse(success bool, err *JsonApiError, resource JsonApiResourcer, r render.Render) {
	if success {
		// TODO: retrieve from the database instead of re-using instance
		r.Header().Set("Location", LinkSelfInstance(resource))
		r.JSON(200, map[string]interface{}{"data": resource}) // given that updated-at is set, a 200 w/ content must be returned
	} else if err != nil {
		// TODO: how do I parse the status code?
		r.JSON(400, map[string]interface{}{"errors": err})
		//r.JSON(412, map[string]interface{}{"errors": err})
	} else {
		r.JSON(422, map[string]interface{}{"errors": resource.Errors()})
	}
}

// HandlePatchResponse => formats appropriate JSON response based on success vs. error
// NOTE: used by both the PUT and PATCH methods
//func HandlePutPatchResponse(success bool, err error, resource JsonApiResourcer, r render.Render) {
//if success {
//// TODO: retrieve from the database instead of re-using instance
//r.JSON(204, map[string]interface{}{})
//} else if err != nil {
//// TODO: how do I parse the status code?
//r.JSON(400, map[string]interface{}{"errors": err})
////r.JSON(412, map[string]interface{}{"errors": err})
//} else {
//r.JSON(422, map[string]interface{}{"errors": resource.Errors()})
//}
//}

func HandleDeleteResponse(err *JsonApiError, r render.Render) {
	if err == nil {
		r.JSON(204, map[string]interface{}{})
	} else {
		r.JSON(400, map[string]interface{}{"errors": err})
	}
}
