package main

import "time"

type GameAction struct {
	Name    string
	Timeout time.Time
}

func (context GameContext) getAction() *GameAction {
	return &context.Action
}

func (context GameContext) getActionString() string {
	return context.Action.Name
}

func (context *GameContext) scheduleAction(action string) {
	timeoutDuration := context.getActionTimeout(action)
	timeOut := time.Now()
	if context.Action.Name == ActionTurn {
		timeOut = context.Action.Timeout
	}
	if action != ActionTurn {
		timeOut = timeOut.Add(timeoutDuration)
	}
	context.Action = GameAction{Name: action, Timeout: timeOut}
}

func (context *GameContext) getActionTimeout(action string) time.Duration {
	timeout := 0
	if context.GameSetting.Speed == 0 {
		switch action {
		case ActionTurn:
			timeout = 30
		case ActionRollDice:
			timeout = 10
		case ActionPlaceSettlement:
			if context.phase == Phase2 || context.phase == Phase3 {
				timeout = 12
			}
		case ActionPlaceRoad:
			if context.phase == Phase2 || context.phase == Phase3 {
				timeout = 15
			}
		case ActionDevPlaceRoad1:
				timeout = 15
		case ActionDevPlaceRoad2:
				timeout = 15
		case ActionDiscardCards:
			timeout = 15
		case ActionPlaceRobber:
			timeout = 10
		case ActionSelectToSteal:
			timeout = 10
		}
	}
	return time.Duration(timeout) * time.Second
}

func (context *GameContext) isActionTimeout() bool {
	diff := time.Now().Sub(context.Action.Timeout).Seconds()
	if diff <= 0 {
		return true
	}
	return false
}
