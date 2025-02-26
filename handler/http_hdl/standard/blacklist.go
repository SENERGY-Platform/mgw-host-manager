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
	"io"
	"net/http"
	"path"
)

type deleteBlacklistValQuery struct {
	Value string `form:"value"`
}

// GetNetItfBlacklistH godoc
// @Summary List network interfaces
// @Description	List blacklisted host network interfaces.
// @Tags Blacklists
// @Produce	json
// @Success	200 {array} string "network interfaces"
// @Failure	500 {string} string "error message"
// @Router /blacklists/net-interfaces [get]
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

// PostNetItfBlacklistValueH godoc
// @Summary Add network interface
// @Description	Add a host network interface to the list.
// @Tags Blacklists
// @Accept plain
// @Param value body string true "interface name"
// @Success	200
// @Failure	400 {string} string "error message"
// @Failure	500 {string} string "error message"
// @Router /blacklists/net-interfaces [post]
func PostNetItfBlacklistValueH(a lib.Api) (string, string, gin.HandlerFunc) {
	return http.MethodPost, path.Join(lib_model.BlacklistsPath, lib_model.NetInterfacesPath), func(gc *gin.Context) {
		defer gc.Request.Body.Close()
		v, err := io.ReadAll(gc.Request.Body)
		if err != nil {
			_ = gc.Error(lib_model.NewInvalidInputError(err))
			return
		}
		err = a.NetItfBlacklistAdd(gc.Request.Context(), string(v))
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.Status(http.StatusOK)
	}
}

// DeleteNetItfBlacklistValueH godoc
// @Summary Delete network interface
// @Description	Remove a host network interface from the list.
// @Tags Blacklists
// @Param value query string true "interface name"
// @Success	200
// @Failure	400 {string} string "error message"
// @Failure	404 {string} string "error message"
// @Failure	500 {string} string "error message"
// @Router /blacklists/net-interfaces [delete]
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

// GetNetRngBlacklistH godoc
// @Summary List network ranges
// @Description	List blacklisted network ranges.
// @Tags Blacklists
// @Produce	json
// @Success	200 {array} string "network rages"
// @Failure	500 {string} string "error message"
// @Router /blacklists/net-ranges [get]
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

// PostNetRngBlacklistValueH godoc
// @Summary Add network range
// @Description	Add a network range to the list.
// @Tags Blacklists
// @Accept plain
// @Param value body string true "network range in CIDR notation"
// @Success	200
// @Failure	400 {string} string "error message"
// @Failure	500 {string} string "error message"
// @Router /blacklists/net-ranges [post]
func PostNetRngBlacklistValueH(a lib.Api) (string, string, gin.HandlerFunc) {
	return http.MethodPost, path.Join(lib_model.BlacklistsPath, lib_model.NetRangesPath), func(gc *gin.Context) {
		defer gc.Request.Body.Close()
		v, err := io.ReadAll(gc.Request.Body)
		if err != nil {
			_ = gc.Error(lib_model.NewInvalidInputError(err))
			return
		}
		if err != nil {
			_ = gc.Error(lib_model.NewInvalidInputError(err))
			return
		}
		err = a.NetRngBlacklistAdd(gc.Request.Context(), string(v))
		if err != nil {
			_ = gc.Error(err)
			return
		}
		gc.Status(http.StatusOK)
	}
}

// DeleteNetRngBlacklistValueH godoc
// @Summary Delete network range
// @Description	Remove a network range from the list.
// @Tags Blacklists
// @Param value query string true "network range"
// @Success	200
// @Failure	400 {string} string "error message"
// @Failure	404 {string} string "error message"
// @Failure	500 {string} string "error message"
// @Router /blacklists/net-ranges [delete]
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
