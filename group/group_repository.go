package group

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// GroupRepository defines an interface for a repository where groups can be
// queried, added, and manipulated. It should be thread safe.
type GroupRepository interface {
	GetAllGroups() (map[string]*Group, error)
	GetGroup(string) (*Group, error)
	SetGroup(*Group) (string, error)
	DeleteGroup(string) error
}

// InmemGroupRepository implements the GroupRepository interface with an inmem
// map of groups. It is thread safe.
type InmemGroupRepository struct {
	sync.Mutex
	groups map[string]*Group
}

// NewInmemGroupRepository instantiates a new InmemGroupRepository
func NewInmemGroupRepository() *InmemGroupRepository {
	return &InmemGroupRepository{
		groups: make(map[string]*Group),
	}
}

// GetAllGroups implements the GroupRepository interface and returns all the
// groups
func (igr *InmemGroupRepository) GetAllGroups() (map[string]*Group, error) {
	igr.Lock()
	defer igr.Unlock()

	return igr.groups, nil
}

// GetGroup implements the GroupRepository interface and returns a group by ID
func (igr *InmemGroupRepository) GetGroup(id string) (*Group, error) {
	igr.Lock()
	defer igr.Unlock()

	g, ok := igr.groups[id]
	if !ok {
		return nil, fmt.Errorf("Group %s not found", id)
	}
	return g, nil
}

// SetGroup implements the GroupRepository interface and inserts or updates a
// group in the local map. If the group's ID is already set and the map already
// contains a corresponding group, then the value is overriden. If the ID is not
// set, we assign a random one and insert the group in the map. In any case we
// return the ID of the group.
func (igr *InmemGroupRepository) SetGroup(group *Group) (string, error) {
	if group.GroupUID == "" {
		group.GroupUID = uuid.New().String()
	}

	group.LastUpdated = time.Now().Unix()

	igr.Lock()
	defer igr.Unlock()

	igr.groups[group.GroupUID] = group

	return group.GroupUID, nil
}

// DeleteGroup implements the GroupRepository interface and removes a group from
// the map
func (igr *InmemGroupRepository) DeleteGroup(id string) error {
	igr.Lock()
	defer igr.Unlock()

	delete(igr.groups, id)
	return nil
}
