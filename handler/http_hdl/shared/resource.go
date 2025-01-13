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

type hostResourcesQuery struct {
}

// GetHostResourcesH godoc
// @Summary List resources
// @Description	List host resources like application sockets or serial adapters.
// @Tags Host Resources
// @Produce	json
// @Success	200 {array} lib_model.HostResource "host resources"
// @Failure	500 {string} string "error message"
// @Router /host-resources [get]
func GetHostResourcesH(a lib.Api) (string, string, gin.HandlerFunc) {
	return http.MethodGet, lib_model.HostResourcesPath, func(gc *gin.Context) {
		query := hostResourcesQuery{}
		if err := gc.ShouldBindQuery(&query); err != nil {
			_ = gc.Error(lib_model.NewInvalidInputError(err))
			return
		}
		resources, err := a.ListHostResources(gc.Request.Context(), lib_model.HostResourceFilter{})
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.JSON(http.StatusOK, resources)
	}
}

// GetHostResourceH godoc
// @Summary Get resource
// @Description	Get a host resource.
// @Tags Host Resources
// @Produce	json
// @Param id path string true "resource id"
// @Success	200 {object} lib_model.HostResource "host resources"
// @Failure	500 {string} string "error message"
// @Router /host-resources/{id} [get]
func GetHostResourceH(a lib.Api) (string, string, gin.HandlerFunc) {
	return http.MethodGet, path.Join(lib_model.HostResourcesPath, ":id"), func(gc *gin.Context) {
		resource, err := a.GetHostResource(gc.Request.Context(), gc.Param("id"))
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.JSON(http.StatusOK, resource)
	}
}
