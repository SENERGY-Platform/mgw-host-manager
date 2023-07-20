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
	"github.com/SENERGY-Platform/mgw-host-manager/handler"
	"github.com/SENERGY-Platform/mgw-host-manager/lib/model"
)

type Api struct {
	hostInfoHdl handler.HostInfoHandler
	resourceHdl handler.ResourceHandler
	mDNSAdvHdl  handler.MDNSAdvHandler
}

func New(hostInfoHandler handler.HostInfoHandler, resourceHandler handler.ResourceHandler, mDNSAdvHdl handler.MDNSAdvHandler) *Api {
	return &Api{
		hostInfoHdl: hostInfoHandler,
		resourceHdl: resourceHandler,
		mDNSAdvHdl:  mDNSAdvHdl,
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

func (a *Api) ListHostResources(ctx context.Context, filter model.ResourceFilter) ([]model.Resource, error) {
	return a.resourceHdl.List(ctx, filter)
}

func (a *Api) GetHostResource(ctx context.Context, rID string) (model.Resource, error) {
	return a.resourceHdl.Get(ctx, rID)
}
