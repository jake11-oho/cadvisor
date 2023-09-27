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

// isulad container state
type ContainerState struct {
	Status     string `json:"Status"`
	Running    bool   `json:"Running"`
	Paused     bool   `json:"Paused"`
	Restarting bool   `json:"Restarting"`
	Pid        int    `json:"Pid"`
	ExitCode   int    `json:"ExitCode"`
	Error      string `json:"Error"`
	StartedAt  string `json:"StartedAt"`
	FinishedAt string `json:"FinishedAt"`
}

// isulad container resource
type ContainerResources struct {
	CPUPeriod  int64 `json:"CPUPeriod"`
	CPUQuota   int64 `json:"CPUQuota"`
	CPUShares  int64 `json:"CPUShares"`
	Memory     int64 `json:"Memory"`
	MemorySwap int64 `json:"MemorySwap"`
	Hugetlbs   []struct {
		PageSize string `json:"PageSize"`
		Limit    uint64 `json:"Limit"`
	} `json:"Hugetlbs"`
}

// isulad container config
type ContainerConfig struct {
	Hostname    string               `json:"Hostname"`
	User        string               `json:"User"`
	Env         []string             `json:"Env"`
	Tty         bool                 `json:"Tty"`
	Cmd         []string             `json:"Cmd,omitempty"`
	Entrypoint  []string             `json:"Entrypoint,omitempty"`
	Labels      map[string]string    `json:"Labels,omitempty"`
	Volumes     map[string]*struct{} `json:"Volumes,omitempty"`
	Annotations map[string]string    `json:"Annotations,omitempty"`
	HealthCheck struct {
		Test []string `json:"Test,omitempty"`
	} `json:"Healthcheck,omitempty"`
	Image      string `json:"Image"`
	ImageRef   string `json:"ImageRef"`
	StopSignal string `json:"StopSignal"`
}

type ContainerGraphDriver struct {
	Data struct {
		LowerDir   string `json:"LowerDir"`
		MergedDir  string `json:"MergedDir"`
		UpperDir   string `json:"UpperDir"`
		WorkDir    string `json:"WorkDir"`
		DeviceId   string `json:"DeviceId"`
		DeviceName string `json:"DeviceName"`
		DeviceSize string `json:"DeviceSize"`
	} `json:"Data"`
	Name string `json:"Name"`
}

type MountPoint struct {
	Type        string `json:"Type"`
	Name        string `json:"Name"`
	Source      string `json:"Source"`
	Destination string `json:"Destination"`
	Driver      string `json:"Driver"`
	Mode        string `json:"Mode"`
	RW          bool   `json:"RW"`
	Propagation string `json:"Propagation"`
}

// isulad container json from inspect[not full]
type ContainerJSON struct {
	Id              string             `json:"Id"`
	Created         string             `json:"Created"`
	Path            string             `json:"Path"`
	Args            []string           `json:"Args"`
	State           ContainerState     `json:"State"`
	Resources       ContainerResources `json:"Resources"`
	Image           string             `json:"Image"`
	ResolvConfPath  string             `json:"ResolvConfPath"`
	HostnamePath    string             `json:"HostnamePath"`
	HostsPath       string             `json:"HostsPath"`
	LogPath         string             `json:"LogPath"`
	Name            string             `json:"Name"`
	RestartCount    int                `json:"RestartCount"`
	MountLabel      string             `json:"MountLabel"`
	ProcessLabel    string             `json:"ProcessLabel"`
	SeccompProfile  string             `json:"SeccompProfile"`
	NoNewPrivileges bool               `json:"NoNewPrivileges"`
	HostConfig      struct {
		NetworkMode string `json:"NeworkMode"`
	} `json:"HostConfig"`
	GraphDriver     ContainerGraphDriver `json:"GraphDriver"`
	Mounts          []MountPoint         `json:"Mounts"`
	Config          ContainerConfig      `json:"Config"`
	NetworkSettings struct {
		IPAddress string `json:"IPAddress"`
	} `json:"NetworkSettings"`
}
