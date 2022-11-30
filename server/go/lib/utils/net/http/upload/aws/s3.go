package aws

import (
	"bufio"
	"context"
	"os"
	"time"

	"github.com/actliboy/hoper/server/go/lib/utils/log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func News3() *s3.S3 {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Endpoint: aws.String("http://oss-cn-hangzhou.aliyuncs.com"),
			Region:   aws.String("oss")},
		Profile: "aliyun",
	}))

	svc := s3.New(sess)
	resp, _ := svc.ListBuckets(&s3.ListBucketsInput{})
	for _, bucket := range resp.Buckets {
		log.Info(*bucket.Name)
	}
	return svc
}

func upload(service *s3.S3) {
	fp, _ := os.Open("s3_test.go")

	defer fp.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(30)*time.Second)
	defer cancel()

	service.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String("hatlonely"),
		Key:    aws.String("hoper/s3_test.go"),
		Body:   fp,
	})

}

func download(service *s3.S3) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(30)*time.Second)
	defer cancel()

	out, _ := service.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String("hatlonely"),
		Key:    aws.String("hoper/s3_test.go"),
	})

	defer out.Body.Close()
	scanner := bufio.NewScanner(out.Body)
	for scanner.Scan() {
		log.Info(scanner.Text())
	}
}

// 目录遍历
func ListObjectsPages(service *s3.S3) {
	var objkeys []string

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(30)*time.Second)
	defer cancel()

	service.ListObjectsPagesWithContext(ctx, &s3.ListObjectsInput{
		Bucket: aws.String("hatlonely"),
		Prefix: aws.String("hoper/"),
	}, func(output *s3.ListObjectsOutput, b bool) bool {
		for _, content := range output.Contents {
			objkeys = append(objkeys, aws.StringValue(content.Key))
		}
		return true
	})
	log.Info(objkeys)
}
