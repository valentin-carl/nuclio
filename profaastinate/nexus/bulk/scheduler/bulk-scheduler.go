package scheduler

import (
	"nexus/bulk/models"
	common "nexus/common/models"
	"time"
)

type BulkScheduler struct {
	common.BaseNexusScheduler

	models.BulkSchedulerConfig
}

func CreateNewScheduler(baseNexusScheduler *common.BaseNexusScheduler, bulkConfig models.BulkSchedulerConfig) *BulkScheduler {
	return &BulkScheduler{
		BaseNexusScheduler:  *baseNexusScheduler,
		BulkSchedulerConfig: bulkConfig,
	}
}

func (ds *BulkScheduler) Start() {
	ds.RunFlag = true

	ds.executeSchedule()
}

func (ds *BulkScheduler) Stop() {
	ds.RunFlag = false
}

func (ds *BulkScheduler) executeSchedule() {
	for ds.RunFlag {
		if ds.Queue.Len() == 0 {
			// TODO: sleep take care of offset due to processing
			time.Sleep(ds.SleepDuration)
			continue
		}

		if indicesToPop := ds.Queue.GetMostCommonEntryIndices(); len(indicesToPop) >= ds.MinAmountOfBulkItems {
			ds.Queue.RemoveAll(indicesToPop)
		}
	}
}
