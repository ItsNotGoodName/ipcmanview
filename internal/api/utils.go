package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/labstack/echo/v4"
)

func useDahuaConn(c echo.Context, db repo.DB, store *dahua.Store) (dahua.Client, error) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return dahua.Client{}, echo.ErrBadRequest.WithInternal(err)
	}

	ctx := c.Request().Context()

	dbDevice, err := db.GetDahuaDevice(ctx, id)
	if err != nil {
		if repo.IsNotFound(err) {
			return dahua.Client{}, echo.ErrNotFound.WithInternal(err)
		}
		return dahua.Client{}, err
	}

	client := store.Client(ctx, dbDevice.Convert().DahuaConn)

	return client, nil
}

// ---------- Stream

func useStream(c echo.Context) *json.Encoder {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response())
}

type StreamPayload struct {
	Data    any     `json:"data,omitempty"`
	Message *string `json:"message,omitempty"`
	OK      bool    `json:"ok"`
}

func sendStreamError(c echo.Context, enc *json.Encoder, err error) error {
	str := err.Error()
	if encodeErr := enc.Encode(StreamPayload{
		OK:      false,
		Message: &str,
	}); encodeErr != nil {
		return errors.Join(encodeErr, err)
	}

	c.Response().Flush()

	return err
}

func sendStream(c echo.Context, enc *json.Encoder, data any) error {
	err := enc.Encode(StreamPayload{
		OK:   true,
		Data: data,
	})
	if err != nil {
		return sendStreamError(c, enc, err)
	}

	c.Response().Flush()

	return nil
}

// ---------- Queries

func queryInts(c echo.Context, key string) ([]int64, error) {
	ids := make([]int64, 0)
	idsStr := c.QueryParams()[key]
	for _, v := range idsStr {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, echo.ErrBadRequest.WithInternal(err)
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func queryIntOptional(c echo.Context, key string) (int, error) {
	str := c.QueryParam(key)
	if str == "" {
		return 0, nil
	}

	number, err := strconv.Atoi(str)
	if err != nil {
		return 0, echo.ErrBadRequest.WithInternal(err)
	}

	return number, nil
}

func queryBoolOptional(c echo.Context, key string) (bool, error) {
	str := c.QueryParam(key)
	if str == "" {
		return false, nil
	}

	bool, err := strconv.ParseBool(str)
	if err != nil {
		return false, echo.ErrBadRequest.WithInternal(err)
	}

	return bool, nil
}
