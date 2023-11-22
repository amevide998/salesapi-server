package Controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	"math/rand"
	"net/http"
	"os"
	"sales-api/Middleware"
	"sales-api/Model"
	"sales-api/Response"
	db "sales-api/config"
	"sales-api/dto"
	"strconv"
	"strings"
	"time"
)

// order controller
// @Description create order
// @Summary create order
// @Tags Order
// @Produce json
// @Param Authorization header string true "authorization"
// @Success 200 {object} Response.WebResponse[dto.OrderResponse]
// @Router /orders [post]
func CreateOrder(c *fiber.Ctx) error {
	//Token authenticate
	headerToken := c.Get("Authorization")
	if headerToken == "" {
		response := Response.NewWebErrorResponse("Unauthorized")
		return c.Status(401).JSON(response)
	}
	if err := Middleware.AuthenticateToken(Middleware.SplitToken(headerToken)); err != nil {
		response := Response.NewWebErrorResponse("Unauthorized")
		return c.Status(401).JSON(response)
	}
	//Token authenticate

	body := struct {
		PaymentId int            `json:"paymentId"`
		TotalPaid int            `json:"totalPaid"`
		Products  []dto.Products `json:"products"`
	}{}

	if err := c.BodyParser(&body); err != nil {
		response := Response.NewWebErrorResponse("Empty Body")
		return c.Status(404).JSON(response)
	}

	Prodresponse := make([]*Model.ProductResponseOrder, 0)

	var TotalInvoicePrice = struct {
		ttprice int
	}{}

	productsIds := ""
	quantities := ""
	for _, v := range body.Products {
		totalPrice := 0
		productsIds = productsIds + "," + strconv.Itoa(v.ProductId)
		quantities = quantities + "," + strconv.Itoa(v.Quantity)

		prods := Model.ProductOrder{}
		var discount Model.Discount
		db.DB.Table("products").Where("id=?", v.ProductId).First(&prods)
		db.DB.Where("id = ?", prods.DiscountId).Find(&discount)
		discCount := 0

		if discount.Type == "BUY_N" {
			totalPrice = prods.Price * v.Quantity

			discCount = totalPrice - discount.Result
			TotalInvoicePrice.ttprice = TotalInvoicePrice.ttprice + discCount

		}

		if discount.Type == "PERCENT" {
			totalPrice = prods.Price * v.Quantity
			percentage := totalPrice * discount.Result / 100
			discCount = totalPrice - percentage
			TotalInvoicePrice.ttprice = TotalInvoicePrice.ttprice + discCount
		}

		Prodresponse = append(Prodresponse,
			&Model.ProductResponseOrder{
				ProductId:        prods.Id,
				Name:             prods.Name,
				Price:            prods.Price,
				Discount:         discount,
				Qty:              v.Quantity,
				TotalNormalPrice: prods.Price,
				TotalFinalPrice:  discCount,
			},
		)

	}
	orderResp := Model.Order{
		CashierID:      1,
		PaymentTypesId: body.PaymentId,
		TotalPrice:     TotalInvoicePrice.ttprice,
		TotalPaid:      body.TotalPaid,
		TotalReturn:    body.TotalPaid - TotalInvoicePrice.ttprice,
		ReceiptId:      "R000" + strconv.Itoa(rand.Intn(1000)),
		ProductId:      productsIds,
		Quantities:     quantities,
		UpdatedAt:      time.Now().UTC(),
		CreatedAt:      time.Now().UTC(),
	}
	db.DB.Create(&orderResp)

	orderResponse := dto.OrderResponse{
		Order:    orderResp,
		Products: Prodresponse,
	}
	response := Response.NewWebResponse(orderResponse)
	return c.Status(http.StatusCreated).JSON(response)
}

