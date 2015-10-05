package gsonapi

import (
	"log"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/modocache/gory"
	"github.com/obieq/gas"
	validations "github.com/obieq/goar-validations"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/guregu/null.v2"

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
	ID          string        `json:"id"`
	Year        int           `json:"year,omitempty"`
	Make        string        `json:"make,omitempty"`
	BodyStyle   *string       `json:"body-style,omitempty"`
	Active      bool          `json:"active,omitempty"`
	Ages        []int         `json:"ages,omitempty"`
	Inspections []interface{} `json:"inspections,omitempty"`
}

// ******************* END MODEL SECTION ****************************** //

// ******************* BEGIN RESOURCE SECTION ************************* //

// Automobile Resource
type AutomobileResource struct {
	Resource    `jsonapi:"-"`
	Year        int              `json:"year,omitempty" jsonapi:"name=year"`
	Make        string           `json:"make,omitempty" jsonapi:"name=make"`
	BodyStyle   null.String      `json:"body-style,omitempty" jsonapi:"name=body-style"`
	Active      bool             `json:"active,omitempty" jsonapi:"name=active"`
	Drivers     []DriverResource `json:"drivers,omitempty" jsonapi:"-"`
	DriversIDs  string           `json:"-" jsonapi:"-"`
	Inspections []interface{}    `json:"inspections,omitempty" jsonapi:"name=inspections"`
	Ages        []interface{}    `json:"ages,omitempty" jsonapi:"name=ages"`
}

type InspectionResource struct {
	Name     string `json:"name,omitempty" jsonapi:"name=name"`
	Location string `json:"location,omitempty" jsonapi:"name=location"`
}

// Driver Resource
type DriverResource struct {
	Resource `jsonapi:"-"`
	Name     *string `json:"name,omitempty" jsonapi:"name=name"`
	Age      *int    `json:"age,omitempty" jsonapi:"name=age"`
	Active   *bool   `json:"active,omitempty" jsonapi:"name=active"`
}

// GetReferences to satisfy the jsonapi.MarshalReferences interface
func (r AutomobileResource) GetReferences() []jsonapi.Reference {
	return []jsonapi.Reference{
		{
			Type: "drivers",
			Name: "drivers",
		},
	}
}

func (r AutomobileResource) GetReferencedIDs() []jsonapi.ReferenceID {
	result := []jsonapi.ReferenceID{}
	for _, driver := range r.Drivers {
		result = append(result, jsonapi.ReferenceID{ID: driver.GetID(), Name: "drivers", Type: "drivers"})
	}

	return result
}

func (r AutomobileResource) GetReferencedStructs() []jsonapi.MarshalIdentifier {
	result := []jsonapi.MarshalIdentifier{}

	for key := range r.Drivers {
		result = append(result, r.Drivers[key])
	}

	return result
}

func (r AutomobileResource) GetName() string {
	return "automobiles"
}

func (r DriverResource) GetName() string {
	return "drivers"
}

// MapFromModel => maps a model to a resource
func (r *AutomobileResource) MapFromModel(model interface{}) (err error) {
	log.Println(model)
	log.Println("obie")
	m := model.(AutomobileModel)

	if !m.HasErrors() {
		r.ID = m.ID
		r.Year = m.Year
		r.Make = m.Make
		r.Active = m.Active
		if m.BodyStyle != nil {
			r.BodyStyle = null.StringFromPtr(m.BodyStyle)
		}
		var ints gas.Ints = m.Ages
		if r.Ages, err = ints.ToInterfaces(); err != nil {
			return err
		}
		r.Inspections = m.Inspections
	} else {
		r.SetErrors(m.ErrorMap())
	}

	return nil
}

