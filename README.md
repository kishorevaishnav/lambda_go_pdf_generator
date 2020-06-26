# lambda_go_pdf_generator

GOOS=linux GOARCH=amd64 go build -o pdfGenerator pdfGenerator2.go

zip pdfGenerator.zip wkhtmltopdf pdfGenerator
