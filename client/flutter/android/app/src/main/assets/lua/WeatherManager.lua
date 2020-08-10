local _weatherManager = {}

_weatherManager.parseWeathers = function (responseStr,callback)
	local t = cjson.decode(responseStr)
	local ret = {}
	if t and t.results and #t.results > 0 and t.results[1].daily then
		local city = t.results[1].location.name
		for i,v in ipairs(t.results[1].daily) do
			local t = {}
			t.wind =v.wind_speed
			t.wind_direction = v.wind_direction
			t.sun_info = v.text_day
			t.low = tonumber(v.low)
			t.high = tonumber(v.high)
			t.id = i
			t.city = city
			table.insert(ret,t)
		end
	end
	if callback then
		callback(ret)
	end
end

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