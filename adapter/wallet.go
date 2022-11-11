package adapter

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Wallet struct {
	Client *Client
}

type RetrieveRemittanceRequest struct {
	RemittanceId int64
}

type RefundWalletTransactionToCardRequest struct {
	WalletTransactionId int64
	RefundPrice         float64 `json:"refundPrice"`
}

type RemittanceRequest struct {
	MemberId             int64   `json:"memberId"`
	Price                float64 `json:"price"`
	Description          string  `json:"description"`
	RemittanceReasonType string  `json:"remittanceReasonType"`
}

type CreateWithdrawRequest struct {
	MemberId    int64    `json:"memberId"`
	Price       float64  `json:"price"`
	Description string   `json:"description"`
	Currency    Currency `json:"currency"`
}

type SearchWalletTransactionsRequest struct {
	WalletId              int64
	WalletTransactionType string `schema:"walletTransactionType,omitempty"`
	Page                  int    `schema:"page,omitempty"`
	Size                  int    `schema:"size,omitempty"`
}

type SearchWithdrawsRequest struct {
	MemberId         int64     `schema:"walletId,omitempty"`
	PayoutStatus     string    `schema:"payoutStatus,omitempty"`
	Currency         Currency  `schema:"currency,omitempty"`
	MinWithdrawPrice float64   `schema:"minWithdrawPrice,omitempty"`
	MaxWithdrawPrice float64   `schema:"maxWithdrawPrice,omitempty"`
	MinCreatedDate   time.Time `schema:"minCreatedDate,omitempty"`
	MaxCreatedDate   time.Time `schema:"maxCreatedDate,omitempty"`
	Page             int       `schema:"page,omitempty"`
	Size             int       `schema:"size,omitempty"`
}

type MemberWalletResponse struct {
	Id               *int64        `json:"id"`
	CreatedDate      *TimeResponse `json:"createdDate"`
	UpdatedDate      *TimeResponse `json:"updatedDate"`
	Amount           *float64      `json:"amount"`
	WithdrawalAmount *float64      `json:"withdrawalAmount"`
	Currency         *Currency     `json:"currency"`
	MemberId         *int64        `json:"memberId"`
}

type RemittanceResponse struct {
	Id                   *int64        `json:"id"`
	CreatedDate          *TimeResponse `json:"createdDate"`
	Active               *int          `json:"active"`
	Price                *float64      `json:"price"`
	MemberId             *int64        `json:"memberId"`
	RemittanceType       *string       `json:"remittanceType"`
	RemittanceReasonType *string       `json:"remittanceReasonType"`
	Description          *string       `json:"description"`
}

type WithdrawResponse struct {
	Id           *int64                   `json:"id"`
	CreatedDate  *TimeResponse            `json:"createdDate"`
	Status       *Status                  `json:"status"`
	MemberId     *int64                   `json:"memberId"`
	PayoutId     *int64                   `json:"payoutId"`
	Price        *float64                 `json:"price"`
	Description  *string                  `json:"description"`
	Currency     *Currency                `json:"currency"`
	PayoutStatus *TransactionPayoutStatus `json:"payoutStatus"`
}

type RefundWalletTransactionToCardResponse struct {
	Id                  *int64                                      `json:"id"`
	CreatedDate         *TimeResponse                               `json:"createdDate"`
	RefundStatus        *string                                     `json:"refundStatus"`
	RefundPrice         *float64                                    `json:"refundPrice"`
	AuthCode            *string                                     `json:"authCode"`
	HostReference       *string                                     `json:"hostReference"`
	TransId             *string                                     `json:"transId"`
	TransactionId       *int64                                      `json:"transactionId"`
	WalletTransactionId *int64                                      `json:"walletTransactionId"`
	PaymentError        *PaymentError                               `json:"paymentError"`
	TransactionType     *WalletTransactionRefundCardTransactionType `json:"transactionType"`
}

type SearchWalletTransactionsResponse struct {
	ID                    *int64        `json:"id"`
	CreatedDate           *TimeResponse `json:"createdDate"`
	WalletTransactionType *string       `json:"walletTransactionType"`
	Amount                *float64      `json:"amount"`
	TransactionID         *int64        `json:"transactionId"`
	WalletID              *int64        `json:"walletId"`
}

type ResetMerchantMemberWalletBalanceRequest struct {
	WalletAmount float64 `json:"walletAmount"`
}

type RetrieveWalletTransactionRefundableAmountResponse struct {
	RefundableAmount *float64 `json:"refundableAmount"`
}

