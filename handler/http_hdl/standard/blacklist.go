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

type deleteBlacklistValQuery struct {
	Value string `form:"value"`
}

func GetNetItfBlacklistH(a lib.Api) (string, string, gin.HandlerFunc) {
	return http.MethodGet, path.Join(lib_model.BlacklistsPath, lib_model.NetInterfacesPath), func(gc *gin.Context) {
		values, err := a.GetNetItfBlacklist(gc.Request.Context())
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.JSON(http.StatusOK, values)
	}
}

func PostNetItfBlacklistValueH(a lib.Api) (string, string, gin.HandlerFunc) {
	return http.MethodPost, path.Join(lib_model.BlacklistsPath, lib_model.NetInterfacesPath), func(gc *gin.Context) {
		var v string
		err := gc.ShouldBindJSON(&v)
		if err != nil {
			_ = gc.Error(lib_model.NewInvalidInputError(err))
			return
		}
		err = a.NetItfBlacklistAdd(gc.Request.Context(), v)
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.Status(http.StatusOK)
	}
}

func DeleteNetItfBlacklistValueH(a lib.Api) (string, string, gin.HandlerFunc) {
	return http.MethodDelete, path.Join(lib_model.BlacklistsPath, lib_model.NetInterfacesPath), func(gc *gin.Context) {
		query := deleteBlacklistValQuery{}
		if err := gc.ShouldBindQuery(&query); err != nil {
			_ = gc.Error(lib_model.NewInvalidInputError(err))
			return
		}
		err := a.NetItfBlacklistRemove(gc.Request.Context(), query.Value)
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.Status(http.StatusOK)
	}
}

func GetNetRngBlacklistH(a lib.Api) (string, string, gin.HandlerFunc) {
	return http.MethodGet, path.Join(lib_model.BlacklistsPath, lib_model.NetRangesPath), func(gc *gin.Context) {
		values, err := a.GetNetRngBlacklist(gc.Request.Context())
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.JSON(http.StatusOK, values)
	}
}

func PostNetRngBlacklistValueH(a lib.Api) (string, string, gin.HandlerFunc) {
	return http.MethodPost, path.Join(lib_model.BlacklistsPath, lib_model.NetRangesPath), func(gc *gin.Context) {
		var v string
		err := gc.ShouldBindJSON(&v)
		if err != nil {
			_ = gc.Error(lib_model.NewInvalidInputError(err))
			return
		}
		err = a.NetRngBlacklistAdd(gc.Request.Context(), v)
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.Status(http.StatusOK)
	}
}

func DeleteNetRngBlacklistValueH(a lib.Api) (string, string, gin.HandlerFunc) {
	return http.MethodDelete, path.Join(lib_model.BlacklistsPath, lib_model.NetRangesPath), func(gc *gin.Context) {
		query := deleteBlacklistValQuery{}
		if err := gc.ShouldBindQuery(&query); err != nil {
			_ = gc.Error(lib_model.NewInvalidInputError(err))
			return
		}
		err := a.NetRngBlacklistRemove(gc.Request.Context(), query.Value)
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.Status(http.StatusOK)
	}
}
