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
	GetAllGroupsByAppID(appID string) (map[string]*Group, error)
	GetGroup(groupID string) (*Group, error)
	SetGroup(group *Group) (string, error)
	DeleteGroup(groupID string) error
}

// InmemGroupRepository implements the GroupRepository interface with an inmem
// map of groups. It is thread safe.
type InmemGroupRepository struct {
	sync.Mutex
	groupsByID    map[string]*Group   // [group ID] => Group
	groupsByAppID map[string][]string // [app ID] => [GroupID,...]
}

// NewInmemGroupRepository instantiates a new InmemGroupRepository
func NewInmemGroupRepository() *InmemGroupRepository {
	return &InmemGroupRepository{
		groupsByID:    make(map[string]*Group),
		groupsByAppID: make(map[string][]string),
	}
}

// GetAllGroups implements the GroupRepository interface and returns all the
// groups
func (igr *InmemGroupRepository) GetAllGroups() (map[string]*Group, error) {
	igr.Lock()
	defer igr.Unlock()

	return igr.groupsByID, nil
}

// GetAllGroupsByAppID implements the GroupRepository interface and returns all
// the groups associated with an AppID
func (igr *InmemGroupRepository) GetAllGroupsByAppID(appID string) (map[string]*Group, error) {
	igr.Lock()
	defer igr.Unlock()

	res := make(map[string]*Group)

	appGroups, ok := igr.groupsByAppID[appID]
	if !ok {
		return res, nil
	}

	for _, gid := range appGroups {
		res[gid] = igr.groupsByID[gid]
	}

	return res, nil
}

// GetGroup implements the GroupRepository interface and returns a group by ID
func (igr *InmemGroupRepository) GetGroup(id string) (*Group, error) {
	igr.Lock()
	defer igr.Unlock()

	g, ok := igr.groupsByID[id]
	if !ok {
		return nil, fmt.Errorf("Group %s not found", id)
	}
	return g, nil
}

// SetGroup implements the GroupRepository interface and inserts or updates a
// group in the local map. The group's AppID must be set. If the group's ID is
// already set and the map already contains a corresponding group, then the
// value is overriden. If the ID is not set, we assign a random one and insert
// the group in the map. In any case we return the ID of the group.
func (igr *InmemGroupRepository) SetGroup(group *Group) (string, error) {
	if group.AppID == "" {
		return "", fmt.Errorf("Group AppID not specified")
	}

	if group.ID == "" {
		group.ID = uuid.New().String()
	}

	group.LastUpdated = time.Now().Unix()

	igr.Lock()
	defer igr.Unlock()

	// If the group does not exist, add it to the AppID index
	if _, gok := igr.groupsByID[group.ID]; !gok {
		appGroups, aok := igr.groupsByAppID[group.AppID]
		if !aok {
			appGroups = []string{}
		}
		appGroups = append(appGroups, group.ID)
		igr.groupsByAppID[group.AppID] = appGroups
	}

	// Set group in main index
	igr.groupsByID[group.ID] = group

	return group.ID, nil
}

// DeleteGroup implements the GroupRepository interface and removes a group from
// the map
func (igr *InmemGroupRepository) DeleteGroup(id string) error {
	igr.Lock()
	defer igr.Unlock()

	// If the group exists, remove it from the AppID index
	if g, gok := igr.groupsByID[id]; gok {
		appGroups, aok := igr.groupsByAppID[g.AppID]
		if aok {
			for i, gid := range appGroups {
				if gid == id {
					// Remove the element at index i from appGroups.
					appGroups[i] = appGroups[len(appGroups)-1] // Copy last element to index i.
					appGroups[len(appGroups)-1] = ""           // Erase last element (write zero value).
					appGroups = appGroups[:len(appGroups)-1]   // Truncate slice.
					break
				}
			}
			igr.groupsByAppID[g.AppID] = appGroups
		}
	}

	delete(igr.groupsByID, id)
	return nil
}
