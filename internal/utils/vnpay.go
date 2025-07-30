package utils

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
)

func CreateVnpayPaymentURL(txnRef string, bookingID string, amount int, clientIP string, orderType string) (string, error) {
	if amount <= 0 {
		return "", errors.New("error.invalid_amount")
	}
	if amount > 999999999 {
		return "", errors.New("error.amount_exceeds_limit")
	}

	clientIP = strings.TrimSpace(clientIP)
	if clientIP == "" {
		return "", errors.New("error.client_ip_empty")
	}
	if clientIP == "::1" {
		clientIP = "127.0.0.1"
	}
	if strings.Count(clientIP, ".") != 3 && !strings.Contains(clientIP, ":") {
		return "", errors.New("error.invalid_ip_address")
	}
	vnp_TmnCode := os.Getenv("VNPAY_TMN_CODE")
	vnp_HashSecret := os.Getenv("VNPAY_HASH_SECRET")
	vnp_Url := os.Getenv("VNPAY_URL")
	vnp_ReturnURL := os.Getenv("VNPAY_RETURN_URL")
	if vnp_TmnCode == "" || vnp_HashSecret == "" || vnp_Url == "" || vnp_ReturnURL == "" {
		return "", errors.New("error.vnpay_configuration_missing")
	}

	now := time.Now()
	createDate := now.Format("20060102150405")
	expireDate := now.Add(15 * time.Minute).Format("20060102150405")

	params := url.Values{}
	params.Set("vnp_Version", "2.1.0")
	params.Set("vnp_Command", "pay")
	params.Set("vnp_TmnCode", vnp_TmnCode)
	params.Set("vnp_Amount", fmt.Sprintf("%d", amount*100))
	params.Set("vnp_CurrCode", "VND")
	params.Set("vnp_TxnRef", txnRef)
	params.Set("vnp_OrderInfo", url.QueryEscape(fmt.Sprintf("Thanh toan dat phong %s", bookingID)))
	params.Set("vnp_OrderType", orderType)
	params.Set("vnp_Locale", "vn")
	params.Set("vnp_ReturnUrl", vnp_ReturnURL)
	params.Set("vnp_IpAddr", clientIP)
	params.Set("vnp_CreateDate", createDate)
	params.Set("vnp_ExpireDate", expireDate)
	params.Set("vnp_SecureHashType", "HMACSHA512")

	var keys []string
	for k := range params {
		if k != "vnp_SecureHash" && k != "vnp_SecureHashType" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	var rawData strings.Builder
	for i, k := range keys {
		if i > 0 {
			rawData.WriteString("&")
		}
		rawData.WriteString(k + "=" + strings.ReplaceAll(url.QueryEscape(params.Get(k)), "+", "%20"))
	}

	h := hmac.New(sha512.New, []byte(vnp_HashSecret))
	h.Write([]byte(rawData.String()))
	secureHash := strings.ToUpper(hex.EncodeToString(h.Sum(nil)))

	query := rawData.String() + "&vnp_SecureHashType=HMACSHA512&vnp_SecureHash=" + secureHash
	return fmt.Sprintf("%s?%s", vnp_Url, query), nil
}
