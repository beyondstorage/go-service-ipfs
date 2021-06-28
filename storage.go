package ipfs

import (
	"context"
	"io"

	. "github.com/beyondstorage/go-storage/v4/types"
)

func (s *Storage) metadata(opt pairStorageMetadata) (meta *StorageMeta) {
	meta = NewStorageMeta()
	meta.Name = s.name
	meta.WorkDir = s.workDir
	return meta
}

func (s *Storage) create(path string, opt pairStorageCreate) (o *Object) {
	panic("not implemented")
}

func (s *Storage) read(ctx context.Context, path string, w io.Writer, opt pairStorageRead) (n int64, err error) {
	f, err := s.ipfs.FilesRead(ctx, path)
	if err != nil {
		return 0, err
	}
	n, err = io.Copy(w, f)
	return
}

func (s *Storage) write(ctx context.Context, path string, r io.Reader, size int64, opt pairStorageWrite) (n int64, err error) {
	err = s.ipfs.FilesWrite(ctx, path, r)
	if err != nil {
		return 0, err
	}
	return size, nil
}

func (s *Storage) delete(ctx context.Context, path string, opt pairStorageDelete) (err error) {
	err = s.ipfs.FilesRm(ctx, path, true)
	return
}

func (s *Storage) stat(ctx context.Context, path string, opt pairStorageStat) (o *Object, err error) {
	stat, err := s.ipfs.FilesStat(ctx, path)
	if err != nil {
		return nil, err
	}
	o = NewObject(s, true)
	o.ID = path
	o.Path = path
	if opt.HasObjectMode && opt.ObjectMode.IsDir() {
		o.Mode |= ModeDir
	} else {
		o.Mode |= ModeRead
	}

	o.SetContentType(stat.Type)
	o.SetContentLength(int64(stat.Size))
	//o.SetContentMd5(stat.Hash) // TODO SHA-1 or SHA-256, not MD5

	var sm ObjectSystemMetadata
	// TODO copy ext metadata
	o.SetSystemMetadata(sm)

	return
}

func (s *Storage) list(ctx context.Context, path string, opt pairStorageList) (oi *ObjectIterator, err error) {
	dir, err := s.ipfs.FilesLs(ctx, path)
	if err != nil {
		return
	}

	finish := false
	nextFn := func(ctx context.Context, page *ObjectPage) error {
		if finish {
			return IterateDone
		}

		for _, f := range dir {
			o := NewObject(s, true)
			o.ID = f.Name
			o.Path = f.Name
			o.Mode |= ModeRead
			o.SetContentLength(int64(f.Size))

			page.Data = append(page.Data, o)
		}
		return nil
	}

	oi = NewObjectIterator(ctx, nextFn, nil)
	return
}
