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
	"flag"
	"fmt"
	"path"
	"regexp"
	"strings"

	"golang.org/x/net/context"
	"k8s.io/klog/v2"

	"github.com/google/cadvisor/container"
	"github.com/google/cadvisor/container/libcontainer"
	"github.com/google/cadvisor/fs"
	info "github.com/google/cadvisor/info/v1"
	"github.com/google/cadvisor/watcher"
)

var ArgIsuladEndpoint = flag.String("isulad", "/var/run/isulad.sock", "isulad endpoint")

const isuladNamespace = "isulad"

// Regexp that identifies isulad cgroups, containers started with
// --cgroup-parent have another prefix than 'isulad'
var isuladCgroupRegexp = regexp.MustCompile(`([a-z0-9]{64})`)

type isuladFactory struct {
	machineInfoFactory info.MachineInfoFactory
	client             IsuladClient
	version            string
	// Information about the mounted cgroup subsystems.
	cgroupSubsystems map[string]string
	// Information about mounted filesystems.
	fsInfo          fs.FsInfo
	includedMetrics container.MetricSet
}

func (f *isuladFactory) String() string {
	return "isulad"
}

func (f *isuladFactory) NewContainerHandler(name string, metadataEnvAllowList []string, inHostNamespace bool) (handler container.ContainerHandler, err error) {
	client, err := Client(*ArgIsuladEndpoint)
	if err != nil {
		return
	}

	return newIsuladContainerHandler(
		client,
		name,
		f.machineInfoFactory,
		f.fsInfo,
		f.cgroupSubsystems,
		inHostNamespace,
		f.includedMetrics,
	)
}

// Returns the isulad ID from the full container name.
func ContainerNameToIsuladID(name string) string {
	id := path.Base(name)
	if matches := isuladCgroupRegexp.FindStringSubmatch(id); matches != nil {
		return matches[1]
	}
	return id
}

// isContainerName returns true if the cgroup with associated name
// corresponds to a isulad container.
func isContainerName(name string) bool {
	// TODO: May be check with HasPrefix IsuladNamespace
	if strings.HasSuffix(name, ".mount") {
		return false
	}
	return isuladCgroupRegexp.MatchString(path.Base(name))
}

// Isulad can handle and accept all isulad created containers
func (f *isuladFactory) CanHandleAndAccept(name string) (bool, bool, error) {
	// if the container is not associated with isulad, we can't handle it or accept it.
	if !isContainerName(name) {
		return false, false, nil
	}
	// Check if the container is known to isulad and it is running.
	id := ContainerNameToIsuladID(name)
	// If container and task lookup in isulad fails then we assume
	// that the container state is not known to isulad
	ctx := context.Background()
	_, err := f.client.LoadContainer(ctx, id)
	if err != nil {
		return false, false, fmt.Errorf("failed to load container: %v", err)
	}

	return true, true, nil
}

func (f *isuladFactory) DebugInfo() map[string][]string {
	return map[string][]string{}
}

// Register root container before running this function!
func Register(factory info.MachineInfoFactory, fsInfo fs.FsInfo, includedMetrics container.MetricSet) error {
	client, err := Client(*ArgIsuladEndpoint)
	if err != nil {
		return fmt.Errorf("unable to create isulad client: %v", err)
	}

	isuladVersion, err := client.Version(context.Background())
	if err != nil {
		return fmt.Errorf("failed to fetch isulad client version: %v", err)
	}

	cgroupSubsystems, err := libcontainer.GetCgroupSubsystems(includedMetrics)
	if err != nil {
		return fmt.Errorf("failed to get cgroup subsystems: %v", err)
	}

	klog.V(1).Infof("Registering isulad factory")
	f := &isuladFactory{
		cgroupSubsystems:   cgroupSubsystems,
		client:             client,
		fsInfo:             fsInfo,
		machineInfoFactory: factory,
		version:            isuladVersion,
		includedMetrics:    includedMetrics,
	}

	container.RegisterContainerHandlerFactory(f, []watcher.ContainerWatchSource{watcher.Raw})
	return nil
}
