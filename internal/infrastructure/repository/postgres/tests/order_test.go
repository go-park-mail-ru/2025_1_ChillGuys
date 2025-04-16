package tests

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	order2 "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/order"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateOrder_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	orderID := uuid.New()
	userID := uuid.New()
	addressID := uuid.New()
	productID := uuid.New()
	itemID := uuid.New()

	req := dto.CreateOrderRepoReq{
		Order: &dto.Order{
			ID:                 orderID,
			UserID:             userID,
			Status:             models.Placed,
			TotalPrice:         100.0,
			TotalPriceDiscount: 90.0,
			AddressID:          addressID,
			Items: []dto.CreateOrderItemDTO{
				{
					ID:        itemID,
					ProductID: productID,
					Price:     50.0,
					Quantity:  2,
				},
			},
		},
		UpdatedQuantities: map[uuid.UUID]uint{
			productID: 5,
		},
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO bazaar.order").
		WithArgs(
			orderID,
			userID,
			"placed",
			float64(100),
			float64(90),
			addressID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("UPDATE bazaar.product SET quantity").
		WithArgs(uint(5), productID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO bazaar.order_item").
		WithArgs(
			itemID,
			orderID,
			productID,
			float64(50),
			uint(2),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := order2.NewOrderRepository(db)
	err = repo.CreateOrder(context.Background(), req)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateOrder_ErrorOnBeginTx(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin().WillReturnError(errors.New("tx error"))

	repo := order2.NewOrderRepository(db)
	err = repo.CreateOrder(context.Background(), dto.CreateOrderRepoReq{})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tx error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateOrder_ErrorOnInsertOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	orderID := uuid.New()
	userID := uuid.New()
	addressID := uuid.New()

	req := dto.CreateOrderRepoReq{
		Order: &dto.Order{
			ID:                 orderID,
			UserID:             userID,
			Status:             models.Placed,
			TotalPrice:         100.0,
			TotalPriceDiscount: 90.0,
			AddressID:          addressID,
		},
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO bazaar.order").
		WithArgs(
			orderID,
			userID,
			"placed",
			float64(100),
			float64(90),
			addressID,
		).
		WillReturnError(errors.New("insert error"))
	mock.ExpectRollback()

	repo := order2.NewOrderRepository(db)
	err = repo.CreateOrder(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insert error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateOrder_ErrorOnUpdateQuantity(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	orderID := uuid.New()
	userID := uuid.New()
	addressID := uuid.New()
	productID := uuid.New()

	req := dto.CreateOrderRepoReq{
		Order: &dto.Order{
			ID:                 orderID,
			UserID:             userID,
			Status:             models.Placed,
			TotalPrice:         100.0,
			TotalPriceDiscount: 90.0,
			AddressID:          addressID,
		},
		UpdatedQuantities: map[uuid.UUID]uint{
			productID: 5,
		},
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO bazaar.order").
		WithArgs(
			orderID,
			userID,
			"placed",
			float64(100),
			float64(90),
			addressID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("UPDATE bazaar.product SET quantity").
		WithArgs(uint(5), productID).
		WillReturnError(errors.New("update error"))
	mock.ExpectRollback()

	repo := order2.NewOrderRepository(db)
	err = repo.CreateOrder(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "update error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateOrder_ErrorOnInsertOrderItem(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	orderID := uuid.New()
	userID := uuid.New()
	addressID := uuid.New()
	productID := uuid.New()
	itemID := uuid.New()

	req := dto.CreateOrderRepoReq{
		Order: &dto.Order{
			ID:                 orderID,
			UserID:             userID,
			Status:             models.Placed,
			TotalPrice:         100.0,
			TotalPriceDiscount: 90.0,
			AddressID:          addressID,
			Items: []dto.CreateOrderItemDTO{
				{
					ID:        itemID,
					ProductID: productID,
					Price:     50.0,
					Quantity:  2,
				},
			},
		},
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO bazaar.order").
		WithArgs(
			orderID,
			userID,
			"placed",
			float64(100),
			float64(90),
			addressID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO bazaar.order_item").
		WithArgs(
			itemID,
			orderID,
			productID,
			float64(50),
			uint(2),
		).
		WillReturnError(errors.New("insert item error"))
	mock.ExpectRollback()

	repo := order2.NewOrderRepository(db)
	err = repo.CreateOrder(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insert item error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateOrder_ErrorOnCommit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	orderID := uuid.New()
	userID := uuid.New()
	addressID := uuid.New()
	productID := uuid.New()
	itemID := uuid.New()

	req := dto.CreateOrderRepoReq{
		Order: &dto.Order{
			ID:                 orderID,
			UserID:             userID,
			Status:             models.Placed,
			TotalPrice:         100.0,
			TotalPriceDiscount: 90.0,
			AddressID:          addressID,
			Items: []dto.CreateOrderItemDTO{
				{
					ID:        itemID,
					ProductID: productID,
					Price:     50.0,
					Quantity:  2,
				},
			},
		},
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO bazaar.order").
		WithArgs(
			orderID,
			userID,
			"placed",
			float64(100),
			float64(90),
			addressID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO bazaar.order_item").
		WithArgs(
			itemID,
			orderID,
			productID,
			float64(50),
			uint(2),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit().WillReturnError(errors.New("commit error"))

	repo := order2.NewOrderRepository(db)
	err = repo.CreateOrder(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "commit error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestProductPrice_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	productID := uuid.New()
	expectedProduct := &models.Product{
		Price:    100.0,
		Status:   models.ProductApproved,
		Quantity: 10,
	}

	rows := sqlmock.NewRows([]string{"price", "status", "quantity"}).
		AddRow(expectedProduct.Price, "approved", expectedProduct.Quantity)

	mock.ExpectQuery("SELECT price, status, quantity FROM bazaar.product").
		WithArgs(productID).
		WillReturnRows(rows)

	repo := order2.NewOrderRepository(db)
	product, err := repo.ProductPrice(context.Background(), productID)

	assert.NoError(t, err)
	assert.Equal(t, expectedProduct, product)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestProductPrice_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	productID := uuid.New()

	mock.ExpectQuery("SELECT price, status, quantity FROM bazaar.product").
		WithArgs(productID).
		WillReturnError(sql.ErrNoRows)

	repo := order2.NewOrderRepository(db)
	product, err := repo.ProductPrice(context.Background(), productID)

	assert.Nil(t, product)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrNotFound))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestProductPrice_InvalidStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	productID := uuid.New()

	rows := sqlmock.NewRows([]string{"price", "status", "quantity"}).
		AddRow(100.0, "invalid_status", 10)

	mock.ExpectQuery("SELECT price, status, quantity FROM bazaar.product").
		WithArgs(productID).
		WillReturnRows(rows)

	repo := order2.NewOrderRepository(db)
	product, err := repo.ProductPrice(context.Background(), productID)

	assert.Nil(t, product)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown product status")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestProductPrice_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	productID := uuid.New()

	mock.ExpectQuery("SELECT price, status, quantity FROM bazaar.product").
		WithArgs(productID).
		WillReturnError(errors.New("database error"))

	repo := order2.NewOrderRepository(db)
	product, err := repo.ProductPrice(context.Background(), productID)

	assert.Nil(t, product)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestProductDiscounts_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	productID := uuid.New()
	now := time.Now()
	expectedDiscounts := []models.ProductDiscount{
		{
			DiscountedPrice:   90.0,
			DiscountStartDate: now.Add(-24 * time.Hour),
			DiscountEndDate:   now.Add(24 * time.Hour),
		},
	}

	rows := sqlmock.NewRows([]string{"discounted_price", "start_date", "end_date"}).
		AddRow(expectedDiscounts[0].DiscountedPrice, expectedDiscounts[0].DiscountStartDate, expectedDiscounts[0].DiscountEndDate)

	mock.ExpectQuery("SELECT discounted_price, start_date, end_date FROM bazaar.discount").
		WithArgs(productID).
		WillReturnRows(rows)

	repo := order2.NewOrderRepository(db)
	discounts, err := repo.ProductDiscounts(context.Background(), productID)

	assert.NoError(t, err)
	assert.Equal(t, expectedDiscounts, discounts)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateProductQuantity_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	productID := uuid.New()
	quantity := uint(5)

	mock.ExpectExec("UPDATE bazaar.product SET quantity").
		WithArgs(quantity, productID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := order2.NewOrderRepository(db)
	err = repo.UpdateProductQuantity(context.Background(), productID, quantity)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateProductQuantity_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	productID := uuid.New()
	quantity := uint(5)

	mock.ExpectExec("UPDATE bazaar.product SET quantity").
		WithArgs(quantity, productID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	repo := order2.NewOrderRepository(db)
	err = repo.UpdateProductQuantity(context.Background(), productID, quantity)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrNotFound))
	assert.NoError(t, mock.ExpectationsWereMet())
}

//func TestGetOrdersByUserID_Success(t *testing.T) {
//  db, mock, err := sqlmock.New()
//  if err != nil {
//    t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//  }
//  defer db.Close()
//
//  userID := uuid.New()
//  orderID := uuid.New()
//  addressID := uuid.New()
//  now := time.Now()
//
//  rows := sqlmock.NewRows([]string{
//    "id", "status", "total_price", "total_price_discount", "address_id",
//    "expected_delivery_at", "actual_delivery_at", "created_at",
//  }).AddRow(
//    orderID, "placed", 100.0, 90.0, addressID, now.Add(24*time.Hour), nil, now,
//  )
//
//  mock.ExpectQuery("SELECT id, status, total_price, total_price_discount, address_id, expected_delivery_at, actual_delivery_at, created_at FROM bazaar.order").
//    WithArgs(userID).
//    WillReturnRows(rows)
//
//  repo := order2.NewOrderRepository(db)
//  orders, err := repo.GetOrdersByUserID(context.Background(), userID)
//
//  assert.NoError(t, err)
//  assert.Len(t, *orders, 1)
//  assert.Equal(t, orderID, (*orders)[0].ID)
//  assert.Equal(t, models.Placed, (*orders)[0].Status)
//  assert.Equal(t, 100.0, (*orders)[0].TotalPrice)
//  assert.Equal(t, 90.0, (*orders)[0].TotalPriceDiscount)
//  assert.Equal(t, addressID, (*orders)[0].AddressID)
//  assert.NotNil(t, (*orders)[0].ExpectedDeliveryAt)
//  assert.Nil(t, (*orders)[0].ActualDeliveryAt)
//  assert.NotNil(t, (*orders)[0].CreatedAt)
//  assert.NoError(t, mock.ExpectationsWereMet())
//}

//func TestGetOrdersByUserID_NotFound(t *testing.T) {
//  db, mock, err := sqlmock.New()
//  if err != nil {
//    t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//  }
//  defer db.Close()
//
//  userID := uuid.New()
//
//  mock.ExpectQuery("SELECT id, status, total_price, total_price_discount, address_id, expected_delivery_at, actual_delivery_at, created_at FROM bazaar.order").
//    WithArgs(userID).
//    WillReturnError(sql.ErrNoRows)
//
//  repo := order2.NewOrderRepository(db)
//  orders, err := repo.GetOrdersByUserID(context.Background(), userID)
//
//  assert.Nil(t, orders)
//  assert.Error(t, err)
//  assert.True(t, errors.Is(err, errs.ErrNotFound))
//  assert.NoError(t, mock.ExpectationsWereMet())
//}

func TestGetOrdersByUserID_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userID := uuid.New()

	mock.ExpectQuery("SELECT id, status, total_price, total_price_discount, address_id, expected_delivery_at, actual_delivery_at, created_at FROM bazaar.order").
		WithArgs(userID).
		WillReturnError(errors.New("database error"))

	repo := order2.NewOrderRepository(db)
	orders, err := repo.GetOrdersByUserID(context.Background(), userID)

	assert.Nil(t, orders)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetOrdersByUserID_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userID := uuid.New()
	orderID := uuid.New()
	addressID := uuid.New()
	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "status", "total_price", "total_price_discount", "address_id",
		"expected_delivery_at", "actual_delivery_at", "created_at",
	}).AddRow(
		orderID, "invalid_status", 100.0, 90.0, addressID, now.Add(24*time.Hour), nil, now,
	)

	mock.ExpectQuery("SELECT id, status, total_price, total_price_discount, address_id, expected_delivery_at, actual_delivery_at, created_at FROM bazaar.order").
		WithArgs(userID).
		WillReturnRows(rows)

	repo := order2.NewOrderRepository(db)
	orders, err := repo.GetOrdersByUserID(context.Background(), userID)

	assert.Nil(t, orders)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown order status")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetOrderProducts_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	orderID := uuid.New()
	productID := uuid.New()

	rows := sqlmock.NewRows([]string{"product_id", "quantity"}).
		AddRow(productID, uint(2))

	mock.ExpectQuery("SELECT product_id, quantity FROM bazaar.order_item").
		WithArgs(orderID).
		WillReturnRows(rows)

	repo := order2.NewOrderRepository(db)
	products, err := repo.GetOrderProducts(context.Background(), orderID)

	assert.NoError(t, err)
	assert.Len(t, *products, 1)
	assert.Equal(t, productID, (*products)[0].ProductID)
	assert.Equal(t, uint(2), (*products)[0].Quantity)
	assert.NoError(t, mock.ExpectationsWereMet())
}

//func TestGetOrderProducts_NotFound(t *testing.T) {
//  db, mock, err := sqlmock.New()
//  if err != nil {
//    t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//  }
//  defer db.Close()
//
//  orderID := uuid.New()
//
//  mock.ExpectQuery("SELECT product_id, quantity FROM bazaar.order_item").
//    WithArgs(orderID).
//    WillReturnError(sql.ErrNoRows)
//
//  repo := order2.NewOrderRepository(db)
//  products, err := repo.GetOrderProducts(context.Background(), orderID)
//
//  assert.Nil(t, products)
//  assert.Error(t, err)
//  assert.True(t, errors.Is(err, errs.ErrNotFound))
//  assert.NoError(t, mock.ExpectationsWereMet())
//}

func TestGetOrderProducts_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	orderID := uuid.New()

	mock.ExpectQuery("SELECT product_id, quantity FROM bazaar.order_item").
		WithArgs(orderID).
		WillReturnError(errors.New("database error"))

	repo := order2.NewOrderRepository(db)
	products, err := repo.GetOrderProducts(context.Background(), orderID)

	assert.Nil(t, products)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetOrderProducts_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	orderID := uuid.New()

	rows := sqlmock.NewRows([]string{"product_id", "quantity"}).
		AddRow("invalid_uuid", "invalid_quantity")

	mock.ExpectQuery("SELECT product_id, quantity FROM bazaar.order_item").
		WithArgs(orderID).
		WillReturnRows(rows)

	repo := order2.NewOrderRepository(db)
	products, err := repo.GetOrderProducts(context.Background(), orderID)

	assert.Nil(t, products)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetProductImage_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	productID := uuid.New()
	imageURL := "http://example.com/image.jpg"

	rows := sqlmock.NewRows([]string{"preview_image_url"}).
		AddRow(imageURL)

	mock.ExpectQuery("SELECT preview_image_url FROM bazaar.product").
		WithArgs(productID).
		WillReturnRows(rows)

	repo := order2.NewOrderRepository(db)
	result, err := repo.GetProductImage(context.Background(), productID)

	assert.NoError(t, err)
	assert.Equal(t, imageURL, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetProductImage_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	productID := uuid.New()

	mock.ExpectQuery("SELECT preview_image_url FROM bazaar.product").
		WithArgs(productID).
		WillReturnError(sql.ErrNoRows)

	repo := order2.NewOrderRepository(db)
	image, err := repo.GetProductImage(context.Background(), productID)

	assert.Empty(t, image)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrNotFound))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetProductImage_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	productID := uuid.New()

	mock.ExpectQuery("SELECT preview_image_url FROM bazaar.product").
		WithArgs(productID).
		WillReturnError(errors.New("database error"))

	repo := order2.NewOrderRepository(db)
	image, err := repo.GetProductImage(context.Background(), productID)

	assert.Empty(t, image)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetOrderAddress_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	addressID := uuid.New()
	expectedAddress := &models.AddressDB{
		Region:        null.StringFrom("Region"),
		City:          null.StringFrom("City"),
		AddressString: null.StringFrom("Street 123"),
		Coordinate:    null.StringFrom("10.0,20.0"),
	}

	rows := sqlmock.NewRows([]string{"region", "city", "address_string", "coordinate"}).
		AddRow(expectedAddress.Region, expectedAddress.City, expectedAddress.AddressString, expectedAddress.Coordinate)

	mock.ExpectQuery("SELECT region, city, address_string, coordinate FROM bazaar.address").
		WithArgs(addressID).
		WillReturnRows(rows)

	repo := order2.NewOrderRepository(db)
	address, err := repo.GetOrderAddress(context.Background(), addressID)

	assert.NoError(t, err)
	assert.Equal(t, expectedAddress, address)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetOrderAddress_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	addressID := uuid.New()

	mock.ExpectQuery("SELECT region, city, address_string, coordinate FROM bazaar.address").
		WithArgs(addressID).
		WillReturnError(sql.ErrNoRows)

	repo := order2.NewOrderRepository(db)
	address, err := repo.GetOrderAddress(context.Background(), addressID)

	assert.Nil(t, address)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrNotFound))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetOrderAddress_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	addressID := uuid.New()

	mock.ExpectQuery("SELECT region, city, address_string, coordinate FROM bazaar.address").
		WithArgs(addressID).
		WillReturnError(errors.New("database error"))

	repo := order2.NewOrderRepository(db)
	address, err := repo.GetOrderAddress(context.Background(), addressID)

	assert.Nil(t, address)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

//func TestProductDiscounts_NotFound(t *testing.T) {
//  db, mock, err := sqlmock.New()
//  if err != nil {
//    t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//  }
//  defer db.Close()
//
//  productID := uuid.New()
//
//  mock.ExpectQuery("SELECT discounted_price, start_date, end_date FROM bazaar.discount").
//    WithArgs(productID).
//    WillReturnError(sql.ErrNoRows)
//
//  repo := order2.NewOrderRepository(db)
//  discounts, err := repo.ProductDiscounts(context.Background(), productID)
//
//  assert.Nil(t, discounts)
//  assert.Error(t, err)
//  assert.True(t, errors.Is(err, errs.ErrNotFound))
//  assert.NoError(t, mock.ExpectationsWereMet())
//}

//func TestProductDiscounts_QueryError(t *testing.T) {
//  db, mock, err := sqlmock.New()
//  if err != nil {
//    t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//  }
//  defer db.Close()
//
//  productID := uuid.New()
//
//  mock.ExpectQuery("SELECT discounted_price, start_date, end_date FROM bazaar.discount").
//    WithArgs(productID).
//    WillReturnError(errors.New("database error"))
//
//  repo := order2.NewOrderRepository(db)
//  discounts, err := repo.ProductDiscounts(context.Background(), productID)
//
//  assert.Nil(t, discounts)
//  assert.Error(t, err)
//  assert.Contains(t, err.Error(), "database error")
//  assert.NoError(t, mock.ExpectationsWereMet())
//}

//func TestProductDiscounts_ScanError(t *testing.T) {
//  db, mock, err := sqlmock.New()
//  if err != nil {
//    t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//  }
//  defer db.Close()
//
//  productID := uuid.New()
//
//  rows := sqlmock.NewRows([]string{"discounted_price", "start_date", "end_date"}).
//    AddRow("invalid_price", time.Now(), time.Now())
//
//  mock.ExpectQuery("SELECT discounted_price, start_date, end_date FROM bazaar.discount").
//    WithArgs(productID).
//    WillReturnRows(rows)
//
//  repo := order2.NewOrderRepository(db)
//  discounts, err := repo.ProductDiscounts(context.Background(), productID)
//
//  assert.Nil(t, discounts)
//  assert.Error(t, err)
//  assert.Contains(t, err.Error(), "scan discount row")
//  assert.NoError(t, mock.ExpectationsWereMet())
//}
