package middlewares

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func validateInt(str string) error {
	if str == "" {
		return fmt.Errorf("field is empty")
	}
	_, err := strconv.ParseInt(str, 10, 64)
	return err
}

func validateBool(str string) error {
	if str == "" {
		return fmt.Errorf("field is empty")
	}
	_, err := strconv.ParseBool(str)
	return err
}

func validateUserID(c echo.Context) error {
	userID := c.Param("user_id")
	if userID == "" {
		return fmt.Errorf("user id is empty")
	}
	_, err := strconv.ParseInt(userID, 10, 64)
	return err
}

func validateOrderID(c echo.Context) error {
	orderID := c.Param("order_id")
	if orderID == "" {
		return fmt.Errorf("order id is empty")
	}
	_, err := uuid.Parse(orderID)
	return err
}

func validateParams(c echo.Context, expect map[string]string) error {
	for k, v := range expect {
		if strings.ToLower(c.QueryParam(k)) != v {
			return fmt.Errorf("query param %q should be %q", k, v)
		}
	}
	return nil
}
