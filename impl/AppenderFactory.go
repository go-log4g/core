package impl

type AppenderFactory func(name string, pending *PendingAppender) Appender
