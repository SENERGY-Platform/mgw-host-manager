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
	"context"
	"fmt"
	"github.com/SENERGY-Platform/mgw-host-manager/lib/model"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func (c *Client) MDNSQueryService(ctx context.Context, service, domain string, timeWindow time.Duration) ([]model.MDNSEntry, error) {
	u, err := url.JoinPath(c.baseUrl, model.MDNSDiscoveryPath)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u+genMDNSQuery(service, domain, timeWindow), nil)
	if err != nil {
		return nil, err
	}
	var entries []model.MDNSEntry
	err = c.baseClient.ExecRequestJSON(req, &entries)
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func genMDNSQuery(srv, dom string, tw time.Duration) string {
	var items []string
	if srv != "" {
		items = append(items, fmt.Sprintf("service=%s", srv))
	}
	if dom != "" {
		items = append(items, fmt.Sprintf("domain=%s", dom))
	}
	if tw > 0 {
		items = append(items, fmt.Sprintf("time_window=%d", tw.Nanoseconds()))
	}
	if len(items) > 0 {
		return "?" + strings.Join(items, "&")
	}
	return ""
}
