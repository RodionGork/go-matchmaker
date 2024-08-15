users = {
    {name='Zaec', skill=10, latency=5, group=-1},
    {name='Volk', skill=12, latency=1, group=-1},
    {name='Ryba', skill=11, latency=5, group=-1},
    {name='Lisa', skill=15, latency=0.5, group=-1},
    {name='Muha', skill=8, latency=2, group=-1},
}

group_size = 2

fname = arg[1]
print('including:', fname)

dofile(fname)

for _, user in pairs(users) do
    print(user['name'], 'goes to group', user['group'])
end
