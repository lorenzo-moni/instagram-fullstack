package main

import (
	"bytes"
    "log"
    "net/http"
    "os"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
)

const defaultPort = "8080"

func init() {

    err := godotenv.Load(".env")

    if err != nil {
        log.Fatal("Error loading .env file")
    }
}

func main() {

	AWS_S3_REGION := os.Getenv("AWS_REGION")
	


    session, err := session.NewSession(&aws.Config{Region: aws.String(AWS_S3_REGION)})
    if err != nil {
        log.Fatal(err)
    }

    // Upload Files
    err = uploadFile(session, "test.jpeg")
    if err != nil {
        log.Fatal(err)
    }
}

func uploadFile(session *session.Session, uploadFileDir string) error {

	AWS_S3_BUCKET := os.Getenv("AWS_S3_BUCKET")
    
    upFile, err := os.Open(uploadFileDir)
    if err != nil {
        return err
    }
    defer upFile.Close()
    
    upFileInfo, _ := upFile.Stat()
    var fileSize int64 = upFileInfo.Size()
    fileBuffer := make([]byte, fileSize)
    upFile.Read(fileBuffer)
    
    _, err = s3.New(session).PutObject(&s3.PutObjectInput{
        Bucket:               aws.String(AWS_S3_BUCKET),
        Key:                  aws.String(uploadFileDir),
        ACL:                  aws.String("private"),
        Body:                 bytes.NewReader(fileBuffer),
        ContentLength:        aws.Int64(fileSize),
        ContentType:          aws.String(http.DetectContentType(fileBuffer)),
        ContentDisposition:   aws.String("attachment"),
        ServerSideEncryption: aws.String("AES256"),
    })
    return err
}