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

package standard

import (
	"github.com/SENERGY-Platform/mgw-host-manager/lib"
	lib_model "github.com/SENERGY-Platform/mgw-host-manager/lib/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
)

// GetHostApplicationsH godoc
// @Summary List applications
// @Description	List host applications.
// @Tags Host Applications
// @Produce	json
// @Success	200 {array} lib_model.HostApplication "host applications"
// @Failure	500 {string} string "error message"
// @Router /applications [get]
func GetHostApplicationsH(a lib.Api) (string, string, gin.HandlerFunc) {
	return http.MethodGet, lib_model.HostAppsPath, func(gc *gin.Context) {
		apps, err := a.ListHostApplications(gc.Request.Context())
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.JSON(http.StatusOK, apps)
	}
}

// PostHostApplicationH godoc
// @Summary Add application
// @Description	Add a host application.
// @Tags Host Applications
// @Accept json
// @Produce	plain
// @Param endpoint body lib_model.HostApplicationBase true "application information"
// @Success	200 {string} string "ID"
// @Failure	400 {string} string "error message"
// @Failure	500 {string} string "error message"
// @Router /applications [post]
func PostHostApplicationH(a lib.Api) (string, string, gin.HandlerFunc) {
	return http.MethodPost, lib_model.HostAppsPath, func(gc *gin.Context) {
		var app lib_model.HostApplicationBase
		err := gc.ShouldBindJSON(&app)
		if err != nil {
			_ = gc.Error(lib_model.NewInvalidInputError(err))
			return
		}
		id, err := a.AddHostApplication(gc.Request.Context(), app)
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.String(http.StatusOK, id)
	}
}

// DeleteHostApplicationH godoc
// @Summary Delete application
// @Description	Remove a host application.
// @Tags Host Applications
// @Param id path string true "resource id"
// @Success	200
// @Failure	404 {string} string "error message"
// @Failure	500 {string} string "error message"
// @Router /applications/{id} [delete]
func DeleteHostApplicationH(a lib.Api) (string, string, gin.HandlerFunc) {
	return http.MethodDelete, path.Join(lib_model.HostAppsPath, ":id"), func(gc *gin.Context) {
		err := a.RemoveHostApplication(gc.Request.Context(), gc.Param("id"))
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.Status(http.StatusOK)
	}
}
