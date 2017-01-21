//
// Copyright 2017 Alsanium, SAS. or its affiliates. All rights reserved.
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
//

package customres

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/eawsy/aws-lambda-go-event/service/lambda/runtime/event/cloudformationevt"
)

// NewPhysicalResourceID generates an AWS CloudFormation like unique physical
// resource id of the forme "StackName-LogicalResourceID-RandomID".
func NewPhysicalResourceID(evt *cloudformationevt.Event) string {
	rns := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	snm := strings.Split(evt.StackID, "/")[1]
	lid := evt.LogicalResourceID
	gen := rand.New(rand.NewSource(time.Now().UnixNano()))
	rnd := make([]byte, 12)
	for i := range rnd {
		rnd[i] = rns[gen.Intn(len(rns))]
	}
	return fmt.Sprintf("%s-%s-%s", snm, lid, rnd)
}
