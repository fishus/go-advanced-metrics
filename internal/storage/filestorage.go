package storage

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/fishus/go-advanced-metrics/internal/metrics"
)

var ErrEmptyFilename = errors.New("filename for store metrics data is empty")

// FileStorage contains a set of values for all metrics and store its in file
type FileStorage struct {
	MemStorage
	filename string
	muSave   sync.Mutex
}

func NewFileStorage(filename string) *FileStorage {
	fs := &FileStorage{
		filename: filename,
	}
	fs.gauges = make(map[string]metrics.Gauge)
	fs.counters = make(map[string]metrics.Counter)
	return fs
}

func (fs *FileStorage) SetFilename(filename string) {
	fs.filename = filename
}

// Save saves metric values to a file.
func (fs *FileStorage) Save() error {
	fs.muSave.Lock()
	defer fs.muSave.Unlock()

	if fs.filename == "" {
		return ErrEmptyFilename
	}

	file, err := os.OpenFile(fs.filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		return err
	}

	defer file.Close()

	encoder := json.NewEncoder(file)

	err = encoder.Encode(&fs)
	if err != nil {
		return err
	}

	return nil
}

// Load reads metric values from a file.
func (fs *FileStorage) Load() error {
	if fs.filename == "" {
		return ErrEmptyFilename
	}

	file, err := os.OpenFile(fs.filename, os.O_RDONLY, 0)
	if err != nil {
		return err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)

	if err = decoder.Decode(&fs); err != nil {
		return err
	}

	return nil
}

var _ MetricsStorager = (*FileStorage)(nil)
var _ json.Marshaler = (*FileStorage)(nil)
var _ json.Unmarshaler = (*FileStorage)(nil)
var _ StorageSaver = (*FileStorage)(nil)
var _ StorageLoader = (*FileStorage)(nil)
