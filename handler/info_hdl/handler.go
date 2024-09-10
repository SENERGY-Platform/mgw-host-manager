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
	"github.com/SENERGY-Platform/mgw-host-manager/handler"
	"net"
)

type Handler struct {
	netInterfaceBlacklist    []string
	netInterfaceBlacklistHdl handler.BlacklistHandler
	netRangeBlacklist        []*net.IPNet
	netRangeBlacklistHdl     handler.BlacklistHandler
}

func New(netInterfaceBlacklist, netRangeBlacklist []string, netInterfaceBlacklistHdl, netRangeBlacklistHdl handler.BlacklistHandler) (*Handler, error) {
	ipNets, err := genIPNets(netRangeBlacklist)
	if err != nil {
		return nil, err
	}
	return &Handler{
		netInterfaceBlacklist:    netInterfaceBlacklist,
		netInterfaceBlacklistHdl: netInterfaceBlacklistHdl,
		netRangeBlacklist:        ipNets,
		netRangeBlacklistHdl:     netRangeBlacklistHdl,
	}, nil
}
