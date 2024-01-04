package scheduler

import (
	"github.com/nuclio/nuclio/pkg/nexus/bulk/models"
	common "github.com/nuclio/nuclio/pkg/nexus/common/scheduler"
	"log"
	"time"
)

type BulkScheduler struct {
	common.BaseNexusScheduler

	models.BulkSchedulerConfig
}

func NewScheduler(baseNexusScheduler *common.BaseNexusScheduler, bulkConfig models.BulkSchedulerConfig) *BulkScheduler {
	return &BulkScheduler{
		BaseNexusScheduler:  *baseNexusScheduler,
		BulkSchedulerConfig: bulkConfig,
	}
}

func NewDefaultScheduler(baseNexusScheduler *common.BaseNexusScheduler) *BulkScheduler {
	return NewScheduler(baseNexusScheduler, *models.NewDefaultBulkSchedulerConfig())
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

		log.Println("Checking for bulking")
		if itemsToPop := ds.Queue.GetMostCommonEntryItems(); len(itemsToPop) >= ds.MinAmountOfBulkItems {
			log.Println("items with name: " + itemsToPop[0].Name)
			ds.Queue.RemoveAll(itemsToPop)
		}
	}
}
