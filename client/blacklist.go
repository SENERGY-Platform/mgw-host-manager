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
	"fmt"
	"github.com/SENERGY-Platform/mgw-host-manager/lib/model"
	"net/http"
	"net/url"
)

func (c *Client) GetNetItfBlacklist(ctx context.Context) ([]string, error) {
	u, err := url.JoinPath(c.baseUrl, model.BlacklistsPath, model.NetInterfaces)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	var values []string
	err = c.baseClient.ExecRequestJSON(req, &values)
	if err != nil {
		return nil, err
	}
	return values, nil
}

func (c *Client) NetItfBlacklistAdd(ctx context.Context, v string) error {
	u, err := url.JoinPath(c.baseUrl, model.BlacklistsPath, model.NetInterfaces)
	if err != nil {
		return err
	}
	body, err := json.Marshal(v)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	return c.baseClient.ExecRequestVoid(req)
}

func (c *Client) NetItfBlacklistRemove(ctx context.Context, v string) error {
	u, err := url.JoinPath(c.baseUrl, model.BlacklistsPath, model.NetInterfaces)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("%s?value=%s", u, v), nil)
	if err != nil {
		return err
	}
	return c.baseClient.ExecRequestVoid(req)
}

func (c *Client) GetNetRngBlacklist(ctx context.Context) ([]string, error) {
	u, err := url.JoinPath(c.baseUrl, model.BlacklistsPath, model.NetRanges)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	var values []string
	err = c.baseClient.ExecRequestJSON(req, &values)
	if err != nil {
		return nil, err
	}
	return values, nil
}

func (c *Client) NetRngBlacklistAdd(ctx context.Context, v string) error {
	u, err := url.JoinPath(c.baseUrl, model.BlacklistsPath, model.NetRanges)
	if err != nil {
		return err
	}
	body, err := json.Marshal(v)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	return c.baseClient.ExecRequestVoid(req)
}

func (c *Client) NetRngBlacklistRemove(ctx context.Context, v string) error {
	u, err := url.JoinPath(c.baseUrl, model.BlacklistsPath, model.NetRanges)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("%s?value=%s", u, v), nil)
	if err != nil {
		return err
	}
	return c.baseClient.ExecRequestVoid(req)
}
