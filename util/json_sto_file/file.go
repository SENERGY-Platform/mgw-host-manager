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
	"encoding/json"
	"errors"
	"github.com/SENERGY-Platform/mgw-host-manager/util"
	"io"
	"os"
)

func Read(p string, t any) error {
	file, err := os.Open(p)
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(t); err != nil {
		return err
	}
	return nil
}

func Write(v any, p string, backup bool) error {
	if backup {
		if err := Copy(p, p+".bk"); err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return err
			}
		}
	}
	file, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	if err = json.NewEncoder(file).Encode(v); err != nil {
		if backup {
			e := Copy(p+".bk", p)
			if e != nil && !errors.Is(e, os.ErrNotExist) {
				util.Logger.Error(e)
			}
		}
		return err
	}
	return nil
}

func Copy(src, dst string) error {
	sFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sFile.Close()
	dFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dFile.Close()
	_, err = io.Copy(dFile, sFile)
	return err
}
