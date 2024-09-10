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
	gin_mw "github.com/SENERGY-Platform/gin-middleware"
	sb_logger "github.com/SENERGY-Platform/go-service-base/logger"
	srv_info_hdl "github.com/SENERGY-Platform/go-service-base/srv-info-hdl"
	sb_util "github.com/SENERGY-Platform/go-service-base/util"
	"github.com/SENERGY-Platform/go-service-base/watchdog"
	"github.com/SENERGY-Platform/mgw-host-manager/api"
	"github.com/SENERGY-Platform/mgw-host-manager/handler/blacklist_hdl"
	"github.com/SENERGY-Platform/mgw-host-manager/handler/http_hdl"
	"github.com/SENERGY-Platform/mgw-host-manager/handler/info_hdl"
	"github.com/SENERGY-Platform/mgw-host-manager/handler/resource_hdl"
	"github.com/SENERGY-Platform/mgw-host-manager/handler/resource_hdl/application_hdl"
	"github.com/SENERGY-Platform/mgw-host-manager/handler/resource_hdl/serial_hdl"
	"github.com/SENERGY-Platform/mgw-host-manager/lib/model"
	"github.com/SENERGY-Platform/mgw-host-manager/util"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
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

	netInterfaceBlacklistHdl, err := blacklist_hdl.New(config.NetItfBlacklistPath)
	if err != nil {
		util.Logger.Error(err)
		ec = 1
		return
	}

	netRangeBlacklistHdl, err := blacklist_hdl.New(config.NetRngBlacklistPath)
	if err != nil {
		util.Logger.Error(err)
		ec = 1
		return
	}
	netRangeBlacklistHdl.SetValidationFunc(info_hdl.ValidateCIDR)

	hostInfoHdl, err := info_hdl.New(config.NetItfBlacklist, config.NetRngBlacklist, netInterfaceBlacklistHdl, netRangeBlacklistHdl)
	if err != nil {
		util.Logger.Error(err)
		ec = 1
		return
	}

	hostAppHdl, err := application_hdl.New(config.ApplicationsPath, config.DockerSocketPath)
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

	hostResourceHdl := resource_hdl.New(map[model.ResourceType]resource_hdl.ResHandler{
		model.SerialDevice: serial_hdl.New(config.SerialDevicePath),
		model.Application:  hostAppHdl,
	})
	util.Logger.Debugf("resource handlers: %s", sb_util.ToJsonStr(hostResourceHdl.Handlers()))

	mApi := api.New(hostInfoHdl, hostResourceHdl, hostAppHdl, netInterfaceBlacklistHdl, netRangeBlacklistHdl, srvInfoHdl)

	gin.SetMode(gin.ReleaseMode)
	httpHandler := gin.New()
	staticHeader := map[string]string{
		model.HeaderApiVer:  srvInfoHdl.GetVersion(),
		model.HeaderSrvName: srvInfoHdl.GetName(),
	}
	httpHandler.Use(gin_mw.StaticHeaderHandler(staticHeader), requestid.New(requestid.WithCustomHeaderStrKey(model.HeaderRequestID)), gin_mw.LoggerHandler(util.Logger, nil, func(gc *gin.Context) string {
		return requestid.Get(gc)
	}), gin_mw.ErrorHandler(util.GetStatusCode, ", "), gin.Recovery())
	httpHandler.UseRawPath = true

	http_hdl.SetRoutes(httpHandler, mApi)
	util.Logger.Debugf("routes: %s", sb_util.ToJsonStr(http_hdl.GetRoutes(httpHandler)))

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
