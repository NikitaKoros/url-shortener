package server

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strings"
	"url-shortener/internal/model"
	"url-shortener/internal/shorten"
	"url-shortener/internal/config"
)

type shortener interface {
	Shorten(context.Context, model.ShortenInput) (*model.Shortening, error)
}

type shortenRequest struct {
	URL        string `json:"url"`
	Identifier string `json:"identifier,omitempty"`
}

type shortenResponse struct {
	ShortURL string `json:"short_url,omitempty"`
	Message  string `json:"message,omitempty"`
}

func HandleShorten(shortener shortener) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req shortenRequest
		if err := c.Bind(&req); err != nil {
			return err
		}

		if err := c.Validate(req); err != nil {
			return err
		}

		identifier := req.Identifier
		if strings.TrimSpace(req.Identifier) == "" {
			identifier = ""
		}

		input := model.ShortenInput{
			RawURL:     req.URL,
			Identifier: identifier,
		}

		shortening, err := shortener.Shorten(c.Request().Context(), input)
		if err != nil {
			if errors.Is(err, model.ErrIdentifierExists) {
				return echo.NewHTTPError(http.StatusConflict, err.Error())
			}

			log.Printf("shortener.Shorten: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		shortURL, err := shorten.PrependBaseURL(config.LoadConfig().BaseURL, shortening.Identifier)
		if err != nil {
			log.Printf("error generating full url for %q: %v", shortening.Identifier, err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		return c.JSON(
			http.StatusOK,
			shortenResponse{ShortURL: shortURL},
		)
	}

}
