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

type Config struct {
	Logger            LoggerConfig `json:"logger" env_var:"LOGGER_CONFIG"`
	Socket            SocketConfig `json:"socket" env_var:"SOCKET_CONFIG"`
	NetItfBlacklist   []string     `json:"net_itf_blacklist" env_var:"NET_ITF_BLACKLIST"`
	NetRngBlacklist   []string     `json:"net_rng_blacklist" env_var:"NET_RNG_BLACKLIST"`
	SerialDevicePath  string       `json:"serial_device_path" env_var:"SERIAL_DEVICE_PATH"`
	ApplicationsPath  string       `json:"applications_path" env_var:"APPLICATIONS_PATH"`
	AvahiServicesPath string       `json:"avahi_services_path" env_var:"AVAHI_SERVICES_PATH"`
	CoreID            string       `json:"core_id" env_var:"CORE_ID"`
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
		SerialDevicePath:  "/dev/serial/by-id",
		ApplicationsPath:  "./applications.json",
		AvahiServicesPath: "/etc/avahi/services",
	}
	err := config_hdl.Load(&cfg, nil, map[reflect.Type]envldr.Parser{reflect.TypeOf(level.Off): sb_logger.LevelParser}, nil, path)
	return &cfg, err
}
