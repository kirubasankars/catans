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

	fmt.Println(gc.Phase, gc.Action.Name, gc.getCurrentPlayer().Id)

	if gc.Phase == Phase2 || gc.Phase == Phase3 {
		switch playerAction.Name {
		case ActionPlaceSettlement:
			{
				HandlePlaceInitialSettlement(gc)
			}
		case ActionPlaceRoad:
			{
				HandlePlaceInitialRoad(gc)
			}
		}
	}

	switch playerAction.Name {
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
