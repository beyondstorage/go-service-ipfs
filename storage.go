package ipfs

import (
	"context"
	"io"

	ipfs "github.com/ipfs/go-ipfs-api"

	"github.com/beyondstorage/go-storage/v4/pkg/iowrap"
	"github.com/beyondstorage/go-storage/v4/services"
	. "github.com/beyondstorage/go-storage/v4/types"
)

func (s *Storage) copy(ctx context.Context, src string, dst string, opt pairStorageCopy) (err error) {
	return s.ipfs.FilesCp(ctx, s.getAbsPath(src), s.getAbsPath(dst))
}

func (s *Storage) create(path string, opt pairStorageCreate) (o *Object) {
	if opt.HasObjectMode && opt.ObjectMode.IsDir() {
		path += "/"
		o = NewObject(s, true)
		o.Mode = ModeDir
	} else {
		o = NewObject(s, false)
		o.Mode = ModeRead
	}
	o.ID = s.getAbsPath(path)
	o.Path = path
	return o
}

// AOS-46: Idempotent Storager Delete Operation
// @see https://github.com/beyondstorage/specs/blob/master/rfcs/46-idempotent-delete.md
func (s *Storage) delete(ctx context.Context, path string, opt pairStorageDelete) (err error) {
	err = s.ipfs.FilesRm(ctx, s.getAbsPath(path), true)
	return
}

func (s *Storage) list(ctx context.Context, path string, opt pairStorageList) (oi *ObjectIterator, err error) {
	rp := s.getAbsPath(path)

	var nextFn NextObjectFunc
	switch {
	case opt.ListMode.IsPart():
	case opt.ListMode.IsDir():
		nextFn = func(ctx context.Context, page *ObjectPage) error {
			dir, err := s.ipfs.FilesLs(ctx, rp, ipfs.FilesLs.Stat(true))
			if err != nil {
				return err
			}
			for _, f := range dir {
				o := NewObject(s, true)
				o.ID = rp + "/" + f.Name
				o.Path = f.Name
				o.Mode |= ModeRead
				o.SetContentLength(int64(f.Size))
				page.Data = append(page.Data, o)
			}
			return IterateDone
		}
	case opt.ListMode.IsPrefix():
	default:
		return nil, services.ListModeInvalidError{Actual: opt.ListMode}
	}
	oi = NewObjectIterator(ctx, nextFn, nil)
	return
}

func (s *Storage) metadata(opt pairStorageMetadata) (meta *StorageMeta) {
	meta = NewStorageMeta()
	meta.Name = s.name
	meta.WorkDir = s.workDir

	// TODO: repo/stat to get total/used size
	return meta
}

func (s *Storage) move(ctx context.Context, src string, dst string, opt pairStorageMove) (err error) {
	return s.ipfs.FilesMv(ctx, s.getAbsPath(src), s.getAbsPath(dst))
}

func (s *Storage) read(ctx context.Context, path string, w io.Writer, opt pairStorageRead) (n int64, err error) {
	fileOpts := make([]ipfs.FilesOpt, 0)
	if opt.HasOffset {
		fileOpts = append(fileOpts, ipfs.FilesRead.Offset(opt.Offset))
	}
	if opt.HasSize {
		fileOpts = append(fileOpts, ipfs.FilesRead.Count(opt.Size))
	}
	f, err := s.ipfs.FilesRead(ctx, s.getAbsPath(path), fileOpts...)
	if err != nil {
		return 0, err
	}
	if opt.HasIoCallback {
		iowrap.CallbackReadCloser(f, opt.IoCallback)
	}
	return io.Copy(w, f)
}

func (s *Storage) stat(ctx context.Context, path string, opt pairStorageStat) (o *Object, err error) {
	rp := s.getAbsPath(path)
	stat, err := s.ipfs.FilesStat(ctx, rp)
	if err != nil {
		return nil, err
	}
	o = NewObject(s, true)
	o.ID = rp
	o.Path = path
	if opt.HasObjectMode && opt.ObjectMode.IsDir() {
		o.Mode |= ModeDir
	} else {
		o.Mode |= ModeRead
	}
	o.SetContentType(stat.Type)
	o.SetContentLength(int64(stat.Size))

	var sm ObjectSystemMetadata
	sm.Hash = stat.Hash
	sm.Blocks = stat.Blocks
	sm.Local = stat.Local
	sm.WithLocality = stat.WithLocality
	sm.CumulativeSize = int(stat.CumulativeSize)
	sm.SizeLocal = int(stat.SizeLocal)
	o.SetSystemMetadata(sm)

	return
}

func (s *Storage) write(ctx context.Context, path string, r io.Reader, size int64, opt pairStorageWrite) (n int64, err error) {
	err = s.ipfs.FilesWrite(
		ctx, s.getAbsPath(path), r,
		ipfs.FilesWrite.Create(true),
		ipfs.FilesWrite.Parents(true),
	)
	if err != nil {
		return 0, err
	}
	return size, nil
}
