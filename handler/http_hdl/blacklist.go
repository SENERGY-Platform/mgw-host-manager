/*
 * Copyright 2024 InfAI (CC SES)
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

type deleteBlacklistValQuery struct {
	Value string `form:"value"`
}

func getNetItfBlacklistH(a lib.Api) gin.HandlerFunc {
	return func(gc *gin.Context) {
		values, err := a.GetNetItfBlacklist(gc.Request.Context())
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.JSON(http.StatusOK, values)
	}
}

func postNetItfBlacklistValueH(a lib.Api) gin.HandlerFunc {
	return func(gc *gin.Context) {
		var v string
		err := gc.ShouldBindJSON(&v)
		if err != nil {
			_ = gc.Error(model.NewInvalidInputError(err))
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

func deleteNetItfBlacklistValueH(a lib.Api) gin.HandlerFunc {
	return func(gc *gin.Context) {
		query := deleteBlacklistValQuery{}
		if err := gc.ShouldBindQuery(&query); err != nil {
			_ = gc.Error(model.NewInvalidInputError(err))
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

func getNetRngBlacklistH(a lib.Api) gin.HandlerFunc {
	return func(gc *gin.Context) {
		values, err := a.GetNetRngBlacklist(gc.Request.Context())
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.JSON(http.StatusOK, values)
	}
}

func postNetRngBlacklistValueH(a lib.Api) gin.HandlerFunc {
	return func(gc *gin.Context) {
		var v string
		err := gc.ShouldBindJSON(&v)
		if err != nil {
			_ = gc.Error(model.NewInvalidInputError(err))
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

func deleteNetRngBlacklistValueH(a lib.Api) gin.HandlerFunc {
	return func(gc *gin.Context) {
		query := deleteBlacklistValQuery{}
		if err := gc.ShouldBindQuery(&query); err != nil {
			_ = gc.Error(model.NewInvalidInputError(err))
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
