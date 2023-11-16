package telegram_bot

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cast"
)

// texas的逻辑
const (
	TEXAS_CARD_NUM = 52       //德州扑克牌的张数
	TEXAS_FACE_MIN = FACE_TWO //德州的小最牌面 为 2
)

type Card byte //牌
type Suit byte //花色
type Face byte //牌面

type CardSlice []Card

func (cs CardSlice) Len() int { return len(cs) }

// 按牌面从大小到排
func (cs CardSlice) Less(i, j int) bool {
	if cs[i].Face() == cs[j].Face() {
		return cs[i].Suit() > cs[j].Suit()
	}
	return cs[i].Face() > cs[j].Face()
}
func (cs CardSlice) Swap(i, j int) { cs[i], cs[j] = cs[j], cs[i] }

const (
	//花色
	SUIT_DIAMOND = 0x10 //方片(花色)
	SUIT_CLUB    = 0x20 //梅花(花色)
	SUIT_HEART   = 0x30 //红桃(花色)
	SUIT_SPADE   = 0x40 //黑桃(花色)
)

const (
	FACE_NONE  = 0
	FACE_TWO   = 0x02
	FACE_THERR = 0x03
	FACE_FOUR  = 0x04
	FACE_FIVE  = 0x05
	FACE_SIX   = 0x06
	FACE_SEVEN = 0x07
	FACE_EIGHT = 0x08
	FACE_NINE  = 0x09
	FACE_TEN   = 0x0A
	FACE_JACK  = 0x0B
	FACE_QUEEN = 0x0C
	FACE_KING  = 0x0D
	FACE_ACE   = 0x0E
)
const CARD_NONE = 0 //牌的无效值

const (
	CTYPE_NONE          = 0  //无牌型
	CTYPE_HIGHCARD      = 1  //高牌
	CTYPE_ONEPAIR       = 2  //一对
	CTYPE_TWOPAIR       = 3  //两对
	CTYPE_THREEKIND     = 4  //三条
	CTYPE_STRAIGHT      = 5  //顺子
	CTYPE_FLUSH         = 6  //同花
	CTYPE_FULLHOUSE     = 7  //葫芦
	CTYPE_KINGKONG      = 8  //四条(金刚)
	CTYPE_FLUSHSTRAIGHT = 9  //同花顺
	CTYPE_ROYALFLUSH    = 10 //皇家同花顺
)

// 获取牌的花色
func (c Card) Suit() Suit {
	return Suit(c & 0xF0)
}

// 获取牌的牌面
func (c Card) Face() Face {
	return Face(c & 0x0F)
}

func (s Suit) String() string {
	switch s {
	case SUIT_DIAMOND:
		return "♦"
	case SUIT_CLUB:
		return "♣"
	case SUIT_HEART:
		return "♥"
	case SUIT_SPADE:
		return "♠"
	default:
		return ""
	}
}

func (f Face) String() string {
	switch f {
	case FACE_TEN:
		return "T"
	case FACE_JACK:
		return "J"
	case FACE_QUEEN:
		return "Q"
	case FACE_KING:
		return "K"
	case FACE_ACE:
		return "A"
	default:
		if f >= FACE_TWO && f <= FACE_NINE {
			return fmt.Sprintf("%d", f)
		}
		return ""
	}
}

// 牌型
type CardType int

func (ct CardType) String() string {
	switch ct {
	case CTYPE_NONE:
		return "无牌型" //
	case CTYPE_HIGHCARD:
		return "高牌" //
	case CTYPE_ONEPAIR:
		return "一对" //
	case CTYPE_TWOPAIR:
		return "两对" //
	case CTYPE_THREEKIND:
		return "三条" //
	case CTYPE_STRAIGHT:
		return "顺子" //
	case CTYPE_FLUSH:
		return "同花" //
	case CTYPE_FULLHOUSE:
		return "葫芦" //
	case CTYPE_KINGKONG:
		return "四条(金刚)" //
	case CTYPE_FLUSHSTRAIGHT:
		return "同花顺" //
	case CTYPE_ROYALFLUSH:
		return "皇家同花顺" //
	default:
		return "未知" //
	}
}

// 牌的字符串表示
func (c Card) String() string {
	if c.Suit().String() == "" || c.Face().String() == "" {
		return fmt.Sprintf("?%d", c)
	}
	return fmt.Sprintf("%s%s", c.Suit(), c.Face())
}

