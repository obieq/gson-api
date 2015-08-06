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
	Year   *int    `json:"year,omitempty"`
	Make   *string `json:"make,omitempty"`
	Active *bool   `json:"active,omitempty"`
}

// BuildLinks => builds JSON API links
func (r *AutomobileResource) BuildLinks() {
	r.Links = &AutomobileLinks{Link: Link{Self: LinkSelfInstance(r)}}
}

// NOTE: the code below is an example of how
//       to customize the URL value on a per resource basis
//func (r *AutomobileResource) URL() string {
//return "https://overridden-url.com/v5/"
//}

func (r *AutomobileResource) URI() string {
	return AUTOMOBILE_RESOURCE_TYPE
}

func (r *AutomobileResource) SelfLink() string {
	return r.Links.Self
}

func (r *AutomobileResource) Atts() interface{} {
	return &r.Attributes
}

func (r *AutomobileResource) SetAtts(atts interface{}) {
	r.Attributes = atts
}

// MapFromModel => maps a model to a resource
func (r *AutomobileResource) MapFromModel(model interface{}) {
	m := model.(AutomobileModel)
	attrs := AutomobileResourceAttributes{}

	if !m.HasErrors() {
		r.ResourceType = AUTOMOBILE_RESOURCE_TYPE
		r.ID = m.ID
		attrs.Year = &m.Year
		attrs.Make = &m.Make
		attrs.Active = &m.Active
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
		if v := attrs.Year; v != nil {
			m.Year = *v
		}
		if v := attrs.Make; v != nil {
			m.Make = *v
		}
		if v := attrs.Active; v != nil {
			m.Active = *v
		}
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

		y := 2010
		m := "Mazda"
		a := true
		attrs := AutomobileResourceAttributes{Year: &y, Make: &m, Active: &a}
		factory["Attributes"] = attrs
	})

	gory.Define("automobileResource2", AutomobileResource{}, func(factory gory.Factory) {
		factory["ID"] = "cccc-3333-dddd-4444"
		factory["ResourceType"] = "automobiles"

		y := 1960
		m := "Austin-Healey"
		a := true
		attrs := AutomobileResourceAttributes{Year: &y, Make: &m, Active: &a}
		factory["Attributes"] = attrs
	})

	gory.Define("automobileModel1", AutomobileModel{}, func(factory gory.Factory) {
		factory["ID"] = "bbbb-2222-eeee-5555"
		factory["Year"] = 1980
		factory["Make"] = "Honda"
		factory["Active"] = true
	})
})

// ******************* END TEST FACTORIES SECTION ********************* //
