package mercadolivre

type Order struct {
	ID           int64   `json:"id"`
	Status       string  `json:"status"`
	StatusDetail *string `json:"status_detail"`

	DateCreated string `json:"date_created"`
	DateClosed  string `json:"date_closed"`
	LastUpdated string `json:"last_updated"`

	PackID      *int64  `json:"pack_id"`
	TotalAmount float64 `json:"total_amount"`
	PaidAmount  float64 `json:"paid_amount"`
	CurrencyID  string  `json:"currency_id"`

	OrderItems []OrderItem `json:"order_items"`
	Payments   []Payment   `json:"payments"`

	Shipping Shipping `json:"shipping"`
	Buyer    Person   `json:"buyer"`
	Seller   Person   `json:"seller"`

	Tags []string `json:"tags"`
}

type OrderItem struct {
	Item Item `json:"item"`

	Quantity          int               `json:"quantity"`
	RequestedQuantity RequestedQuantity `json:"requested_quantity"`
	UnitPrice         float64           `json:"unit_price"`
	CurrencyID        string            `json:"currency_id"`
	SaleFee           float64           `json:"sale_fee"`
	ListingTypeID     string            `json:"listing_type_id"`
}

type Item struct {
	ID                  string               `json:"id"`
	Title               string               `json:"title"`
	CategoryID          string               `json:"category_id"`
	VariationID         *int64               `json:"variation_id"`
	SellerCustomField   *string              `json:"seller_custom_field"`
	SellerSKU           *string              `json:"seller_sku"`
	Warranty            *string              `json:"warranty"`
	Condition           string               `json:"condition"`
	GlobalPrice         *float64             `json:"global_price"`
	NetWeight           *float64             `json:"net_weight"`
	VariationAttributes []VariationAttribute `json:"variation_attributes"`
}

type VariationAttribute struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	ValueID   string `json:"value_id"`
	ValueName string `json:"value_name"`
}

type RequestedQuantity struct {
	Value   int    `json:"value"`
	Measure string `json:"measure"`
}

type Payment struct {
	ID                        int64    `json:"id"`
	OrderID                   int64    `json:"order_id"`
	PayerID                   int64    `json:"payer_id"`
	SiteID                    string   `json:"site_id"`
	Reason                    string   `json:"reason"`
	PaymentMethodID           string   `json:"payment_method_id"`
	PaymentType               string   `json:"payment_type"`
	Status                    string   `json:"status"`
	StatusDetail              string   `json:"status_detail"`
	CurrencyID                string   `json:"currency_id"`
	Installments              int      `json:"installments"`
	TransactionAmount         float64  `json:"transaction_amount"`
	TransactionAmountRefunded float64  `json:"transaction_amount_refunded"`
	ShippingCost              float64  `json:"shipping_cost"`
	TotalPaidAmount           float64  `json:"total_paid_amount"`
	DateApproved              *string  `json:"date_approved"`
	DateCreated               string   `json:"date_created"`
	DateLastModified          string   `json:"date_last_modified"`
	AvailableActions          []string `json:"available_actions"`
}

type Shipping struct {
	ID *int64 `json:"id"`
}

type Person struct {
	ID int64 `json:"id"`
}

func (o *Order) IsPaid() bool {
	if o.Status == "paid" {
		return true
	}

	for _, p := range o.Payments {
		if p.Status == "approved" {
			return true
		}
	}

	return false
}

func (o *Order) IsNoShipping() bool {
	for _, tag := range o.Tags {
		if tag == "no_shipping" {
			return true
		}
	}
	return false
}

func (o *Order) IsDeliveredOrPickedUpManually() bool {
	return o.Status == "paid" && o.IsNoShipping()
}
