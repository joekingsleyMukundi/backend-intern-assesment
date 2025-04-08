package mpesa

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type STKPushResponse struct {
	MerchantRequestID   string `json:"MerchantRequestID"`
	CheckoutRequestID   string `json:"CheckoutRequestID"`
	ResponseCode        string `json:"ResponseCode"`
	ResponseDescription string `json:"ResponseDescription"`
	CustomerMessage     string `json:"CustomerMessage"`
}

func InitializeSTKPush(accessToken, phoneNumber string, amount int64) (*STKPushResponse, error) {
	url := "https://sandbox.safaricom.co.ke/mpesa/stkpush/v1/processrequest"
	auth := "Bearer " + accessToken
	timestamp := time.Now().Format("20060102150405")

	businessShortCode := "174379"
	passkey := ""
	password := base64.StdEncoding.EncodeToString([]byte(businessShortCode + passkey + timestamp))

	if len(phoneNumber) > 0 && phoneNumber[0] == '+' {
		phoneNumber = phoneNumber[1:]
	}
	if len(phoneNumber) == 9 {
		phoneNumber = "254" + phoneNumber
	} else if len(phoneNumber) == 10 && phoneNumber[0] == '0' {
		phoneNumber = "254" + phoneNumber[1:]
	}
	payload := map[string]any{
		"BusinessShortCode": businessShortCode,
		"Password":          password,
		"Timestamp":         timestamp,
		"TransactionType":   "CustomerPayBillOnline",
		"Amount":            amount,
		"PartyA":            phoneNumber,
		"PartyB":            businessShortCode,
		"PhoneNumber":       phoneNumber,
		"CallBackURL":       "https://your-callback-url.com/mpesa/callback",
		"AccountReference":  "test",
		"TransactionDesc":   "Payment",
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("Authorization", auth)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send STK push: %w", err)
	}
	defer resp.Body.Close()
	var stkResp *STKPushResponse
	if err := json.NewDecoder(resp.Body).Decode(&stkResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return stkResp, nil
}
