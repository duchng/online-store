package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	http2 "net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/suite"
	"github.com/uptrace/bun"

	"store-management/internal/adapter/product/http"
	"store-management/internal/adapter/server"
	"store-management/internal/assets"
	"store-management/internal/config"
	"store-management/internal/dependencies"
	"store-management/pkg/dbtest"
)

type ProductTestSuite struct {
	suite.Suite
	ctx      context.Context
	injector do.Injector
	db       *bun.DB

	productHandler *http.ProductHandler
}

func TestProducts(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ProductTestSuite))
}

func (suite *ProductTestSuite) SetupSuite() {
	db := dbtest.NewDatabase(suite.T(), assets.EmbeddedFiles)
	suite.db = db
	injector := NewInjector[config.AppConfig](
		dependencies.NewInjector(), WithDb(suite.db),
	)
	suite.injector = injector
	suite.ctx = context.Background()
	suite.productHandler = do.MustInvoke[*http.ProductHandler](injector)
}

func (suite *ProductTestSuite) TearDownSuite() {
	suite.NoError(suite.db.Close())
}

func (suite *ProductTestSuite) TestProducts() {
	createProductRequest := http.CreateProductRequest{
		Name:          "Apple",
		Description:   "Some description",
		Price:         5.99,
		StockQuantity: 100,
	}
	productJson, _ := json.Marshal(createProductRequest)
	e := echo.New()
	e.Validator = &server.CustomValidator{Validator: validator.New()}
	req := httptest.NewRequest(http2.MethodPost, "/", bytes.NewReader(productJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	suite.NoError(suite.productHandler.CreateProduct(ctx))
	suite.Equal(http2.StatusCreated, rec.Code)
	respBody, err := io.ReadAll(rec.Result().Body)
	suite.NoError(err)
	suite.NotEmpty(respBody)

	q := make(url.Values)
	q.Add("name", "App")
	q.Add("cursor", "0")
	q.Add("limit", "10")
	req = httptest.NewRequest(http2.MethodGet, "/?"+q.Encode(), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	ctx = e.NewContext(req, rec)
	suite.NoError(suite.productHandler.ListProducts(ctx))
	suite.Equal(http2.StatusOK, rec.Code)
	respBody, err = io.ReadAll(rec.Result().Body)
	suite.NoError(err)
	suite.NotEmpty(respBody)
}
