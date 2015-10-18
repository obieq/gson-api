package gsonapi

import (
	"log"

	"github.com/obieq/gas"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	BeforeEach(func() {
	})

	Context("Errors", func() {
		It("should panic when loading the config.json file fails", func() {
			defer func() {
				if r := recover(); r != nil {
					log.Println("Recovered from .env file test panic: ", r)
					Ω(r).ShouldNot(BeNil())
				}
			}()
			c := &config{}
			c.ParseConfigFile("invalid_file_name")
		})

		It("should panic when URL config value is not in the config.json file", func() {
			defer func() {
				if r := recover(); r != nil {
					log.Println("Recovered from URL test panic: ", r)
					Ω(r).ShouldNot(BeNil())
				}
			}()
			c := &config{}
			c.Validate()
		})

		It("should panic when loading a non-existent env value from that config.json file", func() {
			defer func() {
				if r := recover(); r != nil {
					log.Println("Recovered from non-existent env value test panic: ", r)
					Ω(r).ShouldNot(BeNil())
				}
			}()
			c := newConfig()
			Ω(c).ShouldNot(BeNil())
			gas.GetString("non_existing_env_test")
		})
	})
})
