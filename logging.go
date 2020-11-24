package eth

import (
	"github.com/dfuse-io/logging"
	"go.uber.org/zap"
)

var traceEnabled = logging.IsTraceEnabled("eth-go", "github.com/dfuse-io/eth-go")
var zlog *zap.Logger

func init() {
	logging.Register("github.com/dfuse-io/eth-go", &zlog)
}
