// Package settings is a drop-in replacement for the unpublished internal
// chaitin.cn/dev/go/settings module. It implements only what this codebase
// actually calls: New loads a YAML config file into a Setting, whose
// Unmarshal method decodes a single top-level key into a target struct or
// scalar. A package-level Unmarshal operates on the most recently loaded
// Setting, matching how callers use it as a process-wide default.
package settings

import (
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

// Setting holds the parsed top-level keys of a YAML config file.
type Setting struct {
	mu   sync.RWMutex
	root map[string]yaml.Node
}

var (
	defaultMu   sync.RWMutex
	defaultInst *Setting
)

// New reads and parses the YAML file at path.
func New(path string) (*Setting, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	root := map[string]yaml.Node{}
	if err := yaml.Unmarshal(data, &root); err != nil {
		return nil, err
	}

	s := &Setting{root: root}

	defaultMu.Lock()
	defaultInst = s
	defaultMu.Unlock()

	return s, nil
}

// Unmarshal decodes the value at the given top-level key into out. If the
// key is absent, out is left untouched and nil is returned, so callers can
// pre-populate defaults before loading.
func (s *Setting) Unmarshal(key string, out interface{}) error {
	s.mu.RLock()
	node, ok := s.root[key]
	s.mu.RUnlock()
	if !ok {
		return nil
	}
	return node.Decode(out)
}

// Unmarshal decodes the value at key from the most recently loaded Setting.
func Unmarshal(key string, out interface{}) error {
	defaultMu.RLock()
	s := defaultInst
	defaultMu.RUnlock()
	if s == nil {
		return nil
	}
	return s.Unmarshal(key, out)
}
