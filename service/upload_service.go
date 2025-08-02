package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type IUploadService interface {
	GeneratePresignedURL(fileType string) (*PresignedURLResponse, error)
}

type uploadService struct{}

func NewUploadService() IUploadService {
	return &uploadService{}
}

type PresignedURLResponse struct {
	UploadURL string `json:"upload_url"`
	FileURL   string `json:"file_url"`
	FileName  string `json:"file_name"`
	ExpiresIn int64  `json:"expires_in"`
}

func (s *uploadService) GeneratePresignedURL(fileType string) (*PresignedURLResponse, error) {
	// 从环境变量获取R2配置
	var bucketName = "emo-note"
	var accountId = "64249e4e4c53300c8c7444417a5e42b1"
	var accessKeyId = "3a38dcfc50c13767ebaa1ff790c84c3e"
	var accessKeySecret = "c803a1ba9b0fb10dc29035f906bb5853fd8c137d03997f741f41d7195587a25e"

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		log.Fatal(err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountId))
	})

	// 生成唯一文件名
	fileId := uuid.New().String()
	var fileName string
	switch fileType {
	case "image/jpeg", "image/jpg":
		fileName = fmt.Sprintf("avatars/%s.jpg", fileId)
	case "image/png":
		fileName = fmt.Sprintf("avatars/%s.png", fileId)
	case "image/webp":
		fileName = fmt.Sprintf("avatars/%s.webp", fileId)
	default:
		return nil, fmt.Errorf("不支持的文件类型: %s", fileType)
	}

	presignClient := s3.NewPresignClient(client)

	presignResult, err := presignClient.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
	})

	if err != nil {
		panic("Couldn't get presigned URL for PutObject")
	}

	return &PresignedURLResponse{
		UploadURL: presignResult.URL,
		FileURL:   fmt.Sprintf("https://pub-db741f3c485548dfa928869e27c25913.r2.dev/%s", fileName),
		FileName:  fileName,
	}, nil
}

func (s *uploadService) generatePresignedPutURL(endpoint, bucket, key, accessKeyId, secretAccessKey, region string, expiration time.Time, contentType string) (string, error) {
	// AWS Signature Version 4
	method := "PUT"
	service := "s3"

	// 时间格式
	amzDate := expiration.Format("20060102T150405Z")
	dateStamp := expiration.Format("20060102")

	// 构建凭证范围
	credentialScope := fmt.Sprintf("%s/%s/%s/aws4_request", dateStamp, region, service)

	// 构建URL
	u, err := url.Parse(fmt.Sprintf("%s/%s/%s", endpoint, bucket, key))
	if err != nil {
		return "", err
	}

	// 查询参数
	query := url.Values{}
	query.Set("X-Amz-Algorithm", "AWS4-HMAC-SHA256")
	query.Set("X-Amz-Credential", fmt.Sprintf("%s/%s", accessKeyId, credentialScope))
	query.Set("X-Amz-Date", amzDate)
	query.Set("X-Amz-Expires", strconv.FormatInt(900, 10))
	query.Set("X-Amz-SignedHeaders", "content-type;host")

	u.RawQuery = query.Encode()

	// 构建规范请求
	canonicalHeaders := fmt.Sprintf("content-type:%s\nhost:%s\n", contentType, u.Host)
	signedHeaders := "content-type;host"

	canonicalRequest := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\nUNSIGNED-PAYLOAD",
		method,
		u.Path,
		u.RawQuery,
		canonicalHeaders,
		signedHeaders,
	)

	// 构建字符串以签名
	stringToSign := fmt.Sprintf("AWS4-HMAC-SHA256\n%s\n%s\n%s",
		amzDate,
		credentialScope,
		sha256Hex(canonicalRequest),
	)

	// 计算签名
	signature := s.calculateSignature(secretAccessKey, dateStamp, region, service, stringToSign)

	// 添加签名到查询参数
	query.Set("X-Amz-Signature", signature)
	u.RawQuery = query.Encode()

	return u.String(), nil
}

func (s *uploadService) calculateSignature(secretAccessKey, dateStamp, region, service, stringToSign string) string {
	kDate := hmacSHA256([]byte("AWS4"+secretAccessKey), dateStamp)
	kRegion := hmacSHA256(kDate, region)
	kService := hmacSHA256(kRegion, service)
	kSigning := hmacSHA256(kService, "aws4_request")
	signature := hmacSHA256(kSigning, stringToSign)
	return hex.EncodeToString(signature)
}

func hmacSHA256(key []byte, data string) []byte {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return h.Sum(nil)
}

func sha256Hex(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
