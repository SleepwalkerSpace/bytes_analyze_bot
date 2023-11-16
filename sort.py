
players = [
    {"id": 1, "hand_cards":[{"point": 10, "suit": 4}, {"point": 10, "suit": 3}, {"point": 10, "suit": 2}]},
    {"id": 2, "hand_cards":[{"point": 11, "suit": 4}, {"point": 11, "suit": 3}, {"point": 11, "suit": 2}]},
    {"id": 3, "hand_cards":[{"point": 12, "suit": 4}, {"point": 12, "suit": 3}, {"point": 12, "suit": 2}]},
]

"""
1.手牌类型

"""


def c_hand_cards_kind(hand_cards) -> int:
    pass

for i, player in enumerate(players):
    hand_cards_kind = c_hand_cards_kind(player["hand_cards"])
    players[i]["hand_cards_kind"] = hand_cards_kind


def pk_kind_6(i, j) -> bool:

    pass

def pk_kind_5_3(i, j) -> bool:
    j == qka
    return True
    i  == qka
    return False

    a23
    if "a" in i or "a" in j:
        ap[suit]
    
    
    pass


def pks(players):
    for i, pi in enumerate(players):
        for j, pj in enumerate(players):

            if i == j:
                continue

            if pi[i]["hand_cards_kind"] > pj[i]["hand_cards_kind"]:
                continue

            elif pi[i]["hand_cards_kind"] < pj[i]["hand_cards_kind"]:
                players[i], players[j] = players[j], players[i]
            
            else:
                flag = False
                
                kind = pi[i]["hand_cards_kind"]
                if kind == 6:
                    if pk_kind_6(pi[i]["hand_cards_kind"],  pj[i]["hand_cards_kind"]):
                        players[i], players[j] = players[j], players[i]
                        continue

                if kind == 5 and kind == 3:
                     if pk_kind_5_3(pi[i]["hand_cards_kind"],  pj[i]["hand_cards_kind"]):
                        players[i], players[j] = players[j], players[i]
                        continue