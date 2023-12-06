package main

import "github.com/labstack/echo/v4"

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.JSON(200, echo.Map{
			"message": "Hello World",
		})
	})
	e.Logger.Fatal(e.Start(":8080"))
}
