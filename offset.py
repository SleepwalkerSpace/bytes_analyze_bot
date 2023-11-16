addrs = [824635448264,
          824684972776,
            824698819496,
              824709355896,
                824782945160, 824791658568, 824807010136]

poker_10007 = 3193
offset_0 = 484
offset_1 = 488
offset_2 = 344
def poker_offset(pid, addrs, offset_0, offset_1, offset_2):
    for i, addr in enumerate(addrs):
        print("{} - {}".format(i, addr))
        print("\t/debug2 {} {} 1".format(pid, addr - offset_0))
        print("\t/debug2 {} {} 3".format(pid, addr - offset_1 + offset_2))

poker_offset(poker_10007, addrs, offset_0, offset_1, offset_2)