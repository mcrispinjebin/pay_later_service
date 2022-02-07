package models

type Merchant struct {
	MerchantID      int     `json:"merchant_id"`
	MerchantName    string  `json:"merchant_name"`
	DiscountPercent float32 `json:"discount_percent"`
	CreatedAt       int     `json:"created_at"`
	UpdatedAt       int     `json:"updated_at"`
}

type User struct {
	UserID               int     `json:"user_id"`
	UserName             string  `json:"user_name"`
	UserEmail            string  `json:"user_email"`
	CreditLimitOffered   float32 `json:"credit_limit_offered"`
	AvailableCreditLimit float32 `json:"available_credit_limit"`
	CreatedAt            int     `json:"created_at"`
}

type Order struct {
	OrderID     int     `json:"order_id"`
	UserID      int     `json:"user_id"`
	MerchantID  int     `json:"merchant_id"`
	OrderAmount float32 `json:"order_amount"`
	OrderStatus float32 `json:"order_status"`
	CreatedAt   int     `json:"created_at"`
}

type Ledger struct {
	LedgerID  int     `json:"ledger_id"`
	UserID    int     `json:"user_id"`
	Amount    float32 `json:"amount"`
	Status    string  `json:"status"`
	CreatedAt int     `json:"created_at"`
}

type Payout struct {
	PayoutID     int     `json:"payout_id"`
	OrderID      int     `json:"order_id"`
	PayoutAmount float32 `json:"payout_amount"`
	PayoutStatus string  `json:"payout_status"`
	CreatedAt    int     `json:"created_at"`
}
