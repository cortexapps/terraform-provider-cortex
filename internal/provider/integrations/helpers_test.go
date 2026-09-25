package integrations_test

import (
	"os"
	"testing"

	"github.com/cortexapps/terraform-provider-cortex/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// The helpers below are copies of the ones in the provider_test package, which other packages cannot import.

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"cortex": providerserver.NewProtocol6WithError(provider.New("acctest")()),
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("CORTEX_API_TOKEN"); v == "" {
		t.Fatalf("Missing required environment variable: %s", "CORTEX_API_TOKEN")
	}
}
