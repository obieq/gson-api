package gsonapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/go-martini/martini"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/martini-contrib/render"
	"github.com/modocache/gory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const EXPECTED_CONTENT_TYPE = "application/vnd.api+json; charset=UTF-8"

var TEST_SERVER_INFO = JSONApiServerInfo{BaseURL: "http://my.domain", Prefix: "v1"}
var AUTOMOBILE_ID = "aaaa-bbbb-cccc-dddd"

var autoModel1 AutomobileModel
var autoResource1, autoResource2, autoResource3 AutomobileResource

func MapErrorParam(server *martini.ClassicMartini, err error) {
	server.Map(err)
}

func MapSuccessParam(server *martini.ClassicMartini, success bool) {
	server.Map(success)
}

func HandleGetAutomobiles(r render.Render, err error) {
	var jsonApiError *JsonApiError

	if err.Error() != "" {
		jsonApiError = &JsonApiError{Status: "404", Detail: err.Error()}
	}

	automobiles := []AutomobileResource{autoResource1, autoResource2, autoResource3}
	HandleIndexResponse(TEST_SERVER_INFO, jsonApiError, automobiles, r)
}

func HandleGetAutomobile(r render.Render, err error) {
	var jsonApiError *JsonApiError

	if err.Error() != "" {
		jsonApiError = &JsonApiError{Status: "404", Detail: err.Error()}
	}

	auto := autoResource1
	HandleGetResponse(TEST_SERVER_INFO, jsonApiError, auto, r)
}

// func HandleCreateAutomobile(request JsonApiPayload, r render.Render, success bool, err error) {
func HandleCreateAutomobile(request *http.Request, r render.Render, success bool, err error) {
	var resource AutomobileResource
	var jsonApiError *JsonApiError

	if success {
		// map the resource to the model
		m := AutomobileModel{}
		log.Println("**************** IN HANDLE CREATE AUTOMOBILE ****************")
		log.Println(request)
		defer request.Body.Close()

		body, _ := ioutil.ReadAll(request.Body)
		log.Println(body)
		log.Println(string(body))
		err = jsonapi.UnmarshalFromJSON(body, &resource)
		log.Println(err)
		log.Println(resource)

		resource.MapToModel(&m)

		// map the model to the resource
		m.ID = AUTOMOBILE_ID // given that this is a create request, generate an id
		resource = AutomobileResource{}
		resource.MapFromModel(m)
	} else if err != nil && err.Error() != "" {
		jsonApiError = &JsonApiError{Status: "400", Detail: err.Error()}
	} else { // success = false, i.e., business rule validation errors
		resource.SetErrors(BuildErrors())
	}

	HandlePostResponse(TEST_SERVER_INFO, success, jsonApiError, &resource, r)
}

func HandlePatchAutomobile(args martini.Params, request *http.Request, r render.Render, success bool, err error) {
	var resource AutomobileResource
	var jsonApiError *JsonApiError

	if success {
		// map the resource to the model
		// NOTE: should perform a partial update
		log.Println("**************** IN HANDLE PATCH AUTOMOBILE ****************")
		log.Println(request)
		defer request.Body.Close()

		body, _ := ioutil.ReadAll(request.Body)
		log.Println(body)
		log.Println(string(body))
		err = jsonapi.UnmarshalFromJSON(body, &resource)
		log.Println(err)
		log.Println(resource)

		rID := resource.GetID()
		m := *gory.Build("automobileModel" + rID[len(rID)-1:]).(*AutomobileModel)
		resource.MapToModel(&m)

		// map the model to the resource
		resource = AutomobileResource{}
		resource.MapFromModel(m)
	} else if err != nil && err.Error() != "" {
		jsonApiError = &JsonApiError{Status: "400", Detail: err.Error()}
	} else { // success = false, i.e., business rule validation errors
		resource.SetErrors(BuildErrors())
	}

	HandlePatchResponse(TEST_SERVER_INFO, success, jsonApiError, &resource, r)
}

func HandleDeleteAutomobile(r render.Render, err error) {
	var jsonApiError *JsonApiError

	if err.Error() != "" {
		jsonApiError = &JsonApiError{Status: "400", Detail: err.Error()}
	}

	HandleDeleteResponse(jsonApiError, r)
}

