package actions_test

import (
	cli "github.com/docker/docker/client"
	"github.com/jm9e/watchtower/internal/actions"
	"github.com/jm9e/watchtower/pkg/container"
	"github.com/jm9e/watchtower/pkg/container/mocks"
	"github.com/jm9e/watchtower/pkg/types"
	"time"

	. "github.com/jm9e/watchtower/internal/actions/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("the update action", func() {
	var dockerClient cli.CommonAPIClient
	var client MockClient

	BeforeEach(func() {
		server := mocks.NewMockAPIServer()
		dockerClient, _ = cli.NewClientWithOpts(
			cli.WithHost(server.URL),
			cli.WithHTTPClient(server.Client()))
	})

	When("watchtower has been instructed to clean up", func() {
		BeforeEach(func() {
			pullImages := false
			removeVolumes := false
			client = CreateMockClient(
				&TestData{
					NameOfContainerToKeep: "test-container-02",
					Containers: []container.Container{
						CreateMockContainer(
							"test-container-01",
							"test-container-01",
							"fake-image:latest",
							time.Now().AddDate(0, 0, -1)),
						CreateMockContainer(
							"test-container-02",
							"test-container-02",
							"fake-image:latest",
							time.Now()),
						CreateMockContainer(
							"test-container-02",
							"test-container-02",
							"fake-image:latest",
							time.Now()),
					},
				},
				dockerClient,
				pullImages,
				removeVolumes,
			)
		})

		When("there are multiple containers using the same image", func() {
			It("should only try to remove the image once", func() {

				err := actions.Update(client, types.UpdateParams{Cleanup: true})
				Expect(err).NotTo(HaveOccurred())
				Expect(client.TestData.TriedToRemoveImageCount).To(Equal(1))
			})
		})
		When("there are multiple containers using different images", func() {
			It("should try to remove each of them", func() {
				client.TestData.Containers = append(
					client.TestData.Containers,
					CreateMockContainer(
						"unique-test-container",
						"unique-test-container",
						"unique-fake-image:latest",
						time.Now(),
					),
				)
				err := actions.Update(client, types.UpdateParams{Cleanup: true})
				Expect(err).NotTo(HaveOccurred())
				Expect(client.TestData.TriedToRemoveImageCount).To(Equal(2))
			})
		})
	})
})
