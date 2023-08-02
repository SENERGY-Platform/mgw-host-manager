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
	"errors"
	"fmt"
	gin_mw "github.com/SENERGY-Platform/gin-middleware"
	"github.com/SENERGY-Platform/go-service-base/srv-base"
	srv_base_types "github.com/SENERGY-Platform/go-service-base/srv-base/types"
	"github.com/SENERGY-Platform/mgw-host-manager/api"
	"github.com/SENERGY-Platform/mgw-host-manager/handler/http_hdl"
	"github.com/SENERGY-Platform/mgw-host-manager/handler/info_hdl"
	"github.com/SENERGY-Platform/mgw-host-manager/handler/mdns_hdl/avahi_adv_hdl"
	"github.com/SENERGY-Platform/mgw-host-manager/handler/resource_hdl"
	"github.com/SENERGY-Platform/mgw-host-manager/handler/resource_hdl/application_hdl"
	"github.com/SENERGY-Platform/mgw-host-manager/handler/resource_hdl/serial_hdl"
	"github.com/SENERGY-Platform/mgw-host-manager/lib/model"
	"github.com/SENERGY-Platform/mgw-host-manager/util"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

var version string

func main() {
	srv_base.PrintInfo(model.ServiceName, version)

	util.ParseFlags()

	config, err := util.NewConfig(util.Flags.ConfPath)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	logFile, err := util.InitLogger(config.Logger)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		var logFileError *srv_base.LogFileError
		if errors.As(err, &logFileError) {
			os.Exit(1)
		}
	}
	if logFile != nil {
		defer logFile.Close()
	}

	util.Logger.Debugf("config: %s", srv_base.ToJsonStr(config))

	hostInfoHdl := info_hdl.New(config.NetItfBlacklist)

	resourceHandlers := make(map[model.ResourceType]resource_hdl.ResHandler)

	resourceHandlers[model.SerialDevice] = serial_hdl.New(config.SerialDevicePath)

	apps, err := application_hdl.LoadApps(config.ApplicationsPath)
	if err != nil {
		var pathErr *os.PathError
		if !errors.As(err, &pathErr) {
			util.Logger.Fatal(err)
		}
	} else {
		if len(apps) > 0 {
			resourceHandlers[model.Application] = application_hdl.New(apps)
		}
	}

	hostResourceHdl := resource_hdl.New(resourceHandlers)
	util.Logger.Debugf("resource handlers: %s", srv_base.ToJsonStr(hostResourceHdl.Handlers()))

	mDNSAdvHdl := avahi_adv_hdl.New(config.AvahiServicesPath)
	err = mDNSAdvHdl.Init()
	if err != nil {
		util.Logger.Fatal(err)
	}

	mApi := api.New(hostInfoHdl, hostResourceHdl, mDNSAdvHdl)

	gin.SetMode(gin.ReleaseMode)
	httpHandler := gin.New()
	staticHeader := map[string]string{
		model.HeaderApiVer:  version,
		model.HeaderSrvName: model.ServiceName,
	}
	httpHandler.Use(gin_mw.StaticHeaderHandler(staticHeader), requestid.New(requestid.WithCustomHeaderStrKey(model.HeaderRequestID)), gin_mw.LoggerHandler(util.Logger, func(gc *gin.Context) string {
		return requestid.Get(gc)
	}), gin_mw.ErrorHandler(http_hdl.GetStatusCode, ", "), gin.Recovery())
	httpHandler.UseRawPath = true

	http_hdl.SetRoutes(httpHandler, mApi)
	util.Logger.Debugf("routes: %s", srv_base.ToJsonStr(http_hdl.GetRoutes(httpHandler)))

	listener, err := srv_base.NewUnixListener(config.Socket.Path, os.Getuid(), config.Socket.GroupID, config.Socket.FileMode)
	if err != nil {
		util.Logger.Fatal(err)
	}

	srv_base.StartServer(&http.Server{Handler: httpHandler}, listener, srv_base_types.DefaultShutdownSignals, util.Logger)
}
