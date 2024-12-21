package utils

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"payment-gateway/models"
	"payment-gateway/pkg/constants"
)

// decodes the incoming request based on content type
func DecodeRequest(r *http.Request, request *models.SendTransactionRequest) error {
	contentType := r.Header.Get("Content-Type")

	switch contentType {
	case "application/json":
		return json.NewDecoder(r.Body).Decode(request)
	case "text/xml":
		return xml.NewDecoder(r.Body).Decode(request)
	case "application/xml":
		return xml.NewDecoder(r.Body).Decode(request)
	default:
		return fmt.Errorf("unsupported content type")
	}
}

func BuildExternalTransactionRequest(dataFormatSupported, encryptedRequest string) (models.BuildExternalTransaction, error) {
	switch dataFormatSupported {
	case constants.JSON:
		return models.BuildExternalTransaction{
			Request:     dataFormatSupported,
			ContentType: "application/json",
		}, nil

	case constants.SOAP:
		// Construct a SOAP envelope with the encrypted request
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
