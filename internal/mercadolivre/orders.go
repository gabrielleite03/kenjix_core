package mercadolivre

import "context"

type Order struct {
	ID     int64   `json:"id"`
	Status string  `json:"status"`
	Total  float64 `json:"total_amount"`
}

func (c *Client) GetOrders(ctx context.Context, sellerID string) ([]Order, error) {
	var resp struct {
		Results []Order `json:"results"`
	}

	err := c.Get(ctx, "/orders/search?seller="+sellerID, &resp)
	return resp.Results, err
}

func (c *Client) GetOrderByID(ctx context.Context, id string) (*Order, error) {
	var order Order
	err := c.Get(ctx, "/orders/"+id, &order)
	return &order, err
}
