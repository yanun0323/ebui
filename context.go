package ebui

import "context"

type resourcePathKeyInstance struct{}

var (
	resourcePathKey = resourcePathKeyInstance{}
)

func setResourcePath(ctx context.Context, path string) context.Context {
	return context.WithValue(ctx, resourcePathKey, path)
}

func getResourcePath(ctx context.Context) string {
	res, _ := ctx.Value(resourcePathKey).(string)
	return res
}
