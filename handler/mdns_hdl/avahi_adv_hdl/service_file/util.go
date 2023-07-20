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

package service_file

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/SENERGY-Platform/mgw-host-manager/lib/model"
	"os"
	"regexp"
	"strconv"
)

func NewServiceGroup(name string, nameReplaceWildcards bool, services []Service) (ServiceGroup, error) {
	re, err := regexp.Compile(`^[a-zA-Z0-9\.\-_ #()\[\]!"\$%\?='\*\+:,\|@~]+$`)
	if err != nil {
		return ServiceGroup{}, model.NewInternalError(err)
	}
	if !re.MatchString(name) {
		return ServiceGroup{}, model.NewInvalidInputError(fmt.Errorf("invalid name format '%s'", name))
	}
	replaceWildcards := "no"
	if nameReplaceWildcards {
		replaceWildcards = "yes"
	}
	if len(services) == 0 {
		return ServiceGroup{}, model.NewInvalidInputError(errors.New("no services defined"))
	}
	return ServiceGroup{
		Name: Name{
			ReplaceWildcards: replaceWildcards,
			Value:            name,
		},
		Services: services,
	}, nil
}

func NewService(sType string, subTypes []Subtype, port uint, ipVer IPVer, domainName, hostName string, txtRecords []TxtRecord) (Service, error) {
	re, err := regexp.Compile(`^(?:_[a-z0-9-]+\.)(?:_tcp|_udp)$`)
	if err != nil {
		return Service{}, model.NewInternalError(err)
	}
	if !re.MatchString(sType) {
		return Service{}, model.NewInvalidInputError(fmt.Errorf("invalid type format '%s'", sType))
	}
	proto := IPvAny
	if ipVer != "" {
		if !(ipVer == IPvAny || ipVer == IPv4 || ipVer == IPv6) {
			return Service{}, model.NewInvalidInputError(fmt.Errorf("invalid ip version '%s'", ipVer))
		}
		proto = ipVer
	}
	return Service{
		Protocol:   proto,
		Type:       sType,
		Subtypes:   subTypes,
		DomainName: domainName,
		HostName:   hostName,
		Port:       strconv.FormatInt(int64(port), 10),
		TxtRecords: txtRecords,
	}, nil
}

func NewSubtype(sType string) (Subtype, error) {
	re, err := regexp.Compile(`^(?:_[a-z0-9-]+\._sub\.)?(?:_[a-z0-9-]+\.)(?:_tcp|_udp)$`)
	if err != nil {
		return Subtype{}, model.NewInternalError(err)
	}
	if !re.MatchString(sType) {
		return Subtype{}, model.NewInvalidInputError(fmt.Errorf("invalid subtype format '%s'", sType))
	}
	return Subtype{Value: sType}, nil
}

func NewTxtRecord(value string, valueFormat ValueFormat) (TxtRecord, error) {
	valFormat := ValueFormatText
	if valueFormat != "" {
		if !(valueFormat == ValueFormatText || valueFormat == ValueFormatBinaryHex || valueFormat == ValueFormatBinaryBase64) {
			return TxtRecord{}, model.NewInvalidInputError(fmt.Errorf("invalid value format '%s'", valueFormat))
		}
		valFormat = valueFormat
	}
	return TxtRecord{
		ValueFormat: valFormat,
		Value:       value,
	}, nil
}

func Write(sg ServiceGroup, path string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(Header)
	if err != nil {
		return err
	}
	encoder := xml.NewEncoder(file)
	defer encoder.Close()
	encoder.Indent("", "  ")
	err = encoder.Encode(sg)
	if err != nil {
		return err
	}
	return nil
}

func Read(path string) (ServiceGroup, error) {
	file, err := os.Open(path)
	if err != nil {
		return ServiceGroup{}, err
	}
	defer file.Close()
	decoder := xml.NewDecoder(file)
	var sg ServiceGroup
	err = decoder.Decode(&sg)
	if err != nil {
		return ServiceGroup{}, err
	}
	return sg, nil
}
