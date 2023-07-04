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

package application_hdl

import (
	"context"
	"encoding/json"
	"github.com/SENERGY-Platform/mgw-host-manager/lib/model"
	"github.com/SENERGY-Platform/mgw-host-manager/util"
	"os"
)

type Handler struct {
	resources map[string]model.ResourceMeta
}

type Application struct {
	Name   string `json:"name"`
	Socket string `json:"socket"`
}

func New(applications []Application) *Handler {
	resources := make(map[string]model.ResourceMeta)
	for _, app := range applications {
		resources[util.GenHash(app.Socket)] = model.ResourceMeta{
			Name: app.Name,
			Tags: nil,
			Path: app.Socket,
		}
	}
	return &Handler{
		resources: resources,
	}
}

func LoadApps(path string) ([]Application, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	var d []Application
	if err = decoder.Decode(&d); err != nil {
		return nil, err
	}
	return d, nil
}

func (h *Handler) Get(ctx context.Context) (map[string]model.ResourceMeta, error) {
	resources := make(map[string]model.ResourceMeta)
	for id, res := range h.resources {
		if ctx.Err() != nil {
			return nil, model.NewInternalError(ctx.Err())
		}
		resources[id] = res
	}
	return resources, nil
}
