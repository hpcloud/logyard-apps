package drain

import (
	"github.com/ActiveState/log"
	"strings"
	"logyard"
)

// RemoveOrphanedDrains removes all drains created by applog_endpoint.
func RemoveOrphanedDrains() {
	// Note that this is tricky to do when horizontally scalling
	// applog_endpoint. Could be solved easily by using nodeID or ip
	// addr in the drain name.
	logyardConfig := logyard.GetConfig()
	for name, _ := range logyardConfig.Drains {
		if strings.HasPrefix(name, DRAIN_PREFIX) {
			log.Infof("Removing orphaned drain %v", name)
			err := logyard.DeleteDrain(name)
			if err != nil {
				log.Warnf("Failed to delete drain %v -- %v",
					name, err)
			}
		}
	}
}
