package binance

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/thrasher-corp/gocryptotrader/exchanges/request"
)

func TestRateLimit_Limit(t *testing.T) {
	symbol := "BTC-USDT"

	testTable := map[string]struct {
		Expected request.EndpointLimit
		Limit    request.EndpointLimit
		Deadline time.Time
	}{
		"All Orderbooks Ticker": {Expected: spotOrderbookTickerAllRate, Limit: bestPriceLimit("")},
		"Orderbook Ticker":      {Expected: spotDefaultRate, Limit: bestPriceLimit(symbol)},
		"Open Orders":           {Expected: spotOpenOrdersSpecificRate, Limit: openOrdersLimit(symbol)},
		"Orderbook Depth 5":     {Expected: spotDefaultRate, Limit: orderbookLimit(5)},
		"Orderbook Depth 10":    {Expected: spotDefaultRate, Limit: orderbookLimit(10)},
		"Orderbook Depth 20":    {Expected: spotDefaultRate, Limit: orderbookLimit(20)},
		"Orderbook Depth 50":    {Expected: spotDefaultRate, Limit: orderbookLimit(50)},
		"Orderbook Depth 100":   {Expected: spotDefaultRate, Limit: orderbookLimit(100)},
		"Orderbook Depth 500":   {Expected: spotOrderbookDepth500Rate, Limit: orderbookLimit(500)},
		"Orderbook Depth 1000":  {Expected: spotOrderbookDepth1000Rate, Limit: orderbookLimit(1000)},
		"Orderbook Depth 5000":  {Expected: spotOrderbookDepth5000Rate, Limit: orderbookLimit(5000)},
		"Exceeds deadline":      {Expected: spotOrderbookDepth5000Rate, Limit: orderbookLimit(5000), Deadline: time.Now().Add(time.Nanosecond)},
	}
	for name, tt := range testTable {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			exp, got := tt.Expected, tt.Limit
			if exp != got {
				t.Fatalf("incorrect limit applied.\nexp: %v\ngot: %v", exp, got)
			}

			ctx := context.Background()
			if !tt.Deadline.IsZero() {
				var cancel context.CancelFunc
				ctx, cancel = context.WithDeadline(ctx, tt.Deadline)
				defer cancel()
			}

			l := SetRateLimit()
			if err := l.Limit(ctx, tt.Limit); err != nil && !errors.Is(err, context.DeadlineExceeded) {
				t.Fatalf("error applying rate limit: %v", err)
			}
		})
	}
}

func TestRateLimit_LimitStatic(t *testing.T) {
	testTable := map[string]request.EndpointLimit{
		"Default":           spotDefaultRate,
		"Historical Trades": spotHistoricalTradesRate,
		"All Price Changes": spotPriceChangeAllRate,
		"All Orders":        spotAllOrdersRate,
	}
	for name, tt := range testTable {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			l := SetRateLimit()
			if err := l.Limit(context.Background(), tt); err != nil {
				t.Fatalf("error applying rate limit: %v", err)
			}
		})
	}
}
