package aurora

// GrpcServer grpc 协议转发中间件
//func GrpcServer(server *grpc.Server) Middleware {
//	return func(ctx Ctx) bool {
//		if ctx.Request().ProtoMajor == 2 && strings.HasPrefix(ctx.Request().Header.Get("Content-Type"), "application/grpc") {
//			server.ServeHTTP(ctx.Response(), ctx.Request())
//			return false
//		}
//		return true
//	}
//}