// order controller
// @Description subtotal orders
// @Summary subtotal orders
// @Tags Order
// @Produce json
// @Param Authorization header string true "authorization"
// @Success 200 {object} Response.WebResponse[dto.OrderResponse]
// @Router /orders/subtotal [post]
func SubTotalOrder(c *fiber.Ctx) error {
	//Token authenticate
	headerToken := c.Get("Authorization")
	if headerToken == "" {
		response := Response.NewWebErrorResponse("Unauthorized")
		return c.Status(401).JSON(response)
	}
	if err := Middleware.AuthenticateToken(Middleware.SplitToken(headerToken)); err != nil {
		response := Response.NewWebErrorResponse("Unauthorized")
		return c.Status(401).JSON(response)
	}
	//Token authenticate

	type products struct {
		ProductId int `json:"productId"`
		Quantity  int `json:"qty"`
	}

	body := struct {
		Products []products `json:"products"`
	}{}

	if err := c.BodyParser(&body.Products); err != nil {
		response := Response.NewWebErrorResponse(err.Error(), "empty body")
		return c.Status(400).JSON(response)
	}

	Prodresponse := make([]*Model.ProductResponseOrder, 0)

	var TotalInvoicePrice = struct {
		ttprice int
	}{}

	for _, v := range body.Products {
		totalPrice := 0

		prods := Model.ProductOrder{}
		var discount Model.Discount
		db.DB.Table("products").Where("id=?", v.ProductId).First(&prods)
		db.DB.Where("id = ?", prods.DiscountId).Find(&discount)

		disc := 0
		if discount.Type == "PERCENT" {
			totalPrice = prods.Price * v.Quantity //5000*3=15000
			percentage := totalPrice * discount.Result / 100
			disc = totalPrice - percentage
			TotalInvoicePrice.ttprice = TotalInvoicePrice.ttprice + disc
		}
		if discount.Type == "BUY_N" {
			totalPrice = prods.Price * v.Quantity //5000*3=15000
			disc = totalPrice - discount.Result
			TotalInvoicePrice.ttprice = TotalInvoicePrice.ttprice + disc

		}

		Prodresponse = append(Prodresponse,
			&Model.ProductResponseOrder{
				ProductId:        prods.Id,
				Name:             prods.Name,
				Price:            prods.Price,
				Discount:         discount,
				Qty:              v.Quantity,
				TotalNormalPrice: prods.Price,
				TotalFinalPrice:  disc,
			},
		)

	}

	subTotalResponse := dto.SubTotalResponse{
		SubTotal: TotalInvoicePrice.ttprice,
		Products: Prodresponse,
	}

	response := Response.NewWebResponse(subTotalResponse)
	return c.Status(200).JSON(response)
}

// order controller
// @Description checkout order
// @Summary checkout order
// @Tags Order
// @Produce json
// @Param Authorization header string true "authorization"
// @Param orderId path string true "order id"
// @Success 200 {object} Response.WebResponse[dto.CheckOrder]
// @Router /orders/{orderId}/check-download [get]
func CheckOrder(c *fiber.Ctx) error {
	param := c.Params("orderId")

	var order Model.Order
	db.DB.Where("id=?", param).First(&order)
	if order.Id == 0 {
		response := Response.NewWebErrorResponse("order not exists")
		return c.Status(404).JSON(response)
	}

	if order.IsDownload == 0 {
		data := dto.CheckOrder{
			IsDownloaded: false,
		}
		response := Response.NewWebResponse(data)
		return c.Status(200).JSON(response)
	}

	if order.IsDownload == 1 {
		data := dto.CheckOrder{
			IsDownloaded: true,
		}
		response := Response.NewWebResponse(data)
		return c.Status(200).JSON(response)
	}

	return nil

}