// 检测相同花色的牌已按face从大到小排好序的牌是否成顺
// 返回成顺的牌
func checkFlushStraight(sc CardSlice, face_min Face) (CardSlice, bool) {
	//查找是否有5张及以上的同花
	if lsc := len(sc); lsc > 4 {
		//有五张以上相同的花色
		for i := 0; i < lsc-4; i++ {
			//成顺
			if sc[i].Face()-sc[i+4].Face() == 4 {
				return sc[i : i+5], true
			}
		}

		//A 2 3 4 5    A 6 7 8 9 T
		if sc[0].Face() == FACE_ACE && sc[lsc-1].Face() == face_min &&
			sc[lsc-4].Face() == (face_min+3) {
			bestCards := append(sc[lsc-4:], sc[0])
			return bestCards, true
		}
	}
	return sc, false
}

// 检测已按face从大到小排好序的牌是否成顺
// 返回成顺的牌
func checkStraight(sortCards CardSlice, face_min Face) (CardSlice, bool) {
	st := CardSlice{sortCards[0]}
	lensc := len(sortCards) //总牌数组的长度
	for i := 1; i < lensc; i++ {
		last := st[len(st)-1]
		card := sortCards[i]
		if last.Face()-card.Face() == 1 {
			st = append(st, card)
			if len(st) == 4 && card.Face() == face_min && sortCards[0].Face() == FACE_ACE {
				// A 2 3 4 5   A 6 7 8 9 T
				st = append(st, sortCards[0])
				return st, true
			} else if len(st) == 5 {
				return st, true
			}
		} else {
			//当前顺子的长度 + 剩余牌数 <4 就不能成顺
			if len(st)+(lensc-i-1) < 4 {
				return st, false
			} else if last.Face() != card.Face() {
				st = CardSlice{card}
			}
		}
	}
	return st, false
}

// 找出最好的牌型的5张牌: 四条 葫芦 三条 两对 一对
// sortCards 排序好的牌
// 按要求的牌面值的排序数组
// 返回
func findFaceBestCards(sortCards CardSlice, f1, f2 Face) CardSlice {
	var cs1, cs2, cs3 CardSlice //组成上述牌型的3部分
	for _, card := range sortCards {
		switch card.Face() {
		case f1:
			cs1 = append(cs1, card)
		case f2:
			cs2 = append(cs2, card)
		default:
			cs3 = append(cs3, card)
		}
	}
	best := cs1
	if len(cs2) > 0 {
		best = append(best, cs2...)
	}
	if len(cs3) > 0 {
		best = append(best, cs3...)
	}
	return best[:5]
}

// 根据公共牌和手牌 算出牌型和最佳手牌
func CalcTexasCards(public, hands CardSlice) (cardType CardType, bestCards CardSlice) {
	cards := make(CardSlice, 0, 7) //
	cards = append(cards, public...)
	cards = append(cards, hands...)
	if !sort.IsSorted(cards) {
		sort.Sort(cards)
	}

	suitCards := make(map[Suit]CardSlice) //相同花的的牌
	faceNum := make(map[Face]int)         //每种牌面的数量

	var four, three, two []Face //四条，三条， 一对的 牌面值
	for _, card := range cards {
		s := card.Suit()
		f := card.Face()

		suitCards[s] = append(suitCards[s], card)
		faceNum[f]++

		switch faceNum[f] {
		case 2:
			two = append(two, f)
		case 3:
			three = append(three, f)
		case 4:
			four = append(four, f)
		}
	}

	for _, sc := range suitCards {
		if lsc := len(sc); lsc > 4 {
			if bc, ok := checkFlushStraight(sc, TEXAS_FACE_MIN); ok {
				bestCards = bc
				if bc[0].Face() == FACE_ACE {
					cardType = CTYPE_ROYALFLUSH
				} else {
					cardType = CTYPE_FLUSHSTRAIGHT
				}
				return
			}
			//同花
			cardType = CTYPE_FLUSH
			bestCards = sc[0:5]
			return
		}
	}

	if len(four) > 0 {
		//有四条
		cardType = CTYPE_KINGKONG
		bestCards = findFaceBestCards(cards, four[0], 0)
	} else if len(three) > 0 && len(two) > 1 {
		//葫芦
		f2 := two[0]
		if f2 == three[0] {
			f2 = two[1]
		}
		cardType = CTYPE_FULLHOUSE
		bestCards = findFaceBestCards(cards, three[0], f2)
	} else if st, ok := checkStraight(cards, TEXAS_FACE_MIN); ok {
		cardType = CTYPE_STRAIGHT
		bestCards = st
	} else if len(three) > 0 {
		cardType = CTYPE_THREEKIND
		bestCards = findFaceBestCards(cards, three[0], 0)
	} else if len(two) > 1 {
		cardType = CTYPE_TWOPAIR
		bestCards = findFaceBestCards(cards, two[0], two[1])
	} else if len(two) > 0 {
		cardType = CTYPE_ONEPAIR
		bestCards = findFaceBestCards(cards, two[0], 0)
	} else {
		cardType = CTYPE_HIGHCARD
		bestCards = cards[:5]
	}
	return
}

