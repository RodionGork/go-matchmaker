assert(users, "users table undefined") -- with entries {skill, latency, waiting}
assert(group_size, "group_size undefined")

group_count = math.floor(#users / group_size)

for i = 1, group_count do
    for j = 1, group_size do
        table.insert(users[(i-1) * group_size + j], i)
    end
end

for i = group_count*group_size+1,#users do
    table.insert(users[i], -1)
end

