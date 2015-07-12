package gsonapi

import (
	validations "github.com/obieq/goar-validations"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

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
func (p *AutomobileResource) BuildLinks(automobileModel interface{}) {
	root := "https://carz.com" + "/v1/" + AUTOMOBILE_RESOURCE_TYPE + "/" + p.ID
	p.Links = &AutomobileLinks{Link: Link{Self: root}}
}

func (p *AutomobileResource) SelfLink() string {
	return p.Links.Self
}

func (p *AutomobileResource) MapFromModel(model interface{}) {
	m := model.(AutomobileModel)
	attrs := AutomobileResourceAttributes{}

	if !m.HasErrors() {
		p.ResourceType = AUTOMOBILE_RESOURCE_TYPE
		p.ID = m.ID
		attrs.Year = m.Year
		attrs.Make = m.Make
		attrs.Active = m.Active
		p.Attributes = attrs

		// build links
		p.BuildLinks(m)
	} else {
		p.SetErrors(m.ErrorMap())
	}
}

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
// ******************* END RESOURCE SECTION *************************** //
func BuildErrors() map[string]*validations.ValidationError {
	errors := map[string]*validations.ValidationError{}
	veYear := validations.ValidationError{Key: "year", Message: "cannot be greater than 2016"}
	veMake := validations.ValidationError{Key: "make", Message: "cannot be blank"}
	errors[veYear.Key] = &veYear
	errors[veMake.Key] = &veMake
	return errors
}

// ******************* END RESOURCE SECTION *************************** //

func TestGsonApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GsonApi Suite")
}
