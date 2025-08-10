package applications

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"cse190-welp/proto/mydatabase"
)

var (
	ErrRecordNotFound    = errors.New("storage: item not found")
	ErrInvalidDeviceType = errors.New("invalid device type")
)

// EmulatedStorageApp is an in-memory emulated storage layer.
type EmulatedStorageApp struct {
	data    map[string]*mydatabase.DatabaseRecord
	mu      sync.Mutex
	dist    string
	latency time.Duration
}

// NewEmulatedStorageApp creates a new instance of EmulatedStorage. DeviceType must be 'disk' or 'ssd'
func NewEmulatedStorageApp(deviceType string) (*EmulatedStorageApp, error) {
	log.Printf("device type: %v", deviceType)
	// latency constants specified in microseconds
	validDevice := map[string]int{
		"ssd":   500,   // order of magnitude latency for consumer grade SSD
		"disk":  5000,  // order of magnitude latency for commodity disk
		"cloud": 50000, // order of magnitude latency for cloud storage service
	}
	// TODO: add dist options

	// check for valid device type
	latency, ok := validDevice[deviceType]
	if !ok {
		return nil, ErrInvalidDeviceType
	}
	return &EmulatedStorageApp{
		data:    make(map[string]*mydatabase.DatabaseRecord),
		dist:    "uniform",
		latency: time.Duration(latency) * time.Microsecond,
	}, nil
}

func (s *EmulatedStorageApp) sleep() {
	time.Sleep(s.latency)
}

func (s *EmulatedStorageApp) Get(key string) (*mydatabase.DatabaseRecord, bool) {
	s.sleep()
	s.mu.Lock()
	defer s.mu.Unlock()

	value, ok := s.data[key]
	return value, ok
}

func (s *EmulatedStorageApp) Set(record *mydatabase.DatabaseRecord) {
	s.sleep()
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[record.Key] = record
}

func (s *EmulatedStorageApp) Delete(key string) {
	s.sleep()
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, key)
}

type PersistentStorageApp struct {
	data      map[string]*mydatabase.DatabaseRecord
	dataMutex sync.RWMutex
	filePath  string
}

func NewPersistentStorageApp(filePath string) (*PersistentStorageApp, error) {
	kvs := &PersistentStorageApp{
		data:     make(map[string]*mydatabase.DatabaseRecord),
		filePath: filePath,
	}

	err := kvs.loadFromFile()
	if err != nil {
		return nil, err
	}

	return kvs, nil
}

func (s *PersistentStorageApp) Get(key string) (*mydatabase.DatabaseRecord, bool) {
	s.dataMutex.RLock()
	defer s.dataMutex.RUnlock()

	record, ok := s.data[key]
	return record, ok
}

func (kvs *PersistentStorageApp) Set(record *mydatabase.DatabaseRecord) {
	kvs.dataMutex.Lock()
	defer kvs.dataMutex.Unlock()

	kvs.data[record.Key] = record

	err := kvs.saveToFile()
	if err != nil {
		log.Println("Error saving key-value store to file:", err)
	}
}

func (kvs *PersistentStorageApp) loadFromFile() error {
	_, err := os.Stat(kvs.filePath)
	if os.IsNotExist(err) {
		// Create the file if it doesn't exist
		_, err := os.Create(kvs.filePath)
		if err != nil {
			return err
		}
	}

	data, err := ioutil.ReadFile(kvs.filePath)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		// File is empty, initialize kvs.data as an empty map
		kvs.data = make(map[string]*mydatabase.DatabaseRecord)
		return nil
	}

	err = json.Unmarshal(data, &kvs.data)
	if err != nil {
		return err
	}

	return nil
}

func (kvs *PersistentStorageApp) saveToFile() error {
	data, err := json.Marshal(kvs.data)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(kvs.filePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
