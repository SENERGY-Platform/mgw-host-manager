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

package avahi_adv_hdl

import (
	"context"
	"fmt"
	"github.com/SENERGY-Platform/mgw-host-manager/handler/mdns_hdl/avahi_adv_hdl/service_file"
	"github.com/SENERGY-Platform/mgw-host-manager/lib/model"
	"github.com/SENERGY-Platform/mgw-host-manager/util"
	"io/fs"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

const extension = "service"

type Handler struct {
	srvPath   string
	srvGroups map[string]model.ServiceGroup
	mu        sync.RWMutex
}

func New(servicePath string) *Handler {
	return &Handler{
		srvPath: servicePath,
	}
}

func (h *Handler) Init() error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if !path.IsAbs(h.srvPath) {
		return fmt.Errorf("path must be absolute")
	}
	dirEntries, err := fs.ReadDir(os.DirFS(h.srvPath), ".")
	if err != nil {
		return err
	}
	h.srvGroups = make(map[string]model.ServiceGroup)
	for _, entry := range dirEntries {
		if !entry.IsDir() {
			sg, err := read(path.Join(h.srvPath, entry.Name()))
			if err != nil {
				util.Logger.Error(err)
				continue
			}
			h.srvGroups[sg.ID] = sg
		}
	}
	return nil
}

func (h *Handler) List(ctx context.Context) ([]model.ServiceGroup, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	var serviceGroups []model.ServiceGroup
	for _, sg := range h.srvGroups {
		if ctx.Err() != nil {
			return nil, model.NewInternalError(ctx.Err())
		}
		serviceGroups = append(serviceGroups, sg)
	}
	return serviceGroups, nil
}

func (h *Handler) Add(_ context.Context, sg model.ServiceGroup) error {
	re, err := regexp.Compile(`(?m)^[a-z0-9-_]+$`)
	if err != nil {
		return model.NewInternalError(err)
	}
	if !re.MatchString(sg.ID) {
		return model.NewInvalidInputError(fmt.Errorf("invalid id format '%s'", sg.ID))
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.srvGroups[sg.ID]; ok {
		return model.NewInvalidInputError(fmt.Errorf("id '%s' already exisits", sg.ID))
	}
	return h.add(sg)
}

func (h *Handler) Get(_ context.Context, id string) (model.ServiceGroup, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	sg, ok := h.srvGroups[id]
	if !ok {
		return model.ServiceGroup{}, model.NewNotFoundError(fmt.Errorf("service group '%s' not found", id))
	}
	return sg, nil
}

func (h *Handler) Update(_ context.Context, sg model.ServiceGroup) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.srvGroups[sg.ID]; !ok {
		return model.NewNotFoundError(fmt.Errorf("service-group '%s' not found", sg.ID))
	}
	return h.add(sg)
}

func (h *Handler) Delete(_ context.Context, id string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.srvGroups[id]; !ok {
		return model.NewNotFoundError(fmt.Errorf("service-group '%s' not found", id))
	}
	err := os.Remove(path.Join(h.srvPath, getFilename(id)))
	if err != nil {
		return model.NewInternalError(err)
	}
	delete(h.srvGroups, id)
	return nil
}

func (h *Handler) add(sg model.ServiceGroup) error {
	err := write(path.Join(h.srvPath, getFilename(sg.ID)), sg)
	if err != nil {
		return err
	}
	h.srvGroups[sg.ID] = sg
	return nil
}

func read(pth string) (model.ServiceGroup, error) {
	xmlSG, err := service_file.Read(pth)
	if err != nil {
		return model.ServiceGroup{}, err
	}
	var services []model.Service
	for _, xmlSrv := range xmlSG.Services {
		var subtypes []string
		for _, xmlSubtype := range xmlSrv.Subtypes {
			subtypes = append(subtypes, xmlSubtype.Value)
		}
		port, err := strconv.ParseInt(xmlSrv.Port, 10, 0)
		if err != nil {
			return model.ServiceGroup{}, err
		}
		var txtRecords []model.TxtRecord
		for _, xmlTxtRecord := range xmlSrv.TxtRecords {
			txtRecords = append(txtRecords, model.TxtRecord{
				Format: xmlTxtRecord.ValueFormat,
				Value:  xmlTxtRecord.Value,
			})
		}
		services = append(services, model.Service{
			IPVer:      xmlSrv.Protocol,
			Type:       xmlSrv.Type,
			Subtypes:   subtypes,
			DomainName: xmlSrv.DomainName,
			HostName:   xmlSrv.HostName,
			Port:       uint(port),
			TxtRecords: txtRecords,
		})
	}
	return model.ServiceGroup{
		ID:               getID(path.Base(pth)),
		Name:             xmlSG.Name.Value,
		ReplaceWildcards: xmlSG.Name.ReplaceWildcards == "yes",
		Services:         services,
	}, nil
}

func write(path string, sg model.ServiceGroup) error {
	var xmlServices []service_file.Service
	for _, service := range sg.Services {
		var xmlSubtypes []service_file.Subtype
		for _, subtype := range service.Subtypes {
			xmlSubtype, err := service_file.NewSubtype(subtype)
			if err != nil {
				return err
			}
			xmlSubtypes = append(xmlSubtypes, xmlSubtype)
		}
		var xmlTxtRecords []service_file.TxtRecord
		for _, txtRecord := range service.TxtRecords {
			xmlTxtRecord, err := service_file.NewTxtRecord(txtRecord.Value, txtRecord.Format)
			if err != nil {
				return err
			}
			xmlTxtRecords = append(xmlTxtRecords, xmlTxtRecord)
		}
		xmlService, err := service_file.NewService(service.Type, xmlSubtypes, service.Port, service.IPVer, service.DomainName, service.HostName, xmlTxtRecords)
		if err != nil {
			return err
		}
		xmlServices = append(xmlServices, xmlService)
	}
	xmlServiceGroup, err := service_file.NewServiceGroup(sg.Name, sg.ReplaceWildcards, xmlServices)
	if err != nil {
		return err
	}
	err = service_file.Write(xmlServiceGroup, path)
	if err != nil {
		return model.NewInternalError(err)
	}
	return nil
}

func getFilename(id string) string {
	return id + "." + extension
}

func getID(fileName string) string {
	return strings.TrimRight(fileName, "."+extension)
}
