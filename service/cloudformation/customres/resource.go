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
	"github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
	"github.com/eawsy/aws-lambda-go-event/service/lambda/runtime/event/cloudformationevt"
)

// Handler responds to a Lambda invocation whenever a custom resource is
// created, updated, or deleted.
//
// If any sub handler panics, it is automatically recovered, stack trace is
// logged in AWS CloudWatch Logs stream and AWS CloudFormation is informed and
// reacts accordingly.
type Handler interface {
	// Create can either return a physical id and optionaly a map of name-value
	// pairs, or an error, in which case AWS CloudFormation fails the resource
	// creation.
	//
	// The physical id is a name or a unique identifier that corresponds to the
	// physical instance id of the custom resource. See NewPhysicalResourceID, for
	// a strong and unique identifier.
	// The map of name-value pairs can be accessed by name in the template
	// with Fn::GetAtt intrinsic function.
	// If an error is returned, AWS CloudFormation recognizes the failure and
	// reacts accordingly. The detailed information about the error will be
	// written in the Lambda function's AWS CloudWatch Logs stream.
	Create(*cloudformationevt.Event, *runtime.Context) (string, interface{}, error)

	// Update can either return a physical id and optionaly a map of name-value
	// pairs, or an error, in which case AWS CloudFormation fails the resource
	// update.
	//
	// The physical id is usually the same as the one returned in response to the
	// resource creation. However, you can also update custom resources that
	// requires a replacement of the underlying physical resource. In this case,
	// the new custom resource must send a new physical id. When AWS
	// CloudFormation receives the response, it compares the PhysicalResourceID
	// between the old and new custom resources. If they are different,
	// AWS CloudFormation recognizes the update as a replacement and send a delete
	// request to the old resource.
	// Notice also that the returned map of name-value pairs must contain all the
	// pairs returned in response to the resource creation, with updated values;
	// otherwise Fn::GetAtt intrinsic functions fails.
	Update(*cloudformationevt.Event, *runtime.Context) (string, interface{}, error)

	// Delete can optionaly return an error, in which case AWS CloudFormation
	// fails the resource deletion.
	Delete(*cloudformationevt.Event, *runtime.Context) error
}
