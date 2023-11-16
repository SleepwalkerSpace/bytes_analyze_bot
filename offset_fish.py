addrs = []


fish_pid= 10376

offset_0 = 592


for i, addr in enumerate(addrs):
    print("{} - {}".format(i, addr))
    print("\t/debug2 {} {} 3".format(fish_pid, addr - offset_0))
