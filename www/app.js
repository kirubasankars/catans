
function Grid(map) {

    var r = 70;

    function makeGrid(h, w, tileConfig) {
        var x = r;
        var y = r + 20;
        var b = 0.27 * r;
        var grid = []
        var smallRow = true
        var idx = 0

        for (var i = 0; i < h; i++) {
            for (var j = 0; j < w; j++) {
                var tc = tileConfig[idx]
                hx = { x: x, y: y, r: r, index: idx, resource: tc, intersections: [] };
                if (tc.length == 3) {
                    hx.port = true
                    hx.resource = tc[1]
                    hx.direction = parseInt(tc[2])
                }
                grid[idx] = hx;
                x = x + (r * 2) - b;
                idx++;
            }

            if (!smallRow) {
                x = r;
                smallRow = true;
            } else {
                x = (r * 2) - b / 2;
                smallRow = false;
            }
            y = y + (r / 2) + r
        }
        return grid
    }

    function parse(m) {
        var o = m.replace(/[\t ]/g, '')
        segments = o.split(/\n/)
        var output = []
        for (var i = 0; i < segments.length; i++) {
            var row = segments[i].split(',')
            if (row.length > 1) {
                var newRow = []
                for (var r = 0; r < row.length; r++) {
                    var v = row[r]
                    if (v.length == 0) {
                        continue
                    }
                    newRow.push(v)
                }
                output.push(newRow)
            }
        }

        var tiles = []
        for (var i = 0; i < output.length; i++) {
            var row = output[i]
            for (var j = 0; j < row.length; j++) {
                tiles.push(row[j])
            }
        }

        var h = output.length
        var w = output[0].length

        return { tiles, h, w }
    }

    function circlesIntersect(x1, y1, r1, x2, y2, r2) {
        return (x1 - x2) * (x1 - x2) + (y1 - y2) * (y1 - y2) <= (r1 + r2) * (r1 + r2)
    }

    function makeIntersections(grid) {
    
        var a = 2 * Math.PI / 6;
        var ir = r * 0.20

        var id = 0
        var intersections = []
        var intersectionsMap = {}

        grid.forEach((hx) => {
            if (hx.resource == "-" || hx.resource == "s" || hx.port) {
                return
            }
        
            //var neighbors = []
            
            for (var i = 0; i < 6; i++) {
                var x1 = hx.x + r * Math.cos((a * i) + 11)
                var y1 = hx.y + r * Math.sin((a * i) + 11)

                var ins = null
                var nodes = []
                var key = []

                for (var hi = 0; hi < grid.length; hi++) {
                    var h1 = grid[hi]
                    
                    if (circlesIntersect(x1, y1, ir, h1.x, h1.y, h1.r)) {
                        nodes.push(h1)
                        key.push(h1.index)
                    }
                }
                
                var key = key.join("#")

                if (intersectionsMap[key] == undefined) {
                    ins = { index: id, x: x1, y: y1, r: ir, key: key, hasPort: false }
                    intersectionsMap[key] = ins
                    intersections.push(ins)
                    id++
                } else {
                    ins = intersectionsMap[key]
                }                
                hx.intersections.push(ins)

                ins.nodes = nodes
                //neighbors.concat(nodes)                    
            }

            // var neighborMap = {}
            // neighbors.forEach((n) =>{
            //     if (n.index == h.index) {
            //         return
            //     }

            //     if (neighbors[n.index] == undefined) {
            //         h.neighbors.push(n)
            //         neighborMap[n.index] = true
            //     }
            // })            
        })

        function getNextSide(s) {
            if (s == 5) {
                return 0
            }
            return s + 1
        }

        grid.forEach((hx) => {
            if (!hx.port) {
                return
            }
                            
            var portIntersections = []
            intersections.forEach((ins) => {
                if (circlesIntersect(hx.x, hx.y, hx.r, ins.x, ins.y, ins.r)) {
                    portIntersections.push(ins)
                }
            })

            portIntersections.forEach((ins) => {
                var px = hx.x + r * Math.cos((a * hx.direction) + 11)
                var py = hx.y + r * Math.sin((a * hx.direction) + 11)                            
                if (circlesIntersect(ins.x, ins.y, ins.r, px, py, ins.r)) {
                    ins.hasPort = true
                    ins.port = hx
                }

                var ns = getNextSide(hx.direction)
                px = hx.x + r * Math.cos((a * ns) + 11)
                py = hx.y + r * Math.sin((a * ns) + 11)                
                if (circlesIntersect(ins.x, ins.y, ins.r, px, py, ins.r)) {
                    ins.hasPort = true
                    ins.port = hx
                }
            })                            
        })

        intersections.forEach(ins1 => {
            var neighbors = []
            for (var j = 0; j < intersections.length; j++) {
                var ins2 = intersections[j]
                if (ins1.index != ins2.index) {
                    if (circlesIntersect(ins1.x, ins1.y, r, ins2.x, ins2.y, ins2.r)) {
                        neighbors.push(ins2)
                    }
                }
            }
            ins1.neighbors = neighbors
        })

        return { nodes: grid, intersections: intersections }
    }

    o = parse(map)
    grid = makeGrid(o.h, o.w, o.tiles)
    
    return makeIntersections(grid)
}

m = `
-,-,sw2,s,s?3,s,-
-,s,m,p,t,sg3,-
-,s?1,f,h,p,h,s
s,f,t,d,t,m,sl4
-,s?1,t,m,f,p,s
-,s,h,f,p,s?5,-
-,-,so0,s,sb5,s,-`

// grid = Grid(m)

// console.log(grid)

// var canvas = document.getElementById('board')
// var ctx = canvas.getContext('2d');
// var inses = {}
// for (var idx in grid) {
//     var hex = grid[idx]

//     if (hex.resource == "-" || hex.resource == "s") {
//         continue
//     }

//     ctx.beginPath()
//     ctx.arc(hex.x, hex.y, hex.r, 0, 2 * Math.PI, false)
//     ctx.stroke()
//     ctx.font = '15px serif';
//     ctx.fillText(hex.index, hex.x, hex.y);


//     for (var j = 0; j < hex.intersections.length; j++) {
//         var ins = hex.intersections[j]
//         if (inses[ins.index] == undefined) {
//             ctx.beginPath()
//             ctx.arc(ins.x, ins.y, ins.r, 0, 2 * Math.PI, false)
//             ctx.stroke()
//             inses[ins.index] = 1
//         }
//     }

// }