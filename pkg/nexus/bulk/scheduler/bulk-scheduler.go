package scheduler

import (
	models2 "github.com/nuclio/nuclio/pkg/nexus/common/models"
	"time"

	"github.com/nuclio/nuclio/pkg/nexus/bulk/models"
	"github.com/nuclio/nuclio/pkg/nexus/common/models/interfaces"
	"github.com/nuclio/nuclio/pkg/nexus/common/models/structs"
	"github.com/nuclio/nuclio/pkg/nexus/common/scheduler"
)

// BulkScheduler waits for a function call (request) to be queued several times and then calls the function several
// times with the respective attributes
// Purpose: to reduce the amount of colds-starts
// A detailed model of the scheduler can be found here: profaastinate/docs/diagrams/uml/activity/bulk-schedule.puml
type BulkScheduler struct {
	// BaseNexusScheduler is the base scheduler
	scheduler.BaseNexusScheduler

	// BulkSchedulerConfig is the config of the scheduler
	models.BulkSchedulerConfig
}

// NewScheduler creates a new bulk scheduler
func NewScheduler(baseNexusScheduler *scheduler.BaseNexusScheduler, bulkConfig models.BulkSchedulerConfig) *BulkScheduler {
	baseScheduler := baseNexusScheduler
	baseScheduler.Name = models2.BULK_SCHEDULER_NAME

	return &BulkScheduler{
		BaseNexusScheduler:  *baseScheduler,
		BulkSchedulerConfig: bulkConfig,
	}
}

// NewDefaultScheduler creates a new bulk scheduler with default values
func NewDefaultScheduler(baseNexusScheduler *scheduler.BaseNexusScheduler) *BulkScheduler {
	return NewScheduler(baseNexusScheduler, *models.NewDefaultBulkSchedulerConfig())
}

// Start starts the scheduler
func (ds *BulkScheduler) Start() {
	ds.RunFlag = true

	ds.executeSchedule()
}

// Stop stops the scheduler
func (ds *BulkScheduler) Stop() {
	ds.RunFlag = false
}

// GetStatus returns the running status of the scheduler
func (ds *BulkScheduler) GetStatus() interfaces.SchedulerStatus {
	if ds.RunFlag {
		return interfaces.Running
	} else {
		return interfaces.Stopped
	}
}

// executeSchedule checks if any items are ready to be called
func (ds *BulkScheduler) executeSchedule() {
	for ds.RunFlag {
		if ds.Queue.Len() == 0 || ds.BaseNexusScheduler.MaxParallelRequests.Load() < int32(ds.MinAmountOfBulkItems) {
			// TODO: sleep take care of offset due to processing
			time.Sleep(ds.SleepDuration)
			continue
		}

		// log.Println("Checking for bulking")
		if itemsToPop := ds.Queue.GetMostCommonEntryItems(); len(itemsToPop) >= ds.MinAmountOfBulkItems &&
			ds.BaseNexusScheduler.MaxParallelRequests.Load() > (ds.CurrentParallelRequests.Load()+int32(len(itemsToPop))) {
			ds.callAndRemoveItems(itemsToPop)
		} else if ds.BaseNexusScheduler.MaxParallelRequests.Load() >= int32(len(itemsToPop)) {
			time.Sleep(ds.SleepDuration)
		}
	}
}

// callAndRemoveItems calls the items synchronously on the default nuclio endpoint
// then they are removed them from the nexus queue
func (ds *BulkScheduler) callAndRemoveItems(items []*structs.NexusItem) {
	removedItems := ds.Queue.RemoveAll(items)
	if len(removedItems) == 0 {
		return
	}

	ds.Unpause(removedItems[0].Name)
	ds.CurrentParallelRequests.Add(int32(len(removedItems)))

	for _, item := range removedItems {

		go func(item *structs.NexusItem) {
			defer ds.CurrentParallelRequests.Add(-1)

			ds.SendToExecutionChannel(item.Name)
			ds.CallSynchronized(item)
		}(item)

	}
}
