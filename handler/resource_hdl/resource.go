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

package resource_hdl

import (
	"context"
	"errors"
	"fmt"
	"github.com/SENERGY-Platform/mgw-host-manager/lib/model"
	"strings"
)

type Handler struct {
	handlers map[model.ResourceType]ResHandler
}

func New(handlers map[model.ResourceType]ResHandler) *Handler {
	return &Handler{
		handlers: handlers,
	}
}

func (h *Handler) List(ctx context.Context, filter model.ResourceFilter) ([]model.Resource, error) {
	var resources []model.Resource
	for t, handler := range h.handlers {
		res, err := handler.Get(ctx)
		if err != nil {
			return nil, err
		}
		for id, meta := range res {
			resources = append(resources, model.Resource{
				ID:           genID(t, id),
				ResourceMeta: meta,
			})
		}
	}
	return resources, nil
}

func (h *Handler) Get(ctx context.Context, rID string) (model.Resource, error) {
	t, id, err := parseID(rID)
	if err != nil {
		return model.Resource{}, model.NewInvalidInputError(err)
	}
	handler, ok := h.handlers[t]
	if !ok {
		return model.Resource{}, model.NewInvalidInputError(fmt.Errorf("unknown resource type '%s'", t))
	}
	res, err := handler.Get(ctx)
	if err != nil {
		return model.Resource{}, err
	}
	meta, ok := res[id]
	if !ok {
		return model.Resource{}, model.NewNotFoundError(fmt.Errorf("resource '%s' not found", rID))
	}
	return model.Resource{
		ID:           rID,
		ResourceMeta: meta,
	}, nil
}

func (h *Handler) Handlers() []string {
	var handlers []string
	for t := range h.handlers {
		handlers = append(handlers, t)
	}
	return handlers
}

func genID(t model.ResourceType, id string) string {
	return t + ":" + id
}

func parseID(rID string) (model.ResourceType, string, error) {
	parts := strings.Split(rID, ":")
	if len(parts) != 2 {
		return "", "", errors.New("invalid ID format")
	}
	return parts[0], parts[1], nil
}
