package main

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/attribute"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func fibonacciHandler(ctx *gin.Context) {
	traceContext, span := tp.Tracer(serviceName).Start(ctx.Request.Context(), "fibonacci")
	defer span.End()
	number, err := strconv.Atoi(ctx.Query("number"))
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	fib := calculateFibonacci(number, traceContext)
	ctx.String(http.StatusOK, fmt.Sprintf("%d", fib))
}

func calculateFibonacci(n int, ctx context.Context) int {
	traceContext, span := tp.Tracer(serviceName).Start(ctx, "fibonacci-"+strconv.Itoa(n))
	span.SetAttributes(attribute.String("number", strconv.Itoa(n)))
	defer span.End()
	if n <= 1 {
		return n
	}
	return calculateFibonacci(n-1, traceContext) + calculateFibonacci(n-2, traceContext)
}
