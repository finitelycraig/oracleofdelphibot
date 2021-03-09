# Oracle of Delphi Bot

This Twitch chat bot provides short chunks of information about Nethack, a
classic and certainly very pretty roguelike game.

## Accessing the bot

To use the bot in your own chat simply visit [the bot's channel] and type
`!join` in chat to ask the bot to join your channel's chat.  You may need to
send the bot a whispered message and ask to be added to the list of allowed
channels.

## Using the bot

The bot currently accepts input that queries items, monsters and properties.

For example, `!long-sword` returns 

```
A long-sword does d8/d12 damage. It is made of iron, weighs 40, and
is valued at 15zm. Works your skill with long swords.
```

## Running it yourself at a channel of your choosing

The code for running the bot is in bot.go and all the information given by
the bot is contained in the various yaml files included with this repository.

The bot gets the initial channel to join, it's channel and it's oauth token from
environment variables named TWITCHBOT, TWITCHCHANNEL and TWITCHOAUTH: ensure
that these are set.  

## Todos

- [ ] Complete the yaml files with info from the wiki
- [x] Include a mapping from alternative item names to proper item names
- [ ] Include a price ID helper command, like `!idwand 100`

[the bot's channel]: https://twitch.tv/oracleofdeplhibot
