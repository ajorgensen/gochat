package static

import (
	"embed"
	"fmt"
)

//go:embed assets
var AssetsFS embed.FS

func Asset(name string) string {
	return fmt.Sprintf("/public/assets/%s", name)
}
