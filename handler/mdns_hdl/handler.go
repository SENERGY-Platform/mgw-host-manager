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

package mdns_hdl

import (
	"context"
	"errors"
	"github.com/SENERGY-Platform/mgw-host-manager/lib/model"
	"github.com/libp2p/zeroconf/v2"
	"strings"
	"time"
)

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Query(ctx context.Context, service, domain string, window time.Duration) ([]model.MDNSEntry, error) {
	var entries []model.MDNSEntry
	results := make(chan *zeroconf.ServiceEntry)
	go func() {
		for result := range results {
			entries = append(entries, newMDNSEntry(result))
		}
	}()
	ctxWt, cancel := context.WithTimeout(ctx, window)
	defer cancel()
	var err error
	if err = zeroconf.Browse(ctxWt, service, domain, results, zeroconf.SelectIPTraffic(zeroconf.IPv4)); err == nil {
		err = ctxWt.Err()
	}
	if err != nil && !errors.Is(err, context.DeadlineExceeded) {
		return nil, err
	}
	return entries, nil
}

func newMDNSEntry(se *zeroconf.ServiceEntry) model.MDNSEntry {
	var IPv4Addr string
	if len(se.AddrIPv4) > 0 {
		IPv4Addr = se.AddrIPv4[0].String()
	}
	hostname := se.HostName
	hostnameParts := strings.Split(hostname, ".")
	if len(hostnameParts) > 0 {
		hostname = hostnameParts[0]
	}
	return model.MDNSEntry{
		Name:       se.Instance,
		Type:       se.Service,
		Subtypes:   se.Subtypes,
		Domain:     se.Domain,
		Hostname:   hostname,
		Port:       se.Port,
		IPv4Addr:   IPv4Addr,
		TxtRecords: se.Text,
		Expiry:     se.Expiry,
	}
}
