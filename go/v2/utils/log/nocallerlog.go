package log

var  Default = (&LoggerInfo{Product: true}).initConfig(false).Sugar()