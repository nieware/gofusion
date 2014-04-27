# GoFusion - a Go-QML port of "2048"

![GoFusion Screenshot](https://raw.githubusercontent.com/nieware/gofusion/master/_data/screenshot.gif "GoFusion Screenshot")

What is this?
-------------

This is yet another clone of that annoyingly addictive sliding tile game [2048](http://gabrielecirulli.github.io/2048/)
by Gabriele Cirulli (which in turn is a clone of another game, which is a near clone of yet another one...), using Go and go-qml.

This was my entry for the [Go QML contest](http://blog.labix.org/2014/03/13/go-qml-contest) held at [GopherCon 2014](http://www.gophercon.com), and, much to my surprise, it [actually won](http://blog.labix.org/2014/04/25/qml-contest-results)! With this motivation boost, I'll continue work on the project, first by documenting the code and cleaning it up a bit, then I'll start implementing the ideas described below.

For historic reference, this is a link to the project as it looked like when I submitted it for the contest: [https://github.com/nieware/gofusion/tree/6a636c37989464275dd758e66c29256af6cffe8a](https://github.com/nieware/gofusion/tree/6a636c37989464275dd758e66c29256af6cffe8a)


What can I do with it?
----------------------

Basically, you can do almost everything you can with the original - play the game using the keyboard (in a hopefully pleasant-looking way) and keep 
a local highscore.


How do I compile and run it?
----------------------------

You should be able to compile and run it if you are able to compile and run the examples from the go-qml package. See [here](https://github.com/go-qml/qml) for
instructions on how to set up your system for doing that.

Once the "gofusion" binary is compiled successfully, place the files "gofusion.qml", "Button.qml" and "particle.png" into the same directory as the binary. The
"model" subdirectory (and its content) should be a subdirectory of the directory where the binary is located.


Any ideas for expanding it?
---------------------------

Lots...

First, you may ask yourself why all the tiles are 3D models. Well, I was hoping to be able to display them with perspective projection and
do some cool animation effects. However, I found out that this is only possible with Qt 3D, which is still experimental, and I wasn't able 
to get into it (yet) because of lack of time. So, for now, every tile is in its own scene graph (I was at least able to use that to vary the
colors of the tiles by changing the lighting color) and has orthogonal projection, so there is not much 3D to be seen at the moment.

Second, one could write a web service (and host it), so people could store their highscores online. That part is trivial. The non-trivial part
would be making it hack-proof and DOS-proof...

Third, the application only supports keyboard input for now, mouse gestures (a.k.a. swiping) should also be supported for mobile devices.

