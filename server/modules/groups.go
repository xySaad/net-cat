package modules

import "sync"

type groupsMap map[string]map[string]*struct{} // nil value

type SafeGroups struct {
	sync.Mutex
	List groupsMap
}

var Groups = SafeGroups{List: make(groupsMap)}
