package nfs

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alist-org/alist/v3/internal/fs"
	"github.com/alist-org/alist/v3/internal/model"
	"github.com/go-git/go-billy/v5"
	log "github.com/sirupsen/logrus"
)

// AlistFS 是将 alist 存储系统适配到 billy.Filesystem 接口的适配器
type AlistFS struct {
	root string
}

// NewAlistFS 创建一个新的 alist 文件系统适配器
func NewAlistFS(root string) billy.Filesystem {
	if root == "" {
		root = "/"
	}
	return &AlistFS{
		root: root,
	}
}

// Create 创建一个新文件
func (fs *AlistFS) Create(filename string) (billy.File, error) {
	return fs.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
}

// Open 打开一个文件进行读取
func (fs *AlistFS) Open(filename string) (billy.File, error) {
	return fs.OpenFile(filename, os.O_RDONLY, 0)
}

// OpenFile 以指定模式打开文件
func (fs *AlistFS) OpenFile(filename string, flag int, perm os.FileMode) (billy.File, error) {
	fullPath := fs.getFullPath(filename)

	// 检查是否为创建新文件
	if flag&os.O_CREATE != 0 {
		// 尝试获取文件，如果不存在则创建临时文件句柄
		_, err := fs.getObject(fullPath)
		if err != nil {
			// 文件不存在，创建新文件（目前返回临时文件）
			if flag&(os.O_WRONLY|os.O_RDWR) != 0 {
				return NewWritableAlistFile(fullPath, fs), nil
			}
		}
	}

	// 获取文件信息
	obj, err := fs.getObject(fullPath)
	if err != nil {
		return nil, err
	}

	if obj.IsDir() {
		return nil, os.ErrInvalid
	}

	// 如果是写入模式，返回可写文件句柄
	if flag&(os.O_WRONLY|os.O_RDWR) != 0 {
		return NewWritableAlistFile(fullPath, fs), nil
	}

	// 获取下载链接
	link, err := fs.getDownloadLink(obj, fullPath)
	if err != nil {
		return nil, err
	}

	return NewAlistFile(obj, link), nil
}

// Stat 获取文件信息
func (fs *AlistFS) Stat(filename string) (os.FileInfo, error) {
	fullPath := fs.getFullPath(filename)
	obj, err := fs.getObject(fullPath)
	if err != nil {
		return nil, err
	}

	return &AlistFileInfo{obj: obj}, nil
}

// Rename 重命名文件
func (fs *AlistFS) Rename(oldpath, newpath string) error {
	// TODO: 实现重命名功能
	log.Warnf("Rename operation not yet implemented: %s -> %s", oldpath, newpath)
	return os.ErrPermission
}

// Remove 删除文件
func (fs *AlistFS) Remove(filename string) error {
	// TODO: 实现删除功能
	log.Warnf("Remove operation not yet implemented: %s", filename)
	return os.ErrPermission
}

// Join 连接路径
func (fs *AlistFS) Join(elem ...string) string {
	return filepath.Join(elem...)
}

// TempFile 创建临时文件
func (fs *AlistFS) TempFile(dir, prefix string) (billy.File, error) {
	// 使用时间戳生成临时文件名
	tempName := fmt.Sprintf("%s%d.tmp", prefix, time.Now().UnixNano())
	tempPath := fs.Join(dir, tempName)
	return NewWritableAlistFile(tempPath, fs), nil
}

// ReadDir 读取目录
func (fs *AlistFS) ReadDir(path string) ([]os.FileInfo, error) {
	fullPath := fs.getFullPath(path)

	// 获取目录对象
	obj, err := fs.getObject(fullPath)
	if err != nil {
		return nil, err
	}

	if !obj.IsDir() {
		return nil, os.ErrInvalid
	}

	// 列出目录内容
	objs, err := fs.listDirectory(obj, fullPath)
	if err != nil {
		return nil, err
	}

	var infos []os.FileInfo
	for _, o := range objs {
		infos = append(infos, &AlistFileInfo{obj: o})
	}

	return infos, nil
}

// MkdirAll 创建目录（包括父目录）
func (fs *AlistFS) MkdirAll(filename string, perm os.FileMode) error {
	// TODO: 实现目录创建功能
	log.Warnf("MkdirAll operation not yet implemented: %s", filename)
	return os.ErrPermission
}

// Lstat 获取链接信息
func (fs *AlistFS) Lstat(filename string) (os.FileInfo, error) {
	return fs.Stat(filename)
}

// Symlink 创建符号链接
func (fs *AlistFS) Symlink(target, link string) error {
	return os.ErrPermission
}

// Readlink 读取符号链接
func (fs *AlistFS) Readlink(link string) (string, error) {
	return "", os.ErrInvalid
}

// Chmod 修改文件权限
func (fs *AlistFS) Chmod(name string, mode os.FileMode) error {
	return os.ErrPermission
}

// Lchown 修改文件所有者
func (fs *AlistFS) Lchown(name string, uid, gid int) error {
	return os.ErrPermission
}

// Chown 修改文件所有者
func (fs *AlistFS) Chown(name string, uid, gid int) error {
	return os.ErrPermission
}

// Chtimes 修改文件时间
func (fs *AlistFS) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return os.ErrPermission
}

// Chroot 切换根目录
func (fs *AlistFS) Chroot(path string) (billy.Filesystem, error) {
	newRoot := fs.getFullPath(path)
	return NewAlistFS(newRoot), nil
}

// Root 返回根目录
func (fs *AlistFS) Root() string {
	return fs.root
}

// getFullPath 获取完整路径
func (fs *AlistFS) getFullPath(filename string) string {
	if strings.HasPrefix(filename, "/") {
		return filename
	}
	return filepath.Join(fs.root, filename)
}

// getObject 通过路径获取 alist 对象
func (afs *AlistFS) getObject(path string) (model.Obj, error) {
	ctx := context.Background()
	obj, err := afs.Get(ctx, path)
	if err != nil {
		log.Errorf("Failed to get object %s: %v", path, err)
		return nil, os.ErrNotExist
	}
	return obj, nil
}

// listDirectory 列出目录内容
func (afs *AlistFS) listDirectory(dir model.Obj, path string) ([]model.Obj, error) {
	ctx := context.Background()
	objs, err := afs.List(ctx, path)
	if err != nil {
		log.Errorf("Failed to list directory %s: %v", path, err)
		return nil, err
	}
	return objs, nil
}

// getDownloadLink 获取下载链接
func (afs *AlistFS) getDownloadLink(obj model.Obj, path string) (*model.Link, error) {
	ctx := context.Background()
	link, _, err := afs.Link(ctx, path, model.LinkArgs{})
	if err != nil {
		log.Errorf("Failed to get download link for %s: %v", path, err)
		return nil, err
	}
	return link, nil
}

// 这些方法通过 alist 的 fs 包来实现

func (afs *AlistFS) Get(ctx context.Context, path string) (model.Obj, error) {
	return fs.Get(ctx, path, &fs.GetArgs{})
}

func (afs *AlistFS) List(ctx context.Context, path string) ([]model.Obj, error) {
	return fs.List(ctx, path, &fs.ListArgs{})
}

func (afs *AlistFS) Link(ctx context.Context, path string, args model.LinkArgs) (*model.Link, model.Obj, error) {
	return fs.Link(ctx, path, args)
}