// // MapToModel => maps a resource to a model
func (r *AutomobileResource) MapToModel(model interface{}) (err error) {
	m := model.(*AutomobileModel)

	m.Year = r.Year
	m.Make = r.Make
	m.Active = r.Active
	if !r.BodyStyle.IsZero() {
		bs := r.BodyStyle.String
		m.BodyStyle = &bs
	}
	var interfaces gas.Interfaces = r.Ages
	if m.Ages, err = interfaces.ToInts(); err != nil {
		return err
	}
	m.Inspections = r.Inspections

	return nil
}

// // ******************* END RESOURCE SECTION *************************** //
//
// // ******************* BEGIN TEST HELPERS SECTION ********************* //
func BuildErrors() map[string]*validations.ValidationError {
	errors := map[string]*validations.ValidationError{}
	veYear := validations.ValidationError{Key: "year", Message: "cannot be greater than 2016"}
	veMake := validations.ValidationError{Key: "make", Message: "cannot be blank"}
	errors[veYear.Key] = &veYear
	errors[veMake.Key] = &veMake
	return errors
}

//
// // ******************* END TEST HELPERS SECTION *********************** //
//
// // ******************* BEGIN TEST FACTORIES SECTION ******************* //
var _ = BeforeSuite(func() {
	// INSPECTIONS
	gory.Define("inspectionResource1", InspectionResource{}, func(factory gory.Factory) {
		factory["Name"] = "inspection #1"
		factory["Location"] = "216 broad ave, richmond va 23226"
	})

	gory.Define("inspectionResource2", InspectionResource{}, func(factory gory.Factory) {
		factory["Name"] = "inspection #2"
		factory["Location"] = "2201 stoddard ct, arlington va 22202"
	})

	// DRIVERS
	gory.Define("driverResource1", DriverResource{}, func(factory gory.Factory) {
		n := "paul walker"
		age := 40
		a := true

		factory["ID"] = "driver-id-1"
		factory["Name"] = &n
		factory["Age"] = &age
		factory["Active"] = &a
	})

	gory.Define("driverResource2", DriverResource{}, func(factory gory.Factory) {
		n := "steve mcqueen"
		age := 45
		a := false

		factory["ID"] = "driver-id-2"
		factory["Name"] = &n
		factory["Age"] = &age
		factory["Active"] = &a
	})

	// AUTOMOBILES
	gory.Define("automobileResource1", AutomobileResource{}, func(factory gory.Factory) {
		factory["ID"] = "aaaa-1111-bbbb-2222"
		factory["Year"] = 2010
		factory["Make"] = "Mazda"
		factory["BodyStyle"] = null.StringFrom("4 door sedan")
		factory["Active"] = true

		inspection1 := *gory.Build("inspectionResource1").(*InspectionResource)
		inspection2 := *gory.Build("inspectionResource2").(*InspectionResource)
		// factory["Inspections"] = []InspectionResource{inspection1, inspection2}
		factory["Inspections"] = []interface{}{inspection1, inspection2}

		driver1 := *gory.Build("driverResource1").(*DriverResource)
		driver2 := *gory.Build("driverResource2").(*DriverResource)
		factory["Drivers"] = []DriverResource{driver1, driver2}
	})

	gory.Define("automobileResource2", AutomobileResource{}, func(factory gory.Factory) {
		factory["ID"] = "cccc-3333-dddd-4444"
		factory["Year"] = 1960
		factory["Make"] = "Austin-Healey"
		factory["Active"] = true

		inspection1 := *gory.Build("inspectionResource1").(*InspectionResource)
		// factory["Inspections"] = []InspectionResource{inspection1}
		factory["Inspections"] = []interface{}{inspection1}
	})

	gory.Define("automobileResource3", AutomobileResource{}, func(factory gory.Factory) {
		factory["ID"] = "bbbb-2222-eeee-5555"
		factory["Year"] = 1980
		factory["Make"] = "Honda"
		factory["Active"] = false

		driver1 := *gory.Build("driverResource1").(*DriverResource)
		factory["Drivers"] = []DriverResource{driver1}
	})
})

// ******************* END TEST FACTORIES SECTION ********************* //
