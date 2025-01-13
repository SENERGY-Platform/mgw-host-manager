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

package manager

import (
	"context"
	"github.com/SENERGY-Platform/go-service-base/srv-info-hdl"
	srv_info_lib "github.com/SENERGY-Platform/go-service-base/srv-info-hdl/lib"
	lib_model "github.com/SENERGY-Platform/mgw-host-manager/lib/model"
	"time"
)

type Manager struct {
	hostInfoHdl        HostInfoHandler
	hostResourceHdl    HostResourceHandler
	hostAppHdl         HostApplicationHandler
	netItfBlacklistHdl BlacklistHandler
	netRngBlacklistHdl BlacklistHandler
	mdnsDiscoveryHdl   MDNSDiscoveryHandler
	srvInfoHdl         srv_info_hdl.SrvInfoHandler
}

func New(hostInfoHandler HostInfoHandler, hostResourceHandler HostResourceHandler, hostAppHdl HostApplicationHandler, netItfBlacklistHdl, netRngBlacklistHdl BlacklistHandler, mdnsDiscoveryHdl MDNSDiscoveryHandler, srvInfoHandler srv_info_hdl.SrvInfoHandler) *Manager {
	return &Manager{
		hostInfoHdl:        hostInfoHandler,
		hostResourceHdl:    hostResourceHandler,
		hostAppHdl:         hostAppHdl,
		netItfBlacklistHdl: netItfBlacklistHdl,
		netRngBlacklistHdl: netRngBlacklistHdl,
		mdnsDiscoveryHdl:   mdnsDiscoveryHdl,
		srvInfoHdl:         srvInfoHandler,
	}
}

func (m *Manager) GetHostInfo(ctx context.Context) (lib_model.HostInfo, error) {
	netInfo, err := m.hostInfoHdl.GetNet(ctx)
	if err != nil {
		return lib_model.HostInfo{}, err
	}
	return lib_model.HostInfo{
		Network: netInfo,
	}, nil
}

func (m *Manager) GetHostNet(ctx context.Context) (lib_model.HostNet, error) {
	netInfo, err := m.hostInfoHdl.GetNet(ctx)
	if err != nil {
		return lib_model.HostNet{}, err
	}
	return netInfo, nil
}

func (m *Manager) ListHostResources(ctx context.Context, filter lib_model.HostResourceFilter) ([]lib_model.HostResource, error) {
	return m.hostResourceHdl.List(ctx, filter)
}

func (m *Manager) GetHostResource(ctx context.Context, rID string) (lib_model.HostResource, error) {
	return m.hostResourceHdl.Get(ctx, rID)
}

func (m *Manager) ListHostApplications(ctx context.Context) ([]lib_model.HostApplication, error) {
	return m.hostAppHdl.List(ctx)
}

func (m *Manager) AddHostApplication(ctx context.Context, appResBase lib_model.HostApplicationBase) (string, error) {
	return m.hostAppHdl.Add(ctx, appResBase)
}

func (m *Manager) RemoveHostApplication(ctx context.Context, aID string) error {
	return m.hostAppHdl.Remove(ctx, aID)
}

func (m *Manager) GetNetItfBlacklist(ctx context.Context) ([]string, error) {
	return m.netItfBlacklistHdl.List(ctx)
}

func (m *Manager) NetItfBlacklistAdd(ctx context.Context, v string) error {
	return m.netItfBlacklistHdl.Add(ctx, v)
}

func (m *Manager) NetItfBlacklistRemove(ctx context.Context, v string) error {
	return m.netItfBlacklistHdl.Remove(ctx, v)
}

func (m *Manager) GetNetRngBlacklist(ctx context.Context) ([]string, error) {
	return m.netRngBlacklistHdl.List(ctx)
}

func (m *Manager) NetRngBlacklistAdd(ctx context.Context, v string) error {
	return m.netRngBlacklistHdl.Add(ctx, v)
}

func (m *Manager) NetRngBlacklistRemove(ctx context.Context, v string) error {
	return m.netRngBlacklistHdl.Remove(ctx, v)
}

func (m *Manager) MDNSQueryService(ctx context.Context, service, domain string, window time.Duration) ([]lib_model.MDNSEntry, error) {
	return m.mdnsDiscoveryHdl.Query(ctx, service, domain, window)
}

func (m *Manager) GetSrvInfo(_ context.Context) srv_info_lib.SrvInfo {
	return m.srvInfoHdl.GetInfo()
}
