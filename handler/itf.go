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

package handler

import (
	"context"
	"github.com/SENERGY-Platform/mgw-host-manager/lib/model"
)

type HostInfoHandler interface {
	GetNet(ctx context.Context) (model.HostNet, error)
	GetCPU(ctx context.Context) error
	GetRAM(ctx context.Context) error
	GetOS(ctx context.Context) error
}

type HostResourceHandler interface {
	List(ctx context.Context, filter model.HostResourceFilter) ([]model.HostResource, error)
	Get(ctx context.Context, rID string) (model.HostResource, error)
}

type MDNSAdvHandler interface {
	List(ctx context.Context, filter model.ServiceGroupFilter) ([]model.ServiceGroup, error)
	Add(ctx context.Context, sg model.ServiceGroup) error
	Get(ctx context.Context, id string) (model.ServiceGroup, error)
	Update(ctx context.Context, sg model.ServiceGroup) error
	Delete(ctx context.Context, id string) error
}
