-- naive algorithm, to work as placeholder when nothing better is loaded
-- it builds as many groups as possible
-- by simply marking sequential "group_size" chunks from beginning

assert(users, "global variable 'users' should be defined")
assert(group_size, "global variable 'group_size' should be defined")

group_count = math.floor(#users / group_size)

for i = 1, group_count do
    for j = 1, group_size do
        cur_user = users[(i-1) * group_size + j]
        cur_user['group'] = i
    end
end

