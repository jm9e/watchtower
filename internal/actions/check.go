package actions

import (
	"errors"
	"fmt"
	"github.com/jm9e/watchtower/pkg/filters"
	"github.com/jm9e/watchtower/pkg/sorter"
	"sort"
	"strings"
	"time"

	"github.com/opencontainers/runc/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	log "github.com/sirupsen/logrus"

	"github.com/jm9e/watchtower/pkg/container"
)

// CheckForMultipleWatchtowerInstances will ensure that there are not multiple instances of the
// watchtower running simultaneously. If multiple watchtower containers are detected, this function
// will stop and remove all but the most recently started container.
func CheckForMultipleWatchtowerInstances(client container.Client, cleanup bool, scope string) error {
	awaitDockerClient()
	containers, err := client.ListContainers(filters.WatchtowerContainersFilter)

	if err != nil {
		log.Fatal(err)
		return err
	}

	sameScopeContainers := []container.Container{}
	for _, c := range containers {
		if c.GetScope() == scope {
			log.Info("Found watchtower on same scope")
			sameScopeContainers = append(sameScopeContainers, c)
		} else {
			log.Info("Ignore watchtower on different scope " + c.GetScope())
		}
	}

	if len(sameScopeContainers) <= 1 {
		log.Debug("There are no additional watchtower containers")
		return nil
	}

	log.Info("Found multiple running watchtower instances. Cleaning up.")
	return cleanupExcessWatchtowers(sameScopeContainers, client, cleanup)
}

func cleanupExcessWatchtowers(containers []container.Container, client container.Client, cleanup bool) error {
	var cleanupErrors int
	var stopErrors int

	sort.Sort(sorter.ByCreated(containers))
	allContainersExceptLast := containers[0 : len(containers)-1]

	for _, c := range allContainersExceptLast {
		if err := client.StopContainer(c, 60); err != nil {
			// logging the original here as we're just returning a count
			logrus.Error(err)
			stopErrors++
			continue
		}

		if cleanup {
			if err := client.RemoveImageByID(c.ImageID()); err != nil {
				// logging the original here as we're just returning a count
				logrus.Error(err)
				cleanupErrors++
			}
		}
	}

	return createErrorIfAnyHaveOccurred(stopErrors, cleanupErrors)
}

func createErrorIfAnyHaveOccurred(c int, i int) error {
	if c == 0 && i == 0 {
		return nil
	}

	var output strings.Builder

	if c > 0 {
		output.WriteString(fmt.Sprintf("%d errors while stopping containers", c))
	}
	if i > 0 {
		output.WriteString(fmt.Sprintf("%d errors while cleaning up images", c))
	}
	return errors.New(output.String())
}

func awaitDockerClient() {
	log.Debug("Sleeping for a second to ensure the docker api client has been properly initialized.")
	time.Sleep(1 * time.Second)
}
