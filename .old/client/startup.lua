
local success, result = pcall(function()
    load(http.get('http://127.0.0.1:1847/client.lua').readAll)('ws://127.0.0.1:1847/ws')
end)

if not success then
    term.write(result)
end

os.sleep(60)
os.reboot()