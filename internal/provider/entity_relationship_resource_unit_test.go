package provider_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/cortexapps/terraform-provider-cortex/internal/provider"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// fakeRelationshipApi stores destinations like the backend, but its GET keeps returning the previous destinations
// for staleReads reads after each PUT, as the real endpoint briefly does.
type fakeRelationshipApi struct {
	mu         sync.Mutex
	current    []string
	previous   []string
	staleLeft  int
	staleReads int
}

func newFakeRelationshipApi(t *testing.T, staleReads int) string {
	f := &fakeRelationshipApi{current: []string{}, previous: []string{}, staleReads: staleReads}
	server := httptest.NewServer(f)
	t.Cleanup(server.Close)
	return server.URL
}

func (f *fakeRelationshipApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if !strings.HasSuffix(r.URL.Path, "/destinations") {
		fakeCatalogFail(w, http.StatusInternalServerError, "unexpected request "+r.Method+" "+r.URL.Path)
		return
	}
	switch r.Method {
	case http.MethodPut:
		var body struct {
			Destinations []string `json:"destinations"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			fakeCatalogFail(w, http.StatusBadRequest, "bad body")
			return
		}
		f.previous, f.current, f.staleLeft = f.current, body.Destinations, f.staleReads
		writeDestinations(w, f.current)
	case http.MethodGet:
		if f.staleLeft > 0 {
			f.staleLeft--
			writeDestinations(w, f.previous)
			return
		}
		writeDestinations(w, f.current)
	default:
		fakeCatalogFail(w, http.StatusInternalServerError, "unexpected request "+r.Method+" "+r.URL.Path)
	}
}

func writeDestinations(w http.ResponseWriter, tags []string) {
	type destination struct {
		Tag string `json:"tag"`
	}
	out := struct {
		Destinations []destination `json:"destinations"`
	}{Destinations: []destination{}}
	for _, tag := range tags {
		out.Destinations = append(out.Destinations, destination{Tag: tag})
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
}

func entityRelationshipUnitConfig(url string, destinations ...string) string {
	quoted := make([]string, len(destinations))
	for i, d := range destinations {
		quoted[i] = fmt.Sprintf("%q", d)
	}
	return catalogUnitConfig(url, fmt.Sprintf(`
resource "cortex_entity_relationship" "test" {
  source_entity_tag       = "source"
  relationship_type_tag   = "depends-on"
  destination_entity_tags = [%s]
}
`, strings.Join(quoted, ", ")))
}

func shrinkDestinationsTestCase(url string) resource.TestCase {
	return resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: entityRelationshipUnitConfig(url, "dest-1", "dest-2", "dest-3"),
				Check:  resource.TestCheckResourceAttr("cortex_entity_relationship.test", "destination_entity_tags.#", "3"),
			},
			{
				Config: entityRelationshipUnitConfig(url, "dest-1", "dest-2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cortex_entity_relationship.test", "destination_entity_tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("cortex_entity_relationship.test", "destination_entity_tags.*", "dest-1"),
					resource.TestCheckTypeSetElemAttr("cortex_entity_relationship.test", "destination_entity_tags.*", "dest-2"),
				),
			},
		},
	}
}

// The GET lags for about 3.5 seconds after each write (14 reads, one every 250 ms). The write waits it out instead
// of storing the old destinations, which would fail the apply with "Provider produced inconsistent result after
// apply".
func TestEntityRelationshipResource_WaitsOutSlowRead(t *testing.T) {
	resource.UnitTest(t, shrinkDestinationsTestCase(newFakeRelationshipApi(t, 14)))
}

// The GET still lags when the wait runs out. State keeps the applied destinations, so the apply still succeeds.
func TestEntityRelationshipResource_KeepsAppliedTagsWhenReadStaysStale(t *testing.T) {
	t.Cleanup(provider.SetEntityRelationshipConsistencyWait(300 * time.Millisecond))
	resource.UnitTest(t, shrinkDestinationsTestCase(newFakeRelationshipApi(t, 3)))
}
