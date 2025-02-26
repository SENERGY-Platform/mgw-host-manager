/*
 * Copyright 2023 InfAI (CC SES)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"errors"
	"fmt"
	sb_logger "github.com/SENERGY-Platform/go-service-base/logger"
	srv_info_hdl "github.com/SENERGY-Platform/go-service-base/srv-info-hdl"
	sb_util "github.com/SENERGY-Platform/go-service-base/util"
	"github.com/SENERGY-Platform/go-service-base/watchdog"
	"github.com/SENERGY-Platform/mgw-host-manager/handler/blacklist_hdl"
	"github.com/SENERGY-Platform/mgw-host-manager/handler/http_hdl"
	"github.com/SENERGY-Platform/mgw-host-manager/handler/info_hdl"
	"github.com/SENERGY-Platform/mgw-host-manager/handler/mdns_hdl"
	"github.com/SENERGY-Platform/mgw-host-manager/handler/resource_hdl"
	"github.com/SENERGY-Platform/mgw-host-manager/handler/resource_hdl/application_hdl"
	"github.com/SENERGY-Platform/mgw-host-manager/handler/resource_hdl/serial_hdl"
	lib_model "github.com/SENERGY-Platform/mgw-host-manager/lib/model"
	"github.com/SENERGY-Platform/mgw-host-manager/manager"
	"github.com/SENERGY-Platform/mgw-host-manager/util"
	"net/http"
	"os"
	"syscall"
	"time"
)

var version string

func main() {
	srvInfoHdl := srv_info_hdl.New("host-manager", version)

	ec := 0
	defer func() {
		os.Exit(ec)
	}()

	util.ParseFlags()

	config, err := util.NewConfig(util.Flags.ConfPath)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		ec = 1
		return
	}

	logFile, err := util.InitLogger(config.Logger)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		var logFileError *sb_logger.LogFileError
		if errors.As(err, &logFileError) {
			ec = 1
			return
		}
	}
	if logFile != nil {
		defer logFile.Close()
	}

	util.Logger.Printf("%s %s", srvInfoHdl.GetName(), srvInfoHdl.GetVersion())

	util.Logger.Debugf("config: %s", sb_util.ToJsonStr(config))

	watchdog.Logger = util.Logger
	wtchdg := watchdog.New(syscall.SIGINT, syscall.SIGTERM)

	netInterfaceBlacklistHdl, err := blacklist_hdl.New(config.Blacklist.NetInterfaceListPath)
	if err != nil {
		util.Logger.Error(err)
		ec = 1
		return
	}
	netInterfaceBlacklistHdl.SetValidationFunc(info_hdl.ValidateNetItfName)
	if err = netInterfaceBlacklistHdl.Init(); err != nil {
		util.Logger.Error(err)
		ec = 1
		return
	}

	netRangeBlacklistHdl, err := blacklist_hdl.New(config.Blacklist.NetRangeListPath)
	if err != nil {
		util.Logger.Error(err)
		ec = 1
		return
	}
	netRangeBlacklistHdl.SetValidationFunc(info_hdl.ValidateCIDR)
	if err = netRangeBlacklistHdl.Init(); err != nil {
		util.Logger.Error(err)
		ec = 1
		return
	}

	hostInfoHdl, err := info_hdl.New(config.Blacklist.NetInterfaceList, config.Blacklist.NetRangeList, netInterfaceBlacklistHdl, netRangeBlacklistHdl)
	if err != nil {
		util.Logger.Error(err)
		ec = 1
		return
	}

	hostAppHdl, err := application_hdl.New(config.ApplicationsPath, config.Blacklist.AppSocketList)
	if err != nil {
		util.Logger.Error(err)
		ec = 1
		return
	}
	if err = hostAppHdl.Init(); err != nil {
		util.Logger.Error(err)
		ec = 1
		return
	}

	hostResourceHdl := resource_hdl.New(map[lib_model.ResourceType]resource_hdl.ResHandler{
		lib_model.SerialDevice: serial_hdl.New(config.SerialDevicePath),
		lib_model.Application:  hostAppHdl,
	})
	util.Logger.Debugf("resource handlers: %s", sb_util.ToJsonStr(hostResourceHdl.Handlers()))

	hm := manager.New(hostInfoHdl, hostResourceHdl, hostAppHdl, netInterfaceBlacklistHdl, netRangeBlacklistHdl, mdns_hdl.New(), srvInfoHdl)

	httpHandler, err := http_hdl.New(hm, map[string]string{
		lib_model.HeaderApiVer:  srvInfoHdl.GetVersion(),
		lib_model.HeaderSrvName: srvInfoHdl.GetName(),
	})
	if err != nil {
		util.Logger.Error(err)
		ec = 1
		return
	}

	listener, err := sb_util.NewUnixListener(config.Socket.Path, os.Getuid(), config.Socket.GroupID, config.Socket.FileMode)
	if err != nil {
		util.Logger.Error(err)
		ec = 1
		return
	}
	server := &http.Server{Handler: httpHandler}
	srvCtx, srvCF := context.WithCancel(context.Background())
	wtchdg.RegisterStopFunc(func() error {
		if srvCtx.Err() == nil {
			ctxWt, cf := context.WithTimeout(context.Background(), time.Second*5)
			defer cf()
			if err := server.Shutdown(ctxWt); err != nil {
				return err
			}
			util.Logger.Info("http server shutdown complete")
		}
		return nil
	})
	wtchdg.RegisterHealthFunc(func() bool {
		if srvCtx.Err() == nil {
			return true
		}
		util.Logger.Error("http server closed unexpectedly")
		return false
	})

	wtchdg.Start()

	go func() {
		defer srvCF()
		util.Logger.Info("starting http server ...")
		if err := server.Serve(listener); !errors.Is(err, http.ErrServerClosed) {
			util.Logger.Error(err)
			ec = 1
			return
		}
	}()

	ec = wtchdg.Join()
}
