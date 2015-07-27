package gsonapi

import (
	"log"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	BeforeEach(func() {
	})

	Context("Parsing", func() {
		It("should load an env value", func() {
			os.Setenv("EXISTING_ENV_VALUE", "test")
			v := getEnv("EXISTING_ENV_VALUE")
			Ω(v).Should(Equal("test"))
		})
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
			getString("non_existing_env_test")
		})
	})
})
