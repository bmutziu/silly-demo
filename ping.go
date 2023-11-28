package main

import (
	"errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

func pingHandler(ctx *gin.Context) {
	req := resty.New().R().SetHeaderMultiValues(ctx.Request.Header).SetHeader("Content-Type", "application/text")

	otelCtx := ctx.Request.Context()
	span := trace.SpanFromContext(otelCtx)
	defer span.End()
	otel.GetTextMapPropagator().Inject(otelCtx, propagation.HeaderCarrier(req.Header))

	url := ctx.Query("url")
	if len(url) == 0 {
		url = os.Getenv("PING_URL")
		if len(url) == 0 {
			httpErrorBadRequest(errors.New("url is empty"), span, ctx)
			return
		}
	}
	log.Printf("Sending a ping to %s", url)
	resp, err := req.Get(url)
	if err != nil {
		httpErrorBadRequest(err, span, ctx)
		return
	}
	log.Println(resp.String())
	ctx.String(http.StatusOK, resp.String())
}
