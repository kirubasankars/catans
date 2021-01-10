package game

import (
	"fmt"
)

func TimeoutHandler(gc *GameContext) {
	if gc.GetAction() == nil {
		return
	}

	playerAction := gc.GetAction()
	if IsActionTimeout(*gc, *playerAction) {
		return
	}
	fmt.Println(gc.Phase, gc.Action.Name, gc.GetCurrentPlayer().Id)
	switch playerAction.Name {
	case ActionPlaceSettlement:
		{
			if gc.Phase == Phase2 || gc.Phase == Phase3 {
				HandlePlaceInitialSettlement(gc)
			}
		}
	case ActionPlaceRoad:
		{
			if gc.Phase == Phase2 || gc.Phase == Phase3 {
				HandlePlaceInitialRoad(gc)
			}
		}
	case ActionDiscardCards:
		{
			HandleDiscardCards(gc)
		}
	case ActionRollDice:
		{
			HandleRollDice(gc)
		}
	}

	err := gc.EndAction()
	if err != nil {
		fmt.Println(err)
	}
}
