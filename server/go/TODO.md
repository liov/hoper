1.自动日志记录，避免runtime.Caller的开销又能记录详细信息，思考设计一下
2.cookie的设置，把这个抽象出来、，在业务层决定具体的返回头设置
2已经实现，目前的做法是gateway转grpc把cookie放头里，本地http直接调用，在返回中加cookie字段，
返回类型实现cookie方法，这么做还是别捏，待改进
