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
	"errors"
	"fmt"
	"github.com/SENERGY-Platform/mgw-host-manager/lib/model"
	"github.com/SENERGY-Platform/mgw-host-manager/util"
	"github.com/google/uuid"
	"io"
	"os"
	"path"
	"sync"
)

type Handler struct {
	apps map[string]model.HostApplication
	path string
	mu   sync.RWMutex
}

func New(p string) (*Handler, error) {
	if !path.IsAbs(p) {
		return nil, fmt.Errorf("path '%s' not absolute", p)
	}
	return &Handler{
		path: p,
		apps: make(map[string]model.HostApplication),
	}, nil
}

func (h *Handler) Init() error {
	apps, err := readStoFile(h.path)
	if err != nil {
		var jutErr *json.UnmarshalTypeError
		switch {
		case errors.Is(err, os.ErrNotExist):
			return nil
		case errors.As(err, &jutErr):
			apps, err = migrateStoFile(h.path)
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					return nil
				}
				return err
			}
		default:
			return err
		}
	}
	h.apps = apps
	return nil
}

func (h *Handler) List(_ context.Context) ([]model.HostApplication, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	apps := make([]model.HostApplication, len(h.apps))
	for _, app := range h.apps {
		apps = append(apps, app)
	}
	return apps, nil
}

func (h *Handler) Add(_ context.Context, appResBase model.HostApplicationBase) (string, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	idObj, err := uuid.NewUUID()
	if err != nil {
		return "", model.NewInternalError(err)
	}
	id := idObj.String()
	h.apps[id] = model.HostApplication{
		ID:                  id,
		HostApplicationBase: appResBase,
	}
	if err := writeStoFile(h.apps, h.path); err != nil {
		delete(h.apps, id)
		return "", err
	}
	return id, nil
}

func (h *Handler) Remove(_ context.Context, id string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.apps[id]; !ok {
		return model.NewNotFoundError(fmt.Errorf("application '%s' does not exist", id))
	}
	newApps := make(map[string]model.HostApplication)
	for i, app := range h.apps {
		if i != id {
			newApps[i] = app
		}
	}
	if err := writeStoFile(newApps, h.path); err != nil {
		return err
	}
	h.apps = newApps
	return nil
}

func (h *Handler) Get(_ context.Context) (map[string]model.HostResourceBase, error) {
	resources := make(map[string]model.HostResourceBase)
	for id, app := range h.apps {
		resources[id] = model.HostResourceBase{
			Name: app.Name,
			Path: app.Socket,
		}
	}
	return resources, nil
}

func migrateStoFile(p string) (map[string]model.HostApplication, error) {
	if err := copyFile(p, p+".migration_bk"); err != nil {
		return nil, err
	}
	file, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	var oldFmt []model.HostApplicationBase
	if err = decoder.Decode(&oldFmt); err != nil {
		return nil, err
	}
	newFmt := make(map[string]model.HostApplication)
	for _, app := range oldFmt {
		newFmt[util.GenHash(app.Socket)] = model.HostApplication{HostApplicationBase: app}
	}
	if err = writeStoFile(newFmt, p); err != nil {
		return nil, err
	}
	return newFmt, nil
}

func readStoFile(p string) (map[string]model.HostApplication, error) {
	file, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	var apps map[string]model.HostApplication
	if err = decoder.Decode(&apps); err != nil {
		return nil, err
	}
	return apps, nil
}

func writeStoFile(apps map[string]model.HostApplication, p string) error {
	if err := copyFile(p, p+".bk"); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}
	file, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	if err = json.NewEncoder(file).Encode(apps); err != nil {
		e := copyFile(p+".bk", p)
		if e != nil && !errors.Is(e, os.ErrNotExist) {
			util.Logger.Error(e)
		}
		return err
	}
	return nil
}

func copyFile(src, dst string) error {
	sFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sFile.Close()
	dFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dFile.Close()
	_, err = io.Copy(dFile, sFile)
	return err
}
