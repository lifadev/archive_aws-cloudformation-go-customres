<a id="top" name="top"></a>
[<img src="_asset/logo_powered-by-aws.png" alt="Powered by Amazon Web Services" align="right">][aws-home]
[<img src="_asset/logo_created-by-eawsy.png" alt="Created by eawsy" align="right">][eawsy-home]

# eawsy/aws-cloudformation-go-customres

> Author your AWS CloudFormation Custom Resources in Go.

[![Api][badge-api]][eawsy-api]
[![Status][badge-status]](#top)
[![License][badge-license]](LICENSE)
[![Help][badge-help]][eawsy-chat]
[![Social][badge-social]][eawsy-twitter]

[AWS Lambda][aws-lambda-home] lets you run code without provisioning or managing servers. With 
[eawsy/aws-lambda-go-shim][eawsy-runtime], you can author your Lambda function code in Go. This project allows you to 
create [AWS Lambda-backed Custom Resources][aws-cloudformation-custom] for 
[AWS CloudFormation][aws-cloudformation-home] in Go.

[<img src="_asset/misc_arrow-up.png" align="right">](#top)
## Quick Hands-On

> For step by step instructions on how to author your AWS Lambda function code in Go, see 
  [eawsy/aws-lambda-go-shim][eawsy-runtime].

[<img src="_asset/misc_arrow-up.png" align="right">](#top)
### Dependencies

```sh
go get -u -d github.com/eawsy/aws-lambda-go-core/...
go get -u -d github.com/eawsy/aws-lambda-go-event/...
go get -u -d github.com/eawsy/aws-cloudformation-go-customres/...
```

[<img src="_asset/misc_arrow-up.png" align="right">](#top)
### Create

```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/eawsy/aws-cloudformation-go-customres/service/cloudformation/customres"
	"github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
	"github.com/eawsy/aws-lambda-go-event/service/lambda/runtime/event/cloudformationevt"
)

// Handle is the exported handler called by AWS Lambda.
var Handle customres.LambdaHandler

type MyResource struct{}

func (r *MyResource) Create(evt *cloudformationevt.Event, ctx *runtime.Context) (string, interface{}, error) {
	var props map[string]string
	if err := json.Unmarshal(evt.ResourceProperties, &props); err != nil {
		return "", nil, err
	}

	id := customres.NewPhysicalResourceID(evt)

	resp := map[string]string{
		"Message": fmt.Sprintf("Hello, %s!", props["Name"]),
	}

	return id, resp, nil
}

func (r *MyResource) Update(evt *cloudformationevt.Event, ctx *runtime.Context) (string, interface{}, error) {
	return r.Create(evt, ctx)
}

func (r *MyResource) Delete(*cloudformationevt.Event, *runtime.Context) error {
	return nil
}

func init() {
	customres.Register("MyCompany@MyResource", new(MyResource))
	Handle = customres.HandleLambda
}
```

[<img src="_asset/misc_arrow-up.png" align="right">](#top)
### Build

> For step by step instructions on how to author your AWS Lambda function code in Go, see 
  [eawsy/aws-lambda-go-shim][eawsy-runtime].

```sh
make
```

[<img src="_asset/misc_arrow-up.png" align="right">](#top)
### Deploy

> AWS CLI is used for the sake of simplicity. You are free to use your favorite deployment tool.

```sh
aws lambda create-function \
  --role arn:aws:iam::AWS_ACCOUNT_ID:role/lambda_basic_execution \
  --function-name cfn-customres \
  --zip-file fileb://package.zip \
  --runtime python2.7 \
  --handler handler.Handle
```

[<img src="_asset/misc_arrow-up.png" align="right">](#top)
### Invoke

```yaml
Resources:
  HelloWorld:
    Type: Custom::MyCompany@MyResource
    Properties:
      ServiceToken: arn:aws:lambda:AWS_REGION:AWS_ACCOUNT_ID:function:cfn-customres
      Name: World2

Outputs:
  Message:
    Value: !GetAtt HelloWorld.Message
```

```sh
aws cloudformation create-stack \
  --template-body file://example.cfn.yaml
  --stack-name <YOUR STACK NAME>

aws cloudformation wait stack-create-complete \
  --stack-name <YOUR STACK NAME>

aws cloudformation describe-stacks \
  --stack-name <YOUR STACK NAME>

# {
# ...
#   "Outputs": [
#     {
#       "OutputKey": "Message", 
#       "OutputValue": "Hello, World!"
#     }
#   ]
# ...
# }
```

[<img src="_asset/misc_arrow-up.png" align="right">](#top)
## About

[![eawsy](_asset/logo_eawsy.png)][eawsy-home]

This project is maintained and funded by Alsanium, SAS.

[We][eawsy-home] :heart: [AWS][aws-home] and open source software. See [our other projects][eawsy-github], or 
[hire us][eawsy-hire] to help you build modern applications on AWS.

[<img src="_asset/misc_arrow-up.png" align="right">](#top)
## Contact

We want to make it easy for you, users and contributers, to talk with us and connect with each others, to share ideas, 
solve problems and make help this project awesome. Here are the main channels we're running currently and we'd love to 
hear from you on them.

### Twitter 
  
[eawsyhq][eawsy-twitter] 

Follow and chat with us on Twitter. 

Share stories!

### Gitter 

[eawsy/bavardage][eawsy-chat]

This is for all of you. Users, developers and curious. You can find help, links, questions and answers from all the 
community including the core team.

Ask questions!

### GitHub

[pull requests][eawsy-pr] & [issues][eawsy-issues]

You are invited to contribute new features, fixes, or updates, large or small; we are always thrilled to receive pull 
requests, and do our best to process them as fast as we can.

Before you start to code, we recommend discussing your plans through the [eawsy/bavardage channel][eawsy-chat], 
especially for more ambitious contributions. This gives other contributors a chance to point you in the right direction, 
give you feedback on your design, and help you find out if someone else is working on the same thing.

Write code!

[<img src="_asset/misc_arrow-up.png" align="right">](#top)
## License

This product is licensed to you under the Apache License, Version 2.0 (the "License"); you may not use this product 
except in compliance with the License. See [LICENSE](LICENSE) and [NOTICE](NOTICE) for more information.

[<img src="_asset/misc_arrow-up.png" align="right">](#top)
## Trademark

Alsanium, eawsy, the "Created by eawsy" logo, and the "eawsy" logo are trademarks of Alsanium, SAS. or its affiliates in 
France and/or other countries.

Amazon Web Services, the "Powered by Amazon Web Services" logo, and AWS Lambda are trademarks of Amazon.com, Inc. or its 
affiliates in the United States and/or other countries.


[eawsy-home]: https://eawsy.com
[eawsy-github]: https://github.com/eawsy
[eawsy-runtime]: https://github.com/eawsy/aws-lambda-go-shim
[eawsy-chat]: https://gitter.im/eawsy/bavardage
[eawsy-twitter]: https://twitter.com/@eawsyhq
[eawsy-api]: https://godoc.org/github.com/eawsy/aws-cloudformation-go-customres/service/cloudformation/customres
[eawsy-hire]: https://docs.google.com/forms/d/e/1FAIpQLSfPvn1Dgp95DXfvr3ClPHCNF5abi4D1grveT5btVyBHUk0nXw/viewform
[eawsy-pr]: https://github.com/eawsy/aws-lambda-go-net/issues?q=is:pr%20is:open
[eawsy-issues]: https://github.com/eawsy/aws-lambda-go-net/issues?q=is:issue%20is:open

[aws-home]: https://aws.amazon.com/
[aws-lambda-home]: https://aws.amazon.com/lambda/
[aws-cloudformation-home]: https://aws.amazon.com/cloudformation/
[aws-cloudformation-custom]: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/template-custom-resources.html

[badge-api]: http://img.shields.io/badge/api-godoc-3F51B5.svg?style=flat-square
[badge-status]: http://img.shields.io/badge/status-stable-4CAF50.svg?style=flat-square
[badge-license]: http://img.shields.io/badge/license-apache-FF5722.svg?style=flat-square
[badge-help]: http://img.shields.io/badge/help-gitter-E91E63.svg?style=flat-square
[badge-social]: http://img.shields.io/badge/social-twitter-03A9F4.svg?style=flat-square
