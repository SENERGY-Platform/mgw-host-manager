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

package api

import (
	"context"
	"github.com/SENERGY-Platform/go-service-base/srv-info-hdl"
	srv_info_lib "github.com/SENERGY-Platform/go-service-base/srv-info-hdl/lib"
	"github.com/SENERGY-Platform/mgw-host-manager/handler"
	"github.com/SENERGY-Platform/mgw-host-manager/lib/model"
)

type Api struct {
	hostInfoHdl     handler.HostInfoHandler
	hostResourceHdl handler.HostResourceHandler
	hostAppHdl      handler.HostApplicationHandler
	srvInfoHdl      srv_info_hdl.SrvInfoHandler
}

func New(hostInfoHandler handler.HostInfoHandler, hostResourceHandler handler.HostResourceHandler, hostAppHdl handler.HostApplicationHandler, srvInfoHandler srv_info_hdl.SrvInfoHandler) *Api {
	return &Api{
		hostInfoHdl:     hostInfoHandler,
		hostResourceHdl: hostResourceHandler,
		hostAppHdl:      hostAppHdl,
		srvInfoHdl:      srvInfoHandler,
	}
}

func (a *Api) GetHostInfo(ctx context.Context) (model.HostInfo, error) {
	netInfo, err := a.hostInfoHdl.GetNet(ctx)
	if err != nil {
		return model.HostInfo{}, err
	}
	return model.HostInfo{
		Network: netInfo,
	}, nil
}

func (a *Api) GetHostNet(ctx context.Context) (model.HostNet, error) {
	netInfo, err := a.hostInfoHdl.GetNet(ctx)
	if err != nil {
		return model.HostNet{}, err
	}
	return netInfo, nil
}

func (a *Api) ListHostResources(ctx context.Context, filter model.HostResourceFilter) ([]model.HostResource, error) {
	return a.hostResourceHdl.List(ctx, filter)
}

func (a *Api) GetHostResource(ctx context.Context, rID string) (model.HostResource, error) {
	return a.hostResourceHdl.Get(ctx, rID)
}

func (a *Api) GetSrvInfo(_ context.Context) srv_info_lib.SrvInfo {
	return a.srvInfoHdl.GetInfo()
}
