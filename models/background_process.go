package models

import (
	"time"

	"gorm.io/gorm"
)

type ProcessType string

const (
	KEY_GENERATION ProcessType = "KEY_GENERATION"
	DECRYPTION     ProcessType = "DECRYPTION"
	ENCRYPTION     ProcessType = "ENCRYPTION"
)

type ProcessStatus string

const (
	COMPLETED ProcessStatus = "COMPLETED"
	FAILED    ProcessStatus = "FAILED"
)

type Process struct {
	Id          string      `json:"id" gorm:"type:uuid;default:gen_random_uuid()"`
	ProcessType ProcessType `json:"processType" gorm:"type:process_type"`
	StartedAt   time.Time   `json:"startedAt" gorm:"Default:CURRENT_TIMESTAMP"`
	EndedAt     time.Time   `json:"endedAt"`
	Status      string
}

type ProcessService struct {
	DB      *gorm.DB
	Process *Process
}

func (processService *ProcessService) StartProcess(processType ProcessType) (string, error) {
	db := processService.DB
	process := processService.Process
	err := db.Create(process).Error
	if err != nil {
		return "", err
	}
	return process.Id, nil
}

func (processService *ProcessService) EndProcess() error {
	db := processService.DB
	process := processService.Process

	err := db.Find(process, process.Id).Update("status", COMPLETED).Update("endedat", time.Now()).Error
	return err
}

func (processService *ProcessService) MarkFailed(processId string) error {
	db := processService.DB
	process := processService.Process

	err := db.Find(process, process.Id).Update("status", FAILED).Update("endedat", time.Now()).Error
	return err
}
