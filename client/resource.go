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

func (c *Client) GetResources(ctx context.Context, filter model.ResourceFilter) ([]model.Resource, error) {
	u, err := url.JoinPath(c.baseUrl, model.ResourcesPath)
	if err != nil {
		return nil, err
	}
	u += genGetResourcesQuery(filter)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	var resources []model.Resource
	err = c.execRequestJSON(req, &resources)
	if err != nil {
		return nil, err
	}
	return resources, nil
}

func (c *Client) GetResource(ctx context.Context, id string) (model.Resource, error) {
	u, err := url.JoinPath(c.baseUrl, model.ResourcesPath, id)
	if err != nil {
		return model.Resource{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return model.Resource{}, err
	}
	var resource model.Resource
	err = c.execRequestJSON(req, &resource)
	if err != nil {
		return model.Resource{}, err
	}
	return resource, nil
}

func genGetResourcesQuery(filter model.ResourceFilter) string {
	return ""
}
