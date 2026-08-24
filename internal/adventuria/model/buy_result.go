package model

type BuyResultData struct {
	Price int
}

type BuyResult struct {
	data BuyResultData
}

func NewBuyResult(price int) *BuyResult {
	return &BuyResult{
		data: BuyResultData{
			Price: price,
		},
	}
}

func (b *BuyResult) Price() int {
	return b.data.Price
}

func (b *BuyResult) SetPrice(n int) {
	b.data.Price = n
}
