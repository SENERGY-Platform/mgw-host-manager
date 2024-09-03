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

const (
	hostResIdParam = "r"
	mDNSAdvParam   = "a"
)

type hostResourcesQuery struct {
}

func getHostInfoH(a lib.Api) gin.HandlerFunc {
	return func(gc *gin.Context) {
		hostInfo, err := a.GetHostInfo(gc.Request.Context())
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.JSON(http.StatusOK, hostInfo)
	}
}

func getHostNetH(a lib.Api) gin.HandlerFunc {
	return func(gc *gin.Context) {
		hostNet, err := a.GetHostNet(gc.Request.Context())
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.JSON(http.StatusOK, hostNet)
	}
}

func getHostResourcesH(a lib.Api) gin.HandlerFunc {
	return func(gc *gin.Context) {
		query := hostResourcesQuery{}
		if err := gc.ShouldBindQuery(&query); err != nil {
			_ = gc.Error(model.NewInvalidInputError(err))
			return
		}
		resources, err := a.ListHostResources(gc.Request.Context(), model.HostResourceFilter{})
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.JSON(http.StatusOK, resources)
	}
}

func getHostResourceH(a lib.Api) gin.HandlerFunc {
	return func(gc *gin.Context) {
		resource, err := a.GetHostResource(gc.Request.Context(), gc.Param(hostResIdParam))
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.JSON(http.StatusOK, resource)
	}
}

func getSrvInfoH(a lib.Api) gin.HandlerFunc {
	return func(gc *gin.Context) {
		gc.JSON(http.StatusOK, a.GetSrvInfo(gc.Request.Context()))
	}
}