func (api *Wallet) RetrieveMemberWallet(ctx context.Context, memberId int64) (*MemberWalletResponse, error) {
	newRequest, err := api.Client.NewRequest(ctx, http.MethodGet, fmt.Sprintf("/wallet/v1/members/%d/wallet", memberId), nil)
	if err != nil {
		return nil, err
	}

	response := &Response[MemberWalletResponse]{}
	err = api.Client.Do(ctx, newRequest, response)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

func (api *Wallet) RetrieveMerchantMemberWallet(ctx context.Context) (*MemberWalletResponse, error) {
	newRequest, err := api.Client.NewRequest(ctx, http.MethodGet, "/wallet/v1/merchants/me/wallet", nil)
	if err != nil {
		return nil, err
	}

	response := &Response[MemberWalletResponse]{}
	err = api.Client.Do(ctx, newRequest, response)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

func (api *Wallet) ResetMerchantMemberWalletBalance(ctx context.Context, request ResetMerchantMemberWalletBalanceRequest) (*MemberWalletResponse, error) {
	newRequest, err := api.Client.NewRequest(ctx, http.MethodPost, "/wallet/v1/merchants/me/wallet/reset-balance", request)
	if err != nil {
		return nil, err
	}

	response := &Response[MemberWalletResponse]{}
	err = api.Client.Do(ctx, newRequest, response)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

//
//func (api *Wallet) SearchWalletTransactions(request SearchWalletTransactionsRequest) (interface{}, error) {
//	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/wallet/v1/wallets/%d/wallet-transactions", api.Opts.BaseURL, request.WalletId), nil)
//	req.URL.RawQuery, _ = QueryParams(request)
//	res := Response[DataResponse[SearchWalletTransactionsResponse]]{}
//	resErr := rest.SendRequest(req, &res, api.Opts)
//	return &res, resErr
//}
//
//func (api *Wallet) RetrieveRefundableAmountOfWalletTransaction(walletTransactionId int64) (interface{}, error) {
//	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/payment/v1/wallet-transactions/%d/refundable-amount", api.Opts.BaseURL, walletTransactionId), nil)
//	res := Response[RetrieveWalletTransactionRefundableAmountResponse]{}
//	resErr := rest.SendRequest(req, &res, api.Opts)
//	return &res, resErr
//}
//
//func (api *Wallet) RefundWalletTransactionToCard(request RefundWalletTransactionToCardRequest) (interface{}, error) {
//	body, _ := PrepareBody(request)
//	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/payment/v1/wallet-transactions/%d/refunds", api.Opts.BaseURL, request.WalletTransactionId), body)
//	res := Response[RefundWalletTransactionToCardResponse]{}
//	resErr := rest.SendRequest(req, &res, api.Opts)
//	return &res, resErr
//}
//
//func (api *Wallet) RetrieveRefundWalletTransactionToCard(walletTransactionId int64) (interface{}, error) {
//	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/payment/v1/wallet-transactions/%d/refunds", api.Opts.BaseURL, walletTransactionId), nil)
//	res := Response[DataResponse[RefundWalletTransactionToCardResponse]]{}
//	resErr := rest.SendRequest(req, &res, api.Opts)
//	return &res, resErr
//}
//
//func (api *Wallet) SendRemittance(request RemittanceRequest) (interface{}, error) {
//	body, _ := PrepareBody(request)
//	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/wallet/v1/remittances/send", api.Opts.BaseURL), body)
//	res := Response[RemittanceResponse]{}
//	resErr := rest.SendRequest(req, &res, api.Opts)
//	return &res, resErr
//}
//
//func (api *Wallet) ReceiveRemittance(request RemittanceRequest) (interface{}, error) {
//	body, _ := PrepareBody(request)
//	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/wallet/v1/remittances/receive", api.Opts.BaseURL), body)
//	res := Response[RemittanceResponse]{}
//	resErr := rest.SendRequest(req, &res, api.Opts)
//	return &res, resErr
//}
//
//func (api *Wallet) RetrieveRemittance(request RetrieveRemittanceRequest) (interface{}, error) {
//	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/wallet/v1/remittances/%d", api.Opts.BaseURL, request.RemittanceId), nil)
//	req.URL.RawQuery, _ = QueryParams(request)
//	res := Response[RemittanceResponse]{}
//	resErr := rest.SendRequest(req, &res, api.Opts)
//	return &res, resErr
//}
//
//func (api *Wallet) CreateWithdraw(request CreateWithdrawRequest) (interface{}, error) {
//	body, _ := PrepareBody(request)
//	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/wallet/v1/withdraws", api.Opts.BaseURL), body)
//	res := Response[WithdrawResponse]{}
//	resErr := rest.SendRequest(req, &res, api.Opts)
//	return &res, resErr
//}
//
//func (api *Wallet) CancelWithdraw(withdrawId int64) (interface{}, error) {
//	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/wallet/v1/withdraws/%d/cancel", api.Opts.BaseURL, withdrawId), nil)
//	res := Response[WithdrawResponse]{}
//	resErr := rest.SendRequest(req, &res, api.Opts)
//	return &res, resErr
//}
//
//func (api *Wallet) RetrieveWithdraw(withdrawId int64) (interface{}, error) {
//	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/wallet/v1/withdraws/%d", api.Opts.BaseURL, withdrawId), nil)
//	res := Response[WithdrawResponse]{}
//	resErr := rest.SendRequest(req, &res, api.Opts)
//	return &res, resErr
//}
//
//func (api *Wallet) SearchWithdraws(request SearchWithdrawsRequest) (interface{}, error) {
//	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/wallet/v1/withdraws", api.Opts.BaseURL), nil)
//	req.URL.RawQuery, _ = QueryParams(request)
//	res := Response[DataResponse[WithdrawResponse]]{}
//	resErr := rest.SendRequest(req, &res, api.Opts)
//	return &res, resErr
//}
