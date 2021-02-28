package main

const ErrInvalidOperation = "invalid operation"

// Phases
const Phase1 = "SETUP"
const Phase2 = "INITIAL_SETTLEMENT1"
const Phase3 = "INITIAL_SETTLEMENT2"
const Phase4 = "IN_GAME"

// Actions
const ActionPlaceSettlement = "PLACE_SETTLEMENT"
const ActionPlaceRoad = "PLACE_ROAD"
const ActionDiscardCards = "DISCARD_CARDS"
const ActionPlaceRobber = "PLACE_ROBBER"
const ActionSelectToSteal = "SELECT_TO_ROB"
const ActionRollDice = "ROLL_DICE"
const ActionDevPlaceRoad1 = "PLACE_ROAD1"
const ActionDevPlaceRoad2 = "PLACE_ROAD2"
const ActionTurn = "TURN"

//Development Cards
const DevCardKnight = 0
const DevCardVPPoint = 1
const DevCardMonopoly = 2
const DevCard2Resource = 3
const DevCard2Road = 4

const CardLumber = 0
const CardBrick = 1
const CardWool = 2
const CardGrain = 3
const CardOre = 4
