package modules

import "sync"

type groupsMap map[string][]string

type SafeGroups struct {
	sync.Mutex
	List groupsMap
}

var Groups = SafeGroups{List: make(groupsMap)}
