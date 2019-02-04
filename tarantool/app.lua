box.cfg{
    listen = 3311
}

box.once("bootstrap", function()
    l = box.schema.space.create('locations', {engine='vinyl'})
    l:format({{name='id', type='unsigned'}, {name='location_data', type='string'}})
    l:create_index('primary', {type = 'tree', parts = {1, 'unsigned'}})

    il = box.schema.space.create('item_locations', {engine='vinyl'})
    il:format({{name='item_id', type='unsigned'}, {name='location_ids', type='array'}})
    il:create_index('primary', {type = 'tree', parts = {1, 'unsigned'}})
end)

-------------------------------------
-- @param item_id (number)
-------------------------------------
function get_item_locations(item_id)
    assert(type(item_id) == 'number', 'item_id must be a number')
    local locations = {}
    local il = box.space.item_locations
    local l = box.space.locations
    local location_ids = il:get{item_id}
    if location_ids then
        local ids = location_ids['location_ids']
        for _, id in pairs(ids) do
            local location = l:get{id}
            if location then
                table.insert(locations, location['location_data'])
            else
                -- XXX: inconsistency, no location
            end
        end
    end
    return locations
end
-------------------------------------
-- @param item_id (number)
-- @param location_ids_array (table)
-------------------------------------
function put_item_locations(item_id, location_ids_array)
    assert(type(item_id) == 'number', 'item_id must be a number')
    assert(type(location_ids_array) == 'table', 'location_ids_array must be a array')
    local il = box.space.item_locations
    if #location_ids_array == 0 then
        il:delete(item_id)
    else
        il:upsert({item_id, location_ids_array}, {{'=', 2, location_ids_array}})
    end
end

-------------------------------------
-- @param location_id (number)
-- @param location_data (string)
-------------------------------------
function add_location(location_id, location_data)
    assert(type(location_id) == 'number', 'location_id must be a number')
    assert(type(location_data) == 'string', 'location_data must be a string')
    local l = box.space.locations
    l:insert({location_id, location_data})
end

-------------------------------------
-- @param locations (table) [ [location_id:number, location_data:string] ]
-------------------------------------
function add_locations(locations)
    assert(type(locations) == 'table', 'locations must be an array')
    for _, loc in pairs(locations) do
        assert(type(loc[1]) == 'number', "locaion id must be a number")
        assert(type(loc[2]) == 'string', "location_data must be a string")
    end
    for _, loc in pairs(locations) do
        add_location(loc[1], loc[2])
    end
end
