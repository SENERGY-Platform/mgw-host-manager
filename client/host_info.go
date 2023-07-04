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

package client

import (
	"context"
	"github.com/SENERGY-Platform/mgw-host-manager/lib/model"
	"net/http"
	"net/url"
)

func (c *Client) GetHostInfo(ctx context.Context) (model.HostInfo, error) {
	u, err := url.JoinPath(c.baseUrl, model.HostInfoPath)
	if err != nil {
		return model.HostInfo{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return model.HostInfo{}, err
	}
	var hostInfo model.HostInfo
	err = c.execRequestJSON(req, &hostInfo)
	if err != nil {
		return model.HostInfo{}, err
	}
	return hostInfo, nil
}

func (c *Client) GetHostNet(ctx context.Context) (model.HostNet, error) {
	u, err := url.JoinPath(c.baseUrl, model.HostInfoPath, model.HostNetPath)
	if err != nil {
		return model.HostNet{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return model.HostNet{}, err
	}
	var hostNet model.HostNet
	err = c.execRequestJSON(req, &hostNet)
	if err != nil {
		return model.HostNet{}, err
	}
	return hostNet, nil
}
