package cron

import (
	"dapp/repo"
	"dapp/service/utils"
	"github.com/go-co-op/gocron"
	"log"
	"time"
)

// ISvcEventLog EventLog request service interface
type ISvcEventLog interface {
	MeinerCronJob() error
}

type svcEventLogReqs struct {
	svcConf  *utils.SvcConfig
	repoDapp *repo.RepoDapp
}

// endregion =============================================================================

// NewSvcRepoEventLog instantiate the Dapp request services
func NewSvcRepoEventLog(svcConf *utils.SvcConfig) ISvcEventLog {
	repoDapp := repo.NewRepoDapp(svcConf)
	return &svcEventLogReqs{svcConf, repoDapp}
}

// MeinerCronJob periodic task
func (e svcEventLogReqs) MeinerCronJob() error {
	// cron job is started only if it is active in configuration
	if e.svcConf.CronEnabled {
		log.Printf("schedules a new periodic Job with an interval: %d seconds", e.svcConf.EveryTime)
		cron := gocron.NewScheduler(time.UTC)

		_, err := cron.Every(e.svcConf.EveryTime).Seconds().WaitForSchedule().Do(e.doFunc)
		if err != nil {
			return err
		}
		// starts the scheduler asynchronously
		cron.StartAsync()
	}
	return nil
}

func (e svcEventLogReqs) doFunc() {
	log.Println("cron job executing")

	log.Println("cron job ending")
}
