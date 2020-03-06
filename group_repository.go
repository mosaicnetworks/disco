package disco

import (
	"fmt"

	"github.com/google/uuid"
)

// GroupRepository defines an interface for a repository where groups can be
// queried, added, and manipulated
type GroupRepository interface {
	GetAllGroups() (map[string]*Group, error)
	GetGroup(string) (*Group, error)
	SetGroup(*Group) (string, error)
	DeleteGroup(string) error
}

// InmemGroupRepository implements the GroupRepository interface with an inmem
// map of groups
type InmemGroupRepository struct {
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
	return igr.groups, nil
}

// GetGroup implements the GroupRepository interface and returns a group by ID
func (igr *InmemGroupRepository) GetGroup(id string) (*Group, error) {
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
	if group.ID == "" {
		group.ID = uuid.New().String()
	}

	igr.groups[group.ID] = group

	return group.ID, nil
}

// DeleteGroup implements the GroupRepository interface and removes a group from
// the map
func (igr *InmemGroupRepository) DeleteGroup(id string) error {
	delete(igr.groups, id)
	return nil
}
