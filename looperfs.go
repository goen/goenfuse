// this is the go-fuse glue
// +build !bazil

package main

import (
	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"

	"path/filepath"
	"syscall"
	"os"
	"io"
	"log"
	"time"
)

func loopcontext(v interface{}) interface{} {
	return v
}

type LooperFileSystem struct {
	// TODO - this should need default fill in.
	d dump
	pathfs.FileSystem
	Root string
}

// A FUSE filesystem that shunts all request to an underlying file
// system. Its main purpose is to provide test coverage without
// having to build a synthetic filesystem.
func NewLooperFileSystem(root string) pathfs.FileSystem {
	return &LooperFileSystem{
		FileSystem: pathfs.NewDefaultFileSystem(),
		Root:       root,
	}
}
func (fs *LooperFileSystem) OnMount(nodeFs *pathfs.PathNodeFs) {
}
func (fs *LooperFileSystem) OnUnmount() {}
func (fs *LooperFileSystem) GetPath(relPath string) string {
	return filepath.Join(fs.Root, relPath)
}
func (fs *LooperFileSystem) GetAttr(name string, context *fuse.Context) (a *fuse.Attr, code fuse.Status) {
	fs.d.write(Fileop{Code: OP_GETATTR, File: name})

	fullPath := fs.GetPath(name)
	var err error = nil
	st := syscall.Stat_t{}
	if name == "" {
		// When GetAttr is called for the toplevel directory, we always want
		// to look through symlinks.
		err = syscall.Stat(fullPath, &st)
	} else {
		err = syscall.Lstat(fullPath, &st)
	}
	if err != nil {
		return nil, fuse.ToStatus(err)
	}
	if st.Ino == 1 {
		return nil, fuse.ENOENT
	}

	a = &fuse.Attr{}
	a.FromStat(&st)
	return a, fuse.OK
}
func (fs *LooperFileSystem) OpenDir(name string, context *fuse.Context) (stream []fuse.DirEntry, status fuse.Status) {
	fs.d.write(Fileop{Code: OP_OPENDIR, File: name})

	// What other ways beyond O_RDONLY are there to open
	// directories?
	f, err := os.Open(fs.GetPath(name))
	if err != nil {
		return nil, fuse.ToStatus(err)
	}
	want := 500
	output := make([]fuse.DirEntry, 0, want)
	for {
		infos, err := f.Readdir(want)
		for i := range infos {
			// workaround forhttps://code.google.com/p/go/issues/detail?id=5960
			if infos[i] == nil {
				continue
			}



			n := infos[i].Name()
			d := fuse.DirEntry{
				Name: n,
			}
			if s := fuse.ToStatT(infos[i]); s != nil {

				if s.Ino == 1 {
					// workaround for finding a nested root folder
					continue
				}

				d.Mode = uint32(s.Mode)
			} else {
				log.Printf("ReadDir entry %q for %q has no stat info", n, name)
			}
			output = append(output, d)
		}
		if len(infos) < want || err == io.EOF {
			break
		}
		if err != nil {
			log.Println("Readdir() returned err:", err)
			break
		}
	}
	f.Close()
	return output, fuse.OK
}
func (fs *LooperFileSystem) Open(name string, flags uint32, context *fuse.Context) (fuseFile nodefs.File, status fuse.Status) {
	fs.d.write(Fileop{Code: OP_OPEN, File: name})

	f, err := os.OpenFile(fs.GetPath(name), int(flags), 0)
	if err != nil {
		return nil, fuse.ToStatus(err)
	}
	return nodefs.NewLoopbackFile(f), fuse.OK
}
func (fs *LooperFileSystem) Chmod(path string, mode uint32, context *fuse.Context) (code fuse.Status) {
	fs.d.write(Fileop{Code: OP_CHMOD, File: path})

	err := os.Chmod(fs.GetPath(path), os.FileMode(mode))
	return fuse.ToStatus(err)
}
func (fs *LooperFileSystem) Chown(path string, uid uint32, gid uint32, context *fuse.Context) (code fuse.Status) {
	fs.d.write(Fileop{Code: OP_CHOWN, File: path})

	return fuse.ToStatus(os.Chown(fs.GetPath(path), int(uid), int(gid)))
}
func (fs *LooperFileSystem) Truncate(path string, offset uint64, context *fuse.Context) (code fuse.Status) {
	fs.d.write(Fileop{Code: OP_TRUNC, File: path})

	return fuse.ToStatus(os.Truncate(fs.GetPath(path), int64(offset)))
}
func (fs *LooperFileSystem) Utimens(path string, Atime *time.Time, Mtime *time.Time, context *fuse.Context) (code fuse.Status) {
	fs.d.write(Fileop{Code: OP_UTIMENS, File: path})

	var a time.Time
	if Atime != nil {
		a = *Atime
	}
	var m time.Time
	if Mtime != nil {
		m = *Mtime
	}
	return fuse.ToStatus(os.Chtimes(fs.GetPath(path), a, m))
}
func (fs *LooperFileSystem) Readlink(name string, context *fuse.Context) (out string, code fuse.Status) {
	fs.d.write(Fileop{Code: OP_READLNK, File: name})

	f, err := os.Readlink(fs.GetPath(name))
	return f, fuse.ToStatus(err)
}
func (fs *LooperFileSystem) Mknod(name string, mode uint32, dev uint32, context *fuse.Context) (code fuse.Status) {
	fs.d.write(Fileop{Code: OP_MKNOD, File: name})

	return fuse.ToStatus(syscall.Mknod(fs.GetPath(name), mode, int(dev)))
}
func (fs *LooperFileSystem) Mkdir(path string, mode uint32, context *fuse.Context) (code fuse.Status) {
	fs.d.write(Fileop{Code: OP_MKDIR, File: path})

	return fuse.ToStatus(os.Mkdir(fs.GetPath(path), os.FileMode(mode)))
}