func MarshalAutomobileResource(auto AutomobileResource) []byte {
	// Set up a new POST request before every test
	j := JsonApiResource{Data: auto}

	body, err := json.Marshal(j)
	Ω(err).NotTo(HaveOccurred())

	return body
}

func BuildGetListRoute(server *martini.ClassicMartini) {
	server.Group("/v1", func(r martini.Router) {
		r.Get("/automobiles", HandleGetAutomobiles)
	})
}

func BuildGetSingleRoute(server *martini.ClassicMartini) {
	server.Group("/v1", func(r martini.Router) {
		r.Get("/automobiles/:id", HandleGetAutomobile)
	})
}

func BuildPostRoute(server *martini.ClassicMartini) {
	server.Group("/v1", func(r martini.Router) {
		r.Post("/automobiles", HandleCreateAutomobile)
	})
}

func BuildPatchRoute(server *martini.ClassicMartini) {
	server.Group("/v1", func(r martini.Router) {
		r.Patch("/automobiles", HandlePatchAutomobile)
	})
}

func BuildDeleteRoute(server *martini.ClassicMartini) {
	server.Group("/v1", func(r martini.Router) {
		r.Delete("/automobiles/:id", HandleDeleteAutomobile)
	})
}

var _ = Describe("Controller", func() {
	var (
		server   *martini.ClassicMartini
		request  *http.Request
		recorder *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		// Configure Martini
		server = martini.Classic()
		server.Use(render.Renderer())

		// Record HTTP responses
		recorder = httptest.NewRecorder()

		// reset global vars
		autoModel1 = *gory.Build("automobileModel1").(*AutomobileModel)
		autoResource1 = *gory.Build("automobileResource1").(*AutomobileResource)
		autoResource2 = *gory.Build("automobileResource2").(*AutomobileResource)
		autoResource3 = *gory.Build("automobileResource3").(*AutomobileResource)
	})

	Context("HTTP GET (List)", func() {
		BeforeEach(func() {
			server.Group("/v1", func(r martini.Router) {
				r.Get("/automobiles", HandleGetAutomobiles)
			})
		})

		It("should return a 200 Status Code", func() {
			MapErrorParam(server, errors.New(""))
			BuildGetListRoute(server)

			request, _ = http.NewRequest("GET", "/v1/automobiles", nil)

			// send request to server
			server.ServeHTTP(recorder, request)

			// verify
			Ω(recorder.Code).Should(Equal(200))
			// minified via http://www.webtoolkitonline.com/json-minifier.html
			expectedResponse := `{"data":[{"attributes":{"active":true,"ages":null,"body-style":"4 door sedan","inspections":[{"location":"216 broad ave, richmond va 23226","name":"inspection #1"},{"location":"2201 stoddard ct, arlington va 22202","name":"inspection #2"}],"make":"Mazda","year":2010},"id":"aaaa-1111-bbbb-2222","relationships":{"drivers":{"data":[{"id":"driver-id-1","type":"drivers"},{"id":"driver-id-2","type":"drivers"}],"links":{"related":"http://my.domain/v1/automobiles/aaaa-1111-bbbb-2222/drivers","self":"http://my.domain/v1/automobiles/aaaa-1111-bbbb-2222/relationships/drivers"}}},"type":"automobiles"},{"attributes":{"active":true,"ages":null,"body-style":null,"inspections":[{"location":"216 broad ave, richmond va 23226","name":"inspection #1"}],"make":"Austin-Healey","year":1960},"id":"cccc-3333-dddd-4444","relationships":{"drivers":{"data":[],"links":{"related":"http://my.domain/v1/automobiles/cccc-3333-dddd-4444/drivers","self":"http://my.domain/v1/automobiles/cccc-3333-dddd-4444/relationships/drivers"}}},"type":"automobiles"},{"attributes":{"active":false,"ages":null,"body-style":null,"inspections":null,"make":"Honda","year":1980},"id":"bbbb-2222-eeee-5555","relationships":{"drivers":{"data":[{"id":"driver-id-1","type":"drivers"}],"links":{"related":"http://my.domain/v1/automobiles/bbbb-2222-eeee-5555/drivers","self":"http://my.domain/v1/automobiles/bbbb-2222-eeee-5555/relationships/drivers"}}},"type":"automobiles"}],"included":[{"attributes":{"active":true,"age":40,"name":"paul walker"},"id":"driver-id-1","type":"drivers"},{"attributes":{"active":false,"age":45,"name":"steve mcqueen"},"id":"driver-id-2","type":"drivers"}]}`
			Ω(recorder.Body.String()).Should(MatchJSON(expectedResponse))
		})

		It("should return a 404 Status Code", func() {
			MapErrorParam(server, errors.New("not found"))
			BuildGetListRoute(server)

			request, _ = http.NewRequest("GET", "/v1/automobiles", nil)

			// send request to server
			server.ServeHTTP(recorder, request)

			// verify
			Ω(recorder.Code).Should(Equal(404))
			expectedResponse := `{"errors":{"status":"404","detail":"not found"}}`
			Ω(recorder.Body.String()).Should(MatchJSON(expectedResponse))
		})

	})

	Context("HTTP GET (Single)", func() {
		It("should return a 200 Status Code", func() {
			MapErrorParam(server, errors.New(""))
			BuildGetSingleRoute(server)

			request, _ = http.NewRequest("GET", "/v1/automobiles/aaaa-1111-bbbb-2222", nil)

			// send request to server
			server.ServeHTTP(recorder, request)

			// verify
			Ω(recorder.Code).Should(Equal(200))
			Ω(recorder.Header().Get("Content-Type")).Should(Equal(EXPECTED_CONTENT_TYPE))
			// minified via http://www.webtoolkitonline.com/json-minifier.html
			expectedResponse := `{"data":{"attributes":{"active":true,"ages":null,"body-style":"4 door sedan","inspections":[{"location":"216 broad ave, richmond va 23226","name":"inspection #1"},{"location":"2201 stoddard ct, arlington va 22202","name":"inspection #2"}],"make":"Mazda","year":2010},"id":"aaaa-1111-bbbb-2222","relationships":{"drivers":{"data":[{"id":"driver-id-1","type":"drivers"},{"id":"driver-id-2","type":"drivers"}],"links":{"related":"http://my.domain/v1/automobiles/aaaa-1111-bbbb-2222/drivers","self":"http://my.domain/v1/automobiles/aaaa-1111-bbbb-2222/relationships/drivers"}}},"type":"automobiles"},"included":[{"attributes":{"active":true,"age":40,"name":"paul walker"},"id":"driver-id-1","type":"drivers"},{"attributes":{"active":false,"age":45,"name":"steve mcqueen"},"id":"driver-id-2","type":"drivers"}]}`
			Ω(recorder.Header().Get("Content-Type")).Should(Equal(EXPECTED_CONTENT_TYPE))
			Ω(recorder.Body.String()).Should(MatchJSON(expectedResponse))
		})

		It("should return a 404 Status Code", func() {
			MapErrorParam(server, errors.New("not found"))
			BuildGetSingleRoute(server)

			request, _ = http.NewRequest("GET", "/v1/automobiles/aaaa-1111-bbbb-2222", nil)

			// send request to server
			server.ServeHTTP(recorder, request)

			// verify
			Ω(recorder.Code).Should(Equal(404))
			expectedResponse := `{"errors":{"status":"404","detail":"not found"}}`
			Ω(recorder.Body.String()).Should(MatchJSON(expectedResponse))
		})
	})

	Context("HTTP POST", func() {
		It("should return a 201 Status Code", func() {
			MapErrorParam(server, errors.New(""))
			MapSuccessParam(server, true)
			BuildPostRoute(server)

			body := []byte(`{` +
				`"data":{"type":"automobiles",` +
				`"attributes":{"body-style":"4 door sedan","year":2010,"make":"Mazda","active":true,"ages":[4,6],` +
				`"inspections":[{"location":"216 broad ave, richmond va 23226","name":"inspection #1"},{"location":"2201 stoddard ct, arlington va 22202","name":"inspection #2"}]},` +
				`"relationships":{"drivers":{"data":[{ "type": "drivers", "id": "driver-id-1" },{ "type": "drivers", "id": "driver-id-2" }]}}}}`)

			request, _ = http.NewRequest("POST", "/v1/automobiles", bytes.NewReader(body))

			// send request to server
			server.ServeHTTP(recorder, request)

			// verify
			Ω(recorder.Code).Should(Equal(201))
			// minified via http://www.webtoolkitonline.com/json-minifier.html
			responseBody := `{"data":{"attributes":{"active":true,"ages":[4,6],"body-style":"4 door sedan","inspections":[{"location":"216 broad ave, richmond va 23226","name":"inspection #1"},{"location":"2201 stoddard ct, arlington va 22202","name":"inspection #2"}],"make":"Mazda","year":2010},"id":"aaaa-bbbb-cccc-dddd","relationships":{"drivers":{"data":[{"id":"driver-id-1","type":"drivers"},{"id":"driver-id-2","type":"drivers"}],"links":{"related":"http://my.domain/v1/automobiles/aaaa-bbbb-cccc-dddd/drivers","self":"http://my.domain/v1/automobiles/aaaa-bbbb-cccc-dddd/relationships/drivers"}}},"type":"automobiles"},"included":[{"attributes":{"active":true,"age":40,"name":"paul walker"},"id":"driver-id-1","type":"drivers"},{"attributes":{"active":false,"age":45,"name":"steve mcqueen"},"id":"driver-id-2","type":"drivers"}]}`
			Ω(recorder.Body.String()).Should(MatchJSON(responseBody))
		})

		It("should return a 400 Status Code", func() {
			MapErrorParam(server, errors.New("oops"))
			MapSuccessParam(server, false)
			BuildPostRoute(server)

			// prepare request
			body := MarshalAutomobileResource(autoResource1)
			request, _ = http.NewRequest("POST", "/v1/automobiles", bytes.NewReader(body))

			// send request to server
			server.ServeHTTP(recorder, request)

			// verify
			Ω(recorder.Code).Should(Equal(400))
			responseBody := `{"errors":{"status":"400","detail":"oops"}}`
			Ω(recorder.Body.String()).Should(Equal(responseBody))
		})

		It("should return a 422 Status Code", func() {
			MapErrorParam(server, errors.New(""))
			MapSuccessParam(server, false)
			BuildPostRoute(server)

			// prepare request
			body := MarshalAutomobileResource(autoResource1)
			request, _ = http.NewRequest("POST", "/v1/automobiles", bytes.NewReader(body))

			// send request to server
			server.ServeHTTP(recorder, request)

			// verify
			Ω(recorder.Code).Should(Equal(422))

			// NOTE: cannot perform deep equal on errors array, so have to take an alternate approach
			responseBody1 := `{` +
				`"errors":[{"status":"422","detail":"cannot be blank","source":{"pointer":"data/attributes/make"}},` +
				`{"status":"422","detail":"cannot be greater than 2016","source":{"pointer":"data/attributes/year"}}]}`
			responseBody2 := `{` +
				`"errors":[{"status":"422","detail":"cannot be greater than 2016","source":{"pointer":"data/attributes/year"}},` +
				`{"status":"422","detail":"cannot be blank","source":{"pointer":"data/attributes/make"}}]}`

			Ω([]string{responseBody1, responseBody2}).Should(ContainElement(recorder.Body.String()))
		})
	}) // Context "HTTP POST"

	Context("HTTP PATCH", func() {
		It("should return a 200 Status Code", func() {
			MapErrorParam(server, errors.New(""))
			MapSuccessParam(server, true)
			BuildPatchRoute(server)

			// verify model values are what we expect
			Ω(autoModel1.Year).Should(Equal(1980))                 // should update
			Ω(autoModel1.Make).Should(Equal("Honda"))              // should remain unchanged
			Ω(*autoModel1.BodyStyle).Should(Equal("4 door sedan")) // should update
			Ω(autoModel1.Active).Should(Equal(true))               // should update

			// prepare request
			// update the following properties: Year, BodyStyle, and Active
			body := []byte(`{` +
				`"data":{"type":"automobiles","id":"` + autoModel1.ID + `",` +
				`"attributes":{"body-style":"2 door coupe","year":2010,"active":false}}}`)

			request, _ = http.NewRequest("PATCH", "/v1/automobiles", bytes.NewReader(body))

			// send request to server
			server.ServeHTTP(recorder, request)

			// verify response
			Ω(recorder.Code).Should(Equal(200))
			responseBody :=
				`{
          "data": {
            "attributes": {
              "active": false,
              "ages": [],
              "body-style": "2 door coupe",
              "inspections": null,
              "make": "Honda",
              "year": 2010
            },
            "id": "automobile-model-1",
            "relationships": {
              "drivers": {
                "data": [],
                "links": {
                  "related": "http://my.domain/v1/automobiles/automobile-model-1/drivers",
                  "self": "http://my.domain/v1/automobiles/automobile-model-1/relationships/drivers"
                }
              }
            },
            "type": "automobiles"
          }
        }`
			Ω(recorder.Body.String()).Should(MatchJSON(responseBody))
		})

		It("should return a 400 Status Code", func() {
			MapErrorParam(server, errors.New("oops"))
			MapSuccessParam(server, false)
			BuildPatchRoute(server)

			// prepare request
			body := MarshalAutomobileResource(autoResource1)
			request, _ = http.NewRequest("PATCH", "/v1/automobiles", bytes.NewReader(body))

			// send request to server
			server.ServeHTTP(recorder, request)

			// verify
			Ω(recorder.Code).Should(Equal(400))
			responseBody := `{"errors":{"status":"400","detail":"oops"}}`
			Ω(recorder.Body.String()).Should(Equal(responseBody))
		})

		It("should return a 422 Status Code", func() {
			MapErrorParam(server, errors.New(""))
			MapSuccessParam(server, false)
			BuildPatchRoute(server)

			// prepare request
			body := MarshalAutomobileResource(autoResource1)
			request, _ = http.NewRequest("PATCH", "/v1/automobiles", bytes.NewReader(body))

			// send request to server
			server.ServeHTTP(recorder, request)

			// verify
			Ω(recorder.Code).Should(Equal(422))

			// NOTE: cannot perform deep equal on errors array, so have to take an alternate approach
			responseBody1 := `{` +
				`"errors":[{"status":"422","detail":"cannot be blank","source":{"pointer":"data/attributes/make"}},` +
				`{"status":"422","detail":"cannot be greater than 2016","source":{"pointer":"data/attributes/year"}}]}`
			responseBody2 := `{` +
				`"errors":[{"status":"422","detail":"cannot be greater than 2016","source":{"pointer":"data/attributes/year"}},` +
				`{"status":"422","detail":"cannot be blank","source":{"pointer":"data/attributes/make"}}]}`

			Ω([]string{responseBody1, responseBody2}).Should(ContainElement(recorder.Body.String()))
		})
	}) // Context "HTTP PATCH"

	Context("HTTP DELETE", func() {
		It("should return a 204 Status Code", func() {
			MapErrorParam(server, errors.New(""))
			BuildDeleteRoute(server)

			request, _ = http.NewRequest("DELETE", "/v1/automobiles/aaaa-1111-bbbb-2222", nil)

			// send request to server
			server.ServeHTTP(recorder, request)

			// verify
			Ω(recorder.Code).Should(Equal(204))
			expectedResponse := `{}`
			Ω(recorder.Body.String()).Should(MatchJSON(expectedResponse))
		})

		It("should return a 400 Status Code", func() {
			MapErrorParam(server, errors.New("oops"))
			BuildDeleteRoute(server)

			request, _ = http.NewRequest("DELETE", "/v1/automobiles/aaaa-1111-bbbb-2222", nil)

			// send request to server
			server.ServeHTTP(recorder, request)

			// verify
			Ω(recorder.Code).Should(Equal(400))
			log.Println(recorder.Body.String())
			expectedResponse := `{"errors":{"status":"400","detail":"oops"}}`
			Ω(recorder.Body.String()).Should(MatchJSON(expectedResponse))
		})
	}) // Context "HTTP DELETE"
})
