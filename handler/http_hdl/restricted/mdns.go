/*
 * Copyright 2025 InfAI (CC SES)
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

package restricted

import (
	"github.com/SENERGY-Platform/mgw-host-manager/lib"
	lib_model "github.com/SENERGY-Platform/mgw-host-manager/lib/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type mdnsQuery struct {
	Service    string `form:"service"`
	Domain     string `form:"domain"`
	TimeWindow int64  `form:"time_window"`
}

// GetMDNSQueryH godoc
// @Summary MDNS query
// @Description	Query MDNS devices on attached networks.
// @Tags MDNS
// @Produce	json
// @Param service query string true "MDNS service string (e.g.: '_services._dns-sd._udp' for all available services)"
// @Param domain query string false "limit the query to a domain"
// @Param time_window query int false "set the maximum duration for the query (defaults to 1s)"
// @Success	200 {array} lib_model.MDNSEntry "list of services"
// @Failure	400 {string} string "error message"
// @Failure	500 {string} string "error message"
// @Router /endpoints [get]
func GetMDNSQueryH(a lib.Api) (string, string, gin.HandlerFunc) {
	return http.MethodGet, lib_model.MDNSDiscoveryPath, func(gc *gin.Context) {
		var query mdnsQuery
		if err := gc.ShouldBindQuery(&query); err != nil {
			_ = gc.Error(lib_model.NewInvalidInputError(err))
			return
		}
		if query.TimeWindow == 0 {
			query.TimeWindow = int64(time.Second)
		}
		results, err := a.MDNSQueryService(gc.Request.Context(), query.Service, query.Domain, time.Duration(query.TimeWindow))
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.JSON(http.StatusOK, results)
	}
}
