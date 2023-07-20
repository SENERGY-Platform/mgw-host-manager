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

const Header = "<?xml version=\"1.0\" standalone='no'?><!--*-nxml-*-->\n<!DOCTYPE service-group SYSTEM \"avahi-service.dtd\">\n"

const (
	IPvAny IPVer = "any"
	IPv4   IPVer = "ipv4"
	IPv6   IPVer = "ipv6"
)

const (
	ProtocolTcp Protocol = "tcp"
	ProtocolUdp Protocol = "udp"
)

const (
	ValueFormatText         ValueFormat = "text"
	ValueFormatBinaryHex    ValueFormat = "binary-hex"
	ValueFormatBinaryBase64 ValueFormat = "binary-base64"
)
