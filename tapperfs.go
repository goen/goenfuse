// this is the go-fuse glue
// +build !bazil

package main

import (
	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"

	"fmt"
	"os"
	"time"
)

func tapcontext(i [][]string, z *self) nodefs.Node {
	return tapper_root{itemz: i, self: z}
}

type tapper_root struct {
	nodefs.Node
	itemz [][]string // len(itemz) = 1 + maximum name
	*self
}

type tappertrackernode struct {
	nodefs.Node
	*self
	f *os.File
}

type tapperdirnode struct {
	nodefs.Node
	i     uint64 //name = i, inode = i + inodeoffset
	itemz []string
}

type tapperbinlink struct {
	nodefs.Node
	inode uint64
}

func nrn(s *self) *tappertrackernode {
	return &tappertrackernode{Node: nodefs.NewDefaultNode(), self: s}
}

func ndn(i int, itemz []string) *tapperdirnode {
	return &tapperdirnode{Node: nodefs.NewDefaultNode(), i: uint64(i), itemz: itemz}
}

func nln(inode uint64) *tapperbinlink {
	return &tapperbinlink{Node: nodefs.NewDefaultNode(), inode: uint64(inode)}
}

func (s tapperdirnode) OpenDir(context *fuse.Context) ([]fuse.DirEntry, fuse.Status) {
	var foobar []fuse.DirEntry

	//	ibase := (s.i << 18) + 128

	for i := range s.itemz {
		item := fuse.DirEntry{Name: s.itemz[i], Mode: 0555}
		foobar = append(foobar, item)
	}
	return foobar, fuse.OK
}

func (s tapperdirnode) Lookup(out *fuse.Attr, name string, context *fuse.Context) (node *nodefs.Inode, code fuse.Status) {
	ibase := (s.i << 18) + 128

	// TODO: binary search
	for i := range s.itemz {
		if name == s.itemz[i] {

			ch := s.Inode().NewChild(name, false, nln(ibase+uint64(i)))

			return ch, fuse.OK
		}
	}

	return nil, fuse.ENOSYS
}

func (r tapper_root) OpenDir(context *fuse.Context) ([]fuse.DirEntry, fuse.Status) {
	var dirz [103]fuse.DirEntry
	foffset := 3

	dirz[0] = fuse.DirEntry{Name: "tracker", Mode: 0555}
	dirz[1] = fuse.DirEntry{Name: ".track", Mode: 0555}
	dirz[2] = fuse.DirEntry{Name: ".untrack", Mode: 0555}

	end := int(len(r.itemz))
	if end >= 100 {
		end = 100
	}

	for i := 0; i < end; i++ {
		dirz[i+foffset].Mode = uint32(os.ModeDir | 0555)
		dirz[i+foffset].Name = fmt.Sprintf("%02d", i)
	}

	sdirs := dirz[0 : end+foffset]

	return sdirs, fuse.OK
}

func (r tapper_root) Lookup(out *fuse.Attr, name string, context *fuse.Context) (node *nodefs.Inode, code fuse.Status) {
	if name == "tracker" || name == ".track" || name == ".untrack" {
		out.Mode = fuse.S_IFREG | 0644
		ch := r.Inode().NewChild(name, false, nrn(r.self))
		return ch, fuse.OK
	}

	var i int

	n, err := fmt.Sscanf(name, "%02d", &i)

	if (err != nil) || (n != 1) || (i >= len(r.itemz)) {
		return nil, fuse.ENOENT
	}

	out.Mode = fuse.S_IFDIR | 0755
	out.Size = 4096
	ch := r.Inode().NewChild(name, true, ndn(i, r.itemz[i]))

	return ch, fuse.OK
}
func (tapperdirnode) GetAttr(out *fuse.Attr, file nodefs.File, context *fuse.Context) (code fuse.Status) {

	out.Mode = fuse.S_IFDIR | 0755

	return fuse.OK
}
func (tappertrackernode) GetAttr(out *fuse.Attr, file nodefs.File, context *fuse.Context) (code fuse.Status) {

	out.Mode = fuse.S_IFREG | 0644

	return fuse.OK
}

func (tapper_root) OnUnmount() {
	fmt.Println("001")
}
func (tapper_root) OnMount(conn *nodefs.FileSystemConnector) {
	fmt.Println("002")
}
func (tapper_root) StatFs() *fuse.StatfsOut {
	fmt.Println("003")
	return nil
}

