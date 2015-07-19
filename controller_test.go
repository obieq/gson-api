package gsonapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"github.com/modocache/gory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var AUTOMOBILE_ID = "aaaa-bbbb-cccc-dddd"
var CARZ_URL = "https://carz.com/v1"

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

	automobiles := make([]AutomobileResource, 2)
	automobiles[0] = *gory.Build("automobileResource1").(*AutomobileResource)
	automobiles[1] = *gory.Build("automobileResource2").(*AutomobileResource)

	// build links
	automobiles[0].BuildLinks()
	automobiles[1].BuildLinks()

	resource := AutomobileResource{URL: CARZ_URL}
	HandleIndexResponse(jsonApiError, Link{Self: resource.LinkSelfCollection()}, automobiles, r)
}

func HandleGetAutomobile(r render.Render, err error) {
	var jsonApiError *JsonApiError

	if err.Error() != "" {
		jsonApiError = &JsonApiError{Status: "404", Detail: err.Error()}
	}

	auto := *gory.Build("automobileResource1").(*AutomobileResource)

	// build links
	auto.BuildLinks()

	HandleGetResponse(jsonApiError, auto, r)
}

func HandleCreateAutomobile(request JsonApiResource, r render.Render, success bool, err error) {
	var resource AutomobileResource
	var jsonApiError *JsonApiError

	if success {
		// map the resource to the model
		m := AutomobileModel{}
		UnmarshalJsonApiData(request.Data, &resource)
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

	HandlePostResponse(success, jsonApiError, &resource, r)
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
	auto.Attributes = AutomobileResourceAttributes{Year: 2010, Make: "Acura"}

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
		r.Post("/automobiles", binding.Json(JsonApiResource{}), HandleCreateAutomobile)
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
			expectedResponse := `{` +
				`"data":[{"type":"automobiles","id":"aaaa-1111-bbbb-2222","attributes":{"year":2010,"make":"Mazda"},` +
				`"links":{"self":"https://carz.com/v1/automobiles/aaaa-1111-bbbb-2222"}},` +
				`{"type":"automobiles","id":"cccc-3333-dddd-4444","attributes":{"year":1960,"make":"Austin-Healey"},` +
				`"links":{"self":"https://carz.com/v1/automobiles/cccc-3333-dddd-4444"}}],` +
				`"links":{"self":"https://carz.com/v1/automobiles"}}`
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
			expectedResponse := `{` +
				`"data":{"type":"automobiles","id":"aaaa-1111-bbbb-2222","attributes":{"year":2010,"make":"Mazda"},` +
				`"links":{"self":"https://carz.com/v1/automobiles/aaaa-1111-bbbb-2222"}}}`
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
		var (
			auto1 *AutomobileResource = &AutomobileResource{}
		)

		It("should return a 201 Status Code", func() {
			MapErrorParam(server, errors.New(""))
			MapSuccessParam(server, true)
			BuildPostRoute(server)

			// prepare request
			body := MarshalAutomobileResource(*auto1)
			request, _ = http.NewRequest("POST", "/v1/automobiles", bytes.NewReader(body))

			// send request to server
			server.ServeHTTP(recorder, request)

			// verify
			Ω(recorder.Code).Should(Equal(201))
			responseBody :=
				`{` +
					`"data":{"type":"automobiles","id":"` + AUTOMOBILE_ID + `",` +
					`"attributes":{"year":2010,"make":"Acura"},` +
					`"links":{"self":"https://carz.com/v1/automobiles/` + AUTOMOBILE_ID + `"}}}`
			Ω(recorder.Body.String()).Should(MatchJSON(responseBody))
		})

		It("should return a 400 Status Code", func() {
			MapErrorParam(server, errors.New("oops"))
			MapSuccessParam(server, false)
			BuildPostRoute(server)

			// prepare request
			body := MarshalAutomobileResource(*auto1)
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
			body := MarshalAutomobileResource(*auto1)
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
