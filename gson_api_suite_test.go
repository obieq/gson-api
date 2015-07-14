package gsonapi

import (
	gory "github.com/modocache/gory"
	validations "github.com/obieq/goar-validations"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGsonApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GsonApi Suite")
}

// ******************* BEGIN MARTINI SERVER SECTION ******************* //
// Wrap the Martini server struct.
//type Server *martini.ClassicMartini
// ******************* END MARTINI SERVER SECTION ********************* //

// ******************* BEGIN MODEL SECTION **************************** //
type AutomobileModel struct {
	validations.Validation
	ID     string `json:"id"`
	Year   int    `json:"year,omitempty"`
	Make   string `json:"make,omitempty"`
	Active bool   `json:"active,omitempty"`
}

// ******************* END MODEL SECTION ****************************** //

// ******************* BEGIN RESOURCE SECTION ************************* //
const AUTOMOBILE_RESOURCE_TYPE string = "automobiles"

// AutomobileLinks => JSON API links
type AutomobileLinks struct {
	Link
}

// Automobile Resource
type AutomobileResource struct {
	Resource
	Links *AutomobileLinks `json:"links,omitempty"`
}

type AutomobileResourceAttributes struct {
	Year   int    `json:"year,omitempty"`
	Make   string `json:"make,omitempty"`
	Active bool   `json:"active,omitempty"`
}

// BuildLinks => builds JSON API links
func (r *AutomobileResource) BuildLinks() {
	root := "https://carz.com" + "/v1/" + AUTOMOBILE_RESOURCE_TYPE + "/" + r.ID
	r.Links = &AutomobileLinks{Link: Link{Self: root}}
}

func (r *AutomobileResource) SelfLink() string {
	return r.Links.Self
}

// MapFromModel => maps a model to a resource
func (r *AutomobileResource) MapFromModel(model interface{}) {
	m := model.(AutomobileModel)
	attrs := AutomobileResourceAttributes{}

	if !m.HasErrors() {
		r.ResourceType = AUTOMOBILE_RESOURCE_TYPE
		r.ID = m.ID
		attrs.Year = m.Year
		attrs.Make = m.Make
		attrs.Active = m.Active
		r.Attributes = attrs

		// build links
		r.BuildLinks()
	} else {
		r.SetErrors(m.ErrorMap())
	}
}

// MapToModel => maps a resource to a model
func (r *AutomobileResource) MapToModel(model interface{}) error {
	var attrs AutomobileResourceAttributes
	var err error
	m := model.(*AutomobileModel)

	if err = UnmarshalJsonApiData(r.Attributes, &attrs); err == nil {
		m.Year = attrs.Year
		m.Make = attrs.Make
		m.Active = attrs.Active
	}

	return err
}

// ******************* END RESOURCE SECTION *************************** //

// ******************* BEGIN TEST HELPERS SECTION ********************* //
func BuildErrors() map[string]*validations.ValidationError {
	errors := map[string]*validations.ValidationError{}
	veYear := validations.ValidationError{Key: "year", Message: "cannot be greater than 2016"}
	veMake := validations.ValidationError{Key: "make", Message: "cannot be blank"}
	errors[veYear.Key] = &veYear
	errors[veMake.Key] = &veMake
	return errors
}

// ******************* END TEST HELPERS SECTION *********************** //

// ******************* BEGIN TEST FACTORIES SECTION ******************* //
var _ = BeforeSuite(func() {
	gory.Define("automobileResource1", AutomobileResource{}, func(factory gory.Factory) {
		factory["ID"] = "aaaa-1111-bbbb-2222"
		factory["ResourceType"] = "automobiles"

		attrs := AutomobileResourceAttributes{Year: 2010, Make: "Mazda"}
		factory["Attributes"] = attrs
	})

	gory.Define("automobileResource2", AutomobileResource{}, func(factory gory.Factory) {
		factory["ID"] = "cccc-3333-dddd-4444"
		factory["ResourceType"] = "automobiles"

		attrs := AutomobileResourceAttributes{Year: 1960, Make: "Austin-Healey"}
		factory["Attributes"] = attrs
	})
})

// ******************* END TEST FACTORIES SECTION ********************* //
