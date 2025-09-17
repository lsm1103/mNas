package nfs

import (
	"fmt"
	"net"
	"sync"

	nfs "github.com/willscott/go-nfs"
	nfshelper "github.com/willscott/go-nfs/helpers"
	log "github.com/sirupsen/logrus"
)

// Server NFS 服务器实例
type Server struct {
	listener net.Listener
	handler  nfs.Handler
	address  string
	port     int
	running  bool
	mutex    sync.RWMutex
}

// Config NFS 服务器配置
type Config struct {
	Address string // 监听地址，默认 "0.0.0.0"
	Port    int    // 监听端口，默认 2049
	Enable  bool   // 是否启用 NFS 服务
}

// NewServer 创建新的 NFS 服务器
func NewServer(config Config) *Server {
	if config.Address == "" {
		config.Address = "0.0.0.0"
	}
	if config.Port == 0 {
		config.Port = 2049
	}

	return &Server{
		address: config.Address,
		port:    config.Port,
		running: false,
	}
}

// Start 启动 NFS 服务器
func (s *Server) Start() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.running {
		return fmt.Errorf("NFS 服务器已在运行")
	}

	// 创建监听器
	addr := fmt.Sprintf("%s:%d", s.address, s.port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("无法启动 NFS 服务器监听 %s: %v", addr, err)
	}

	s.listener = listener

	// 创建文件系统适配器
	filesystem := NewAlistFS("/")

	// 创建认证处理器（目前使用无认证）
	authHandler := nfshelper.NewNullAuthHandler(filesystem)

	// 添加缓存层以提高性能
	cacheHandler := nfshelper.NewCachingHandler(authHandler, 1024)

	s.handler = cacheHandler
	s.running = true

	log.Infof("NFS 服务器启动成功，监听地址: %s", listener.Addr())

	// 在 goroutine 中运行服务器
	go func() {
		err := nfs.Serve(listener, cacheHandler)
		if err != nil && s.running {
			log.Errorf("NFS 服务器运行错误: %v", err)
		}
	}()

	return nil
}

// Stop 停止 NFS 服务器
func (s *Server) Stop() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.running {
		return fmt.Errorf("NFS 服务器未在运行")
	}

	if s.listener != nil {
		err := s.listener.Close()
		if err != nil {
			log.Errorf("关闭 NFS 监听器时出错: %v", err)
		}
	}

	s.running = false
	log.Info("NFS 服务器已停止")

	return nil
}

// IsRunning 检查服务器是否正在运行
func (s *Server) IsRunning() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.running
}

// GetAddress 获取服务器监听地址
func (s *Server) GetAddress() string {
	if s.listener != nil {
		return s.listener.Addr().String()
	}
	return fmt.Sprintf("%s:%d", s.address, s.port)
}