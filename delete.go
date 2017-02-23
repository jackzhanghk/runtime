// Copyright (c) 2014,2015,2016 Docker, Inc.
// Copyright (c) 2017 Intel Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"

	vc "github.com/containers/virtcontainers"
	"github.com/containers/virtcontainers/pkg/oci"
	"github.com/urfave/cli"
)

var deleteCommand = cli.Command{
	Name:  "delete",
	Usage: "Delete any resources held by one or more containers",
	ArgsUsage: `<container-id> [container-id...]

   <container-id> is the name for the instance of the container.

EXAMPLE:
   If the container id is "ubuntu01" and ` + name + ` list currently shows the
   status of "ubuntu01" as "stopped" the following will delete resources held
   for "ubuntu01" removing "ubuntu01" from the ` + name + ` list of containers:

       # ` + name + ` delete ubuntu01`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "force, f",
			Usage: "Forcibly deletes the container if it is still running (uses SIGKILL)",
		},
	},
	Action: func(context *cli.Context) error {
		return delete(context.String("container-id"), context.Bool("force"))
	},
}

func delete(containerID string, force bool) error {
	// Checks the MUST and MUST NOT from OCI runtime specification
	if err := validContainer(containerID); err != nil {
		return err
	}

	if force == false {
		podStatus, err := vc.StatusPod(containerID)
		if err != nil {
			return err
		}

		state, err := oci.StatusToOCIState(podStatus)
		if err != nil {
			return err
		}

		running, err := processRunning(state.Pid)
		if err != nil {
			return err
		}

		if running == true {
			return fmt.Errorf("Container still running, should be stopped")
		}
	}

	if _, err := vc.StopPod(containerID); err != nil {
		return err
	}

	if _, err := vc.DeletePod(containerID); err != nil {
		return err
	}

	return nil
}