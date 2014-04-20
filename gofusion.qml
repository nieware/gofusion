import QtQuick 2.0
import QtQuick.Particles 2.0
import QtGraphicalEffects 1.0
//import Qt3D 1.0
import GoExtensions 1.0

Rectangle {
    id: screen
	width: 600; height: 675
	color: "navy"
	focus: true

    gradient: Gradient {
        GradientStop { position: 0.0; color: "#001166"; }
        //GradientStop { position: 0.8; color: "#875864"; }
        GradientStop { position: 1.0; color: "#001133"; }
    }

    SystemPalette { id: activePalette }
    
    Keys.onPressed: ctrl.handleKey(event.key)
    
    Rectangle {
        id: toolBar
        width: parent.width; height: 30
        color: "yellow"
        anchors.top: screen.top

        gradient: Gradient {
            GradientStop { position: 0.0; color: "#001133"; }
            GradientStop { position: 1.0; color: "#001166"; }
        }

        Button {
            anchors { left: parent.left; verticalCenter: parent.verticalCenter }
            text: "Restart"
            onClicked: ctrl.handleRestartButton()
        }

        Text {
            id: score
            objectName: "score"
            color: "white"
            anchors { right: parent.right; verticalCenter: parent.verticalCenter }
            text: "Score: 0"
        }
    }   
    
    Item {
        width: parent.width
        anchors { top: parent.top; bottom: toolBar.top }

        Item {
            id: gameCanvas
            objectName: "gameCanvas"

            property int score: 0
            //property int blockSize: 40

            width: parent.width //- (parent.width % blockSize)
            height: parent.height //- (parent.height % blockSize)
            anchors.centerIn: parent

            MouseArea {
                z: 100
                anchors.fill: parent
                onClicked: game.handleClick(mouse.x, mouse.y)
                
                Text {
                    id: message
                    objectName: "message"
                    font.pointSize: 24
                    color: "white"
                    y: 300 // verticalCenter doesn't work?!
                    z: 100
                    //verticalAlignment: Text.AlignVCenter
                    anchors {
                        //verticalCenter: parent.verticalCenter
                        horizontalCenter: parent.horizontalCenter
                    }
                    text: "GoFusion"
                }

                Text {
                    id: submessage
                    objectName: "submessage"
                    font.pointSize: 14
                    color: "white"
                    y: 350 // verticalCenter doesn't work?!
                    z: 100
                    //verticalAlignment: Text.AlignVCenter
                    anchors {
                        //verticalCenter: parent.verticalCenter
                        horizontalCenter: parent.horizontalCenter
                    }
                    text: "a '2048' clone by nieware"
                }

                Glow {
                    anchors.fill: message
                    radius: 8
                    samples: 16
                    color: "white"
                    source: message
                }

                Glow {
                    anchors.fill: submessage
                    radius: 4
                    samples: 16
                    color: "white"
                    source: submessage
                }
            }
        }
    }

	property var tileComponent: Component {
		id: tileComponent
		Tile {
			id: tile
            property int nvalue: 1
            property int zOrder: 0

            property real bounceY0: 0
            property real bounceY1: 0
            property bool bounceEnable: false    
            property real bounceDuration: 700

            property bool fallEnable: false    
            property real fallDuration: 2000

            x: 300; y: 300; z: zOrder
			width: 0; height: 0
            Behavior on x  {
                NumberAnimation  { duration: 500; easing.type: Easing.OutBounce; 
                    onRunningChanged: {
                        if (!running) {
                            ctrl.handleMoveAnimationDone();
                        }
                    } }
            }
            Behavior on y  {
                NumberAnimation  { duration: 500; easing.type: Easing.OutBounce }
            }
            Behavior on width  {
                NumberAnimation  { duration: 500; easing.type: Easing.OutBounce }
            }
            Behavior on height  {
                NumberAnimation  { duration: 500; easing.type: Easing.OutBounce }
            }

            SequentialAnimation on y {
                running: bounceEnable
                loops: Animation.Infinite
                NumberAnimation {to: bounceY1; duration: bounceDuration; easing.type: "OutQuad"}
                NumberAnimation {to: bounceY0; duration: bounceDuration; easing.type: "InQuad"}
            }

            NumberAnimation on y {
                running: fallEnable
                to: 2000; duration: fallDuration; easing.type: "OutQuad"
            }
    	}
	}


    ParticleSystem { id: sys }

    ImageParticle {
        system: sys
        source: "particle.png"
        color: "white"
        colorVariation: 1.0
        alpha: 0.1
    }

    property var emitterComponent: Component {
        id: emitterComponent
        Emitter {
            id: container
            system: sys
            Emitter {
                system: sys
                emitRate: 128
                lifeSpan: 600
                size: 16
                endSize: 8
                velocity: AngleDirection { angleVariation:360; magnitude: 60 }
            }

            property int life: 2600
            property real targetX: 0
            property real targetY: 0
            emitRate: 128
            lifeSpan: 600
            size: 24
            endSize: 8
            NumberAnimation on x {
                objectName: "xAnim"
                id: xAnim;
                to: targetX
                duration: life
                running: false
            }
            NumberAnimation on y {
                objectName: "yAnim"
                id: yAnim;
                to: targetY
                duration: life
                running: false
            }
            Timer {
                interval: life
                running: true
                onTriggered: ctrl.done(container)
            }
        }
    }
}