// order controller
// @Description  order detail
// @Summary order detail
// @Tags Order
// @Produce json
// @Param Authorization header string true "authorization"
// @Param orderId path string true "order id"
// @Success 200 {object} Response.WebResponse[dto.OrderDetailResponse]
// @Router /orders/{orderId} [get]
func OrderDetail(c *fiber.Ctx) error {
	//Token authenticate
	headerToken := c.Get("Authorization")
	if headerToken == "" {
		response := Response.NewWebErrorResponse("Unauthorized")
		return c.Status(401).JSON(response)
	}
	if err := Middleware.AuthenticateToken(Middleware.SplitToken(headerToken)); err != nil {
		response := Response.NewWebErrorResponse("Unauthorized")
		return c.Status(401).JSON(response)
	}
	//Token authenticate

	param := c.Params("orderId")

	var order Model.Order
	db.DB.Where("id=?", param).First(&order)

	if order.Id == 0 {
		response := Response.NewWebErrorResponse("not found")
		return c.Status(404).JSON(response)
	}
	productIds := strings.Split(order.ProductId, ",")
	TotalProducts := make([]*Model.Product, 0)

	for i := 1; i < len(productIds); i++ {
		prods := Model.Product{}
		db.DB.Where("id = ?", productIds[i]).Find(&prods)
		TotalProducts = append(TotalProducts, &prods)
	}
	cashier := Model.Cashier{}
	db.DB.Where("id = ?", order.CashierID).Find(&cashier)

	paymentType := Model.PaymentType{}
	db.DB.Where("id = ?", order.PaymentTypesId).Find(&paymentType)

	orderTable := Model.Order{}
	db.DB.Where("id = ?", order.Id).Find(&orderTable)

	orderDetail := dto.OrderDetail{
		OrderId:        order.Id,
		CashierId:      order.CashierID,
		PaymentTypesId: order.PaymentTypesId,
		TotalPrice:     order.TotalPrice,
		TotalPaid:      order.TotalPaid,
		TotalReturn:    order.TotalReturn,
		ReceiptId:      order.ReceiptId,
		CreatedAt:      order.CreatedAt,
		Cashier:        cashier,
		PaymentType:    paymentType,
	}

	responseData := dto.OrderDetailResponse{
		Orders:   orderDetail,
		Products: TotalProducts,
	}

	response := Response.NewWebResponse(responseData)

	return c.Status(200).JSON(response)

}

// order controller
// @Description order list
// @Summary order list
// @Tags Order
// @Produce json
// @Param Authorization header string true "authorization"
// @Success 200 {object} Response.WebResponse[dto.OrderListResponse]
// @Router /orders [get]
func OrdersList(c *fiber.Ctx) error {

	//Token authenticate
	headerToken := c.Get("Authorization")
	if headerToken == "" {
		response := Response.NewWebErrorResponse("Unauthorized")
		return c.Status(401).JSON(response)
	}
	if err := Middleware.AuthenticateToken(Middleware.SplitToken(headerToken)); err != nil {
		response := Response.NewWebErrorResponse("Unauthorized")
		return c.Status(401).JSON(response)
	}
	//Token authenticate

	limit, _ := strconv.Atoi(c.Query("limit"))
	skip, _ := strconv.Atoi(c.Query("skip"))
	var count int64
	var order []Model.Order

	db.DB.Select("*").Limit(limit).Offset(skip).Find(&order).Count(&count)

	OrderResponse := make([]*dto.OrderList, 0)

	for _, v := range order {
		cashier := Model.Cashier{}
		db.DB.Where("id = ?", v.CashierID).Find(&cashier)
		paymentType := Model.PaymentType{}
		db.DB.Where("id = ?", v.PaymentTypesId).Find(&paymentType)

		OrderResponse = append(OrderResponse, &dto.OrderList{
			OrderId:        v.Id,
			CashierID:      v.CashierID,
			PaymentTypesId: v.PaymentTypesId,
			TotalPrice:     v.TotalPrice,
			TotalPaid:      v.TotalPaid,
			TotalReturn:    v.TotalReturn,
			ReceiptId:      v.ReceiptId,
			CreatedAt:      v.CreatedAt,
			Payments:       paymentType,
			Cashiers:       cashier,
		})

	}

	orderListResponse := dto.OrderListResponse{
		Order: OrderResponse,
		Meta: dto.Pagination{
			Total: count,
			Limit: limit,
			Skip:  skip,
		},
	}

	response := Response.NewWebResponse(orderListResponse)

	return c.Status(404).JSON(response)
}

