package nfs

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/alist-org/alist/v3/internal/model"
	"github.com/go-git/go-billy/v5"
	log "github.com/sirupsen/logrus"
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
	if f.link.URL == "" {
		return nil, os.ErrInvalid
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("GET", f.link.URL, nil)
	if err != nil {
		log.Errorf("Failed to create HTTP request for %s: %v", f.link.URL, err)
		return nil, err
	}

	// 添加必要的请求头
	if f.link.Header != nil {
		for key, values := range f.link.Header {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Failed to fetch HTTP content from %s: %v", f.link.URL, err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		log.Errorf("HTTP request failed with status %d for %s", resp.StatusCode, f.link.URL)
		return nil, os.ErrPermission
	}

	return resp.Body, nil
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
		return os.ModeDir | 0755  // 目录权限：rwxr-xr-x
	}
	return 0644  // 文件权限：rw-r--r--
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

// WritableAlistFile 实现可写的 billy.File 接口（用于上传）
type WritableAlistFile struct {
	path   string
	fs     *AlistFS
	buffer *bytes.Buffer
	name   string
	closed bool
}

// NewWritableAlistFile 创建一个新的可写 AlistFile
func NewWritableAlistFile(path string, fs *AlistFS) billy.File {
	name := filepath.Base(path)
	return &WritableAlistFile{
		path:   path,
		fs:     fs,
		buffer: bytes.NewBuffer(nil),
		name:   name,
		closed: false,
	}
}

// Name 返回文件名
func (f *WritableAlistFile) Name() string {
	return f.name
}

// Write 写入数据
func (f *WritableAlistFile) Write(p []byte) (n int, err error) {
	if f.closed {
		return 0, os.ErrClosed
	}
	return f.buffer.Write(p)
}

// Read 读取数据（写文件不支持读取）
func (f *WritableAlistFile) Read(p []byte) (n int, err error) {
	return 0, os.ErrInvalid
}

// ReadAt 在指定位置读取数据（写文件不支持）
func (f *WritableAlistFile) ReadAt(p []byte, off int64) (n int, err error) {
	return 0, os.ErrInvalid
}

// Seek 移动文件指针
func (f *WritableAlistFile) Seek(offset int64, whence int) (int64, error) {
	// 对于写入模式，简化 Seek 操作
	switch whence {
	case io.SeekStart:
		if offset == 0 {
			return 0, nil
		}
	case io.SeekCurrent:
		return int64(f.buffer.Len()), nil
	case io.SeekEnd:
		return int64(f.buffer.Len()), nil
	}
	return int64(f.buffer.Len()), nil
}

// Close 关闭文件并上传到 alist
func (f *WritableAlistFile) Close() error {
	if f.closed {
		return nil
	}
	f.closed = true

	// TODO: 实现文件上传到 alist
	log.Warnf("File upload not yet implemented: %s (%d bytes)", f.path, f.buffer.Len())

	return nil
}

// Lock 锁定文件
func (f *WritableAlistFile) Lock() error {
	return nil
}

// Unlock 解锁文件
func (f *WritableAlistFile) Unlock() error {
	return nil
}

// Truncate 截断文件
func (f *WritableAlistFile) Truncate(size int64) error {
	if f.closed {
		return os.ErrClosed
	}

	if size == 0 {
		f.buffer.Reset()
		return nil
	}

	// 简单实现：如果 size 小于当前大小，截断缓冲区
	if size < int64(f.buffer.Len()) {
		data := f.buffer.Bytes()[:size]
		f.buffer.Reset()
		f.buffer.Write(data)
	}

	return nil
}