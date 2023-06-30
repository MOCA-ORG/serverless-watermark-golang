package main

import (
	"bytes"
	"encoding/base64"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"image"
	"image/draw"
	"image/jpeg"
	_ "image/png"
	"os"
)

type Environment struct {
	bucketName         string
	watermarkImageName string
}

// Get environment variables.
var env = Environment{
	bucketName:         os.Getenv("S3_BUCKET_NAME"),
	watermarkImageName: os.Getenv("WATERMARK_IMAGE_NAME"),
}

// Cache the watermark image across lambda invocations.
var watermarkImage image.Image

// Define common error messages.
var (
	internalServerError = generateErrorResponse(500, "Internal Server Error")
	notFoundError       = generateErrorResponse(404, "Not Found")
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var err error
	var imageKey string
	var baseImageBuffer []byte
	var baseImage image.Image
	var processedImageBuffer []byte
	var processedImage image.Image

	if imageKey = request.PathParameters["imageKey"]; imageKey == "" {
		return notFoundError, nil
	}
	if watermarkImage, err = fetchWatermarkImage(); err != nil {
		return internalServerError, nil
	}
	if baseImageBuffer, err = fetchImage(env.bucketName, imageKey); err != nil {
		return notFoundError, nil
	}
	if baseImage, _, err = decodeImage(baseImageBuffer); err != nil {
		return internalServerError, nil
	}
	processedImage = compositeImages(baseImage, watermarkImage)
	if processedImageBuffer, err = encodeImage(processedImage); err != nil {
		return internalServerError, nil
	}

	return formatResponse(processedImageBuffer), nil
}

// fetchImage fetches an image from S3.
func fetchImage(bucketName string, key string) ([]byte, error) {
	sess := session.Must(session.NewSession())
	svc := s3.New(sess)
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	}
	buf := new(bytes.Buffer)
	var result *s3.GetObjectOutput
	var err error
	if result, err = svc.GetObject(input); err != nil {
		return nil, err
	}
	if _, err = buf.ReadFrom(result.Body); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// fetchWatermarkImage fetches the watermark image from S3.
func fetchWatermarkImage() (image.Image, error) {
	var imgData []byte
	var img image.Image
	var err error

	if watermarkImage != nil {
		return watermarkImage, nil
	}
	if imgData, err = fetchImage(env.bucketName, env.watermarkImageName); err != nil {
		return nil, err
	}
	if img, _, err = decodeImage(imgData); err != nil {
		return nil, err
	}
	watermarkImage = img
	return img, nil
}

// compositeImages composites the watermark image onto the given image.
func compositeImages(imgData image.Image, watermarkData image.Image) image.Image {
	imgSize := imgData.Bounds().Size()
	watermarkSize := watermarkData.Bounds().Size()
	offset := image.Pt((imgSize.X-watermarkSize.X)/2, (imgSize.Y-watermarkSize.Y)/2)
	bound := imgData.Bounds()
	rgba := image.NewRGBA(bound)
	draw.Draw(rgba, bound, imgData, image.Point{}, draw.Src)
	draw.Draw(rgba, watermarkData.Bounds().Add(offset), watermarkData, image.Point{}, draw.Over)
	return rgba
}

// decodeImage decodes the given byte array into an image.
func decodeImage(imgData []byte) (image.Image, string, error) {
	return image.Decode(bytes.NewReader(imgData))
}

// encodeImage encodes the given image into a byte array.
func encodeImage(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, nil)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// formatResponse formats the output image into a base64 encoded string.
func formatResponse(image []byte) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       base64.StdEncoding.EncodeToString(image),
		Headers: map[string]string{
			"Content-Type": "image/jpeg",
		},
		IsBase64Encoded: true,
	}
}

// generateErrorResponse generates an error response with the given status code and message.
func generateErrorResponse(statusCode int, message string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       message,
		Headers: map[string]string{
			"Content-Type": "text/plain",
		},
	}
}

func main() {
	lambda.Start(Handler)
}
