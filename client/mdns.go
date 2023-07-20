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
	"bytes"
	"context"
	"encoding/json"
	"github.com/SENERGY-Platform/mgw-host-manager/lib/model"
	"net/http"
	"net/url"
)

func (c *Client) ListMDNSAdv(ctx context.Context, filter model.ServiceGroupFilter) ([]model.ServiceGroup, error) {
	u, err := url.JoinPath(c.baseUrl, model.MDNSAdvPath)
	if err != nil {
		return nil, err
	}
	u += genServiceGroupQuery(filter)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	var serviceGroups []model.ServiceGroup
	err = c.execRequestJSON(req, &serviceGroups)
	if err != nil {
		return nil, err
	}
	return serviceGroups, nil
}

func (c *Client) AddMDNSAdv(ctx context.Context, serviceGroup model.ServiceGroup) error {
	u, err := url.JoinPath(c.baseUrl, model.MDNSAdvPath)
	if err != nil {
		return err
	}
	body, err := json.Marshal(serviceGroup)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, u, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	return c.execRequestVoid(req)
}

func (c *Client) GetMDNSAdv(ctx context.Context, id string) (model.ServiceGroup, error) {
	u, err := url.JoinPath(c.baseUrl, model.MDNSAdvPath, id)
	if err != nil {
		return model.ServiceGroup{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return model.ServiceGroup{}, err
	}
	var serviceGroup model.ServiceGroup
	err = c.execRequestJSON(req, &serviceGroup)
	if err != nil {
		return model.ServiceGroup{}, err
	}
	return serviceGroup, nil
}

func (c *Client) UpdateMDNSAdv(ctx context.Context, serviceGroup model.ServiceGroup) error {
	u, err := url.JoinPath(c.baseUrl, model.MDNSAdvPath, serviceGroup.ID)
	if err != nil {
		return err
	}
	body, err := json.Marshal(serviceGroup)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, u, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	return c.execRequestVoid(req)
}

func (c *Client) DeleteMDNSAdv(ctx context.Context, id string) error {
	u, err := url.JoinPath(c.baseUrl, model.MDNSAdvPath, id)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return err
	}
	return c.execRequestVoid(req)
}

func genServiceGroupQuery(filter model.ServiceGroupFilter) string {
	return ""
}
