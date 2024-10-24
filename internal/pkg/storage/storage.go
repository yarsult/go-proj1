package storage

import (
	"encoding/json"
	"os"
	"proj1/internal/pkg/saving"
	"strconv"

	"go.uber.org/zap"
)

type Value struct {
	kind Kind
	st   string
	d    string
}

type Kind string

const (
	KindString Kind = "S"
	KindInt    Kind = "D"
)

type Storage struct {
	inner  map[string]Value
	logger *zap.Logger
}

func NewStorage() (Storage, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return Storage{}, err
	}
	defer logger.Sync()
	logger.Info("Created new storage")
	return Storage{inner: make(map[string]Value),
		logger: logger}, nil
}

func (r Storage) Set(key, val string) {
	var val1 Value
	if _, err := strconv.Atoi(val); err == nil {
		val1 = Value{kind: KindInt, d: val}
	} else {
		val1 = Value{kind: KindString, st: val}
	}
	r.inner[key] = val1
	r.logger.Info("key has been set")
}

func (r Storage) Get(key string) (string, bool) {
	res, ok := r.inner[key]
	if !ok {
		return "", false
	}
	r.logger.Info("val got")
	if res.kind == KindString {
		return res.st, true
	}
	return res.d, true
}

func (r Storage) GetKind(key string) string {
	res := r.inner[key]
	return string(res.kind)
}

func (r Storage) SaveToFile(filename string) error {
	data, err := json.MarshalIndent(r.inner, "", "  ")
	if err != nil {
		r.logger.Error("Failed to marshal SliceStorage to JSON", zap.Error(err))
		return err
	}
	err = saving.WriteAtomic(filename, data)
	if err != nil {
		r.logger.Error("Failed to write JSON to file", zap.Error(err))
		return err
	}

	r.logger.Info("SliceStorage successfully saved to file", zap.String("filename", filename))
	return nil
}

func (r Storage) LoadFromFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		r.logger.Error("Failed to read file", zap.Error(err))
		return err
	}
	var inner map[string]Value
	if err = json.Unmarshal(data, &inner); err != nil {
		r.logger.Error("Failed to unmarshal JSON", zap.Error(err))
		return err
	}
	r.inner = inner
	r.logger.Info("SliceStorage successfully loaded from file", zap.String("filename", filename))
	return nil
}