func (tapper_root) Deletable() bool {
	fmt.Println("004")
	return true
}

func (tapper_root) OnForget() {
	fmt.Println("005")
}

func (tapper_root) Access(mode uint32, context *fuse.Context) (code fuse.Status) {
	return fuse.ENOSYS
}
func (tapper_root) Readlink(c *fuse.Context) ([]byte, fuse.Status) {
	return nil, fuse.ENOSYS
}
func (tapper_root) Mknod(name string, mode uint32, dev uint32, context *fuse.Context) (newNode *nodefs.Inode, code fuse.Status) {
	return nil, fuse.ENOSYS
}
func (tapper_root) Mkdir(name string, mode uint32, context *fuse.Context) (newNode *nodefs.Inode, code fuse.Status) {
	return nil, fuse.ENOSYS
}
func (tapper_root) Unlink(name string, context *fuse.Context) (code fuse.Status) {
	return fuse.ENOSYS
}
func (tapper_root) Rmdir(name string, context *fuse.Context) (code fuse.Status) {
	return fuse.ENOSYS
}
func (tapper_root) Symlink(name string, content string, context *fuse.Context) (newNode *nodefs.Inode, code fuse.Status) {
	return nil, fuse.ENOSYS
}
func (tapper_root) Rename(oldName string, newParent nodefs.Node, newName string, context *fuse.Context) (code fuse.Status) {
	return fuse.ENOSYS
}
func (tapper_root) Link(name string, existing nodefs.Node, context *fuse.Context) (newNode *nodefs.Inode, code fuse.Status) {
	return nil, fuse.ENOSYS
}
func (tapper_root) Create(name string, flags uint32, mode uint32, context *fuse.Context) (file nodefs.File, newNode *nodefs.Inode, code fuse.Status) {
	return nil, nil, fuse.ENOSYS
}
func (tapper_root) Open(flags uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	return nil, fuse.ENOSYS
}
func (tapper_root) Flush(file nodefs.File, openFlags uint32, context *fuse.Context) (code fuse.Status) {
	return fuse.ENOSYS
}

func (tapper_root) GetXAttr(attribute string, context *fuse.Context) (data []byte, code fuse.Status) {
	return nil, fuse.ENOSYS
}
func (tapper_root) RemoveXAttr(attr string, context *fuse.Context) fuse.Status {
	return fuse.ENOSYS
}
func (tapper_root) SetXAttr(attr string, data []byte, flags int, context *fuse.Context) fuse.Status {
	return fuse.ENOSYS
}
func (tapper_root) ListXAttr(context *fuse.Context) (attrs []string, code fuse.Status) {
	return nil, fuse.ENOSYS
}

func (tapper_root) GetAttr(out *fuse.Attr, file nodefs.File, context *fuse.Context) (code fuse.Status) {

	out.Mode = fuse.S_IFDIR | 0755

	return fuse.OK
}
func (tapper_root) Chmod(file nodefs.File, perms uint32, context *fuse.Context) (code fuse.Status) {
	return fuse.ENOSYS
}
func (tapper_root) Chown(file nodefs.File, uid uint32, gid uint32, context *fuse.Context) (code fuse.Status) {
	return fuse.ENOSYS
}
func (tapper_root) Truncate(file nodefs.File, size uint64, context *fuse.Context) (code fuse.Status) {
	return fuse.ENOSYS
}
func (tapper_root) Utimens(file nodefs.File, atime *time.Time, mtime *time.Time, context *fuse.Context) (code fuse.Status) {
	return fuse.ENOSYS
}
func (tapper_root) Fallocate(file nodefs.File, off uint64, size uint64, mode uint32, context *fuse.Context) (code fuse.Status) {
	return fuse.ENOSYS
}
func (tapper_root) Read(file nodefs.File, dest []byte, off int64, context *fuse.Context) (fuse.ReadResult, fuse.Status) {
	if file != nil {
		return file.Read(dest, off)
	}
	return nil, fuse.ENOSYS
}
func (tapper_root) Write(file nodefs.File, data []byte, off int64, context *fuse.Context) (written uint32, code fuse.Status) {
	if file != nil {
		return file.Write(data, off)
	}
	return 0, fuse.ENOSYS
}
