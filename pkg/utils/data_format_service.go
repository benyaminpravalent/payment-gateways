package utils

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"payment-gateway/models"
	"payment-gateway/pkg/constants"
)

func DecodeRequest(r *http.Request, request interface{}) error {
	if r.Body == nil {
		return fmt.Errorf("request body is empty")
	}
	defer r.Body.Close() // Close the body after reading

	contentType := r.Header.Get("Content-Type")

	switch contentType {
	case "application/json":
		if err := json.NewDecoder(r.Body).Decode(request); err != nil && err != io.EOF {
			return fmt.Errorf("failed to decode JSON: %v", err)
		}
	case "text/xml", "application/xml":
		if err := xml.NewDecoder(r.Body).Decode(request); err != nil && err != io.EOF {
			return fmt.Errorf("failed to decode XML: %v", err)
		}
	default:
		return fmt.Errorf("unsupported content type: %s", contentType)
	}

	return nil
}

func BuildExternalTransactionRequest(dataFormatSupported, encryptedRequest string) (models.BuildExternalTransaction, error) {
	switch dataFormatSupported {
	case constants.JSON:
		return models.BuildExternalTransaction{
			Request:     dataFormatSupported,
			ContentType: "application/json",
		}, nil

	case constants.SOAP:
		soapEnvelope := `
			<?xml version="1.0" encoding="UTF-8"?>
				<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
					<soap:Body>
						<Request>%s</Request>
					</soap:Body>
				</soap:Envelope>`
		formattedSOAP := fmt.Sprintf(soapEnvelope, encryptedRequest)

		return models.BuildExternalTransaction{
			Request:     formattedSOAP,
			ContentType: "text/xml",
		}, nil
	}

	return models.BuildExternalTransaction{}, fmt.Errorf("unsupported data format: %s", dataFormatSupported)
}
