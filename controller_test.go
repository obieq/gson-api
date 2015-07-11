package gsonapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var AUTOMOBILE_ID string = "aaaa-bbbb-cccc-dddd"

func HandleCreateAutomobile(request JsonApiResource, r render.Render) {
	var resource AutomobileResource

	// map the resource to the model
	m := AutomobileModel{}
	err := UnmarshalJsonApiData(request.Data, &resource)
	resource.MapToModel(&m)

	// map the model to the resource
	if err == nil {
		// given that this is a create request, generate an id
		m.ID = AUTOMOBILE_ID

		resource = AutomobileResource{}
		resource.MapFromModel(m)
	}

	success := err == nil
	HandlePostResponse(success, err, &resource, r)
}

func MarshalAutomobileResource(auto AutomobileResource) []byte {
	// Set up a new POST request before every test
	auto.Attributes = AutomobileResourceAttributes{Year: 2010, Make: "Acura"}

	j := JsonApiResource{Data: auto}

	body, err := json.Marshal(j)
	Ω(err).NotTo(HaveOccurred())

	return body
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
			server.Group("/v1", func(r martini.Router) {
				r.Post("/automobiles", binding.Json(JsonApiResource{}), HandleCreateAutomobile)
			})
		})

		It("should return a 201 Status Code", func() {
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
	}) // Context "HTTP POST"
})
