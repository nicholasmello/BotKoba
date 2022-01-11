# Rocket League Bot Koba

## About
Rocket League Bot designed to face Koba in Rocket League. Currently the bot only works in 1v1.

## Create your own

If you are looking to create your own bot using Go take a look at [Trey2K](https://github.com/Trey2k)/[https://github.com/Trey2k/RLBotGo/](RLBotGo) for the Go specific framework as well as my Always Toward Ball Agent ([ATBA](https://github.com/xonmello/RLBotGoATBA)) bot for the basic on movement. 

In addition to Go, many other languages have framework for RLBot. Choose one on the [RLBot website](http://rlbot.org/) or join the [RLBot Discord server](https://discord.gg/zbaAKPt).

## Run this Bot

1. Download and run RLBot for Windows from the [RLBot website](http://rlbot.org/)
2. Add the bot to the Bots directory or add the path to the bot to the sources
3. Compile the bot `go build ./` *This requires [go](https://go.dev/dl/) to be installed*
4. Use the RLBot GUI to start the match

## Development

### States

BotKoba is a state based bot. This means it does different actions depending on the state of the game. The state of the game could refer to if it is during a kickoff, if the ball is in it's own corner, etc. This is compared to using machine learning to train a bot. In my implimentation, when a state is decided it records when it started the state for later and allots some time it has before it can decide if it wants to switch or not. This allows for states such as the kickoff to keep going after the trigger (Game being stopped) ends. 

### Utility Functions

Useful utility functions found in [utils.go](https://github.com/xonmello/BotKoba/blob/main/utils.go) are functions that are used or likely will be used by multiple states and help with control of the bot. This lets the states focus on what should happen rather than how to do it. These functions include steerToward and flipToward. These are likely useful to other bots and can be copied as stated in the LICENSE.