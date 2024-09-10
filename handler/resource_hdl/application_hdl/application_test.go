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

package application_hdl

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/SENERGY-Platform/mgw-host-manager/lib/model"
	"github.com/SENERGY-Platform/mgw-host-manager/util"
	"os"
	"path"
	"reflect"
	"testing"
)

func TestHandler_Init(t *testing.T) {
	tmpFilePath := path.Join(t.TempDir(), "test.json")
	t.Run("file does not exist", func(t *testing.T) {
		h, err := New(tmpFilePath, "")
		if err != nil {
			t.Error(err)
		}
		if err = h.Init(); err != nil {
			t.Error(err)
		}
	})
	apps := map[string]model.HostApplication{
		util.GenHash("/test/socket1"): {
			ID: util.GenHash("/test/socket1"),
			HostApplicationBase: model.HostApplicationBase{
				Name:   "Test 1",
				Socket: "/test/socket1",
			},
		},
	}
	t.Run("file exists migration", func(t *testing.T) {
		oldFmt := []model.HostApplicationBase{
			{
				Name:   "Test 1",
				Socket: "/test/socket1",
			},
		}
		f, err := os.Create(tmpFilePath)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		err = json.NewEncoder(f).Encode(oldFmt)
		if err != nil {
			t.Fatal(err)
		}
		h, err := New(tmpFilePath, "")
		if err != nil {
			t.Error(err)
		}
		if err = h.Init(); err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(apps, h.apps) {
			t.Errorf("got %+v, expected %+v", h.apps, apps)
		}
	})
	t.Run("file exists", func(t *testing.T) {
		h, err := New(tmpFilePath, "")
		if err != nil {
			t.Error(err)
		}
		if err = h.Init(); err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(apps, h.apps) {
			t.Errorf("got %+v, expected %+v", h.apps, apps)
		}
	})
	t.Run("error", func(t *testing.T) {
		f, err := os.Create(tmpFilePath)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		h, err := New(tmpFilePath, "")
		if err != nil {
			t.Error(err)
		}
		if err = h.Init(); err == nil {
			t.Error("expected error")
		}
	})
}

func TestHandler_List(t *testing.T) {
	h, err := New(path.Join(t.TempDir(), "test.json"), "")
	if err != nil {
		t.Error(err)
	}
	t.Run("empty", func(t *testing.T) {
		li, err := h.List(context.Background())
		if err != nil {
			t.Error(err)
		}
		if len(li) != 0 {
			t.Error("expected empty list")
		}
	})
	h.apps["123"] = model.HostApplication{
		ID: "123",
		HostApplicationBase: model.HostApplicationBase{
			Name:   "Test 1",
			Socket: "/test/socket1",
		},
	}
	a := []model.HostApplication{{
		ID: "123",
		HostApplicationBase: model.HostApplicationBase{
			Name:   "Test 1",
			Socket: "/test/socket1",
		},
	},
	}
	b, err := h.List(context.Background())
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(a, b) {
		t.Errorf("got %+v, expected %+v", b, a)
	}
}

func TestHandler_Add(t *testing.T) {
	h, err := New(path.Join(t.TempDir(), "test.json"), "")
	if err != nil {
		t.Error(err)
	}
	id, err := h.Add(context.Background(), model.HostApplicationBase{
		Name:   "Test",
		Socket: "/test/socket",
	})
	if err != nil {
		t.Error(err)
	}
	a := model.HostApplication{
		ID: id,
		HostApplicationBase: model.HostApplicationBase{
			Name:   "Test",
			Socket: "/test/socket",
		},
	}
	b, ok := h.apps[id]
	if !ok {
		t.Error("id not in map")
	}
	if !reflect.DeepEqual(a, b) {
		t.Errorf("got %+v, expected %+v", b, a)
	}
}

func TestHandler_Remove(t *testing.T) {
	h, err := New(path.Join(t.TempDir(), "test.json"), "")
	if err != nil {
		t.Error(err)
	}
	t.Run("doesn't exist", func(t *testing.T) {
		err = h.Remove(context.Background(), "123")
		if err == nil {
			t.Error("expected error")
		}
		var nf *model.NotFoundError
		if !errors.As(err, &nf) {
			t.Error("expected NotFoundError")
		}
	})
	h.apps["123"] = model.HostApplication{}
	err = h.Remove(context.Background(), "123")
	if err != nil {
		t.Error(err)
	}
	if len(h.apps) != 0 {
		t.Error("expected empty map")
	}
}

func TestHandler_Get(t *testing.T) {
	h, err := New(path.Join(t.TempDir(), "test.json"), "")
	if err != nil {
		t.Error(err)
	}
	t.Run("empty", func(t *testing.T) {
		li, err := h.Get(context.Background())
		if err != nil {
			t.Error(err)
		}
		if len(li) != 0 {
			t.Error("expected empty list")
		}
	})
	h.apps["123"] = model.HostApplication{
		ID: "123",
		HostApplicationBase: model.HostApplicationBase{
			Name:   "Test 1",
			Socket: "/test/socket1",
		},
	}
	a := map[string]model.HostResourceBase{
		"123": {
			Name: "Test 1",
			Path: "/test/socket1",
		},
	}
	b, err := h.Get(context.Background())
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(a, b) {
		t.Errorf("got %+v, expected %+v", b, a)
	}
}

func TestStorageMigration(t *testing.T) {
	tmpFilePath := path.Join(t.TempDir(), "test.json")
	oldFmt := []model.HostApplicationBase{
		{
			Name:   "Test 1",
			Socket: "/test/socket1",
		},
		{
			Name:   "Test 2",
			Socket: "/test/socket2",
		},
	}
	newFmt := map[string]model.HostApplication{
		util.GenHash("/test/socket1"): {
			ID: util.GenHash("/test/socket1"),
			HostApplicationBase: model.HostApplicationBase{
				Name:   "Test 1",
				Socket: "/test/socket1",
			},
		},
		util.GenHash("/test/socket2"): {
			ID: util.GenHash("/test/socket2"),
			HostApplicationBase: model.HostApplicationBase{
				Name:   "Test 2",
				Socket: "/test/socket2",
			},
		},
	}
	t.Run("file does not exist", func(t *testing.T) {
		_, err := migrateStoFile(tmpFilePath)
		if err == nil {
			t.Error("error should not be nil")
		}
	})
	t.Run("file exists empty", func(t *testing.T) {
		f, err := os.Create(tmpFilePath)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		_, err = migrateStoFile(tmpFilePath)
		if err == nil {
			t.Error("error should not be nil")
		}
	})
	t.Run("file exists old format", func(t *testing.T) {
		f, err := os.Create(tmpFilePath)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		err = json.NewEncoder(f).Encode(oldFmt)
		if err != nil {
			t.Fatal(err)
		}
		apps, err := migrateStoFile(tmpFilePath)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(newFmt, apps) {
			t.Errorf("got %+v, expected %+v", apps, newFmt)
		}
		f2, err := os.Open(tmpFilePath)
		if err != nil {
			t.Fatal(err)
		}
		defer f2.Close()
		var apps2 map[string]model.HostApplication
		err = json.NewDecoder(f2).Decode(&apps2)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(newFmt, apps2) {
			t.Errorf("got %+v, expected %+v", apps2, newFmt)
		}
	})
}
