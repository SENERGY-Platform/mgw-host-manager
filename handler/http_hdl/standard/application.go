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
