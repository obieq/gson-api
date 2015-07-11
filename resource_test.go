package gsonapi

import (
	validations "github.com/obieq/goar-validations"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Resource", func() {
	var (
		errors *map[string]*validations.ValidationError
	)

	BeforeEach(func() {
		errorList := map[string]*validations.ValidationError{}
		veYear := validations.ValidationError{Key: "year", Message: "cannot be greater than 2016"}
		veMake := validations.ValidationError{Key: "make", Message: "cannot be blank"}
		errorList[veYear.Key] = &veYear
		errorList[veMake.Key] = &veMake

		errors = &errorList
	})

	Context("Errors", func() {
		It("should set errors", func() {
			r := AutomobileResource{}
			r.SetErrors(*errors)
		})

		It("should get errors", func() {
			r := AutomobileResource{}
			Ω(len(r.Errors())).Should(Equal(0))

			r.SetErrors(*errors)

			// verify
			Ω(len(r.Errors())).Should(Equal(2))
		})
	})
})
