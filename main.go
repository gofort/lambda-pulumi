package main

import (
	"os"

	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		// Create an IAM role.
		role, err := iam.NewRole(ctx, "task-exec-role", &iam.RoleArgs{
			AssumeRolePolicy: pulumi.String(`{
				"Version": "2012-10-17",
				"Statement": [{
					"Sid": "",
					"Effect": "Allow",
					"Principal": {
						"Service": "lambda.amazonaws.com"
					},
					"Action": "sts:AssumeRole"
				}]
			}`),
		})
		if err != nil {
			return err
		}

		// Set arguments for constructing the function resource.
		args := &lambda.FunctionArgs{
			Handler:  pulumi.String("main-" + os.Getenv("GITHUB_SHA")),
			Runtime:  pulumi.String("go1.x"),
			S3Bucket: pulumi.String("hello-world-zip"),
			S3Key:    pulumi.String("main-" + os.Getenv("GITHUB_SHA") + ".zip"),
			Role:     role.Arn,
		}

		// Create the lambda using the args.
		function, err := lambda.NewFunction(
			ctx,
			"hello-world",
			args,
		)
		if err != nil {
			return err
		}

		// Export the lambda ARN.
		ctx.Export("lambda", function.Arn)

		return nil
	})
}
