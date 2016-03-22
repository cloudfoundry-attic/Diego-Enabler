package models_test

import (
	. "github.com/cloudfoundry-incubator/diego-enabler/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Space", func() {
	Describe("Parser", func() {
		jsonBody := `{
  "total_results": 1,
  "total_pages": 1,
  "prev_url": null,
  "next_url": null,
  "resources": [
    {
      "metadata": {
        "guid": "1f7ac3a5-6f4e-4d6c-8edd-ce694fc8c907",
        "url": "/v2/spaces/1f7ac3a5-6f4e-4d6c-8edd-ce694fc8c907",
        "created_at": "2016-03-16T16:36:38Z",
        "updated_at": null
      },
      "entity": {
        "name": "myspace",
        "organization_guid": "94fe9c1a-6bda-483b-bf48-d6fa39d08cb6",
        "space_quota_definition_guid": null,
        "allow_ssh": true,
        "organization_url": "/v2/organizations/94fe9c1a-6bda-483b-bf48-d6fa39d08cb6",
        "organization": {
          "metadata": {
            "guid": "94fe9c1a-6bda-483b-bf48-d6fa39d08cb6",
            "url": "/v2/organizations/94fe9c1a-6bda-483b-bf48-d6fa39d08cb6",
            "created_at": "2016-03-16T16:36:24Z",
            "updated_at": "2016-03-17T22:08:00Z"
          },
          "entity": {
            "name": "myorg",
            "billing_enabled": false,
            "quota_definition_guid": "e1b4ef20-a3a7-434c-bf07-01f09eea9441",
            "status": "active",
            "quota_definition_url": "/v2/quota_definitions/e1b4ef20-a3a7-434c-bf07-01f09eea9441",
            "spaces_url": "/v2/organizations/94fe9c1a-6bda-483b-bf48-d6fa39d08cb6/spaces",
            "domains_url": "/v2/organizations/94fe9c1a-6bda-483b-bf48-d6fa39d08cb6/domains",
            "private_domains_url": "/v2/organizations/94fe9c1a-6bda-483b-bf48-d6fa39d08cb6/private_domains",
            "users_url": "/v2/organizations/94fe9c1a-6bda-483b-bf48-d6fa39d08cb6/users",
            "managers_url": "/v2/organizations/94fe9c1a-6bda-483b-bf48-d6fa39d08cb6/managers",
            "billing_managers_url": "/v2/organizations/94fe9c1a-6bda-483b-bf48-d6fa39d08cb6/billing_managers",
            "auditors_url": "/v2/organizations/94fe9c1a-6bda-483b-bf48-d6fa39d08cb6/auditors",
            "app_events_url": "/v2/organizations/94fe9c1a-6bda-483b-bf48-d6fa39d08cb6/app_events",
            "space_quota_definitions_url": "/v2/organizations/94fe9c1a-6bda-483b-bf48-d6fa39d08cb6/space_quota_definitions"
          }
        },
        "developers_url": "/v2/spaces/1f7ac3a5-6f4e-4d6c-8edd-ce694fc8c907/developers",
        "developers": [
          {
            "metadata": {
              "guid": "4c408dbc-0ddf-4c9b-9fbf-a74fb991c298",
              "url": "/v2/users/4c408dbc-0ddf-4c9b-9fbf-a74fb991c298",
              "created_at": "2016-03-15T19:13:27Z",
              "updated_at": null
            },
            "entity": {
              "admin": false,
              "active": true,
              "default_space_guid": null,
              "spaces_url": "/v2/users/4c408dbc-0ddf-4c9b-9fbf-a74fb991c298/spaces",
              "organizations_url": "/v2/users/4c408dbc-0ddf-4c9b-9fbf-a74fb991c298/organizations",
              "managed_organizations_url": "/v2/users/4c408dbc-0ddf-4c9b-9fbf-a74fb991c298/managed_organizations",
              "billing_managed_organizations_url": "/v2/users/4c408dbc-0ddf-4c9b-9fbf-a74fb991c298/billing_managed_organizations",
              "audited_organizations_url": "/v2/users/4c408dbc-0ddf-4c9b-9fbf-a74fb991c298/audited_organizations",
              "managed_spaces_url": "/v2/users/4c408dbc-0ddf-4c9b-9fbf-a74fb991c298/managed_spaces",
              "audited_spaces_url": "/v2/users/4c408dbc-0ddf-4c9b-9fbf-a74fb991c298/audited_spaces"
            }
          }
        ],
        "managers_url": "/v2/spaces/1f7ac3a5-6f4e-4d6c-8edd-ce694fc8c907/managers",
        "managers": [
          {
            "metadata": {
              "guid": "4c408dbc-0ddf-4c9b-9fbf-a74fb991c298",
              "url": "/v2/users/4c408dbc-0ddf-4c9b-9fbf-a74fb991c298",
              "created_at": "2016-03-15T19:13:27Z",
              "updated_at": null
            },
            "entity": {
              "admin": false,
              "active": true,
              "default_space_guid": null,
              "spaces_url": "/v2/users/4c408dbc-0ddf-4c9b-9fbf-a74fb991c298/spaces",
              "organizations_url": "/v2/users/4c408dbc-0ddf-4c9b-9fbf-a74fb991c298/organizations",
              "managed_organizations_url": "/v2/users/4c408dbc-0ddf-4c9b-9fbf-a74fb991c298/managed_organizations",
              "billing_managed_organizations_url": "/v2/users/4c408dbc-0ddf-4c9b-9fbf-a74fb991c298/billing_managed_organizations",
              "audited_organizations_url": "/v2/users/4c408dbc-0ddf-4c9b-9fbf-a74fb991c298/audited_organizations",
              "managed_spaces_url": "/v2/users/4c408dbc-0ddf-4c9b-9fbf-a74fb991c298/managed_spaces",
              "audited_spaces_url": "/v2/users/4c408dbc-0ddf-4c9b-9fbf-a74fb991c298/audited_spaces"
            }
          }
        ],
        "auditors_url": "/v2/spaces/1f7ac3a5-6f4e-4d6c-8edd-ce694fc8c907/auditors",
        "auditors": [

        ],
        "apps_url": "/v2/spaces/1f7ac3a5-6f4e-4d6c-8edd-ce694fc8c907/apps",
        "apps": [
          {
            "metadata": {
              "guid": "fafd3547-50c7-4a51-be93-6d722fa42886",
              "url": "/v2/apps/fafd3547-50c7-4a51-be93-6d722fa42886",
              "created_at": "2016-03-21T17:19:09Z",
              "updated_at": "2016-03-21T17:20:53Z"
            },
            "entity": {
              "name": "ilovedogs",
              "production": false,
              "space_guid": "1f7ac3a5-6f4e-4d6c-8edd-ce694fc8c907",
              "stack_guid": "f3cecf19-4567-4dca-ad35-2a3af733cbde",
              "buildpack": null,
              "detected_buildpack": "staticfile 1.3.1",
              "environment_json": {

              },
              "memory": 256,
              "instances": 10,
              "disk_quota": 1024,
              "state": "STARTED",
              "version": "abc88442-db04-41f4-a2d6-96ae826fe9be",
              "command": null,
              "console": false,
              "debug": null,
              "staging_task_id": "df27c5755e964654bb1c79b1a6706efb",
              "package_state": "STAGED",
              "health_check_type": "port",
              "health_check_timeout": null,
              "staging_failed_reason": null,
              "staging_failed_description": null,
              "diego": true,
              "docker_image": null,
              "package_updated_at": "2016-03-21T17:19:12Z",
              "detected_start_command": "sh boot.sh",
              "enable_ssh": true,
              "docker_credentials_json": {
                "redacted_message": "[PRIVATE DATA HIDDEN]"
              },
              "ports": [
                8080
              ],
              "space_url": "/v2/spaces/1f7ac3a5-6f4e-4d6c-8edd-ce694fc8c907",
              "stack_url": "/v2/stacks/f3cecf19-4567-4dca-ad35-2a3af733cbde",
              "events_url": "/v2/apps/fafd3547-50c7-4a51-be93-6d722fa42886/events",
              "service_bindings_url": "/v2/apps/fafd3547-50c7-4a51-be93-6d722fa42886/service_bindings",
              "routes_url": "/v2/apps/fafd3547-50c7-4a51-be93-6d722fa42886/routes",
              "route_mappings_url": "/v2/apps/fafd3547-50c7-4a51-be93-6d722fa42886/route_mappings"
            }
          }
        ],
        "routes_url": "/v2/spaces/1f7ac3a5-6f4e-4d6c-8edd-ce694fc8c907/routes",
        "routes": [
          {
            "metadata": {
              "guid": "dcbe808f-3636-463d-aaa3-d326671acaae",
              "url": "/v2/routes/dcbe808f-3636-463d-aaa3-d326671acaae",
              "created_at": "2016-03-16T16:40:43Z",
              "updated_at": null
            },
            "entity": {
              "host": "ilovedogs",
              "path": "",
              "domain_guid": "07ec62ee-47c8-4599-9c16-a683dd52e282",
              "space_guid": "1f7ac3a5-6f4e-4d6c-8edd-ce694fc8c907",
              "service_instance_guid": null,
              "port": 0,
              "domain_url": "/v2/domains/07ec62ee-47c8-4599-9c16-a683dd52e282",
              "space_url": "/v2/spaces/1f7ac3a5-6f4e-4d6c-8edd-ce694fc8c907",
              "apps_url": "/v2/routes/dcbe808f-3636-463d-aaa3-d326671acaae/apps"
            }
          },
          {
            "metadata": {
              "guid": "4b8b370c-6a7c-4a4f-8550-96b2f9ea1449",
              "url": "/v2/routes/4b8b370c-6a7c-4a4f-8550-96b2f9ea1449",
              "created_at": "2016-03-17T22:06:51Z",
              "updated_at": null
            },
            "entity": {
              "host": "myapp",
              "path": "",
              "domain_guid": "07ec62ee-47c8-4599-9c16-a683dd52e282",
              "space_guid": "1f7ac3a5-6f4e-4d6c-8edd-ce694fc8c907",
              "service_instance_guid": null,
              "port": 0,
              "domain_url": "/v2/domains/07ec62ee-47c8-4599-9c16-a683dd52e282",
              "space_url": "/v2/spaces/1f7ac3a5-6f4e-4d6c-8edd-ce694fc8c907",
              "apps_url": "/v2/routes/4b8b370c-6a7c-4a4f-8550-96b2f9ea1449/apps"
            }
          }
        ],
        "domains_url": "/v2/spaces/1f7ac3a5-6f4e-4d6c-8edd-ce694fc8c907/domains",
        "domains": [
          {
            "metadata": {
              "guid": "07ec62ee-47c8-4599-9c16-a683dd52e282",
              "url": "/v2/domains/07ec62ee-47c8-4599-9c16-a683dd52e282",
              "created_at": "2016-03-15T19:06:19Z",
              "updated_at": "2016-03-17T16:55:10Z"
            },
            "entity": {
              "name": "bosh-lite.com",
              "router_group_guid": null
            }
          }
        ],
        "service_instances_url": "/v2/spaces/1f7ac3a5-6f4e-4d6c-8edd-ce694fc8c907/service_instances",
        "service_instances": [

        ],
        "app_events_url": "/v2/spaces/1f7ac3a5-6f4e-4d6c-8edd-ce694fc8c907/app_events",
        "events_url": "/v2/spaces/1f7ac3a5-6f4e-4d6c-8edd-ce694fc8c907/events",
        "security_groups_url": "/v2/spaces/1f7ac3a5-6f4e-4d6c-8edd-ce694fc8c907/security_groups",
        "security_groups": [
          {
            "metadata": {
              "guid": "bef73d8d-34f4-4aa0-b2bb-75eafbba94fe",
              "url": "/v2/security_groups/bef73d8d-34f4-4aa0-b2bb-75eafbba94fe",
              "created_at": "2016-03-15T19:06:19Z",
              "updated_at": null
            },
            "entity": {
              "name": "public_networks",
              "rules": [
                {
                  "destination": "0.0.0.0-9.255.255.255",
                  "protocol": "all"
                },
                {
                  "destination": "11.0.0.0-169.253.255.255",
                  "protocol": "all"
                },
                {
                  "destination": "169.255.0.0-172.15.255.255",
                  "protocol": "all"
                },
                {
                  "destination": "172.32.0.0-192.167.255.255",
                  "protocol": "all"
                },
                {
                  "destination": "192.169.0.0-255.255.255.255",
                  "protocol": "all"
                }
              ],
              "running_default": true,
              "staging_default": true,
              "spaces_url": "/v2/security_groups/bef73d8d-34f4-4aa0-b2bb-75eafbba94fe/spaces"
            }
          },
          {
            "metadata": {
              "guid": "265f01dc-afe3-444c-a126-b175c66ffd2e",
              "url": "/v2/security_groups/265f01dc-afe3-444c-a126-b175c66ffd2e",
              "created_at": "2016-03-15T19:06:19Z",
              "updated_at": null
            },
            "entity": {
              "name": "dns",
              "rules": [
                {
                  "destination": "0.0.0.0/0",
                  "ports": "53",
                  "protocol": "tcp"
                },
                {
                  "destination": "0.0.0.0/0",
                  "ports": "53",
                  "protocol": "udp"
                }
              ],
              "running_default": true,
              "staging_default": true,
              "spaces_url": "/v2/security_groups/265f01dc-afe3-444c-a126-b175c66ffd2e/spaces"
            }
          },
          {
            "metadata": {
              "guid": "a174ea6c-3838-408f-8baf-75b16b126a9d",
              "url": "/v2/security_groups/a174ea6c-3838-408f-8baf-75b16b126a9d",
              "created_at": "2016-03-15T19:06:19Z",
              "updated_at": null
            },
            "entity": {
              "name": "services",
              "rules": [
                {
                  "destination": "10.244.1.0/24",
                  "protocol": "all"
                },
                {
                  "destination": "10.244.3.0/24",
                  "protocol": "all"
                }
              ],
              "running_default": true,
              "staging_default": false,
              "spaces_url": "/v2/security_groups/a174ea6c-3838-408f-8baf-75b16b126a9d/spaces"
            }
          },
          {
            "metadata": {
              "guid": "97a1fa85-ba9e-4cea-9cbd-104073d23360",
              "url": "/v2/security_groups/97a1fa85-ba9e-4cea-9cbd-104073d23360",
              "created_at": "2016-03-15T19:06:19Z",
              "updated_at": null
            },
            "entity": {
              "name": "load_balancer",
              "rules": [
                {
                  "destination": "10.244.0.34",
                  "protocol": "all"
                }
              ],
              "running_default": true,
              "staging_default": false,
              "spaces_url": "/v2/security_groups/97a1fa85-ba9e-4cea-9cbd-104073d23360/spaces"
            }
          },
          {
            "metadata": {
              "guid": "54668902-0f14-4e24-ad7e-c3ebb8a4e7d4",
              "url": "/v2/security_groups/54668902-0f14-4e24-ad7e-c3ebb8a4e7d4",
              "created_at": "2016-03-15T19:06:19Z",
              "updated_at": null
            },
            "entity": {
              "name": "user_bosh_deployments",
              "rules": [
                {
                  "destination": "10.244.4.0-10.254.0.0",
                  "protocol": "all"
                }
              ],
              "running_default": true,
              "staging_default": false,
              "spaces_url": "/v2/security_groups/54668902-0f14-4e24-ad7e-c3ebb8a4e7d4/spaces"
            }
          }
        ]
      }
    }
  ]
}`

		It("parses", func() {
			spaces, err := SpacesParser{}.Parse([]byte(jsonBody))
			Expect(err).NotTo(HaveOccurred())
			Expect(spaces).NotTo(BeEmpty())

			space := spaces[0]
			Expect(space.Name).To(Equal("myspace"))
			Expect(space.Guid).To(Equal("1f7ac3a5-6f4e-4d6c-8edd-ce694fc8c907"))
			Expect(space.OrganizationGuid).To(Equal("94fe9c1a-6bda-483b-bf48-d6fa39d08cb6"))

			org := space.Organization
			Expect(org.Name).To(Equal("myorg"))
			Expect(org.Guid).To(Equal("94fe9c1a-6bda-483b-bf48-d6fa39d08cb6"))
		})
	})
})
