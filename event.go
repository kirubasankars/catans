package main

import "fmt"


func (context *GameContext) EventRolled(dice int) {
	context.Events = append(context.Events, fmt.Sprintf(`{id:%d,"type":"dice",player:%d,dice:%d}`, context.EventID, context.CurrentPlayerID, dice))
	context.publishMessage()
	context.EventID ++
}

func (context *GameContext) EventCardDistributed(playerID, cardType, count int) {
	context.Events = append(context.Events, fmt.Sprintf(`{id:%d,"type":"card_distribution",player:%d,cardtype:%d,count:%d}`, context.EventID, playerID, cardType, count))
	context.publishMessage()
	context.EventID ++
}

func (context *GameContext) EventBoughtDevelopmentCard() {
	context.Events = append(context.Events, fmt.Sprintf(`{id:%d,"type":"bought_dev_card",player:%d}`, context.EventID, context.CurrentPlayerID))
	context.publishMessage()
	context.EventID ++
}

func (context *GameContext) EventPutSettlement(intersection int) {
	context.Events = append(context.Events, fmt.Sprintf(`{id:%d,"type":"settlement",player:%d,"intersection":%d}`, context.EventID, context.CurrentPlayerID, intersection))
	context.publishMessage()
	context.EventID ++
}

func (context *GameContext) EventPutRoad(road [2]int) {
	context.Events = append(context.Events, fmt.Sprintf(`{id:%d,"type":"road",player:%d,"road":[%d,%d]}`, context.EventID, context.CurrentPlayerID, road[0], road[1]))
	context.publishMessage()
	context.EventID ++
}