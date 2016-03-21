package api_test

import (
	. "github.com/cloudfoundry-incubator/diego-enabler/api"

	"github.com/cloudfoundry-incubator/diego-enabler/api/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
)

var _ = Describe("Api", func() {
	Describe("NewGetAppsRequest", func() {
		var (
			apiClient  *ApiClient
			fakeFilter *fakes.FakeFilter
			params     map[string]interface{}
			baseUrl    string
			authToken  string

			request *http.Request
			err     error
		)

		BeforeEach(func() {
			baseUrl = "https://api.my-crazy-domain.com"
			authToken = "some-auth-token"
			fakeFilter = new(fakes.FakeFilter)
			params = map[string]interface{}{}
		})

		JustBeforeEach(func() {
			apiClient, err = NewApiClient(baseUrl, authToken)
			Expect(err).NotTo(HaveOccurred())

			request, err = apiClient.NewGetAppsRequest(fakeFilter, params)
		})

		It("works", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("sets the Authorization header", func() {
			Expect(request.Header.Get("Authorization")).To(Equal(authToken))
		})

		Context("when given filters", func() {
			BeforeEach(func() {
				fakeFilter.ToFilterQueryParamReturns("something")
			})

			It("puts the filter into `q`", func() {
				Expect(request.URL.Query().Get("q")).To(Equal("something"))
			})
		})

		Context("when given params", func() {
			BeforeEach(func() {
				params = map[string]interface{}{"param1": "paramValue", "param2": "paramValue2"}
			})

			It("adds the params to the request", func() {
				Expect(request.URL.Query().Get("param1")).To(Equal("paramValue"))
				Expect(request.URL.Query().Get("param2")).To(Equal("paramValue2"))
			})
		})

		It("hits the appropriate API URL", func() {
			Expect(request.Method).To(Equal("GET"))
			Expect(request.URL.String()).To(Equal("https://api.my-crazy-domain.com/v2/apps"))
		})
	})

	Describe("EqualFilter", func() {
		It("serializes to name:val", func() {
			filter := EqualFilter{
				Name:  "foo",
				Value: true,
			}

			Expect(filter.ToFilterQueryParam()).To(Equal("foo:true"))

			filter = EqualFilter{
				Name:  "something",
				Value: 2,
			}

			Expect(filter.ToFilterQueryParam()).To(Equal("something:2"))

			filter = EqualFilter{
				Name:  "quux",
				Value: "bar",
			}

			Expect(filter.ToFilterQueryParam()).To(Equal("quux:bar"))
		})
	})

	Describe("Filters", func() {
		It("combines its filters together with semicolons", func() {
			filter1 := new(fakes.FakeFilter)
			filter1.ToFilterQueryParamReturns("something>2")

			filter2 := new(fakes.FakeFilter)
			filter2.ToFilterQueryParamReturns("bar::baaz")

			filters := Filters{
				filter1,
				filter2,
			}

			Expect(filters.ToFilterQueryParam()).To(Equal("something>2;bar::baaz"))
		})
	})

	Describe("PageParser", func() {
		It("parses", func() {
			jsonBody := `{
   "total_results": 2,
   "total_pages": 1,
   "prev_url": null,
   "next_url": null,
   "resources": [
      {
         "metadata": {
            "guid": "b2ba6466-23f7-4f90-935b-4da1c87b8943",
            "url": "/v2/apps/b2ba6466-23f7-4f90-935b-4da1c87b8943",
            "created_at": "2016-03-16T16:40:43Z",
            "updated_at": "2016-03-16T16:42:01Z"
         },
         "entity": {
            "name": "ilovedogs",
            "production": false,
            "space_guid": "1f7ac3a5-6f4e-4d6c-8edd-ce694fc8c907",
            "stack_guid": "f3cecf19-4567-4dca-ad35-2a3af733cbde",
            "buildpack": null,
            "detected_buildpack": "staticfile 1.3.1",
            "environment_json": {},
            "memory": 512,
            "instances": 4,
            "disk_quota": 1024,
            "state": "STARTED",
            "version": "7b0f71b8-39e0-4f21-8ed3-3dc287b8f9d2",
            "command": null,
            "console": false,
            "debug": null,
            "staging_task_id": "36e314e2d0dd41d7922e42493e2b7aee",
            "package_state": "STAGED",
            "health_check_type": "port",
            "health_check_timeout": null,
            "staging_failed_reason": null,
            "staging_failed_description": null,
            "diego": false,
            "docker_image": null,
            "package_updated_at": "2016-03-16T16:41:55Z",
            "detected_start_command": "sh boot.sh",
            "enable_ssh": true,
            "docker_credentials_json": {
               "redacted_message": "[PRIVATE DATA HIDDEN]"
            },
            "ports": null,
            "space_url": "/v2/spaces/1f7ac3a5-6f4e-4d6c-8edd-ce694fc8c907",
            "stack_url": "/v2/stacks/f3cecf19-4567-4dca-ad35-2a3af733cbde",
            "events_url": "/v2/apps/b2ba6466-23f7-4f90-935b-4da1c87b8943/events",
            "service_bindings_url": "/v2/apps/b2ba6466-23f7-4f90-935b-4da1c87b8943/service_bindings",
            "routes_url": "/v2/apps/b2ba6466-23f7-4f90-935b-4da1c87b8943/routes",
            "route_mappings_url": "/v2/apps/b2ba6466-23f7-4f90-935b-4da1c87b8943/route_mappings"
         }
      },
      {
         "metadata": {
            "guid": "280daf11-47ef-4121-b31b-73f3ef77cf04",
            "url": "/v2/apps/280daf11-47ef-4121-b31b-73f3ef77cf04",
            "created_at": "2016-03-17T22:06:49Z",
            "updated_at": "2016-03-17T22:07:27Z"
         },
         "entity": {
            "name": "myapp",
            "production": false,
            "space_guid": "1f7ac3a5-6f4e-4d6c-8edd-ce694fc8c907",
            "stack_guid": "f3cecf19-4567-4dca-ad35-2a3af733cbde",
            "buildpack": null,
            "detected_buildpack": "staticfile 1.3.1",
            "environment_json": {},
            "memory": 256,
            "instances": 10,
            "disk_quota": 1024,
            "state": "STARTED",
            "version": "5fae953a-8ef2-4a06-9a6b-e2ff38f2de60",
            "command": null,
            "console": false,
            "debug": null,
            "staging_task_id": "17176d03cf9f44858287b12894b55584",
            "package_state": "STAGED",
            "health_check_type": "port",
            "health_check_timeout": null,
            "staging_failed_reason": null,
            "staging_failed_description": null,
            "diego": false,
            "docker_image": null,
            "package_updated_at": "2016-03-17T22:06:55Z",
            "detected_start_command": "sh boot.sh",
            "enable_ssh": true,
            "docker_credentials_json": {
               "redacted_message": "[PRIVATE DATA HIDDEN]"
            },
            "ports": null,
            "space_url": "/v2/spaces/1f7ac3a5-6f4e-4d6c-8edd-ce694fc8c907",
            "stack_url": "/v2/stacks/f3cecf19-4567-4dca-ad35-2a3af733cbde",
            "events_url": "/v2/apps/280daf11-47ef-4121-b31b-73f3ef77cf04/events",
            "service_bindings_url": "/v2/apps/280daf11-47ef-4121-b31b-73f3ef77cf04/service_bindings",
            "routes_url": "/v2/apps/280daf11-47ef-4121-b31b-73f3ef77cf04/routes",
            "route_mappings_url": "/v2/apps/280daf11-47ef-4121-b31b-73f3ef77cf04/route_mappings"
         }
      }
   ]
}`
			pages, err := PageParser{}.Parse([]byte(jsonBody))
			Expect(err).NotTo(HaveOccurred())
			Expect(pages.TotalPages).To(Equal(1))
		})
	})
})
