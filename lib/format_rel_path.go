package lib

import "os"

func AddRootToPath(path string) string {
	return os.Getenv("root_dir") + path
}
