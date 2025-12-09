package service

import (
	"app/src/utils"
	"errors"
	"runtime"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type HealthCheckService struct {
	log *logrus.Logger
	db  *gorm.DB
}

func NewHealthCheckService(db *gorm.DB) *HealthCheckService {
	return &HealthCheckService{
		log: utils.Log,
		db:  db,
	}
}

func (s *HealthCheckService) GormCheck() error {
	sqlDB, errDB := s.db.DB()
	if errDB != nil {
		s.log.Errorf("failed to access the database connection pool: %v", errDB)
		return errDB
	}

	if err := sqlDB.Ping(); err != nil {
		s.log.Errorf("failed to ping the database: %v", err)
		return err
	}

	return nil
}

// MemoryHeapCheck checks if heap memory usage exceeds a threshold
func (s *HealthCheckService) MemoryHeapCheck() error {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats) // Collect memory statistics

	heapAlloc := memStats.HeapAlloc            // Heap memory currently allocated
	heapThreshold := uint64(300 * 1024 * 1024) // Example threshold: 300 MB

	s.log.Infof("Heap Memory Allocation: %v bytes", heapAlloc)

	// If the heap allocation exceeds the threshold, return an error
	if heapAlloc > heapThreshold {
		s.log.Errorf("Heap memory usage exceeds threshold: %v bytes", heapAlloc)
		return errors.New("heap memory usage too high")
	}

	return nil
}
