package ports

type Clock interface{ NowUnix() int64 }
