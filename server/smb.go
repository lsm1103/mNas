package server

import (
	"github.com/alist-org/alist/v3/internal/conf"
	log "github.com/sirupsen/logrus"
)

// StartSMB 启动 SMB 服务器
func StartSMB() error {
	if !conf.Conf.SMB.Enable {
		log.Debug("SMB 服务未启用")
		return nil
	}

	// TODO: 实现 SMB 服务器
	log.Info("SMB 服务器功能尚未实现")
	return nil
}

// StopSMB 停止 SMB 服务器
func StopSMB() error {
	// TODO: 实现 SMB 服务器停止
	log.Debug("SMB 服务器停止功能尚未实现")
	return nil
}