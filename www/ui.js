window.onload = function () {
    $.get('/create_game', function (g){
        g = JSON.parse(g)
        $.get("/join_game?g=" + g.g + "&user_id=0", function () {
            a(JSON.parse(g.g))
        })
    })

    function a(game) {
        console.log(game)
        $.get('/board?g='+ game, function (output) {

            output = JSON.parse(output)

            console.log(output)
            var printerFriendly = true

            var canvas = document.getElementById('board')
            paper.setup(canvas);

            renderLand()
            renderHexagon()
            renderToken()
            renderLabel()
            renderPort()
            renderIntersection()

            function getColor(r) {
                switch (r) {
                    case 'l':
                        return '#159739'
                    case 'b':
                        return '#e46a2b'
                    case 'w':
                        return '#91b60b'
                    case 'g':
                        return '#f7be2c'
                    case 'o':
                        return '#a8aeaa'
                    case 't':
                        return '#159739'
                    case 'h':
                        return '#e46a2b'
                    case 'p':
                        return '#91b60b'
                    case 'f':
                        return '#f7be2c'
                    case 'm':
                        return '#a8aeaa'
                    case 'd':
                        return '#d7cf91'
                    case '?':
                        return '#fff'
                }
            }

            function getLabel(r) {
                switch (r) {
                    case 'l':
                        return 'lumber'
                    case 'b':
                        return 'brick'
                    case 'w':
                        return 'wool'
                    case 'g':
                        return 'grain'
                    case 'o':
                        return 'ore'
                    case 't':
                        return 'lumber'
                    case 'h':
                        return 'brick'
                    case 'p':
                        return 'wool'
                    case 'f':
                        return 'grain'
                    case 'm':
                        return 'ore'
                    case 'd':
                        return 'desert'
                    case '?':
                        return '?'
                }
            }

            function renderLand() {
                if(!printerFriendly) {
                    var _hexagons = {}
                    output.hexagons.forEach(el => {
                        _hexagons[el.id] = el
                    })

                    output.intersections.forEach(el => {
                        if (el.hexagons.length == 2) {
                            el.hexagons.forEach(h => {
                                var hx = _hexagons[h]
                                var point = new paper.Point(hx.x, hx.y)
                                var cricle = new paper.Path.Circle(point, hx.r + 6);
                                cricle.fillColor = '#20abfeff';
                            })
                            el.hexagons.forEach(h => {
                                var hx = _hexagons[h]
                                var point = new paper.Point(hx.x, hx.y)
                                var cricle = new paper.Path.Circle(point, hx.r + 3);
                                cricle.fillColor = '#53cbffff';
                            })
                            el.hexagons.forEach(h => {
                                var hx = _hexagons[h]
                                var point = new paper.Point(hx.x, hx.y)
                                var cricle = new paper.Path.Circle(point, hx.r);
                                cricle.fillColor = '#cce9feff';
                            })
                        }
                    })

                    output.intersections.forEach(el => {
                        if (el.hexagons.length == 3) {
                            var point = new paper.Point(el.x, el.y)
                            var cricle = new paper.Path.Circle(point, el.r);
                            cricle.fillColor = '#d7cf91';
                        }
                    })

                    output.hexagons.forEach(el => {
                        if (el.port) {
                            return
                        }
                        var point = new paper.Point(el.x, el.y)
                        var cricle = new paper.Path.Circle(point, el.r - 3);
                        cricle.fillColor = '#d7cf91';
                    })
                }
            }

            function renderHexagon() {

                output.hexagons.forEach(el => {
                    if (el.port) {
                        return
                    }

                    var color = "#fff"

                    if (!printerFriendly) {
                        color = getColor(el.terrain)
                    }

                    var point = new paper.Point(el.x, el.y)

                    var hexagon = new paper.Path.RegularPolygon(point, 6, el.r - 7);
                    hexagon.strokeWidth = 4;
                    hexagon.strokeColor = color;
                    hexagon.fillColor = color;
                    hexagon.strokeColor.lightness -= 0.09;
                    hexagon.fillColor.lightness += 0.09;
                })

            }

            function renderToken() {
                var dots = { 2: 1, 3: 3, 4: 3, 5: 4, 6: 5, 8: 5, 9: 4, 10: 3, 11: 2, 12: 1 }
                output.hexagons.forEach(el => {
                    if (el.port || el.terrain == 'd') {
                        return
                    }

                    var rectangle = new paper.Rectangle(new paper.Point(el.x - 16, el.y - 25), new paper.Size(32, 26));
                    var cornerSize = new paper.Size(5, 5);
                    var shape = new paper.Shape.Rectangle(rectangle, cornerSize);
                    shape.strokeWidth = 0.2;
                    shape.fillColor = '#fff'

                    if (printerFriendly) {
                        shape.strokeColor = 'red';
                    }

                    var token = new paper.PointText(new paper.Point(el.x, el.y - 10));
                    token.justification = 'center';
                    if (el.token == 6 || el.token == 8) {
                        token.fillColor = 'red';
                    } else {
                        token.fillColor = 'black';
                    }
                    token.fontSize = '12pt'
                    token.content = el.token;

                    var dot = new paper.PointText(new paper.Point(el.x, el.y - 3));
                    dot.fillColor = 'black';
                    dot.justification = 'center';
                    dot.fontSize = '15pt'
                    var s = ""
                    for (var i = 0; i < dots[el.token]; i++) {
                        s += "."
                    }
                    dot.content = s;
                })
            }

            function renderLabel() {
                output.hexagons.forEach(el => {
                    if (el.port) {
                        return
                    }
                    var text = new paper.PointText(new paper.Point(el.x, el.y + 20));
                    text.justification = 'center';
                    text.fillColor = 'black';
                    text.fontSize = '8pt'
                    text.content = getLabel(el.terrain);

                    if (printerFriendly) {
                        var textIndex = new paper.PointText(new paper.Point(el.x, el.y + 35));
                        textIndex.justification = 'center';
                        textIndex.fillColor = 'black';
                        textIndex.fontSize = '6pt'
                        textIndex.content = el.id
                    }
                })
            }

            function showIntersections(path1, path2) {
                var intersections = path1.getIntersections(path2);
                for (var i = 0; i < intersections.length; i++) {
                    new paper.Path.Circle({
                        center: intersections[i].point,
                        radius: 5,
                        fillColor: '#009dec'
                    })
                }
            }

            function renderPort() {
                output.intersections.forEach(el => {
                    if (!el.port) {
                        return
                    }

                    //calculate line length
                    var path = new paper.Path.Line({
                        from: [el.port.x, el.port.y],
                        to: [el.x, el.y]
                    });

                    var portHex = new paper.Path.Circle(new paper.Point(el.port.x, el.port.y), el.r + 15);
                    var pointA = path.getIntersections(portHex);

                    var portIns = new paper.Path.Circle(new paper.Point(el.x, el.y), el.r);
                    var pointB = path.getIntersections(portIns);
                    //calculate line length

                    var path = new paper.Path.Line({
                        from: [pointA[0].point.x, pointA[0].point.y],
                        to: [pointB[0].point.x, pointB[0].point.y]
                        //to: [el.x, el.y]
                    });
                    path.strokeColor = 'black';
                    if (printerFriendly) {
                        path.strokeColor = '#a8aeaa'
                    }
                    path.strokeWidth = 4;
                    path.strokeCap = 'round';

                    var cricle = new paper.Path.Circle(new paper.Point(el.port.x, el.port.y), el.r + 5);
                    cricle.fillColor = getColor(el.port.resource);
                    cricle.strokeColor = getColor(el.port.resource);
                    cricle.fillColor.lightness += 0.09;
                    cricle.strokeColor.lightness -= 0.09;

                    if (printerFriendly) {
                        cricle.fillColor = '#fff'
                        cricle.strokeColor = '#a8aeaa';
                    }

                    var portLabel = new paper.PointText(new paper.Point(el.port.x, el.port.y - 2));
                    portLabel.justification = 'center';
                    portLabel.fillColor = 'black';
                    portLabel.fontSize = '7pt'
                    if (el.port.resource == "?") {
                        portLabel.content = "3:1";
                    } else {
                        portLabel.content = " 2:1";
                    }

                    var portlabel = new paper.PointText(new paper.Point(el.port.x + 0.5, el.port.y + 7));
                    portlabel.justification = 'center';
                    portlabel.fillColor = 'black';
                    portlabel.fontSize = '5pt'
                    portlabel.content = getLabel(el.port.resource)
                })
            }

            function renderIntersection() {
                output.intersections.forEach(el => {
                    if (printerFriendly) {
                        var cricle = new paper.Shape.Circle(new paper.Point(el.x, el.y), el.r);
                        cricle.fillColor = '#a8aeaa';
                        cricle.opacity = 0.4;
                        cricle.strokeColor = 'black';

                        var text = new paper.PointText(new paper.Point(el.x, el.y + 3));
                        text.justification = 'center';
                        text.fillColor = 'black';
                        text.fontSize = '8pt'
                        text.content = el.id;
                    }
                })
            }

            // var t = true
            // var r = 0.099
            // var m = 30

            // paper.view.onFrame = function (event) {
            //     if (t) {
            //         cricles.forEach(c => {
            //             c.radius = 10
            //         })
            //     } else {
            //         cricles.forEach(c => {
            //             c.radius = 8.5
            //         })
            //     }

            //     if (event.count % m == 0) {
            //         t = !t
            //     }
            // }

            paper.view.scale(1.6)
            paper.view.center = new paper.Point(400, 300)

            canvas.onwheel = function (event) {
                var newZoom = paper.view.zoom;
                var oldZoom = paper.view.zoom;

                if (event.deltaY < 0) {
                    newZoom = paper.view.zoom * 1.05;
                } else {
                    newZoom = paper.view.zoom * 0.95;
                }

                var beta = oldZoom / newZoom;

                if (newZoom > 2 || newZoom < 0.5) {
                    return
                }

                var mousePosition = new paper.Point(event.offsetX, event.offsetY);

                //viewToProject: gives the coordinates in the Project space from the Screen Coordinates
                var viewPosition = paper.view.viewToProject(mousePosition);

                var mpos = viewPosition;
                var ctr = paper.view.center;

                var pc = mpos.subtract(ctr);
                var offset = mpos.subtract(pc.multiply(beta)).subtract(ctr);

                paper.view.zoom = newZoom;
                paper.view.center = paper.view.center.add(offset);

                event.preventDefault();
                paper.view.draw();
            }
        })
    }

    var conn = new WebSocket("ws://" + document.location.host + "/ws");
    conn.onclose = function (evt) {

    };
    conn.onmessage = function (evt) {
        var messages = evt.data.split('\n');
        console.log(messages)
    };
}