// Don't use os.Remove, it removes twice (unlink followed by rmdir).
func (fs *LooperFileSystem) Unlink(name string, context *fuse.Context) (code fuse.Status) {
	fs.d.write(Fileop{Code: OP_UNLINK, File: name})

	return fuse.ToStatus(syscall.Unlink(fs.GetPath(name)))
}
func (fs *LooperFileSystem) Rmdir(name string, context *fuse.Context) (code fuse.Status) {
	fs.d.write(Fileop{Code: OP_RMDIR, File: name})

	return fuse.ToStatus(syscall.Rmdir(fs.GetPath(name)))
}
func (fs *LooperFileSystem) Symlink(pointedTo string, linkName string, context *fuse.Context) (code fuse.Status) {
	fs.d.write(Fileop{Code: OP_SLINK, File: linkName, Extra: pointedTo})

	return fuse.ToStatus(os.Symlink(pointedTo, fs.GetPath(linkName)))
}
func (fs *LooperFileSystem) Rename(oldPath string, newPath string, context *fuse.Context) (codee fuse.Status) {
	fs.d.write(Fileop{Code: OP_RENAME, File: newPath, Extra: oldPath})

	err := os.Rename(fs.GetPath(oldPath), fs.GetPath(newPath))
	return fuse.ToStatus(err)
}
func (fs *LooperFileSystem) Link(orig string, newName string, context *fuse.Context) (code fuse.Status) {
	fs.d.write(Fileop{Code: OP_LINK, File: newName, Extra: orig})

	return fuse.ToStatus(os.Link(fs.GetPath(orig), fs.GetPath(newName)))
}
func (fs *LooperFileSystem) Access(name string, mode uint32, context *fuse.Context) (code fuse.Status) {
	fs.d.write(Fileop{Code: OP_ACCESS, File: name})

	return fuse.ToStatus(syscall.Access(fs.GetPath(name), mode))
}
func (fs *LooperFileSystem) Create(path string, flags uint32, mode uint32, context *fuse.Context) (fuseFile nodefs.File, code fuse.Status) {
	fs.d.write(Fileop{Code: OP_CREATE, File: path})

	f, err := os.OpenFile(fs.GetPath(path), int(flags)|os.O_CREATE, os.FileMode(mode))
	return nodefs.NewLoopbackFile(f), fuse.ToStatus(err)
}
