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

package json_sto_file

import (
	"os"
	"path"
	"reflect"
	"testing"
)

func TestReadAndWrite(t *testing.T) {
	tmpFilePath := path.Join(t.TempDir(), "test.json")
	testValue := "test"
	t.Run("read file does not exist", func(t *testing.T) {
		err := Read(tmpFilePath, nil)
		if err == nil {
			t.Error("error should not be nil")
		}
	})
	t.Run("write file does not exist", func(t *testing.T) {
		err := Write(testValue, tmpFilePath, true)
		if err != nil {
			t.Error(err)
		}
		_, err = os.Stat(tmpFilePath)
		if err != nil {
			t.Error(err)
		}
		_, err = os.Stat(tmpFilePath + ".bk")
		if err == nil {
			t.Error("error should not be nil")
		}
	})
	t.Run("read file exists", func(t *testing.T) {
		var s string
		err := Read(tmpFilePath, &s)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(testValue, s) {
			t.Errorf("got %+v, expected %+v", s, testValue)
		}
	})
	t.Run("write file exists", func(t *testing.T) {
		testValue = "test2"
		err := Write(testValue, tmpFilePath, true)
		if err != nil {
			t.Error(err)
		}
		_, err = os.Stat(tmpFilePath + ".bk")
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("read file exists", func(t *testing.T) {
		var s string
		err := Read(tmpFilePath, &s)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(testValue, s) {
			t.Errorf("got %+v, expected %+v", s, testValue)
		}
	})
}
