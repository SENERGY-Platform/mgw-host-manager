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

package blacklist_hdl

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/SENERGY-Platform/mgw-host-manager/lib/model"
	"os"
	"path"
	"reflect"
	"testing"
)

func TestHandler_Init(t *testing.T) {
	tmpFilePath := path.Join(t.TempDir(), "test.json")
	t.Run("file does not exist", func(t *testing.T) {
		h, err := New(tmpFilePath)
		if err != nil {
			t.Error(err)
		}
		if err = h.Init(); err != nil {
			t.Error(err)
		}
	})
	values := []string{"a", "b"}
	t.Run("file exists", func(t *testing.T) {
		f, err := os.Create(tmpFilePath)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		err = json.NewEncoder(f).Encode([]string{"a", "b"})
		if err != nil {
			t.Fatal(err)
		}
		h, err := New(tmpFilePath)
		if err != nil {
			t.Error(err)
		}
		if err = h.Init(); err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(values, h.values) {
			t.Errorf("got %+v, expected %+v", h.values, values)
		}
	})
	t.Run("error", func(t *testing.T) {
		f, err := os.Create(tmpFilePath)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		h, err := New(tmpFilePath)
		if err != nil {
			t.Error(err)
		}
		if err = h.Init(); err == nil {
			t.Error("expected error")
		}
	})
}

func TestHandler_List(t *testing.T) {
	h, err := New(path.Join(t.TempDir(), "test.json"))
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
	h.values = append(h.values, "a")
	a := []string{"a"}
	b, err := h.List(context.Background())
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(a, b) {
		t.Errorf("got %+v, expected %+v", b, a)
	}
}

func TestHandler_Add(t *testing.T) {
	h, err := New(path.Join(t.TempDir(), "test.json"))
	if err != nil {
		t.Error(err)
	}
	val := "a"
	err = h.Add(context.Background(), val)
	if err != nil {
		t.Error(err)
	}
	if !inSlice(val, h.values) {
		t.Error("value not in list")
	}
	t.Run("error", func(t *testing.T) {
		err = h.Add(context.Background(), val)
		if err == nil {
			t.Error("expected error")
		}
		var ii *model.InvalidInputError
		if !errors.As(err, &ii) {
			t.Error("invalid error type")
		}
	})
}

func TestHandler_Remove(t *testing.T) {
	val := "a"
	h, err := New(path.Join(t.TempDir(), "test.json"))
	if err != nil {
		t.Error(err)
	}
	t.Run("doesn't exist", func(t *testing.T) {
		err = h.Remove(context.Background(), val)
		if err == nil {
			t.Error("expected error")
		}
		var nf *model.NotFoundError
		if !errors.As(err, &nf) {
			t.Error("expected NotFoundError")
		}
	})
	h.values = append(h.values, val)
	err = h.Remove(context.Background(), val)
	if err != nil {
		t.Error(err)
	}
	if len(h.values) != 0 {
		t.Error("expected empty map")
	}
}
