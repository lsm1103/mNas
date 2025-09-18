package server

import (
	"github.com/alist-org/alist/v3/internal/conf"
	"github.com/alist-org/alist/v3/server/nfs"
	log "github.com/sirupsen/logrus"
)

var nfsServer *nfs.Server

// StartNFS 启动 NFS 服务器
func StartNFS() error {
	if !conf.Conf.NFS.Enable {
		log.Debug("NFS 服务未启用")
		return nil
	}

	config := nfs.Config{
		Address: conf.Conf.NFS.Address,
		Port:    conf.Conf.NFS.Port,
		Enable:  conf.Conf.NFS.Enable,
	}

	nfsServer = nfs.NewServer(config)

	err := nfsServer.Start()
	if err != nil {
		log.Errorf("启动 NFS 服务器失败: %v", err)
		return err
	}

	log.Infof("NFS 服务器已启动，监听地址: %s", nfsServer.GetAddress())
	return nil
}

// StopNFS 停止 NFS 服务器
func StopNFS() error {
	if nfsServer != nil && nfsServer.IsRunning() {
		err := nfsServer.Stop()
		if err != nil {
			log.Errorf("停止 NFS 服务器失败: %v", err)
			return err
		}
		log.Info("NFS 服务器已停止")
	}
	return nil
}