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
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
	"time"
)

const queryUpdateOrderStatus = `
	UPDATE bazaar.order
	SET 
		status = $1,
		updated_at = now()
	WHERE id = $2`

const queryGetOrders = `
		SELECT id, status, total_price, total_price_discount, 
			address_id, expected_delivery_at, actual_delivery_at, created_at 
		FROM bazaar.order WHERE status = 'placed'`

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

	mock.ExpectQuery(`SELECT a.region, a.city, a.address_string, a.coordinate, ua.label
	FROM bazaar.address a
	LEFT JOIN bazaar.user_address ua ON a.id = ua.address_id
	WHERE a.id = \$1
	LIMIT 1`).
		WithArgs(addressID).
		WillReturnRows(sqlmock.NewRows([]string{
			"region", "city", "address_string", "coordinate", "label",
		}).AddRow(
			expectedAddress.Region,
			expectedAddress.City,
			expectedAddress.AddressString,
			expectedAddress.Coordinate,
			expectedAddress.Label, // <-- добавь поле Label в expectedAddress
		))

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

	mock.ExpectQuery(`SELECT a\.region, a\.city, a\.address_string, a\.coordinate, ua\.label
	FROM bazaar\.address a
	LEFT JOIN bazaar\.user_address ua ON a\.id = ua\.address_id
	WHERE a\.id = \$1
	LIMIT 1`).
		WithArgs(addressID).
		WillReturnError(sql.ErrNoRows) // для TestGetOrderAddress_NotFound

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

	mock.ExpectQuery(`SELECT a\.region, a\.city, a\.address_string, a\.coordinate, ua\.label
	FROM bazaar\.address a
	LEFT JOIN bazaar\.user_address ua ON a\.id = ua\.address_id
	WHERE a\.id = \$1
	LIMIT 1`).
		WithArgs(addressID).
		WillReturnError(errors.New("database error")) // для TestGetOrderAddress_QueryError

	repo := order2.NewOrderRepository(db)
	address, err := repo.GetOrderAddress(context.Background(), addressID)

	assert.Nil(t, address)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetOrdersPlaced_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := order2.NewOrderRepository(db)

	mock.ExpectQuery("SELECT (.+) FROM bazaar.order WHERE status = 'placed'").
		WillReturnError(errors.New("query error"))

	orders, err := repo.GetOrdersPlaced(context.Background())
	assert.Nil(t, orders)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "query error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserIDByOrderID_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := order2.NewOrderRepository(db)
	orderID := uuid.New()

	// Update this line to match the actual query format
	mock.ExpectQuery("SELECT user_id FROM bazaar.order WHERE id = \\$1").
		WithArgs(orderID).
		WillReturnError(errors.New("scan failed"))

	uid, err := repo.GetUserIDByOrderID(context.Background(), orderID)
	assert.Equal(t, uuid.Nil, uid)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "scan failed")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetOrderProducts_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	orderID := uuid.New()
	expectedProducts := []dto.GetOrderProductResDTO{
		{
			ProductID:   uuid.New(),
			Quantity:    2,
			ProductName: "Test Product 1",
		},
		{
			ProductID:   uuid.New(),
			Quantity:    3,
			ProductName: "Test Product 2",
		},
	}

	rows := sqlmock.NewRows([]string{"product_id", "quantity", "product_name"})
	for _, p := range expectedProducts {
		rows.AddRow(p.ProductID, p.Quantity, p.ProductName)
	}

	mock.ExpectQuery(`SELECT oi.product_id, oi.quantity, p.name FROM bazaar.order_item oi JOIN bazaar.product p ON oi.product_id = p.id WHERE oi.order_id = \$1`).
		WithArgs(orderID).
		WillReturnRows(rows)

	repo := order2.NewOrderRepository(db)
	products, err := repo.GetOrderProducts(context.Background(), orderID)

	assert.NoError(t, err)
	assert.Equal(t, len(expectedProducts), len(*products))
	assert.Equal(t, expectedProducts, *products)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetOrderProducts_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	orderID := uuid.New()

	mock.ExpectQuery(`SELECT oi.product_id, oi.quantity, p.name FROM bazaar.order_item oi JOIN bazaar.product p ON oi.product_id = p.id WHERE oi.order_id = \$1`).
		WithArgs(orderID).
		WillReturnRows(sqlmock.NewRows([]string{"product_id", "quantity", "name"}))

	repo := order2.NewOrderRepository(db)
	products, err := repo.GetOrderProducts(context.Background(), orderID)

	assert.Nil(t, products)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateStatus_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	orderID := uuid.New()
	status := models.Placed

	mock.ExpectExec(regexp.QuoteMeta(queryUpdateOrderStatus)).
		WithArgs(status.String(), orderID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := order2.NewOrderRepository(db)
	err = repo.UpdateStatus(context.Background(), orderID, status)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateStatus_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	orderID := uuid.New()
	status := models.Placed

	mock.ExpectExec(regexp.QuoteMeta(queryUpdateOrderStatus)).
		WithArgs(status.String(), orderID).
		WillReturnError(errors.New("database error"))

	repo := order2.NewOrderRepository(db)
	err = repo.UpdateStatus(context.Background(), orderID, status)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetOrdersPlaced_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := order2.NewOrderRepository(db)

	rows := sqlmock.NewRows([]string{
		"id", "status", "total_price", "total_price_discount",
		"address_id", "expected_delivery_at", "actual_delivery_at", "created_at",
	}).AddRow(
		uuid.New(), "placed", 1000, 900,
		uuid.New(), time.Now(), time.Now(), time.Now(),
	)

	mock.ExpectQuery(regexp.QuoteMeta(queryGetOrders)).
		WillReturnRows(rows)

	ctx := context.Background()
	orders, err := repo.GetOrdersPlaced(ctx)

	require.NoError(t, err)
	require.Len(t, *orders, 1)
	assert.Equal(t, models.Placed, (*orders)[0].Status)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetOrdersPlaced_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := order2.NewOrderRepository(db)

	rows := sqlmock.NewRows([]string{
		"id", "status", "total_price", "total_price_discount",
		"address_id", "expected_delivery_at", "actual_delivery_at", "created_at",
	}).AddRow(
		"not-a-uuid", "placed", 1000, 900, // <--- здесь ошибка
		uuid.New(), time.Now(), time.Now(), time.Now(),
	)

	mock.ExpectQuery(regexp.QuoteMeta(queryGetOrders)).
		WillReturnRows(rows)

	ctx := context.Background()
	orders, err := repo.GetOrdersPlaced(ctx)

	require.Error(t, err)
	assert.Nil(t, orders)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetOrdersPlaced_InvalidStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := order2.NewOrderRepository(db)

	rows := sqlmock.NewRows([]string{
		"id", "status", "total_price", "total_price_discount",
		"address_id", "expected_delivery_at", "actual_delivery_at", "created_at",
	}).AddRow(
		uuid.New(), "invalid_status", 1000, 900,
		uuid.New(), time.Now(), time.Now(), time.Now(),
	)

	mock.ExpectQuery(regexp.QuoteMeta(queryGetOrders)).
		WillReturnRows(rows)

	ctx := context.Background()
	orders, err := repo.GetOrdersPlaced(ctx)

	require.Error(t, err)
	assert.Nil(t, orders)
	assert.NoError(t, mock.ExpectationsWereMet())
}
