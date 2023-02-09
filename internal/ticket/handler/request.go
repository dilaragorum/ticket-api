package handler

type CreateTicketOptionRequestBody struct {
	Name       string `json:"name"`
	Desc       string `json:"desc"`
	Allocation int    `json:"allocation"`
}

type CreatePurchaseTicketOptionRequestBody struct {
	Quantity int    `json:"quantity"`
	UserID   string `json:"user_id"`
}
