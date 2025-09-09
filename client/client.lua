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
local file = fs.open("peripherals.json", "w")
file.write(textutils.serializeJSON(peripheral.getNames()))
file.close()


-- Get terminals
local terms = { term.native(), peripheral.find('monitor') }
for _, term in ipairs(terms) do
    if term.setTextScale then term.setTextScale(0.5) end
    term.height = select(2, term.getSize())
end

local function log(packet)
    if packet == nil then
        return
    end
    if type(packet) == 'string' or type(packet) == 'number' then
        packet = {t=packet, c=0}
    end
    for _, term in ipairs(terms) do
        term.scroll(1)
        term.setCursorPos(1, term.height)
        term.setTextColor(bit.blshift(1, packet.color or packet.c or 0))
        term.write(packet.text or packet.t)
    end
end

local function exec(task, result)
    if task.type == "log" then
        log(task.args[1])
    elseif task.type == "peripheral" then
        result.result = peripheral.call(table.unpack(task.args))
    elseif task.type == "turtle" then
        local f = table.remove(task.args, l)
        result.result = {turtle[f](table.unpack(task.args))}
    else
        error('invalid type: ' .. tostring(task.type))
    end
    return 0
end

--
while true do
    http.websocketAsync(url)

    local ws
    while true do
        local event = { os.pullEvent() }
        if event[1] == 'websocket_failure' then
            if event[2] == url then
                log({t="websocket_failure " .. event[3], c=14})
                os.sleep(1)
                break
            end
        elseif event[1] == 'websocket_success' then
            if event[2] == url then
                log({t="websocket_success", c=5})
                ws = event[3]
                break
            end
        end
    end

    if ws then
        local out = {}
        local tasks = {}
        table.insert(out, { client = label })

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
                log({ text = filter, color = 14 })
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
        end

        while true do
            local success = true
            while #out > 0 do
                local to_send = textutils.serialiseJSON(table.remove(out, 1))
                log(to_send)
                local result
                success, result = pcall(ws.send, to_send, true)
                if not success then
                    log(result[0])
                    break
                end
            end
            if not success then break end

            local event = { os.pullEvent() }
            if event[1] == 'websocket_message' then
                if event[2] == url then
                    log("websocket_message " .. event[3])
                    handler(event[3])
                end
            elseif event[1] == 'websocket_failure' then
                if event[2] == url then
                    log({t="websocket_failure", c=14})
                    os.sleep(1)
                    break
                end
            elseif event[1] == 'websocket_closed' then
                if event[2] == url then
                    log({t="websocket_closed", c=14})
                    os.sleep(1)
                    break
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
    end
end

log("aaaaaa")
