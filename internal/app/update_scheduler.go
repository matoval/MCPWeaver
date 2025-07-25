package app

import (
	"sync"
	"time"
)

// UpdateScheduler manages scheduled update operations
type UpdateScheduler struct {
	schedule     *UpdateSchedule
	scheduledJob *ScheduledJob
	mutex        sync.RWMutex
	ticker       *time.Ticker
	stopChan     chan bool
	callback     ScheduleCallback
}

// ScheduleCallback defines the callback function for scheduled operations
type ScheduleCallback func(jobType ScheduledJobType) error

// ScheduledJob represents a scheduled update operation
type ScheduledJob struct {
	ID         string             `json:"id"`
	Type       ScheduledJobType   `json:"type"`
	Schedule   *UpdateSchedule    `json:"schedule"`
	UpdateInfo *UpdateInfo        `json:"updateInfo,omitempty"`
	CreatedAt  time.Time          `json:"createdAt"`
	LastRun    *time.Time         `json:"lastRun,omitempty"`
	NextRun    time.Time          `json:"nextRun"`
	RunCount   int                `json:"runCount"`
	Status     ScheduledJobStatus `json:"status"`
	Error      *APIError          `json:"error,omitempty"`
}

// ScheduledJobType defines the type of scheduled job
type ScheduledJobType string

const (
	ScheduledJobTypeCheck    ScheduledJobType = "check"
	ScheduledJobTypeDownload ScheduledJobType = "download"
	ScheduledJobTypeInstall  ScheduledJobType = "install"
)

// ScheduledJobStatus defines the status of a scheduled job
type ScheduledJobStatus string

const (
	ScheduledJobStatusActive   ScheduledJobStatus = "active"
	ScheduledJobStatusPaused   ScheduledJobStatus = "paused"
	ScheduledJobStatusComplete ScheduledJobStatus = "complete"
	ScheduledJobStatusFailed   ScheduledJobStatus = "failed"
	ScheduledJobStatusCanceled ScheduledJobStatus = "canceled"
)

// NewUpdateScheduler creates a new update scheduler
func NewUpdateScheduler(callback ScheduleCallback) *UpdateScheduler {
	return &UpdateScheduler{
		callback: callback,
		stopChan: make(chan bool),
	}
}

// ScheduleJob schedules a new update job
func (s *UpdateScheduler) ScheduleJob(jobType ScheduledJobType, schedule *UpdateSchedule, updateInfo *UpdateInfo) (*ScheduledJob, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Validate schedule
	if schedule == nil {
		return nil, &APIError{
			Type:    ErrorTypeValidation,
			Code:    "INVALID_SCHEDULE",
			Message: "Schedule is required",
		}
	}

	// Calculate next run time
	nextRun, err := s.calculateNextRun(schedule)
	if err != nil {
		return nil, err
	}

	// Create scheduled job
	job := &ScheduledJob{
		ID:         generateID(),
		Type:       jobType,
		Schedule:   schedule,
		UpdateInfo: updateInfo,
		CreatedAt:  time.Now(),
		NextRun:    nextRun,
		RunCount:   0,
		Status:     ScheduledJobStatusActive,
	}

	// Cancel existing job if any
	if s.scheduledJob != nil {
		s.cancelCurrentJob()
	}

	s.scheduledJob = job
	s.schedule = schedule

	// Start scheduler if not already running
	s.startScheduler()

	return job, nil
}

// GetScheduledJob returns the current scheduled job
func (s *UpdateScheduler) GetScheduledJob() *ScheduledJob {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.scheduledJob
}

// CancelScheduledJob cancels the current scheduled job
func (s *UpdateScheduler) CancelScheduledJob() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.scheduledJob == nil {
		return &APIError{
			Type:    ErrorTypeSystem,
			Code:    "NO_SCHEDULED_JOB",
			Message: "No scheduled job to cancel",
		}
	}

	s.cancelCurrentJob()
	return nil
}

// PauseScheduledJob pauses the current scheduled job
func (s *UpdateScheduler) PauseScheduledJob() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.scheduledJob == nil {
		return &APIError{
			Type:    ErrorTypeSystem,
			Code:    "NO_SCHEDULED_JOB",
			Message: "No scheduled job to pause",
		}
	}

	if s.scheduledJob.Status != ScheduledJobStatusActive {
		return &APIError{
			Type:    ErrorTypeSystem,
			Code:    "JOB_NOT_ACTIVE",
			Message: "Scheduled job is not active",
		}
	}

	s.scheduledJob.Status = ScheduledJobStatusPaused
	s.stopScheduler()

	return nil
}

