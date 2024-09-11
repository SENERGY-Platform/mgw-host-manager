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
	"sort"
)

func SetRoutes(e *gin.Engine, a lib.Api) {
	standardGrp := e.Group("")
	restrictedGrp := e.Group(model.RestrictedPath)
	setSrvInfoRoutes(a, standardGrp, restrictedGrp)
	setHostInfoRoutes(a, standardGrp, restrictedGrp)
	setHostResourcesRoutes(a, standardGrp, restrictedGrp)
	setHostApplicationsRoutes(a, standardGrp)
	setBlacklistRoutes(a, standardGrp.Group(model.BlacklistsPath))
}

func GetRoutes(e *gin.Engine) [][2]string {
	routes := e.Routes()
	sort.Slice(routes, func(i, j int) bool {
		return routes[i].Path < routes[j].Path
	})
	var rInfo [][2]string
	for _, info := range routes {
		rInfo = append(rInfo, [2]string{info.Method, info.Path})
	}
	return rInfo
}

func setHostApplicationsRoutes(a lib.Api, rg *gin.RouterGroup) {
	rg.GET(model.HostAppsPath, getHostApplicationsH(a))
	rg.POST(model.HostAppsPath, postHostApplicationH(a))
	rg.DELETE(model.HostAppsPath+"/:"+hostAppIdParam, deleteHostApplicationH(a))
}

func setHostResourcesRoutes(a lib.Api, rGroups ...*gin.RouterGroup) {
	for _, rg := range rGroups {
		rg.GET(model.HostResourcesPath, getHostResourcesH(a))
		rg.GET(model.HostResourcesPath+"/:"+hostResIdParam, getHostResourceH(a))
	}
}

func setHostInfoRoutes(a lib.Api, rGroups ...*gin.RouterGroup) {
	for _, rg := range rGroups {
		rg.GET(model.HostInfoPath, getHostInfoH(a))
		rg.GET(model.HostInfoPath+"/"+model.HostNetPath, getHostNetH(a))
	}
}

func setBlacklistRoutes(a lib.Api, rg *gin.RouterGroup) {
	rg.GET(model.NetInterfaces, getNetItfBlacklistH(a))
	rg.POST(model.NetInterfaces, postNetItfBlacklistValueH(a))
	rg.DELETE(model.NetInterfaces, deleteNetItfBlacklistValueH(a))
	rg.GET(model.NetRanges, getNetRngBlacklistH(a))
	rg.POST(model.NetRanges, postNetRngBlacklistValueH(a))
	rg.DELETE(model.NetRanges, deleteNetRngBlacklistValueH(a))
}

func setSrvInfoRoutes(a lib.Api, rGroups ...*gin.RouterGroup) {
	for _, rg := range rGroups {
		rg.GET(model.SrvInfoPath, getSrvInfoH(a))
	}
}
