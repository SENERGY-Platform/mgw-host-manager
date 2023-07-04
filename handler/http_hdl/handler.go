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

package http_hdl

import (
	"github.com/SENERGY-Platform/mgw-host-manager/lib"
	"github.com/SENERGY-Platform/mgw-host-manager/lib/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

const resIdParam = "r"

type resourcesQuery struct {
}

func getHostInfo(a lib.Api) gin.HandlerFunc {
	return func(gc *gin.Context) {
		hostInfo, err := a.GetHostInfo(gc.Request.Context())
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.JSON(http.StatusOK, hostInfo)
	}
}

func getHostNet(a lib.Api) gin.HandlerFunc {
	return func(gc *gin.Context) {
		hostNet, err := a.GetHostNet(gc.Request.Context())
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.JSON(http.StatusOK, hostNet)
	}
}

func getResources(a lib.Api) gin.HandlerFunc {
	return func(gc *gin.Context) {
		query := resourcesQuery{}
		if err := gc.ShouldBindQuery(&query); err != nil {
			_ = gc.Error(model.NewInvalidInputError(err))
			return
		}
		resources, err := a.GetResources(gc.Request.Context(), model.ResourceFilter{})
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.JSON(http.StatusOK, resources)
	}
}

func getResource(a lib.Api) gin.HandlerFunc {
	return func(gc *gin.Context) {
		resource, err := a.GetResource(gc.Request.Context(), gc.Param(resIdParam))
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.JSON(http.StatusOK, resource)
	}
}
