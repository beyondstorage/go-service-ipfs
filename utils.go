package ipfs

import (
	"errors"
	"fmt"
	"strings"

	"github.com/beyondstorage/go-storage/v4/types"
	ipfs "github.com/ipfs/go-ipfs-api"
)

// Storage is the example client.
type Storage struct {
	ipfs *ipfs.Shell

	defaultPairs DefaultStoragePairs
	features     StorageFeatures

	name    string
	workDir string

	types.UnimplementedStorager
}

// String implements Storager.String
func (s *Storage) String() string {
	return fmt.Sprintf(
		"Storager IPFS {Name: %s, WorkDir: %s}",
		s.name, s.workDir,
	)
}

// NewStorager will create Storager only.
func NewStorager(pairs ...types.Pair) (types.Storager, error) {
	opt, err := parsePairStorageNew(pairs)
	if err != nil {
		return nil, err
	}

	st := &Storage{
		name:    opt.Name,
		workDir: "/",
	}
	if opt.HasWorkDir {
		st.workDir = opt.WorkDir
	}

	// @see https://beyondstorage.io/zh-CN/docs/go-storage/pairs/endpoint/
	endpointParts := strings.SplitN(opt.Endpoint, ":", 2)
	if len(endpointParts) < 2 {
		return nil, errors.New("endpoint format error")
	}
	sh := ipfs.NewShell(endpointParts[1])
	if !sh.IsUp() {
		return nil, errors.New("ipfs not online")
	}
	st.ipfs = sh

	return st, nil
}

func (s *Storage) formatError(op string, err error, path ...string) error {
	panic("implement me")
}
