package main

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/url"

	wkhtml "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func PDFGenerator(url string) {

	// Create new PDF generator
	pdfg, err := wkhtml.NewPDFGenerator()
	if err != nil {
		log.Fatal(err)
	}

	// Set global options
	pdfg.MarginTop.Set(10)
	pdfg.MarginBottom.Set(10)
	pdfg.MarginLeft.Set(10)
	pdfg.MarginRight.Set(10)

	// pdfg.Dpi.Set(300)
	pdfg.Orientation.Set(wkhtml.OrientationPortrait)
	// pdfg.Grayscale.Set(false)

	// Create a new input page from an URL
	page := wkhtml.NewPage(url)

	// html := "<html>Hi</html>"
	// pdfgen.AddPage(NewPageReader(strings.NewReader(html)))

	// Set options for this page

	page.FooterRight.Set("Page [page] #{footer}")
	page.FooterFontSize.Set(9)
	page.FooterSpacing.Set(3)

	// Add to document
	pdfg.AddPage(page)

	// Create PDF document in internal buffer
	err = pdfg.Create()
	if err != nil {
		log.Fatal(err)
	}

	// Write buffer contents to file on disk
	err = pdfg.WriteFile("/tmp/hello.pdf")
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println("Done")
	// Output: Done
}

func main() {
	lambda.Start(LambdaHandler)
}

func LambdaHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	param1, found := request.QueryStringParameters["url1"]
	if found {
		// query parameters are typically URL encoded so to get the value
		value, err := url.QueryUnescape(param1)
		if nil != err {
			return events.APIGatewayProxyResponse{
				Body:       string(err.Error()),
				StatusCode: 500,
			}, err
		}
		// ... now use the value as needed
		PDFGenerator(value)
	}
	data, err := ioutil.ReadFile("/tmp/hello.pdf")
	if err != nil {
		panic(err)
	}

	return events.APIGatewayProxyResponse{
		Body:            base64.StdEncoding.EncodeToString(data),
		IsBase64Encoded: true,
		MultiValueHeaders: map[string][]string{
			"Content-Type": []string{"application/pdf"},
		},
		StatusCode: 200,
	}, nil
}
