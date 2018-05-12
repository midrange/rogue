The goal of the Rogue project is to build an AI to play Magic: The Gathering at a better-than-human level.

The initial goal is to develop similar functionality as https://github.com/andrewljohnson/CardAI while making it run faster.

## Running Hello World

If you are a developer, clone this repo into your `$GOPATH` in the `github.com/midrange/rogue` directory.

From the `rogue` directory:

```
go install ./... && play
```

You should see it print out something like:

```
 ~~~~~~ Welcome to Rogue ~~~~~~

1) Human vs AI
2) AI vs AI

Enter a number:
```