// ResumeScheduledJob resumes a paused scheduled job
func (s *UpdateScheduler) ResumeScheduledJob() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.scheduledJob == nil {
		return &APIError{
			Type:    ErrorTypeSystem,
			Code:    "NO_SCHEDULED_JOB",
			Message: "No scheduled job to resume",
		}
	}

	if s.scheduledJob.Status != ScheduledJobStatusPaused {
		return &APIError{
			Type:    ErrorTypeSystem,
			Code:    "JOB_NOT_PAUSED",
			Message: "Scheduled job is not paused",
		}
	}

	s.scheduledJob.Status = ScheduledJobStatusActive

	// Recalculate next run time
	nextRun, err := s.calculateNextRun(s.scheduledJob.Schedule)
	if err != nil {
		return err
	}
	s.scheduledJob.NextRun = nextRun

	s.startScheduler()

	return nil
}

// UpdateSchedule updates the schedule of the current job
func (s *UpdateScheduler) UpdateSchedule(newSchedule *UpdateSchedule) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.scheduledJob == nil {
		return &APIError{
			Type:    ErrorTypeSystem,
			Code:    "NO_SCHEDULED_JOB",
			Message: "No scheduled job to update",
		}
	}

	if newSchedule == nil {
		return &APIError{
			Type:    ErrorTypeValidation,
			Code:    "INVALID_SCHEDULE",
			Message: "New schedule is required",
		}
	}

	// Calculate next run time with new schedule
	nextRun, err := s.calculateNextRun(newSchedule)
	if err != nil {
		return err
	}

	// Update job schedule
	s.scheduledJob.Schedule = newSchedule
	s.scheduledJob.NextRun = nextRun
	s.schedule = newSchedule

	// Restart scheduler with new schedule
	if s.scheduledJob.Status == ScheduledJobStatusActive {
		s.stopScheduler()
		s.startScheduler()
	}

	return nil
}

// Stop stops the scheduler
func (s *UpdateScheduler) Stop() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.stopScheduler()
}

// Private methods

func (s *UpdateScheduler) startScheduler() {
	if s.ticker != nil {
		s.ticker.Stop()
	}

	if s.scheduledJob == nil || s.scheduledJob.Status != ScheduledJobStatusActive {
		return
	}

	// Calculate duration until next run
	now := time.Now()
	duration := s.scheduledJob.NextRun.Sub(now)

	if duration <= 0 {
		// Should run immediately
		go s.executeJob()
		return
	}

	// Create ticker for the next run
	s.ticker = time.NewTicker(duration)

	go func() {
		select {
		case <-s.ticker.C:
			s.executeJob()
		case <-s.stopChan:
			return
		}
	}()
}

func (s *UpdateScheduler) stopScheduler() {
	if s.ticker != nil {
		s.ticker.Stop()
		s.ticker = nil
	}

	select {
	case s.stopChan <- true:
	default:
	}
}

func (s *UpdateScheduler) cancelCurrentJob() {
	if s.scheduledJob != nil {
		s.scheduledJob.Status = ScheduledJobStatusCanceled
	}
	s.stopScheduler()
	s.scheduledJob = nil
	s.schedule = nil
}

func (s *UpdateScheduler) executeJob() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.scheduledJob == nil || s.scheduledJob.Status != ScheduledJobStatusActive {
		return
	}

	// Execute the callback
	err := s.callback(s.scheduledJob.Type)

	now := time.Now()
	s.scheduledJob.LastRun = &now
	s.scheduledJob.RunCount++

	if err != nil {
		s.scheduledJob.Status = ScheduledJobStatusFailed
		s.scheduledJob.Error = err.(*APIError)
		s.stopScheduler()
		return
	}

	// Check if this was a one-time job
	if s.scheduledJob.Schedule.Type == ScheduleTypeImmediate {
		s.scheduledJob.Status = ScheduledJobStatusComplete
		s.stopScheduler()
		return
	}

	// Calculate next run time for recurring jobs
	nextRun, err := s.calculateNextRun(s.scheduledJob.Schedule)
	if err != nil {
		s.scheduledJob.Status = ScheduledJobStatusFailed
		s.scheduledJob.Error = err.(*APIError)
		s.stopScheduler()
		return
	}

	s.scheduledJob.NextRun = nextRun

	// Restart scheduler for next run
	s.stopScheduler()
	s.startScheduler()
}

