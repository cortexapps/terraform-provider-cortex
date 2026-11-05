package provider

import "time"

// SetEntityRelationshipConsistencyWait shortens the post-write wait in tests
// and returns a func that restores it.
func SetEntityRelationshipConsistencyWait(d time.Duration) func() {
	old := entityRelationshipConsistencyWait
	entityRelationshipConsistencyWait = d
	return func() { entityRelationshipConsistencyWait = old }
}
