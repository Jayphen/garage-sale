# Buy my things

This is a little project I am building to learn how to use Pocketbase as a framework with Go. Also I need to sell my things so that I can move across the world.

This is turning out to be a little more difficult than expected, as the documentation for Pocketbase as a framework is very sparse, and I am new to Go. Let's give it a shot anyway.

## What is it tho?

It's a little twist on an auction. The products all have a maximum and a minimum price. When the store is open, the prices will drop every 2 minutes until they reach the minimum price at the end of trading. The prices then reset to the max in the evening. Since there is only 1 of each item, customers have to decide whether to risk waiting until the price drops or not.

## Running this

`task watch_all` will watch templ for changes to the htmx templates and recompile them, and also generate a new build if anything in the go application changes. It works most of the time.
