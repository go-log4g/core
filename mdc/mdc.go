package mdc

import "context"

type contextKey struct{}

func Put(ctx context.Context, key, value string) context.Context {
	return PutAll(ctx, map[string]string{key: value})
}

func PutAll(ctx context.Context, values map[string]string) context.Context {
	result := valuesCopy(ctx)

	for key, value := range values {
		result[key] = value
	}

	return context.WithValue(ctx, contextKey{}, result)
}

func Get(ctx context.Context, key string) (string, bool) {
	value, ok := values(ctx)[key]
	return value, ok
}

func Values(ctx context.Context) map[string]string {
	return valuesCopy(ctx)
}

func Remove(ctx context.Context, key string) context.Context {
	result := valuesCopy(ctx)
	delete(result, key)
	return context.WithValue(ctx, contextKey{}, result)
}

func Clear(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKey{}, map[string]string{})
}

func values(ctx context.Context) map[string]string {
	if ctx == nil {
		return nil
	}

	values, _ := ctx.Value(contextKey{}).(map[string]string)
	return values
}

func valuesCopy(ctx context.Context) map[string]string {
	source := values(ctx)
	result := make(map[string]string, len(source))

	for key, value := range source {
		result[key] = value
	}

	return result
}
