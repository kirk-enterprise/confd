package api

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/vault"
)

func TestSysMountConfig(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	config := DefaultConfig()
	config.Address = addr

	client, err := NewClient(config)
	if err != nil {
		t.Fatal(err)
	}
	client.SetToken(token)

	// Set up a test mount
	path, err := testMount(client)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Sys().Unmount(path)

	// Get config info for this mount
	mountConfig, err := client.Sys().MountConfig(path)
	if err != nil {
		t.Fatal(err)
	}

	expectedDefaultTTL := 2592000
	if mountConfig.DefaultLeaseTTL != expectedDefaultTTL {
		t.Fatalf("Expected default lease TTL: %d, got %d",
			expectedDefaultTTL, mountConfig.DefaultLeaseTTL)
	}

	expectedMaxTTL := 2592000
	if mountConfig.MaxLeaseTTL != expectedMaxTTL {
		t.Fatalf("Expected default lease TTL: %d, got %d",
			expectedMaxTTL, mountConfig.MaxLeaseTTL)
	}
}

// testMount sets up a test mount of a generic backend w/ a random path; caller
// is responsible for unmounting
func testMount(client *Client) (string, error) {
	rand.Seed(time.Now().UTC().UnixNano())
	randInt := rand.New(rand.NewSource(time.Now().UnixNano())).Int()
	path := fmt.Sprintf("testmount-%d", randInt)
	err := client.Sys().Mount(path, &MountInput{Type: "generic"})
	return path, err
}
