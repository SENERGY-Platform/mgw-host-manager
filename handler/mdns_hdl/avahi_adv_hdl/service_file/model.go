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

import "encoding/xml"

type ServiceGroup struct {
	XMLName  xml.Name `xml:"service-group"`
	Name     Name
	Services []Service `xml:"service"`
}

type Name struct {
	XMLName          xml.Name `xml:"name"`
	ReplaceWildcards string   `xml:"replace-wildcards,attr"` // (yes|no) default: "no"
	Value            string   `xml:",chardata"`
}

type Service struct {
	XMLName    xml.Name    `xml:"service"`
	Protocol   string      `xml:"protocol,attr"` // (ipv4|ipv6|any) default: "any"
	Type       string      `xml:"type"`
	Subtypes   []Subtype   `xml:"subtype,omitempty"`
	DomainName string      `xml:"domain-name,omitempty"`
	HostName   string      `xml:"host-name,omitempty"`
	Port       string      `xml:"port"`
	TxtRecords []TxtRecord `xml:"txt-record,omitempty"`
}

type Subtype struct {
	XMLName xml.Name `xml:"subtype"`
	Value   string   `xml:",chardata"`
}

type TxtRecord struct {
	XMLName     xml.Name `xml:"txt-record"`
	ValueFormat string   `xml:"value-format,attr"` // (text|binary-hex|binary-base64) default: "text"
	Value       string   `xml:",chardata"`
}

type IPVer = string

type ValueFormat = string

type Protocol = string
