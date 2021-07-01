package tests

import (
	"os"
	"testing"

	tests "github.com/beyondstorage/go-integration-test/v4"
)

func TestStorage(t *testing.T) {
	if os.Getenv("STORAGE_IPFS_INTEGRATION_TEST") != "on" {
		t.Skipf("STORAGE_IPFS_INTEGRATION_TEST is not 'on', skipped")
	}
	s := setupTest(t)
	tests.TestStorager(t, s)
	tests.TestCopier(t, s)
	tests.TestMover(t, s)
}
