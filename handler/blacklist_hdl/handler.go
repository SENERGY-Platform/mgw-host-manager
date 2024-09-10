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

package blacklist_hdl

import (
	"context"
	"errors"
	"fmt"
	"github.com/SENERGY-Platform/mgw-host-manager/lib/model"
	"github.com/SENERGY-Platform/mgw-host-manager/util/json_sto_file"
	"os"
	"path"
	"sync"
)

type Handler struct {
	values map[string]struct{}
	path   string
	mu     sync.RWMutex
}

func New(p string) (*Handler, error) {
	if !path.IsAbs(p) {
		return nil, fmt.Errorf("path '%s' not absolute", p)
	}
	return &Handler{
		values: make(map[string]struct{}),
		path:   p,
	}, nil
}

func (h *Handler) Init() error {
	var values []string
	if err := json_sto_file.Read(h.path, &values); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	for _, value := range values {
		h.values[value] = struct{}{}
	}
	return nil
}

func (h *Handler) List(_ context.Context) ([]string, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	var values []string
	for val := range h.values {
		values = append(values, val)
	}
	return values, nil
}

func (h *Handler) Add(_ context.Context, v string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.values[v]; ok {
		return model.NewInvalidInputError(fmt.Errorf("value '%s' already in list", v))
	}
	h.values[v] = struct{}{}
	if err := json_sto_file.Write(h.values, h.path, true); err != nil {
		delete(h.values, v)
		return model.NewInternalError(err)
	}
	return nil
}

func (h *Handler) Remove(_ context.Context, v string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.values[v]; !ok {
		return model.NewNotFoundError(fmt.Errorf("value '%s' not in list", v))
	}
	newValues := make(map[string]struct{})
	for val := range h.values {
		if val != v {
			newValues[val] = struct{}{}
		}
	}
	if err := json_sto_file.Write(newValues, h.path, true); err != nil {
		return model.NewInternalError(err)
	}
	h.values = newValues
	return nil
}

func (h *Handler) Has(_ context.Context, v string) (ok bool, err error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok = h.values[v]
	return
}