// order controller
// @Description order list
// @Summary order list
// @Tags Order
// @Produce json
// @Param orderId path string true "order id"
// @Param Authorization header string true "authorization"
// @Success 200 {object} Response.WebResponse[dto.OrderListResponse]
// @Router /orders/{orderId}/download [get]
func DownloadOrder(c *fiber.Ctx) error {
	//Token authenticate
	headerToken := c.Get("Authorization")
	if headerToken == "" {
		response := Response.NewWebErrorResponse("Unauthorized")
		return c.Status(401).JSON(response)
	}
	if err := Middleware.AuthenticateToken(Middleware.SplitToken(headerToken)); err != nil {
		response := Response.NewWebErrorResponse("Unauthorized")
		return c.Status(401).JSON(response)
	}
	//Token authenticate

	param := c.Params("orderId")

	var order Model.Order
	db.DB.Where("id=?", param).First(&order)

	if order.Id == 0 {
		response := Response.NewWebErrorResponse("Order not found")
		return c.Status(404).JSON(response)
	}
	productIds := strings.Split(order.ProductId, ",")

	TotalProducts := make([]*Model.Product, 0)

	for i := 1; i < len(productIds); i++ {
		prods := Model.Product{}
		db.DB.Where("id = ?", productIds[i]).Find(&prods)
		TotalProducts = append(TotalProducts, &prods)
	}
	cashier := Model.Cashier{}
	db.DB.Where("id = ?", order.CashierID).Find(&cashier)
	paymentType := Model.PaymentType{}

	db.DB.Where("id = ?", order.PaymentTypesId).Find(&paymentType)
	orderTable := Model.Order{}
	db.DB.Where("id = ?", order.Id).Find(&orderTable)

	///pdf Generating
	twoDarray := [][]string{{}}
	quantities := strings.Split(order.Quantities, ",")
	quantities = quantities[1:]
	for i := 0; i < len(TotalProducts); i++ {

		s_array := []string{}
		s_array = append(s_array, TotalProducts[i].Sku)

		s_array = append(s_array, TotalProducts[i].Name)
		s_array = append(s_array, quantities[i])
		s_array = append(s_array, strconv.Itoa(TotalProducts[i].Price))
		twoDarray = append(twoDarray, s_array)

	}

	begin := time.Now()
	grayColor := getGrayColor()
	whiteColor := color.NewWhite()
	header := getHeader()
	contents := twoDarray

	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(10, 15, 10)
	//m.SetBorder(true)

	//Top Heading
	m.SetBackgroundColor(grayColor)
	m.Row(10, func() {
		m.Col(12, func() {
			m.Text("Order Invoice #"+strconv.Itoa(order.Id), props.Text{
				Top:   3,
				Style: consts.Bold,
				Align: consts.Center,
			})
		})
	})
	m.SetBackgroundColor(whiteColor)

	//Table setting
	m.TableList(header, contents, props.TableList{
		HeaderProp: props.TableListContent{
			Size:      9,
			GridSizes: []uint{3, 4, 2, 3},
		},
		ContentProp: props.TableListContent{
			Size:      8,
			GridSizes: []uint{3, 4, 2, 3},
		},
		Align:                consts.Center,
		AlternatedBackground: &grayColor,
		HeaderContentSpace:   1,
		Line:                 false,
	})
	//Total price
	m.Row(20, func() {
		m.ColSpace(7)
		m.Col(2, func() {
			m.Text("Total:", props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  8,
				Align: consts.Right,
			})
		})
		m.Col(3, func() {
			m.Text("RS. "+strconv.Itoa(order.TotalPrice), props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  8,
				Align: consts.Center,
			})
		})
	})
	m.Row(21, func() {
		m.ColSpace(7)
		m.Col(2, func() {
			m.Text("Total Paid:", props.Text{
				Top:   0.5,
				Style: consts.Bold,
				Size:  8,
				Align: consts.Right,
			})
		})
		m.Col(3, func() {
			m.Text("RS. "+strconv.Itoa(order.TotalPaid), props.Text{
				Top:   0.5,
				Style: consts.Bold,
				Size:  8,
				Align: consts.Center,
			})
		})
	})

	m.Row(22, func() {
		m.ColSpace(7)
		m.Col(2, func() {
			m.Text("Total Return", props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  8,
				Align: consts.Right,
			})
		})
		m.Col(3, func() {
			m.Text("RS. "+strconv.Itoa(order.TotalReturn), props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  8,
				Align: consts.Center,
			})
		})
	})

	//Invoice creation
	currentTime := time.Now()
	pdfFileName := "invoice-" + currentTime.Format("2006-Jan-02")
	err := m.OutputFileAndClose(pdfFileName + ".pdf")
	if err != nil {
		fmt.Println("Could not save PDF:", err)
		os.Exit(1)
	}

	end := time.Now()
	fmt.Println(end.Sub(begin))

	//update recepit is downloaded to 1 means true
	db.DB.Table("orders").Where("id=?", order.Id).Update("is_download", 1)
	response := Response.NewWebResponse("")
	return c.Status(200).JSON(response)

}

func getHeader() []string {
	return []string{"Product Sku", "Name", "Qty", "Price"}
}

func getGrayColor() color.Color {
	return color.Color{
		Red:   200,
		Green: 200,
		Blue:  200,
	}
}
