package limiter

import "time"

type RateLimiter struct {
	rate       int           // Number of actions allowed per time window
	interval   time.Duration // Time window duration
	tokenCount int           // Available tokens at the moment
	tokens     chan struct{} // Channel to hold tokens
}

func NewRateLimiter(rate int, interval time.Duration) *RateLimiter {
	if rate == 0 {
		rate = 1
	}

	if interval == 0 {
		interval = time.Second
	}

	limiter := &RateLimiter{
		rate:       rate,
		interval:   interval,
		tokenCount: rate,
		tokens:     make(chan struct{}, rate),
	}

	go limiter.refillTokens()

	return limiter
}

func (limiter *RateLimiter) refillTokens() {
	for {
		select {
		case <-time.After(limiter.interval / time.Duration(limiter.rate)):
			for i := 0; i < limiter.rate; i++ {
				select {
				case limiter.tokens <- struct{}{}:
				default:
					//just next the request, don't block it
				}
			}
		}
	}
}

func (limiter *RateLimiter) Allow() bool {
	select {
	case <-limiter.tokens:
		return true
	default:
		return false
	}
}
