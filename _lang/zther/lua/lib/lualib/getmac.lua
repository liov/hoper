
function getmac()

	local ip="192.168.99.166"
	local name ="张三"

	for line in io.lines("E:\\a.txt") do
		if string.match(line,ip)~=nil then --如果存在字符串‘error’
			local state = string.match(line,"%w+",27)
			local mac = string.match(line,"[%w+:]+",33)
			print(state)
			print(mac)

			local flag =true
			local file = io.open ("E:\\b.txt","a+")

				for line1 in file:lines() do --检查mac是否已经存在
					if string.match(line1,mac) then
						flag = false
					end
				end

				if flag then
					io.output(file)
					io.write(mac.." "..name.."\n")
				end
				io.close(file)
		end
	end
end

function qiandao()
	for line in io.lines("E:\\b.txt") do
			local mac = string.match(line,"[%w:]+")
			local name =string.match(line,".+",19)
		if mac~=nil then
			for line1 in io.lines("E:\\a.txt") do
				local state = string.match(line1,"%w+",27)
				if (string.match(line1,mac) and state=="0x2") then
					print("<p>"..name.."  已到aa</p><br>")
				elseif string.match(line1,mac) then
					print("<p>"..name.."  未到bb</p><br>")
				end
			end

		else
			print("无人注册cc")
		end
	end
end

function split(szFullString, szSeparator)
    local nFindStartIndex = 1
    local nSplitIndex = 1
    local nSplitArray = {}
    while true do
       local nFindLastIndex = string.find(szFullString, szSeparator, nFindStartIndex)
       if not nFindLastIndex then
        nSplitArray[nSplitIndex] = string.sub(szFullString, nFindStartIndex, string.len(szFullString))
        break
       end
       nSplitArray[nSplitIndex] = string.sub(szFullString, nFindStartIndex, nFindLastIndex - 1)
       nFindStartIndex = nFindLastIndex + string.len(szSeparator)
       nSplitIndex = nSplitIndex + 1
    end
    return nSplitArray
    end
function test1()
for line in io.lines("E:\\test1.tpl") do
    if line~=nil then
        local list = split(line,"|")
        if list[4]=="t" then
            print(list[2]..":<input type='text' name='"..list[3].."' required='required'><br>")
        end
        if list[4]=="r" then
            local rlist=split(list[5],",")
            for k,v in ipairs(rlist) do
                print(list[2]..":<input type='radio' name='"..list[3].."' value='"..v.."'>"..v.."<br>")
            end
        end
        if list[4]=="c" then
            local clist=split(list[5],",")
            for k,v in ipairs(clist) do
                print(list[2]..":<input type='checkbox' name='"..list[3].."' value='"..v.."'>"..v.."<br>")
            end
        end
        if list[4]=="s" then
            local clist=split(list[5],",")
            print(list[2]..":<select name='"..list[3].."'><br>")
            for k,v in ipairs(clist) do
                print("<option value='"..v.."'>"..v.."</option>")
            end
            print("</select><br>")
        end
    end
end
end

local list =string.gsub("%2Ffile%2F","%%2F","/")
math.randomseed(os.time())

print(math.random(1000))

