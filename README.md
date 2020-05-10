# guesswho-discord

A simple discord bot that enables you to play "guess who".

## Running the program

Run the program using `go run main.go -t <bot_token>` where `<bot_token>` can be found in the bot's settings at https://discord.com/developers.

## Playing the game

Start the game by messaging the bot with `new`. Get players to join the game by also messaging the bot `join`. When all players are in, send `start` to the bot to begin the game. It will randomize the players and ask everyone to choose a character for another player in the game. Message the bot with the chosen character. When everyone has sent a character, the bot will distribute everyone else's characters.
