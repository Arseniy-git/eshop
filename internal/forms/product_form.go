package forms

type ProductForm struct {
	Title       string  `form:"title" binding:"required"`
	Description string  `form:"description" binding:"required"`
	Price       float64 `form:"price" binding:"required"`
	Quantity    int     `form:"quantity" binding:"required"`
}
