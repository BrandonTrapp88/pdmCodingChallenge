package main

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

type Repository struct {
	mu      sync.RWMutex
	parts   map[string][]PartVersion
	nextID  int
	nextVer int
}

func NewRepository() *Repository {
	return &Repository{
		parts:   make(map[string][]PartVersion),
		nextID:  1,
		nextVer: 1,
	}
}

func (r *Repository) CreatePart(part Part) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	part.ID = r.generateID()
	version := r.nextVer
	r.nextVer++
	timestamp := time.Now()

	partVersion := PartVersion{
		Version:   version,
		Timestamp: timestamp,
		Part:      part,
	}
	r.parts[part.ID] = append(r.parts[part.ID], partVersion)

	return part.ID, nil
}

func (r *Repository) GetPart(id string) (Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	versions, exists := r.parts[id]
	if !exists || len(versions) == 0 {
		return Part{}, errors.New("part not found")
	}

	return versions[len(versions)-1].Part, nil
}

func (r *Repository) UpdatePart(id string, part Part) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.parts[id]
	if !exists {
		return errors.New("part not found")
	}

	part.ID = id
	version := r.nextVer
	r.nextVer++
	timestamp := time.Now()

	partVersion := PartVersion{
		Version:   version,
		Timestamp: timestamp,
		Part:      part,
	}
	r.parts[id] = append(r.parts[id], partVersion)

	return nil
}

func (r *Repository) DeletePart(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.parts[id]
	if !exists {
		return errors.New("part not found")
	}

	delete(r.parts, id)
	return nil
}

func (r *Repository) ListParts() []Part {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var parts []Part
	for _, versions := range r.parts {
		if len(versions) > 0 {
			parts = append(parts, versions[len(versions)-1].Part)
		}
	}
	return parts
}

func (r *Repository) GetPartVersion(id string, version int) (Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	versions, exists := r.parts[id]
	if !exists {
		return Part{}, errors.New("part not found")
	}

	for _, v := range versions {
		if v.Version == version {
			return v.Part, nil
		}
	}
	return Part{}, errors.New("version not found")
}

func (r *Repository) generateID() string {
	id := r.nextID
	r.nextID++
	return fmt.Sprintf("P%04d", id)
}

func (r *Repository) SearchPartsByName(name string) []Part {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []Part
	for _, versions := range r.parts {
		if len(versions) > 0 {
			latestVersion := versions[len(versions)-1].Part
			if strings.Contains(strings.ToLower(latestVersion.Name), strings.ToLower(name)) {
				results = append(results, latestVersion)
			}
		}
	}
	return results
}
