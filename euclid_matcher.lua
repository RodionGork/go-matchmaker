assert(users, "global variable 'users' should be defined")
assert(group_size, "global variable 'group_size' should be defined")

score = {}

-- first calculate score for each user based on skill and latency (time not used)
for i, user in ipairs(users) do
    s = math.sqrt(math.pow(user['skill'], 2) / 100 + 3 / (user['latency'] + 1))
    table.insert(score, {i, s})
end

table.sort(score, function(a, b) return a[2] < b[2] end)

group_count = math.floor(#users / group_size)

for g = 1, group_count do
    for j = 1, group_size do
        i, s = table.unpack(table.remove(score))
        users[i]['group'] = g
    end
end