func DataStringBytesToCardsSlice(dataStr string) CardSlice {
	var dataStrs = strings.SplitAfter(dataStr, " ")
	fmt.Printf("dataStrs: %+v\n", dataStrs)
	var dataStrArr []string
	for _, dataStr := range dataStrs {
		dataStrArr = append(dataStrArr, strings.TrimSpace(dataStr))
	}
	fmt.Printf("dataStrArr: %+v\n", dataStrArr)
	var dataBytes []byte
	for _, dataStr := range dataStrArr {
		dataBytes = append(dataBytes, byte(cast.ToInt8(dataStr)))
	}
	fmt.Printf("dataBytes: %+v\n", dataBytes)
	var cardSlice CardSlice
	for _, dataByte := range dataBytes {
		cardSlice = append(cardSlice, Card(dataByte))
	}
	fmt.Printf("cardSlice: %+v\n", cardSlice)
	return cardSlice
}

func AnalyzePoker(dataStr string, seatNum int) string {
	cardSlice := DataStringBytesToCardsSlice(dataStr)

	var (
		hCards CardSlice = cardSlice[:seatNum*2]
		pCards CardSlice = cardSlice[seatNum*2 : seatNum*2+5]
		seats  []CardSlice
	)

	for i := 0; i < seatNum; i++ {
		var seat CardSlice
		seat = append(seat, hCards[i], hCards[i+seatNum])
		seats = append(seats, seat)
	}

	var output string
	output += fmt.Sprintf("公:%v\n", pCards)
	for i, seat := range seats {
		cardType, best := CalcTexasCards(pCards, seat)
		output += fmt.Sprintf("[%v]%v -> %v = %v\n", i, seat, cardType, best)
	}

	return output
}

func AnalyzeFish(dataStr string) string {
	cardSlice := DataStringBytesToCardsSlice(dataStr)

	var (
		pCards CardSlice = cardSlice[:5]
		rCards CardSlice
		lCards CardSlice
	)
	rCards = append(rCards, cardSlice[6])
	rCards = append(rCards, cardSlice[8])

	lCards = append(lCards, cardSlice[5])
	lCards = append(lCards, cardSlice[7])

	rCardType, rB := CalcTexasCards(pCards, rCards)
	lCardType, lB := CalcTexasCards(pCards, lCards)

	var winner string
	if rCardType > lCardType {
		winner = rCardType.String()

	} else if rCardType < lCardType {
		winner = lCardType.String()

	} else if rCardType == lCardType {
		sam := true
		for i, c := range rB {
			if c.Face() != lB[i].Face() {
				sam = false
			}
		}
		if sam {
			winner = "平"
		} else {
			winner = rCardType.String()
		}
	}

	if rCards[0].Face() == rCards[1].Face() || lCards[0].Face() == lCards[1].Face() {
		winner = fmt.Sprintf("%v - %v", winner, "手牌成对")
	}

	if (rCards[0].Face() == rCards[1].Face() && rCards[0].Face() == FACE_ACE) ||
		(lCards[0].Face() == lCards[1].Face() && lCards[0].Face() == FACE_ACE) {
		winner = fmt.Sprintf("%v - %v", winner, "手牌AA")
	}

	if rCards[0].Face() == 2 && rCards[1].Face() == FACE_TWO {
		winner = fmt.Sprintf("%v - %v", winner, "R手牌22")
	}

	return fmt.Sprintf("公:%v\n左:%v -> %v = %v\n右:%v -> %v = %v\n获胜牌型:%v\n",
		pCards, lCards, lCardType, lB, rCards, rCardType, rB, winner)
}
