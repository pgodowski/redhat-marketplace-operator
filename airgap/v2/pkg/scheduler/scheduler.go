// Copyright 2021 IBM Corp.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package scheduler

import (
	"time"

	"github.com/go-co-op/gocron"
	"github.com/go-logr/logr"
	v1 "github.com/redhat-marketplace/redhat-marketplace-operator/airgap/v2/apis/model/v1"
	"github.com/redhat-marketplace/redhat-marketplace-operator/airgap/v2/pkg/database"
	"github.com/redhat-marketplace/redhat-marketplace-operator/airgap/v2/pkg/dqlite"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SchedulerConfig struct {
	Log            logr.Logger
	Fs             *database.Database
	DBConfig       dqlite.DatabaseConfig
	CleanAfter     string
	PurgeAfter     string
	CronExpression string
}

// createScheduler return gocron.scheduler with job(s)
func (sfg *SchedulerConfig) createScheduler() *gocron.Scheduler {
	s := gocron.NewScheduler(time.UTC)
	s.SetMaxConcurrentJobs(1, gocron.WaitMode)

	if sfg.CleanAfter != "" {
		sfg.createJob(s, sfg.CleanAfter, false, "cleanAfter")
	}

	if sfg.PurgeAfter != "" {
		sfg.createJob(s, sfg.PurgeAfter, true, "purgeAfter")
	}

	if len(s.Jobs()) == 0 {
		return nil
	}

	return s
}

// createJob creates scheduler job
func (sfg *SchedulerConfig) createJob(s *gocron.Scheduler, before string, purge bool, tag string) {
	sfg.Log.Info("creating new job", "tag", tag, "before", before, "purge", purge, "cron", sfg.CronExpression)

	_, err := time.ParseDuration(before)
	if err != nil {
		sfg.Log.Error(err, "error parsing time duration")
	} else {
		_, err = s.Cron(sfg.CronExpression).Tag(tag).Do(
			func() {
				// run handler only for leader node
				if isLeader, _ := sfg.DBConfig.IsLeader(); isLeader {
					fileIds, _ := sfg.handler(before, purge)
					sfg.Log.Info("result", "fileIds", fileIds)
				}
			},
		)
		if err != nil {
			sfg.Log.Error(err, "error creating job")
		}
	}
}

// handler cleans/purges files based on given time duration and purge flag
func (sfg *SchedulerConfig) handler(before string, purge bool) ([]*v1.FileID, error) {
	now := time.Now()
	bf, _ := time.ParseDuration(before)
	t1 := now.Add(bf).Unix()
	t := &timestamppb.Timestamp{Seconds: t1}
	sfg.Log.Info("Job", "time", time.Now().Unix(), "before", t, "purge", purge)

	fileIds, err := sfg.Fs.CleanTombstones(t, purge)
	if err != nil {
		sfg.Log.Error(err, "failed to clean tombstoned files from database")
		return nil, err
	}
	return fileIds, err
}

// StartScheduler starts all job(s) for created scheduler
func (sfg *SchedulerConfig) StartScheduler() {
	s := sfg.createScheduler()
	if s != nil {
		sfg.Log.Info("starting scheduler")
		s.StartAsync()
	} else {
		sfg.Log.Info("no scheduler to start")
	}
}
