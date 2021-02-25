package native

import (
	"github.com/dfuse-io/logging"
	"go.uber.org/zap"
)

var zlog = zap.NewNop()

func setupLogger() {
	logging.Register("github.com/dfuse-io/eth-go/signer/native", &zlog)
}
