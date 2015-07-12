package gsonapi

import (
	"github.com/martini-contrib/render"
	goar "github.com/obieq/goar"
)

func HandleIndexResponse(resultError error, link Link, result interface{}, r render.Render) {
	if resultError == nil {
		r.JSON(200, map[string]interface{}{"links": link, "data": result}) // TODO: return links before data
	} else {
		r.JSON(404, map[string]interface{}{"errors": resultError})
	}
}

func HandleGetResponse(resultError error, result interface{}, r render.Render) {
	if resultError == nil {
		r.JSON(200, map[string]interface{}{"data": result})
	} else {
		r.JSON(404, map[string]interface{}{"errors": resultError})
	}
}

// HandlePostResponse => formats appropriate JSON response based on success vs. error
func HandlePostResponse(success bool, resultError error, resource JsonApiResourcer, r render.Render) {
	if success {
		// TODO: retrieve from the database instead of re-using instance
		r.Header().Set("Location", resource.SelfLink())
		r.JSON(201, map[string]interface{}{"data": resource})
	} else if resultError != nil {
		// TODO: how do I parse the status code?
		r.JSON(400, map[string]interface{}{"errors": resultError})
		//r.JSON(412, map[string]interface{}{"errors": err})
	} else {
		r.JSON(422, map[string]interface{}{"errors": resource.Errors()})
	}
}

// HandlePatchResponse => formats appropriate JSON response based on success vs. error
// NOTE: used by both the PUT and PATCH methods
//func HandlePutPatchResponse(success bool, resultError error, resource JsonApiResourcer, r render.Render) {
//if success {
//// TODO: retrieve from the database instead of re-using instance
//r.JSON(204, map[string]interface{}{})
//} else if resultError != nil {
//// TODO: how do I parse the status code?
//r.JSON(400, map[string]interface{}{"errors": resultError})
////r.JSON(412, map[string]interface{}{"errors": err})
//} else {
//r.JSON(422, map[string]interface{}{"errors": resource.Errors()})
//}
//}

func HandleDeleteResponse(model goar.ActiveRecordInterfacer, r render.Render) {
	goar.ToAR(model)

	if err := model.Delete(); err != nil {
		r.JSON(400, map[string]interface{}{"errors": err})
	} else {
		r.JSON(204, map[string]interface{}{})
	}
}
