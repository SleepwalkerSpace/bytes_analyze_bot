balance = 60000
single_bet_amount = 100
bet_area_oddr = 2.2
empty = 1
total_bet_amount = 0
for i in range(1, 11):
    total_bet_amount += single_bet_amount
    balance -= single_bet_amount
    print("[{}]余额:{}\t下注:{}\t总注:{}\t预计收益:{}\n".format(
        i,
        balance, 
        single_bet_amount,
        total_bet_amount,
        int(single_bet_amount * bet_area_oddr - total_bet_amount)
        ))
    
    if i > empty:
        single_bet_amount *= 2