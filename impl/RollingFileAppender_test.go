package impl_test

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-jang/go/util/optional"
	"github.com/go-log4g/core/impl"
	"github.com/go-log4g/core/impl/rolling"
	"github.com/stretchr/testify/require"
)

func TestRollingFileAppender(test *testing.T) {
	dir := os.TempDir()
	file := filepath.Join(dir, "go-log4g-TestRollingFileAppender.log")
	filePattern := filepath.Join(dir, "go-log4g-TestRollingFileAppender-%i.log")

	_ = os.Remove(file)
	_ = os.Remove(filepath.Join(dir, "go-log4g-TestRollingFileAppender-1.log"))
	_ = os.Remove(filepath.Join(dir, "go-log4g-TestRollingFileAppender-2.log"))
	_ = os.Remove(filepath.Join(dir, "go-log4g-TestRollingFileAppender-3.log"))

	statusLogger := impl.NewStatusLogger()
	layout := impl.NewPatternLayout("%m%n", impl.NewPatternParser(statusLogger))
	policy := rolling.NewSizeBasedTriggeringPolicy(12)
	strategy := rolling.NewDefaultRolloverStrategy(3, nil)
	appender := impl.NewRollingFileAppender(file, filePattern, false, 1024, true, layout, nil, statusLogger, nil, policy, strategy)

	callerContext := &impl.CallerContext{}

	appender.Append(impl.NewLogEvent(context.Background(), slog.NewRecord(time.Now(), slog.LevelInfo, "11111", 0), callerContext, nil, nil))
	appender.Append(impl.NewLogEvent(context.Background(), slog.NewRecord(time.Now(), slog.LevelInfo, "22222", 0), callerContext, nil, nil))
	appender.Append(impl.NewLogEvent(context.Background(), slog.NewRecord(time.Now(), slog.LevelInfo, "33333", 0), callerContext, nil, nil))
	appender.Append(impl.NewLogEvent(context.Background(), slog.NewRecord(time.Now(), slog.LevelInfo, "44444", 0), callerContext, nil, nil))

	require.Equal(test, "11111\n22222\n", string(optional.OfCommaErr(os.ReadFile(filepath.Join(dir, "go-log4g-TestRollingFileAppender-1.log"))).OrElsePanic("cannot read file")))
	require.Equal(test, "33333\n44444\n", string(optional.OfCommaErr(os.ReadFile(file)).OrElsePanic("cannot read file")))
}
