package gsonapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var AUTOMOBILE_ID string = "aaaa-bbbb-cccc-dddd"

func MapHandleCreateAutomobileParams(server *martini.ClassicMartini, success bool, err error) {
	server.Map(success)
	server.Map(err)
}

func HandleCreateAutomobile(request JsonApiResource, r render.Render, success bool, err error) {
	var resource AutomobileResource

	if success {
		err = nil
		// map the resource to the model
		m := AutomobileModel{}
		err = UnmarshalJsonApiData(request.Data, &resource)
		resource.MapToModel(&m)

		// map the model to the resource
		if err == nil {
			// given that this is a create request, generate an id
			m.ID = AUTOMOBILE_ID

			resource = AutomobileResource{}
			resource.MapFromModel(m)
		}

		//success = err == nil
		HandlePostResponse(true, err, &resource, r)
	} else if err != nil {
		HandlePostResponse(false, err, &resource, r)
	} else { // success = false
		HandlePostResponse(success, nil, &resource, r)
	}
}

func MarshalAutomobileResource(auto AutomobileResource) []byte {
	// Set up a new POST request before every test
	auto.Attributes = AutomobileResourceAttributes{Year: 2010, Make: "Acura"}

	j := JsonApiResource{Data: auto}

	body, err := json.Marshal(j)
	Ω(err).NotTo(HaveOccurred())

	return body
}

func BuildPostRoute(server *martini.ClassicMartini) {
	server.Group("/v1", func(r martini.Router) {
		r.Post("/automobiles", binding.Json(JsonApiResource{}), HandleCreateAutomobile)
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

	//It("should unmarshal a resource to a model", func() {
	//var m AutomobileModel
	//r := automobileResource1
	//j, err := json.Marshal(r)
	//Ω(err).NotTo(HaveOccurred())

	//log.Println(r)
	//err = UnmarshalJsonApiData(j, &m)
	//Ω(err).NotTo(HaveOccurred())
	//Ω(m).ShouldNot(BeNil())
	//log.Println(r)
	//log.Println(m.Year)
	//log.Println(m.Make)
	//Ω(m.Year).Should(Equal(r.Attributes.(AutomobileResourceAttributes).Year))
	//Ω(m.Make).Should(Equal(r.Attributes.(AutomobileResourceAttributes).Make))
	//})

	Context("HTTP POST", func() {
		var (
			auto1 *AutomobileResource = &AutomobileResource{}
		)

		BeforeEach(func() {
			//success := true
			//server.Map(success)
		})

		It("should return a 201 Status Code", func() {
			MapHandleCreateAutomobileParams(server, true, errors.New(""))
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
			Ω(recorder.Body.String()).Should(Equal(responseBody))
		})

		It("should return a 400 Status Code", func() {
			MapHandleCreateAutomobileParams(server, false, errors.New("oops"))
			BuildPostRoute(server)

			// prepare request
			body := MarshalAutomobileResource(*auto1)
			request, _ = http.NewRequest("POST", "/v1/automobiles", bytes.NewReader(body))

			// send request to server
			server.ServeHTTP(recorder, request)

			// verify
			Ω(recorder.Code).Should(Equal(400))
			responseBody := `{"errors":{}}`
			Ω(recorder.Body.String()).Should(Equal(responseBody))
		})
	}) // Context "HTTP POST"
})
