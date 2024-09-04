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

package lib

import (
	"context"
	srv_info_lib "github.com/SENERGY-Platform/go-service-base/srv-info-hdl/lib"
	"github.com/SENERGY-Platform/mgw-host-manager/lib/model"
)

type Api interface {
	GetHostInfo(ctx context.Context) (model.HostInfo, error)
	GetHostNet(ctx context.Context) (model.HostNet, error)
	ListHostResources(ctx context.Context, filter model.HostResourceFilter) ([]model.HostResource, error)
	GetHostResource(ctx context.Context, rID string) (model.HostResource, error)
	ListHostApplications(ctx context.Context) ([]model.HostApplication, error)
	AddHostApplication(ctx context.Context, appResBase model.HostApplicationBase) (string, error)
	RemoveHostApplication(ctx context.Context, aID string) error
	srv_info_lib.Api
}
