package core

import "fmt"

type Middleware func(func(*WindowContext)) func(*WindowContext)

func Compose(mws ...Middleware) Middleware {
	return func(final func(*WindowContext)) func(*WindowContext) {
		for i := len(mws) - 1; i >= 0; i-- {
			final = mws[i](final)
		}
		return final
	}
}

func WithWindowGuard(next func(*WindowContext)) func(*WindowContext) {
	return func(ctx *WindowContext) {
		bounds := ctx.Win.Bounds()
		width := int(bounds.W())
		height := int(bounds.H())
		if width <= 0 || height <= 0 {
			return
		}
		next(ctx)
	}
}

func WithRecover(next func(*WindowContext)) func(*WindowContext) {
	return func(ctx *WindowContext) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("[panic recovered] %v\n", r)
			}
		}()
		next(ctx)
	}
}
