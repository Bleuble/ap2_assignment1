package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type HttpPaymentClient struct {
	client  *http.Client
	baseURL string
}

func NewHttpPaymentClient(baseURL string) *HttpPaymentClient {
	return &HttpPaymentClient{
		client: &http.Client{
			Timeout: 2 * time.Second,
		},
		baseURL: baseURL,
	}
}

type paymentPayload struct {
	OrderID string `json:"order_id"`
	Amount  int64  `json:"amount"`
}

type paymentResponse struct {
	Status        string `json:"status"`
	TransactionID string `json:"transaction_id"`
}

func (c *HttpPaymentClient) AuthorizePayment(orderID string, amount int64) (string, error) {
	url := fmt.Sprintf("%s/payments", c.baseURL)

	payload := paymentPayload{
		OrderID: orderID,
		Amount:  amount,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("payment rejected with status %d", resp.StatusCode)
	}

	var res paymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	if res.Status != "Authorized" {
		return "", fmt.Errorf("payment declined")
	}

	return res.TransactionID, nil
}
