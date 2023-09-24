// Copyright 2017 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package isulad

import (
	"context"
	"fmt"

	containersapi "github.com/google/cadvisor/third_party/isulad/api/services/containers"
)

type isuladClientMock struct {
	cntrs     map[string]*containersapi.Container
	returnErr error
}

func (c *isuladClientMock) LoadContainer(ctx context.Context, id string) (*containersapi.Container, error) {
	if c.returnErr != nil {
		return nil, c.returnErr
	}
	cntr, ok := c.cntrs[id]
	if !ok {
		return nil, fmt.Errorf("unable to find container %q", id)
	}
	return cntr, nil
}

func (c *isuladClientMock) Version(ctx context.Context) (string, error) {
	return "test-v0.0.0", nil
}

func mockIsuladClient(cntrs map[string]*containersapi.Container, returnErr error) IsuladClient {
	return &isuladClientMock{
		cntrs:     cntrs,
		returnErr: returnErr,
	}
}
