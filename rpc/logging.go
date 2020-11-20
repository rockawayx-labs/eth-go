package rpc

import (
	"github.com/dfuse-io/logging"
	"go.uber.org/zap"
)

var traceEnabled = logging.IsTraceEnabled("eth-go", "github.com/dfuse-io/eth-go/rpc")
var zlog = zap.NewNop()

func init() {
	logging.Register("github.com/dfuse-io/eth-go/rpc", &zlog)
}
