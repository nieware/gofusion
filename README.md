# Go-QML port of "2048"


What is this?
-------------

This is yet another clone of that annoyingly addictive sliding tile game [2048](http://gabrielecirulli.github.io/2048/)
by Gabriele Cirulli (which in turn is a clone of another game, which is a near clone of yet another one...), using Go and go-qml.


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

