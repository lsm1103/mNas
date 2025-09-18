package nfs

import (
	"io"
	"os"
	"time"

	"github.com/alist-org/alist/v3/internal/model"
	"github.com/go-git/go-billy/v5"
)

// AlistFile 实现 billy.File 接口
type AlistFile struct {
	obj      model.Obj
	link     *model.Link
	reader   io.ReadCloser
	position int64
	name     string
}

// NewAlistFile 创建一个新的 AlistFile
func NewAlistFile(obj model.Obj, link *model.Link) billy.File {
	return &AlistFile{
		obj:      obj,
		link:     link,
		position: 0,
		name:     obj.GetName(),
	}
}

// Name 返回文件名
func (f *AlistFile) Name() string {
	return f.name
}

// Write 写入数据（暂不支持）
func (f *AlistFile) Write(p []byte) (n int, err error) {
	return 0, os.ErrPermission
}

// Read 读取数据
func (f *AlistFile) Read(p []byte) (n int, err error) {
	if f.reader == nil {
		reader, err := f.getReader()
		if err != nil {
			return 0, err
		}
		f.reader = reader
	}

	n, err = f.reader.Read(p)
	f.position += int64(n)
	return n, err
}

// ReadAt 在指定位置读取数据
func (f *AlistFile) ReadAt(p []byte, off int64) (n int, err error) {
	// 对于网络流，ReadAt 比较复杂，暂时不完全支持
	return 0, os.ErrPermission
}

// Seek 移动文件指针
func (f *AlistFile) Seek(offset int64, whence int) (int64, error) {
	// 对于网络流，Seek 操作比较复杂，暂时不完全支持
	switch whence {
	case io.SeekStart:
		if offset == 0 {
			f.position = 0
			if f.reader != nil {
				f.reader.Close()
				f.reader = nil
			}
			return 0, nil
		}
	case io.SeekCurrent:
		if offset == 0 {
			return f.position, nil
		}
	case io.SeekEnd:
		return f.obj.GetSize(), nil
	}
	return f.position, os.ErrPermission
}

// Close 关闭文件
func (f *AlistFile) Close() error {
	if f.reader != nil {
		return f.reader.Close()
	}
	return nil
}

// Lock 锁定文件
func (f *AlistFile) Lock() error {
	// NFS 通常不需要文件锁
	return nil
}

// Unlock 解锁文件
func (f *AlistFile) Unlock() error {
	// NFS 通常不需要文件锁
	return nil
}

// Truncate 截断文件
func (f *AlistFile) Truncate(size int64) error {
	return os.ErrPermission
}

// getReader 获取文件读取器
func (f *AlistFile) getReader() (io.ReadCloser, error) {
	if f.link == nil {
		return nil, os.ErrInvalid
	}

	// 根据链接类型获取读取器
	switch {
	case f.link.URL != "":
		// HTTP 链接
		return f.getHTTPReader()
	case f.link.MFile != nil:
		// 本地文件路径
		return f.getFileReader()
	default:
		return nil, os.ErrInvalid
	}
}

// getHTTPReader 获取 HTTP 读取器
func (f *AlistFile) getHTTPReader() (io.ReadCloser, error) {
	// 这里需要实现 HTTP 流读取
	// 暂时返回错误，需要后续完善
	return nil, os.ErrPermission
}

// getFileReader 获取文件读取器
func (f *AlistFile) getFileReader() (io.ReadCloser, error) {
	if f.link.MFile == nil {
		return nil, os.ErrInvalid
	}
	return f.link.MFile, nil
}

// AlistFileInfo 实现 os.FileInfo 接口
type AlistFileInfo struct {
	obj model.Obj
}

// Name 返回文件名
func (fi *AlistFileInfo) Name() string {
	return fi.obj.GetName()
}

// Size 返回文件大小
func (fi *AlistFileInfo) Size() int64 {
	return fi.obj.GetSize()
}

// Mode 返回文件模式
func (fi *AlistFileInfo) Mode() os.FileMode {
	if fi.obj.IsDir() {
		return os.ModeDir | 0755
	}
	return 0644
}

// ModTime 返回修改时间
func (fi *AlistFileInfo) ModTime() time.Time {
	return fi.obj.ModTime()
}

// IsDir 判断是否为目录
func (fi *AlistFileInfo) IsDir() bool {
	return fi.obj.IsDir()
}

// Sys 返回底层数据
func (fi *AlistFileInfo) Sys() interface{} {
	return fi.obj
}