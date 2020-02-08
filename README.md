# Planning

0. Evernote API License - seems OK, it's free but they reserve all rights. 
1. We need go SDK - make an adapter over ?python? SDK somehow (or maybe c++, could be easier). Or maybe it's not that hard to write it yourself. See Evernote repos, maybe there's something there. There seem to be some very shady / small ones already (https://github.com/dreampuf/evernote-sdk-golang)
2. 2 parts of the app basically: TUI (terminal widgets) + "backend" (communication with Evernote API)
3. How to structure the app?
4. Which terminal library to use for drawing?
5. Make sure it syncs smoothly

Ad 3)
Start with just main package
https://github.com/golang-standards/project-layout
https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1#.ds38va3pp
REMEMEBER TO MODEL APPLICATION LANGUAGE

Ad 4)
A nice breakdown here https://appliedgo.net/tui/
I'd like to have at least grid / flexbox supported, and a nice "framework" for creating new widgets.
https://github.com/rivo/tview seems promising
https://github.com/jroimartin/gocui also seems nice. Minimalistic.
There is also somewhat not maintained https://github.com/nsf/termbox-go

Og, decided to go with **gocui**

Nice features to have:
* customizable key-bindings (vim like navigation)
* Markdown support
* Maybe 'Edit in vim' feature?

# ROADMAP

To make this fun.

1. Start with writing a terminal app
	* Create a good "data abstraction", that would allow for multiple storage patterns underneath easily (disk storage, evernote API sync)
	* Create UI that allows for basic operations
	* Create "more specific" UI that fits my needs and is only so close to evernote's one as it needs to be
2. Implement the actual storage underneath

