local _weatherManager = {}

_weatherManager.loadWeather = function (callback)
    lua_http.request({ url  = "https://api.seniverse.com/v3/weather/daily.json?key=SNVXTU-TmTj7-AEm_&location=beijing&language=zh-Hans&unit=c&start=0&days=5",
                       onResponse = function (response)
                           if response.http_code ~= 200 then
                               if callback then
                                   callback(nil)
                               end
                           else
                               lua_thread.postToThread(BusinessThreadLOGIC,"WeatherManager","parseWeathers",response.response,function(data)
                                   if callback then
                                       callback(data)
                                   end
                               end)
                           end
                       end})
end

return _weatherManager