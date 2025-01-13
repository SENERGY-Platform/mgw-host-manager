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

package shared

import (
	"github.com/SENERGY-Platform/mgw-host-manager/lib"
	lib_model "github.com/SENERGY-Platform/mgw-host-manager/lib/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
)

// GetHostInfoH godoc
// @Summary Get all
// @Description	Get host information.
// @Tags Host Information
// @Produce	json
// @Success	200 {object} lib_model.HostInfo "host info"
// @Failure	500 {string} string "error message"
// @Router /host-info [get]
func GetHostInfoH(a lib.Api) (string, string, gin.HandlerFunc) {
	return http.MethodGet, lib_model.HostInfoPath, func(gc *gin.Context) {
		hostInfo, err := a.GetHostInfo(gc.Request.Context())
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.JSON(http.StatusOK, hostInfo)
	}
}

// GetHostNetH godoc
// @Summary Get network
// @Description	Get host network information.
// @Tags Host Information
// @Produce	json
// @Success	200 {object} lib_model.HostNet "host network info"
// @Failure	500 {string} string "error message"
// @Router /host-info/network [get]
func GetHostNetH(a lib.Api) (string, string, gin.HandlerFunc) {
	return http.MethodGet, path.Join(lib_model.HostInfoPath, lib_model.HostNetPath), func(gc *gin.Context) {
		hostNet, err := a.GetHostNet(gc.Request.Context())
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.JSON(http.StatusOK, hostNet)
	}
}
