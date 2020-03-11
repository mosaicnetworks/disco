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
	groups       map[string]*Group
	timeout      int
	pruneTimeout int
	nextPrune    int64
}

// NewInmemGroupRepository instantiates a new InmemGroupRepository
func NewInmemGroupRepository(timeout int, pruneTimeout int) *InmemGroupRepository {
	return &InmemGroupRepository{
		groups:       make(map[string]*Group),
		timeout:      timeout,
		pruneTimeout: pruneTimeout,
	}
}

// GetAllGroups implements the GroupRepository interface and returns all the
// groups
func (igr *InmemGroupRepository) GetAllGroups() (map[string]*Group, error) {
	igr.Lock()
	defer igr.Unlock()

	igr.PruneGroups()

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

	g, ok := igr.groups[group.GroupUID]

	if !ok || g.LastBlockIndex <= group.LastBlockIndex {
		igr.groups[group.GroupUID] = group
	}

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

// PruneGroups iterates through the group list removing all items who
// were last updated more than igr.timeout seconds ago.
// If this function has been run within the last igr.pruneTimeout seconds,
// then the function returns without updating any groups.
// This function should only be called from a function that has already
// called igr.Lock().
func (igr *InmemGroupRepository) PruneGroups() error {

	//	igr.Lock()
	//	defer igr.Unlock()

	currentTime := time.Now().Unix()

	if igr.nextPrune > currentTime {
		// Next prune is not yet due

		fmt.Println("Prune Groups - do nothing")
		return nil
	}

	fmt.Println("Prune Groups")

	igr.nextPrune = currentTime + int64(igr.pruneTimeout)

	triggerTime := currentTime - int64(igr.timeout)

	for key, group := range igr.groups {
		if group.LastUpdated < triggerTime {
			fmt.Println("Prune Groups - delete " + key)
			delete(igr.groups, key)
		}
	}

	return nil
}
