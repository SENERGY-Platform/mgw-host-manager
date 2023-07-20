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

package model

type ServiceGroup struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	ReplaceWildcards bool      `json:"replace_wildcards"`
	Services         []Service `json:"services"`
}

type Service struct {
	Type       string      `json:"type"`
	Subtypes   []string    `json:"subtypes"`
	DomainName string      `json:"domain_name"`
	HostName   string      `json:"host_name"`
	Port       uint        `json:"port"`
	IPVer      string      `json:"ip_ver"`
	TxtRecords []TxtRecord `json:"txt_records"`
}

type TxtRecord struct {
	Format string `json:"format"`
	Value  string `json:"value"`
}
