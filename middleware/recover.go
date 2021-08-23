package middleware

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/klintcheng/kim"
	"github.com/klintcheng/kim/logger"
	"github.com/klintcheng/kim/wire/pkt"
)

func Recover() kim.HandlerFunc {
	return func(ctx kim.Context) {
		defer func() {
			if err := recover(); err != nil {
				var callers []string
				for i := 1; ; i++ {
					_, file, line, got := runtime.Caller(i)
					if !got {
						break
					}
					callers = append(callers, fmt.Sprintf("%s:%d", file, line))
				}
				logger.WithFields(logger.Fields{
					"ChannelId": ctx.Header().ChannelId,
					"Command":   ctx.Header().Command,
					"Seq":       ctx.Header().Sequence,
				}).Error(err, strings.Join(callers, "\n"))

				_ = ctx.Resp(pkt.Status_SystemException, &pkt.ErrorResp{Message: "SystemException"})
			}
		}()

		ctx.Next()
	}

}
