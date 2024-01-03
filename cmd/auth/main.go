package auth

//
//var (
//	cfg config.App
//)
//
//func init() {
//	err := env.Parse(&cfg)
//	if err != nil {
//		panic(err)
//	}
//
//	level, err := log.ParseLevel(cfg.ErrorLevel)
//	if err != nil {
//		panic(err)
//	}
//	log.SetLevel(level)
//	log.SetFormatter(&log.JSONFormatter{})
//	log.SetReportCaller(true)
//}
//
//func main() {
//	application, err := internal.NewApplication(cfg)
//	if err != nil {
//		panic(err)
//	}
//
//	application.Run()
//}
