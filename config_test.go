package gsonapi

import (
	"errors"
	"log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	BeforeEach(func() {
	})

	Context("Errors", func() {
		It("should panic when loading the .env file fails", func() {
			defer func() {
				if r := recover(); r != nil {
					log.Println("Recovered from .env file test panic: ", r)
					Ω(r).ShouldNot(BeNil())
				}
			}()
			ct := &ConfigTester{Error: errors.New("oops")}
			c := Config{}
			c.ParseEnvs(ct)
		})

		It("should panic when URL config value is not in the .env file", func() {
			defer func() {
				if r := recover(); r != nil {
					log.Println("Recovered from URL test panic: ", r)
					Ω(r).ShouldNot(BeNil())
				}
			}()
			ct := &ConfigTester{Envs: map[string]string{}}
			c := Config{}
			c.ParseEnvs(ct)
		})
	})
})
