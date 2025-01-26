package handler

import "github.com/labstack/echo/v4"

func bindAndValidate(c echo.Context, i interface{}) error {
	if err := c.Bind(i); err != nil {
		return err
	}

	if err := c.Validate(i); err != nil {
		return err
	}

	return nil
}
