package main

import (
	"fmt"
	larkevent "github.com/larksuite/oapi-sdk-go/v3/event"
	"testing"
)

func TestVerifySign(t *testing.T) {
	encryptBody := "46GNOVa+qjnRuNzIM9jaB+T3He9m2wpZc6eOicI7/YB2qEx5C5IOP6MjfTeIpZ800nsZMcEM+RCNiZ6/h1IEPaGqW1koPjE+2zyqu19dbLMkO+hqaAmer+8Y7Tb+VI/7rm1o9R7KfgM56eAiSEkETWpq4zrtkouyE9a5AO9MvEZDN6VlbaAX+B4KqabUEAyhUNTf2lzStPFdZx9iAqYHVs/IQJSE1y+S6U/4qv8KQ8cJTRTDx6RKTBJiE2aLuLZaqcgsN45VgHF1tgWY8yqFc3p/QXS2H51D8J/TZDCjRFiwYOJJRPpvUs5jD9INL5MrCw0qzkH8sFs5bf83hQSbRwyB7Okw7XLCj0jeL8czxa3l7ExzywGAW8auQxbzZZRHgcqj5+YDUG4ZDKx3f5KYYjAyLssvcaYvFwsLmTtE7wUW1EJg4G8PM1rhzVP1kzb9vsV2KxdX5hRxEhPgjlsIPsjviDVKIMn0ypy7vZv2ZWxQr9JQQhrF3JaLVa/OUdInKtQjJWNfbVnWy23OKYA96s7ZhzAJ6G2Sj0OK/TAthh7KvxkT3vc2M+HE5RXEB6qJsuv2luIz1IM/M1CxUuDq6KMgNLDAk436LZZTBneR2UZ+oOHuZjhd/DBzMAUXxU+EirL0CxzK5gGPz/fGTZuXDGA/QLq5kv86EAEGQRWiZr+PMGZxYpntjpnZ6hh63KkeEh+AG6lJrK8DBNs0CBr3f/AqEQC1wpG9fBm5QjhSvRlS0i+0tfqXY3ansGdTnCFckOoNx6m0+l+/FldofuqS+gkV7aBpnjSs4D8YSUSh42307Gwxt0UGWduHpRipSsLbCIFavCjQ+RFQHryCbWsPKUchV23q2fq6X6BZ2nmOHjGzH66rmWJZvT1eWmJ0G8SNpsnfcgEGEr7NrM5VRStk/Hcad/zxOkR2z5g5akF/vPBvKejiMEMlqPIDs2rVpjcfDefp7Kf/2wNlum1/xAMtXz8tcOVzg9Q3kn+8MH8981uE5vIwSCkN0SgpR8A6T1axVCB/1lcb1+Jgtg5h0OpFuL3VDr6g506gT3k79XdrtNoyP9beg/Or5h+ZU9IUPM9n9TQT+GOBNJVC22MhNn19x4nv/c2xEbGjhNVXqp5vETW3RK+MjKnSp7nZOyPxtkGTKVZ3vL6+EaAgnTpV26/j1klZ7EKHnDCFQhy3e0r7fL7goUxnAHAJUEyME+7KrvvQYhxJhHIp8PzHUKn5iURUnt1hBgaLDoa3MzRTpCUvDB0="

	plainBodyData, err := larkevent.EventDecrypt(encryptBody, larkEventEncryptKey)
	if err != nil {
		t.Fatal(err)
	}

	plainBody := string(plainBodyData)
	fmt.Printf("body: %v\n", plainBody)

	targetSign := larkevent.Signature("1677250880", "15935179", larkEventEncryptKey, plainBody)
	fmt.Printf("target sign: %v\n", targetSign)

	sourceSign := "09983655e57e8ae2959d20fda2e45b3220609e544415812fba9a42444911b425"

	if targetSign != sourceSign {
		t.Fail()
	}
}
