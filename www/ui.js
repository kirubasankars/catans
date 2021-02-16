
	

window.onload = function () {
    output = Grid(m)
    console.log(grid)
    var dedupeMap = {}

    var canvas = document.getElementById('board')
    paper.setup(canvas);

    canvas.onwheel = function(event) {
		var newZoom = paper.view.zoom; 
		var oldZoom = paper.view.zoom;
		
		if (event.deltaY > 0) {			
			newZoom = paper.view.zoom * 1.05;
		} else {
			newZoom = paper.view.zoom * 0.95;
		}
		
		var beta = oldZoom / newZoom;
		
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

    var PlacedNodeCount = 0
    var chits = [10, 2, 9, 12, 6, 4, 10, 9, 11, 3, 8, 8, 3, 4, 5, 5, 6, 11]
    var idx = 0
    output.nodes.forEach(hex => {
        if (hex.resource == "-" || hex.resource == "s" || hex.port) {
            return
        }

        var pr = new paper.Point(Math.floor((Math.random() * paper.view._viewSize.width) + 1), Math.floor((Math.random() * paper.view._viewSize.height) + 1))
        var p = new paper.Point(hex.x, hex.y)
        
        var hexagon = new paper.Path.RegularPolygon(pr, 6, hex.r - 10); 
        hexagon.tween({
            'position.x': hex.x,
            'position.y': hex.y
        }, { duration: 1000 }).then(function(){
            PlacedNodeCount += 1
            if (PlacedNodeCount == 19) {
                output.nodes.forEach(hex => {
                    hex.intersections.forEach(intersection => {
                        if (!intersection.hasPort) {
                            return
                        }
                        var myPath = new paper.Path();
                        myPath.strokeColor = 'black';          
                        myPath.strokeWidth = 4;
                        myPath.strokeCap = 'round';
                        myPath.add(new paper.Point(intersection.port.x, intersection.port.y));
                        myPath.add(new paper.Point(intersection.x, intersection.y));          
                        var text = new paper.PointText(new paper.Point(intersection.port.x, intersection.port.y ));
                        text.justification = 'center';
                        text.fillColor = 'black';
                        text.fontSize = '18pt'
                        text.content = intersection.port.resource;
                    })

                    
                })
            }
        })

      

        var text = new paper.PointText(pr);
        text.tween({
            'position.x': hex.x,
            'position.y': hex.y + 10
        },{  duration: 1000 })
        //text.justification = 'center';
        //text.fillColor = 'black';

       
        if (hex.resource != 'd') {
            // var cricle = new paper.Shape.Rectangle(new paper.Point(hex.x - 10, hex.y - 30), 20);
            // cricle.fillColor = '#fff'
            // cricle.strokeColor = '#fff'

            var token = new paper.PointText(pr);
            token.tween({
                'position.x': hex.x,
                'position.y': hex.y - 20
            },{  duration: 1000 })
            token.justification = "center"
            token.fillColor = 'black';
            token.content = chits[idx]
            idx ++
        }

        var color = ""
        if (hex.resource == 't') {            
            color = '#159739';            
            text.content = "Tree";
        }
        if (hex.resource == 'h') {            
            color = '#e46a2b';            
            text.content = "Brick";
        }
        if (hex.resource == 'p') {            
            color = '#91b60b';            
            text.content = "Whool";
        }
        if (hex.resource == 'f') {            
            color = '#f7be2c';            
            text.content = "Grain";
        }
        if (hex.resource == 'm') {            
            color = '#a8aeaa';            
            text.content = "Ore";
        }
        if (hex.resource == 'd') {            
            color = '#d7cf91';            
            text.content = "Desert";
        }
        if (hex.resource == 's') {            
            color = '#0967a6';                        
        }

        hexagon.fillColor = color;
        hexagon.strokeColor = color;
        hexagon.strokeColor.lightness += 0.09;
        hexagon.fillColor.lightness -= 0.09;
        hexagon.strokeWidth = 5;
    

        
    });



    output.intersections.forEach((ins) => {
        // var p = new paper.Point(ins.x, ins.y)
        // var shape = new paper.Shape.Circle(p, ins.r);
        // shape.strokeColor = 'red';

        // var text = new paper.PointText(p);
        // text.justification = 'center';
        // text.fillColor = 'black';
        // text.fontSize = '8pt'
        // text.content = ins.index;

        
    })

    // var p31 = grid.intersections[31]
    // var p34 = grid.intersections[34]

    // console.log(p31, p34)

    // var myPath = new paper.Path();
    // myPath.strokeColor = 'blue';          
    // myPath.strokeWidth = 8;
    // myPath.strokeCap = 'round';
    // myPath.add(new paper.Point(p31.x, p31.y));
    // myPath.add(new paper.Point(p34.x, p34.y));            

}