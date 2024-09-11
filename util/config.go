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

package util

import (
	"github.com/SENERGY-Platform/go-service-base/config-hdl"
	sb_logger "github.com/SENERGY-Platform/go-service-base/logger"
	envldr "github.com/y-du/go-env-loader"
	"github.com/y-du/go-log-level/level"
	"io/fs"
	"os"
	"reflect"
)

type SocketConfig struct {
	Path     string      `json:"path" env_var:"SOCKET_PATH"`
	GroupID  int         `json:"group_id" env_var:"SOCKET_GROUP_ID"`
	FileMode fs.FileMode `json:"file_mode" env_var:"SOCKET_FILE_MODE"`
}

type LoggerConfig struct {
	Level        level.Level `json:"level" env_var:"LOGGER_LEVEL"`
	Utc          bool        `json:"utc" env_var:"LOGGER_UTC"`
	Path         string      `json:"path" env_var:"LOGGER_PATH"`
	FileName     string      `json:"file_name" env_var:"LOGGER_FILE_NAME"`
	Terminal     bool        `json:"terminal" env_var:"LOGGER_TERMINAL"`
	Microseconds bool        `json:"microseconds" env_var:"LOGGER_MICROSECONDS"`
	Prefix       string      `json:"prefix" env_var:"LOGGER_PREFIX"`
}

type BlacklistConfig struct {
	NetInterfaceList     []string `json:"net_interface_list" env_var:"BLACKLIST_NET_INTERFACE_LIST"`
	NetRangeList         []string `json:"net_range_list" env_var:"BLACKLIST_NET_RANGE_LIST"`
	AppSocketList        []string `json:"app_socket_list" env_var:"BLACKLIST_APP_SOCKET_LIST"`
	NetInterfaceListPath string   `json:"net_interface_list_path" env_var:"BLACKLIST_NET_INTERFACE_LIST_PATH"`
	NetRangeListPath     string   `json:"net_range_list_path" env_var:"BLACKLIST_NET_RANGE_LIST_PATH"`
}

type Config struct {
	Logger           LoggerConfig    `json:"logger" env_var:"LOGGER_CONFIG"`
	Socket           SocketConfig    `json:"socket" env_var:"SOCKET_CONFIG"`
	Blacklist        BlacklistConfig `json:"blacklist" env_var:"BLACKLIST_CONFIG"`
	SerialDevicePath string          `json:"serial_device_path" env_var:"SERIAL_DEVICE_PATH"`
	ApplicationsPath string          `json:"applications_path" env_var:"APPLICATIONS_PATH"`
	CoreID           string          `json:"core_id" env_var:"CORE_ID"`
}

func NewConfig(path string) (*Config, error) {
	cfg := Config{
		Logger: LoggerConfig{
			Level:        level.Warning,
			Utc:          true,
			Path:         "./",
			FileName:     "mgw_host_manager",
			Microseconds: true,
		},
		Socket: SocketConfig{
			Path:     "./h_manager.sock",
			GroupID:  os.Getgid(),
			FileMode: 0660,
		},
		SerialDevicePath: "/dev/serial/by-id",
	}
	err := config_hdl.Load(&cfg, nil, map[reflect.Type]envldr.Parser{reflect.TypeOf(level.Off): sb_logger.LevelParser}, nil, path)
	return &cfg, err
}
