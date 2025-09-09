local url = ...

-- Set Computer Label
local label
if turtle then
    label = "T"
elseif pocket then
    label = "P"
else
    label = "C"
end
label = label .. os.computerID()
os.setComputerLabel(label)

-- Export peripheral names
local file = fs.open("peripherals.json","w")
file.write(textutils.serializeJSON(peripheral.getNames()))
file.close()

-- Get terminals
local terms = { term.native(), peripheral.find('monitor') }
for _, term in ipairs(terms) do
    if term.setTextScale then term.setTextScale(0.5) end
    term.height = select(2, term.getSize())
end

-- log function
local function log(packet)
    for _, term in ipairs(terms) do
        term.scroll(1)
        term.setCursorPos(1, term.height)
        term.setTextColor(bit.blshift(1, packet.color or packet.c or 14))
        term.write(packet.text or packet.t)
    end
end

local function exec(task, result)
    --result.result = 'test'
    if task.action == "log" then
        log(task.args)
    elseif task.action == "peripheral" then
        result.result = {peripheral.call(table.unpack(task.args))}
    elseif task.action == "turtle" then
        local f = table.remove(task.args, l)
        result.result = {turtle[f](table.unpack(task.args))}
    else
        error('invalid action: ' .. tostring(task.action))
    end
    return 0
end

--print("hello", url, clientName)

while true do
    local tid = os.startTimer(3)
    local url = url .. '#' .. tid
    log({ text = 'Connecting to ' .. label .. '@' .. url, color = 4 })
    http.websocketAsync(url)

    local socket
    local event = { os.pullEvent() }
    if event[1] == 'timer' then
        if event[2] == tid then
            log({ t = 'Timed out', c = 14 })
        end
    elseif event[1] == 'websocket_failure' then
        if event[2] == url then
            log({ text = 'Websocket failure', color = 14 })
            os.sleep(5)
        end
    elseif event[1] == 'websocket_success' then
        if event[2] == url then
            socket = event[3]
        end
    end

    if socket then
        log({ text = 'Websocket connected', color = 13 })

        local out = {{addr = label,}}
        local tasks = {}

        local handler = function(data)
            local task = coroutine.create(exec)
            log({ text = data, color = 13 })
            data = textutils.unserialiseJSON(data)
            local result = {
                id = data.id,
                result = {},
                error = nil,
            }
            local success, filter = coroutine.resume(task, data, result)
            if not success then
                result.result = nil
                result.error = filter
                --log({ text = filter, color = 14 })
                table.insert(out, result)
            elseif type(filter) == 'number' then
                table.insert(out, result)
            else
                table.insert(tasks, {
                    task = task,
                    filter = filter,
                    result = result,
                })
            end
            --todo call exec
        end

        while true do
            local success = true
            while #out > 0 do
                local result
                local to_send_obj = table.remove(out, 1)
                if not to_send_obj.result then
                    to_send_obj.result = {}
                end
                local to_send = textutils.serialiseJSON(to_send_obj)
                success, result = pcall(socket.send, to_send, true)
                if not success then
                    log({ text = result, color = 14 })
                    break
                end
            end
            if not success then break end

            local event = { os.pullEvent() }
            if event[1] == 'timer' then
                if event[2] == tid then
                    tid = nil
                end
            elseif event[1] == 'websocket_closed' then
                if event[2] == url then
                    log({ text = "Websocket closed", color = 14 })
                    break
                end
            elseif event[1] == 'websocket_message' then
                if event[2] == url then
                    handler(event[3])
                end
            end

            local newTasks = {}
            for _, v in ipairs(tasks) do
                if not v.filter or v.filter == event[1] then
                    local success, ret = coroutine.resume(v.task, table.unpack(event))
                    if not success then
                        v.result.result = nil
                        v.result.error = ret
                        log({ text = ret, color = 14 })
                        table.insert(out, v.result)
                    elseif type(ret) == 'number' then
                        table.insert(out, v.result)
                    else
                        table.insert(newTasks, {
                            task = v.task,
                            filter = ret,
                            result = v.result,
                        })
                    end
                else
                    table.insert(newTasks, v)
                end
            end
            tasks = newTasks
        end
        if tid then repeat local e = { os.pullEvent() } until e[1] ~= 'timer' or e[2] ~= tid end
    end
end

--while true do
--os.sleep(5)
--end
