package gsonapi

import "github.com/martini-contrib/render"

func HandleIndexResponse(err *JsonApiError, link Link, result interface{}, r render.Render) {
	if err == nil {
		r.JSON(200, map[string]interface{}{"links": link, "data": result}) // TODO: return links before data
	} else {
		r.JSON(404, map[string]interface{}{"errors": err})
	}
}

func HandleGetResponse(err *JsonApiError, result interface{}, r render.Render) {
	if err == nil {
		r.JSON(200, map[string]interface{}{"data": result})
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
