# pokedex - CLI-based REPL pokedex using PokeAPI

## Ideas for Improvements and Expansion
* Update the CLI to support the "up" arrow to cycle through previous commands (constant thought during manual testing that would be very helpful)
* Simulate battles between pokemon
* Add more unit tests for commands
* Refactor code to organize it better and make it more testable
    * specifically, break out commands to individual files
    * make smaller funcs when possible
* Keep pokemon in a "party" and allow them to level up
* Allow for pokemon that are caught to evolve after a set amount of time
* Persist a user's Pokedex to disk so they can save progress between sessions
* Use the PokeAPI to make exploration more interesting. For example, rather than typing the names of areas, maybe give the user a choice of areas and then they can just type "left" or "right"
* Random encounters with wild pokemon
* Adding support for different types of balls (Pokeballs, Great Balls, Ultra Balls, etc), which have different chances of catching pokemon