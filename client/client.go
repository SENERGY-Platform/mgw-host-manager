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

package client

import (
	"github.com/SENERGY-Platform/go-base-http-client"
	"github.com/SENERGY-Platform/mgw-host-manager/lib"
	"github.com/SENERGY-Platform/mgw-host-manager/lib/model"
	"net/http"
)

type HmClient = lib.Api

type Client struct {
	*base_client.Client
	baseUrl string
}

func New(httpClient base_client.HTTPClient, baseUrl string) *Client {
	return &Client{
		Client:  base_client.New(httpClient, customError, model.HeaderRequestID),
		baseUrl: baseUrl,
	}
}

func customError(code int, err error) error {
	switch code {
	case http.StatusInternalServerError:
		err = model.NewInternalError(err)
	case http.StatusNotFound:
		err = model.NewNotFoundError(err)
	case http.StatusBadRequest:
		err = model.NewInvalidInputError(err)
	}
	return err
}
