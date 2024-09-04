package database

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/joho/godotenv"
)

func init() {
	fmt.Println("checking for aws variables")
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
		panic("Error Loading .env File for AWS, Check If It Exists.")
	}
}

const (
	AWS_S3_REGION = "us-east-1"
	AWS_S3_BUCKET = "bucketregion1"
)

var Sess = ClientAWS()

func ClientAWS() *session.Session {
	// region := AWS_S3_REGION
	// if region == "US" || region == "" {
	//     region = "us-east-1"
	// }
	// fmt.Println("region is", region)

	// sess, err := session.NewSession(
	// 	&aws.Config{
	// 		// Region: aws.String(os.Getenv("AWS_S3_REGION")),
	// 		Region: aws.String(region),
	// 	},
	// )
	// if err != nil {
	// 	panic(err)
	// }
	// return sess
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"), // Ensure this is correct
		LogLevel: aws.LogLevel(aws.LogDebugWithHTTPBody),
	})
	if err != nil {
		panic(err)
	}
	return sess

}
