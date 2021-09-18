package tests

import (
	"os"
	"testing"

	tests "github.com/beyondstorage/go-integration-test/v4"
)

func TestStorager(t *testing.T) {
	if os.Getenv("STORAGE_IPFS_INTEGRATION_TEST") != "on" {
		t.Skipf("STORAGE_IPFS_INTEGRATION_TEST is not 'on', skipped")
	}
	tests.TestStorager(t, setupTest(t))
}

func TestCopier(t *testing.T) {
	if os.Getenv("STORAGE_IPFS_INTEGRATION_TEST") != "on" {
		t.Skipf("STORAGE_IPFS_INTEGRATION_TEST is not 'on', skipped")
	}
	tests.TestCopier(t, setupTest(t))
}

func TestMover(t *testing.T) {
	if os.Getenv("STORAGE_IPFS_INTEGRATION_TEST") != "on" {
		t.Skipf("STORAGE_IPFS_INTEGRATION_TEST is not 'on', skipped")
	}
	tests.TestMover(t, setupTest(t))
}

func TestStorageHttpSignerRead(t *testing.T) {
	if os.Getenv("STORAGE_IPFS_INTEGRATION_TEST") != "on" {
		t.Skipf("STORAGE_IPFS_INTEGRATION_TEST is not 'on', skipped")
	}
	tests.TestStorageHTTPSignerRead(t, setupTest(t))
}
