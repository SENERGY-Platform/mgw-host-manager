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
	"github.com/SENERGY-Platform/mgw-host-manager/handler/mdns_hdl/avahi_adv_hdl/service_file"
	"strconv"
)

type Handler struct {
	srvPath string
}

func New(servicePath string) *Handler {
	return &Handler{
		srvPath: servicePath,
	}
}

func read(path string) (ServiceGroup, error) {
	xmlSG, err := service_file.Read(path)
	if err != nil {
		return ServiceGroup{}, err
	}
	var services []Service
	for _, xmlSrv := range xmlSG.Services {
		var subtypes []string
		for _, xmlSubtype := range xmlSrv.Subtypes {
			subtypes = append(subtypes, xmlSubtype.Value)
		}
		port, err := strconv.ParseInt(xmlSrv.Port, 10, 0)
		if err != nil {
			return ServiceGroup{}, err
		}
		var txtRecords []TxtRecord
		for _, xmlTxtRecord := range xmlSrv.TxtRecords {
			txtRecords = append(txtRecords, TxtRecord{
				Format: xmlTxtRecord.ValueFormat,
				Value:  xmlTxtRecord.Value,
			})
		}
		services = append(services, Service{
			IPVer:      xmlSrv.Protocol,
			Type:       xmlSrv.Type,
			Subtypes:   subtypes,
			DomainName: xmlSrv.DomainName,
			HostName:   xmlSrv.HostName,
			Port:       uint(port),
			TxtRecords: txtRecords,
		})
	}
	return ServiceGroup{
		Name:             xmlSG.Name.Value,
		ReplaceWildcards: xmlSG.Name.ReplaceWildcards == "yes",
		Services:         services,
	}, nil
}

func write(path string, sg ServiceGroup) error {
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
	return service_file.Write(xmlServiceGroup, path)
}
