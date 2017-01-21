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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
	"github.com/eawsy/aws-lambda-go-event/service/lambda/runtime/event/cloudformationevt"
)

type response struct {
	Status             string      `json:"Status"`
	StackID            string      `json:"StackId"`
	RequestID          string      `json:"RequestId"`
	PhysicalResourceID string      `json:"PhysicalResourceId"`
	LogicalResourceID  string      `json:"LogicalResourceId"`
	Reason             string      `json:"Reason,omitempty"`
	Data               interface{} `json:"Data,omitempty"`

	url string
	ctx *runtime.Context
}

func newResponse(evt *cloudformationevt.Event, ctx *runtime.Context) *response {
	return &response{
		StackID:            evt.StackID,
		RequestID:          evt.RequestID,
		PhysicalResourceID: evt.PhysicalResourceID,
		LogicalResourceID:  evt.LogicalResourceID,

		url: evt.ResponseURL,
		ctx: ctx,
	}
}

func (r *response) send() error {
	body, err := json.Marshal(r)
	if err != nil {
		return err
	}

	log.Printf("CloudFormation Response:\n%s", body)

	req, err := http.NewRequest(http.MethodPut, r.url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Del("Content-Type")

	cln := &http.Client{}
	res, err := cln.Do(req)
	if err != nil {
		return err
	}

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	res.Body.Close()

	log.Printf("CloudFormation Status: %d %s\n %s", res.StatusCode, res.Status, body)

	return nil
}

func (r *response) success(id string, data interface{}) error {
	r.Status = "SUCCESS"
	r.PhysicalResourceID = id
	r.Data = data
	return r.send()
}

func (r *response) error(id string, err error) error {
	log.Printf("CloudFormation Error: %s", err)
	r.Status = "FAILED"
	r.PhysicalResourceID = id
	r.Reason = fmt.Sprintf("See CloudWatch Logs at %s/%s/%s", r.ctx.LogGroupName, r.ctx.LogStreamName, r.ctx.AWSRequestID)
	return r.send()
}
