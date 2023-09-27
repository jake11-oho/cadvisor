// Copyright 2023 Google Inc. All Rights Reserved.
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
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/google/cadvisor/container/containerd/errdefs"
	"github.com/google/cadvisor/container/containerd/pkg/dialer"
	containersapi "github.com/google/cadvisor/third_party/isulad/api/services/containers"
)

type client struct {
	containerService containersapi.ContainerServiceClient
}

type IsuladClient interface {
	InspectContainer(ctx context.Context, id string) (*ContainerJSON, error)
	Version(ctx context.Context) (string, error)
}

var (
	ErrTaskIsInUnknownState = errors.New("isulad task is in unknown state") // used when process reported in isulad task is in Unknown State
)

var once sync.Once
var ctrdClient IsuladClient = nil

const (
	maxBackoffDelay   = 3 * time.Second
	baseBackoffDelay  = 100 * time.Millisecond
	connectionTimeout = 2 * time.Second
	maxMsgSize        = 16 * 1024 * 1024 // 16MB
)

// Client creates a containerd client
func Client(address string) (IsuladClient, error) {
	var retErr error
	once.Do(func() {
		tryConn, err := net.DialTimeout("unix", address, connectionTimeout)
		if err != nil {
			retErr = fmt.Errorf("isulad: cannot unix dial isulad api service: %v", err)
			return
		}
		tryConn.Close()

		connParams := grpc.ConnectParams{
			Backoff: backoff.DefaultConfig,
		}
		connParams.Backoff.BaseDelay = baseBackoffDelay
		connParams.Backoff.MaxDelay = maxBackoffDelay
		gopts := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithContextDialer(dialer.ContextDialer),
			grpc.WithBlock(),
			grpc.WithConnectParams(connParams),
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMsgSize)),
		}

		ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
		defer cancel()
		conn, err := grpc.DialContext(ctx, dialer.DialAddress(address), gopts...)
		if err != nil {
			retErr = err
			return
		}
		ctrdClient = &client{
			containerService: containersapi.NewContainerServiceClient(conn),
		}
	})
	return ctrdClient, retErr
}

func (c *client) InspectContainer(ctx context.Context, id string) (*ContainerJSON, error) {
	const timeout int32 = 120
	r, err := c.containerService.Inspect(ctx, &containersapi.InspectContainerRequest{
		Id:      id,
		Bformat: false,
		Timeout: timeout,
	})
	if err != nil {
		return nil, errdefs.FromGRPC(err)
	}
	var container ContainerJSON
	err = json.Unmarshal([]byte(r.ContainerJSON), &container)
	if err != nil {
		return nil, err
	}
	return &container, nil
}

func (c *client) Version(ctx context.Context) (string, error) {
	response, err := c.containerService.Version(ctx, &containersapi.VersionRequest{})
	if err != nil {
		return "", errdefs.FromGRPC(err)
	}
	return response.Version, nil
}
