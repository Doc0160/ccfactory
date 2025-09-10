
local success, result = pcall(function()
    load(http.get('http://127.0.0.1:1847/client.lua').readAll)()
end)

if not success then
    term.write(result)
end
term.write("Rebooting in 60s seconds ...")

os.sleep(60)
os.reboot()