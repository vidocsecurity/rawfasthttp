//go:build !amd64 && !arm64 && !ppc64 && !ppc64le
// +build !amd64,!arm64,!ppc64,!ppc64le

package rawfasthttp

const (
	maxHexIntChars = 7
)