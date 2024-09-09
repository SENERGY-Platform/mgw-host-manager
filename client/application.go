/*
 * Copyright 2024 InfAI (CC SES)
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

func (c *Client) ListHostApplications(ctx context.Context) ([]model.HostApplication, error) {
	u, err := url.JoinPath(c.baseUrl, model.HostAppsPath)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	var apps []model.HostApplication
	err = c.baseClient.ExecRequestJSON(req, &apps)
	if err != nil {
		return nil, err
	}
	return apps, nil
}

func (c *Client) AddHostApplication(ctx context.Context, appResBase model.HostApplicationBase) (string, error) {
	u, err := url.JoinPath(c.baseUrl, model.HostAppsPath)
	if err != nil {
		return "", err
	}
	body, err := json.Marshal(appResBase)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	return c.baseClient.ExecRequestString(req)
}

func (c *Client) RemoveHostApplication(ctx context.Context, aID string) error {
	u, err := url.JoinPath(c.baseUrl, model.HostAppsPath, aID)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return err
	}
	return c.baseClient.ExecRequestVoid(req)
}