func (s *UpdateScheduler) calculateNextRun(schedule *UpdateSchedule) (time.Time, error) {
	now := time.Now()

	switch schedule.Type {
	case ScheduleTypeImmediate:
		return now, nil

	case ScheduleTypeDaily:
		// Parse time (HH:MM format)
		targetTime, err := time.Parse("15:04", schedule.Time)
		if err != nil {
			return time.Time{}, &APIError{
				Type:    ErrorTypeValidation,
				Code:    "INVALID_TIME_FORMAT",
				Message: "Invalid time format, expected HH:MM",
			}
		}

		// Calculate next daily run
		nextRun := time.Date(now.Year(), now.Month(), now.Day(),
			targetTime.Hour(), targetTime.Minute(), 0, 0, now.Location())

		// If the time has already passed today, schedule for tomorrow
		if nextRun.Before(now) {
			nextRun = nextRun.Add(24 * time.Hour)
		}

		return nextRun, nil

	case ScheduleTypeWeekly:
		// Parse time
		targetTime, err := time.Parse("15:04", schedule.Time)
		if err != nil {
			return time.Time{}, &APIError{
				Type:    ErrorTypeValidation,
				Code:    "INVALID_TIME_FORMAT",
				Message: "Invalid time format, expected HH:MM",
			}
		}

		// Calculate next weekly run
		daysUntilTarget := (schedule.DayOfWeek - int(now.Weekday()) + 7) % 7
		if daysUntilTarget == 0 {
			// Same day of week, check if time has passed
			todayAtTime := time.Date(now.Year(), now.Month(), now.Day(),
				targetTime.Hour(), targetTime.Minute(), 0, 0, now.Location())
			if todayAtTime.Before(now) {
				daysUntilTarget = 7 // Next week
			}
		}

		nextRun := now.Add(time.Duration(daysUntilTarget) * 24 * time.Hour)
		nextRun = time.Date(nextRun.Year(), nextRun.Month(), nextRun.Day(),
			targetTime.Hour(), targetTime.Minute(), 0, 0, nextRun.Location())

		return nextRun, nil

	case ScheduleTypeMonthly:
		// Parse time
		targetTime, err := time.Parse("15:04", schedule.Time)
		if err != nil {
			return time.Time{}, &APIError{
				Type:    ErrorTypeValidation,
				Code:    "INVALID_TIME_FORMAT",
				Message: "Invalid time format, expected HH:MM",
			}
		}

		// Calculate next monthly run
		year, month := now.Year(), now.Month()

		// Try current month first
		nextRun := time.Date(year, month, schedule.DayOfMonth,
			targetTime.Hour(), targetTime.Minute(), 0, 0, now.Location())

		// If the date doesn't exist in current month or has passed, go to next month
		if nextRun.Day() != schedule.DayOfMonth || nextRun.Before(now) {
			// Go to next month
			month++
			if month > 12 {
				month = 1
				year++
			}
			nextRun = time.Date(year, month, schedule.DayOfMonth,
				targetTime.Hour(), targetTime.Minute(), 0, 0, now.Location())

			// If day doesn't exist in next month, use last day of month
			if nextRun.Day() != schedule.DayOfMonth {
				nextRun = time.Date(year, month+1, 0,
					targetTime.Hour(), targetTime.Minute(), 0, 0, now.Location())
			}
		}

		return nextRun, nil

	case ScheduleTypeManual:
		// Manual scheduling doesn't have automatic next run
		return time.Time{}, &APIError{
			Type:    ErrorTypeValidation,
			Code:    "MANUAL_SCHEDULE_NO_AUTO_RUN",
			Message: "Manual schedules do not have automatic next run times",
		}

	default:
		return time.Time{}, &APIError{
			Type:    ErrorTypeValidation,
			Code:    "INVALID_SCHEDULE_TYPE",
			Message: "Invalid schedule type",
		}
	}
}
