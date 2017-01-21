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
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	goruntime "runtime"
	"strings"

	"github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
	"github.com/eawsy/aws-lambda-go-event/service/lambda/runtime/event/cloudformationevt"
)

var handlers = make(map[string]Handler)

// LambdaHandler responds to a Lambda function invocation.
type LambdaHandler func(json.RawMessage, *runtime.Context) (interface{}, error)

// Register registers the handler for the given custom resource. If a handler
// already exists for the custom resource, Register panics.
//
// Custom resource type names can include alphanumeric characters and the
// following characters: _@-
// You can specify a custom resource type name up to a maximum length of 60
// characters.
// Example: MyCompany@MyCustomResource
func Register(name string, handler Handler) {
	if ok, _ := regexp.MatchString(`^[A-Za-z]+[A-Za-z0-9_@\-]*$`, name); !ok {
		panic("customres: invalid type name")
	}
	if handler == nil {
		panic("customres: nil handler")
	}
	if _, ok := handlers[name]; ok {
		panic("customres: multiple registrations for " + name)
	}
	handlers[name] = handler
}

// HandleLambda dispatches the request to the matching custom resource handler.
// If no handler is registered for the requested custom resource, then an error
// is returned to AWS CloudFormation.
//
// For detailed information about the incoming raw request, outgoing raw
// response and eventual errors, refer to the AWS CloudWatch Logs stream of the
// AWS Lambda function.
func HandleLambda(raw json.RawMessage, ctx *runtime.Context) (interface{}, error) {
	evt := new(cloudformationevt.Event)
	if err := json.Unmarshal(raw, evt); err != nil {
		return nil, err
	}

	log.Printf("CloudFormation Request:\n%s", raw)

	res := newResponse(evt, ctx)

	typ := strings.TrimPrefix(evt.ResourceType, "Custom::")
	hld, ok := handlers[typ]
	if !ok {
		return nil, res.error(NewPhysicalResourceID(evt), fmt.Errorf("handler for '%s' not found", typ))
	}

	var (
		id   string
		data interface{}
		err  error
	)

	func() {
		defer func() {
			info := recover()
			if info == nil {
				return
			}

			buf := make([]byte, 64<<10)
			buf = buf[:goruntime.Stack(buf, false)]
			err = fmt.Errorf("%s\n%s", info, buf)
			id = evt.PhysicalResourceID
			if id == "" {
				id = NewPhysicalResourceID(evt)
			}
		}()

		switch evt.RequestType {
		case "Create":
			id, data, err = hld.Create(evt, ctx)
		case "Update":
			id, data, err = hld.Update(evt, ctx)
		case "Delete":
			id = evt.PhysicalResourceID
			err = hld.Delete(evt, ctx)
		}
	}()

	if err != nil {
		return nil, res.error(id, err)
	}
	return nil, res.success(id, data)
}
