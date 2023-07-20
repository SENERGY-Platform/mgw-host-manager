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

package serial_hdl

import (
	"context"
	"errors"
	"github.com/SENERGY-Platform/mgw-host-manager/lib/model"
	"github.com/SENERGY-Platform/mgw-host-manager/util"
	"io/fs"
	"os"
)

type Handler struct {
	path string
}

func New(path string) *Handler {
	return &Handler{
		path: path,
	}
}

func (h *Handler) Get(ctx context.Context) (map[string]model.HostResourceBase, error) {
	dir := os.DirFS(h.path)
	entries, err := fs.ReadDir(dir, ".")
	if err != nil {
		var pathErr *os.PathError
		if errors.As(err, &pathErr) && pathErr.Op == "open" {
			return nil, nil
		}
		return nil, model.NewInternalError(err)
	}
	resources := make(map[string]model.HostResourceBase)
	for _, entry := range entries {
		if ctx.Err() != nil {
			return nil, model.NewInternalError(ctx.Err())
		}
		if !entry.IsDir() {
			resources[util.GenHash(entry.Name())] = model.HostResourceBase{
				Name: entry.Name(),
				Tags: nil,
				Path: h.path + "/" + entry.Name(),
			}
		}
	}
	return resources, nil
}
