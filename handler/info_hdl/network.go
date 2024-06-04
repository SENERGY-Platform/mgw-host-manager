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

package info_hdl

import (
	"context"
	"github.com/SENERGY-Platform/mgw-host-manager/lib/model"
	"net"
	"os"
	"strings"
)

func (h *Handler) GetNet(ctx context.Context) (model.HostNet, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return model.HostNet{}, model.NewInternalError(err)
	}
	interfaces, err := h.getNetInterfaces(ctx)
	if err != nil {
		return model.HostNet{}, model.NewInternalError(err)
	}
	return model.HostNet{
		Hostname:   hostname,
		Interfaces: interfaces,
	}, nil
}

func (h *Handler) getNetInterfaces(ctx context.Context) ([]model.NetInterface, error) {
	addrMap := make(map[string][2]string)
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	for _, addr := range addrs {
		ip, n, err := net.ParseCIDR(addr.String())
		if err != nil {
			return nil, err
		}
		if n != nil {
			addrMap[ip.String()] = [2]string{addr.String(), net.IP(n.Mask).String()}
		}
	}
	ifs, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var interfaces []model.NetInterface
	for _, i := range ifs {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		if h.blacklistedInterface(i.Name) {
			continue
		}
		if i.Flags&net.FlagUp == 1 && i.Flags&net.FlagRunning != 0 && i.Flags&net.FlagLoopback == 0 {
			ip, err := getInterfaceAddr(i)
			if err != nil {
				return nil, err
			}
			if ip == nil {
				continue
			}
			values, ok := addrMap[ip.String()]
			if !ok {
				continue
			}
			interfaces = append(interfaces, model.NetInterface{
				Name:        i.Name,
				IPv4Address: ip.String(),
				IPv4NetMask: values[1],
				IPv4CIDR:    values[0],
			})
		}
	}
	return interfaces, nil
}

func (h *Handler) blacklistedInterface(name string) bool {
	for _, s := range h.netInterfaceBlacklist {
		if strings.Contains(name, s) {
			return true
		}
	}
	return false
}

func getInterfaceAddr(i net.Interface) (net.IP, error) {
	addrs, err := i.Addrs()
	if err != nil {
		return nil, err
	}
	for _, addr := range addrs {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}
		if ip == nil || ip.IsLoopback() {
			continue
		}
		ip = ip.To4()
		if ip == nil {
			continue
		}
		return ip, nil
	}
	return nil, nil
}
