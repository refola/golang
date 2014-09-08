golang
======

Mark's stuff written in the Go language


This should include all the interesting/useful code I've written in Go. Here's what I've written in Go so far:

* A minimal wiki server. This needs to be updated for Go 1.0 changes. It's missing a lot of things a "real world" wiki server should have, but I think the wiki syntax to html part is pretty neat.
* A backup program that's a frontend for rsync. It uses a json configuration file to run several backups in sequence depending on which (external) locations are available. I've been using this for several months for my computer's backups. I've been thinking about extending/generalizing it a bunch of ways to make it more flexible (e.g., using other backup backends besides rsync) and maybe even make it easier to make a new config file. I've also been trying to get topological sort working so it can properly sort backups by dependency order.
* Various minor packages to augment Go's standard library for my programs
* A bunch of broken/incomplete/unstarted stuff that isn't too useful (yet?)


I think the following should be pulled out into separate repositories after getting started with this one.

* backup
* server
* telemach
* util and other minor packages used in the above

Then this repo will be for just miscellaneous half-baked ideas.


Dedication and licensing:
======

I cannot survive on my own. I literally owe my life to this world. As long as I'm alive, I must see value in my life that the world enables, so by extension I must trust and value the entities of this world overall. Thus, as long as I value my life enough to keep it, I dedicate my work to maximizing universal utility.

I've seen how copyright and copyleft have kept different peoples' innovations and ideas from being combined into greater wholes. Copyfree licenses are better, but copyright and attribution requirements can get messy over the years as several authors' names get added. The point is, I don't want copyright to get in the way of this code being as useful as it can be. Thus I'm "licensing" my code under the "unlicense" public domain¹ option.

If you have your reasons to make your own changes and not release them, that's great! I'm glad my code is useful to you. If you want to change and release my code under some other license, go for it! I know there are lots of reasons for using a particular license, and I want my code to work with all these cases. However, if you want to contribute code back to this project, it must be public domain, so it can benefit as many people as possible.

Finally, here's what I'm asking² for in return: Do whatever you reasonably can that you think is fair payment. If you're busy, you could be happy to use this code. If you have time, you could get involved by leaving a comment, reporting a bug, or submitting a patch. If you're rich, you could give me a donation or even hire me for your own project.


Notes:

¹ Even though my code is public domain/unlicense, some of it pulls in other code to build the final program (see telemach). I am not a lawyer, but I think the end result in this case must be licensed under some combination of the upstream licenses.

² I'm asking, not requiring. You're free to be sad about my code, but I think you'll find happiness much nicer.


Do stuff and be happy!
- Mark Haferkamp
