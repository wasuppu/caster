package main

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"sort"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	WINDOW_WIDTH     = 1280
	WINDOW_HEIGHT    = 720
	SCREEN_WIDTH     = 640
	SCREEN_HEIGHT    = 480
	MAP_WIDTH        = 24
	MAP_HEIGHT       = 24
	TEX_WIDTH        = 64
	TEX_HEIGHT       = 64
	FLOOR_HORIZONTAL = true
	NUM_SPRITES      = 19
)

var worldMap = [MAP_WIDTH][MAP_HEIGHT]int{
	{8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 4, 4, 6, 4, 4, 6, 4, 6, 4, 4, 4, 6, 4},
	{8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4},
	{8, 0, 3, 3, 0, 0, 0, 0, 0, 8, 8, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 6},
	{8, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 6},
	{8, 0, 3, 3, 0, 0, 0, 0, 0, 8, 8, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4},
	{8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 4, 0, 0, 0, 0, 0, 6, 6, 6, 0, 6, 4, 6},
	{8, 8, 8, 8, 0, 8, 8, 8, 8, 8, 8, 4, 4, 4, 4, 4, 4, 6, 0, 0, 0, 0, 0, 6},
	{7, 7, 7, 7, 0, 7, 7, 7, 7, 0, 8, 0, 8, 0, 8, 0, 8, 4, 0, 4, 0, 6, 0, 6},
	{7, 7, 0, 0, 0, 0, 0, 0, 7, 8, 0, 8, 0, 8, 0, 8, 8, 6, 0, 0, 0, 0, 0, 6},
	{7, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 6, 0, 0, 0, 0, 0, 4},
	{7, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 6, 0, 6, 0, 6, 0, 6},
	{7, 7, 0, 0, 0, 0, 0, 0, 7, 8, 0, 8, 0, 8, 0, 8, 8, 6, 4, 6, 0, 6, 6, 6},
	{7, 7, 7, 7, 0, 7, 7, 7, 7, 8, 8, 4, 0, 6, 8, 4, 8, 3, 3, 3, 0, 3, 3, 3},
	{2, 2, 2, 2, 0, 2, 2, 2, 2, 4, 6, 4, 0, 0, 6, 0, 6, 3, 0, 0, 0, 0, 0, 3},
	{2, 2, 0, 0, 0, 0, 0, 2, 2, 4, 0, 0, 0, 0, 0, 0, 4, 3, 0, 0, 0, 0, 0, 3},
	{2, 0, 0, 0, 0, 0, 0, 0, 2, 4, 0, 0, 0, 0, 0, 0, 4, 3, 0, 0, 0, 0, 0, 3},
	{1, 0, 0, 0, 0, 0, 0, 0, 1, 4, 4, 4, 4, 4, 6, 0, 6, 3, 3, 0, 0, 0, 3, 3},
	{2, 0, 0, 0, 0, 0, 0, 0, 2, 2, 2, 1, 2, 2, 2, 6, 6, 0, 0, 5, 0, 5, 0, 5},
	{2, 2, 0, 0, 0, 0, 0, 2, 2, 2, 0, 0, 0, 2, 2, 0, 5, 0, 5, 0, 0, 0, 5, 5},
	{2, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 2, 5, 0, 5, 0, 5, 0, 5, 0, 5},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5},
	{2, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 2, 5, 0, 5, 0, 5, 0, 5, 0, 5},
	{2, 2, 0, 0, 0, 0, 0, 2, 2, 2, 0, 0, 0, 2, 2, 0, 5, 0, 5, 0, 0, 0, 5, 5},
	{2, 2, 2, 2, 1, 2, 2, 2, 2, 2, 2, 1, 2, 2, 2, 5, 5, 5, 5, 5, 5, 5, 5, 5},
}

type Sprite struct {
	x, y    float64
	texture int
}

var sprite = [NUM_SPRITES]Sprite{
	{20.5, 11.5, 10}, //green light in front of playerstart
	//green lights in every room
	{18.5, 4.5, 10},
	{10.0, 4.5, 10},
	{10.0, 12.5, 10},
	{3.5, 6.5, 10},
	{3.5, 20.5, 10},
	{3.5, 14.5, 10},
	{14.5, 20.5, 10},

	//row of pillars in front of wall: fisheye test
	{18.5, 10.5, 9},
	{18.5, 11.5, 9},
	{18.5, 12.5, 9},

	//some barrels around the map
	{21.5, 1.5, 8},
	{15.5, 1.5, 8},
	{16.0, 1.8, 8},
	{16.2, 1.2, 8},
	{3.5, 2.5, 8},
	{9.5, 15.5, 8},
	{10.0, 15.1, 8},
	{10.5, 15.8, 8},
}

type Graphic struct {
	window          *sdl.Window
	surface         *sdl.Surface
	pos, dir, plane vec2f
	time, oldTime   uint64 // time of current frame and time of previous frame
	frameTime       float64
	buffer          [SCREEN_HEIGHT][SCREEN_WIDTH]uint32 // y-coordinate first because it works per scanline
	texture         [11][TEX_WIDTH * TEX_HEIGHT]uint32
	zbuffer         [SCREEN_WIDTH]float64
	spriteOrder     [NUM_SPRITES]int
	spriteDistance  [NUM_SPRITES]float64
}

func (g *Graphic) Init() (err error) {
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		return fmt.Errorf("failed to initializing SDL: %s", err)
	}

	g.window, err = sdl.CreateWindow("raycast", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		WINDOW_WIDTH, WINDOW_HEIGHT, sdl.WINDOW_BORDERLESS|sdl.WINDOW_ALLOW_HIGHDPI)
	if err != nil {
		return fmt.Errorf("failed to create SDL window: %s", err)
	}

	g.surface, err = g.window.GetSurface()
	if err != nil {
		return fmt.Errorf("failed to get surface: %s", err)
	}
	g.surface.FillRect(nil, sdl.MapRGB(g.surface.Format, 0x10, 0x10, 0x10))
	g.window.UpdateSurface()

	g.pos = vec2f{22, 11.5}  // x and y start position
	g.dir = vec2f{-1, 0}     // initial direction vector
	g.plane = vec2f{0, 0.66} // the 2d raycaster version of camera plane

	// load some textures
	loadTexture(g.texture[0][:], "pics/eagle.png")
	loadTexture(g.texture[1][:], "pics/redbrick.png")
	loadTexture(g.texture[2][:], "pics/purplestone.png")
	loadTexture(g.texture[3][:], "pics/greystone.png")
	loadTexture(g.texture[4][:], "pics/bluestone.png")
	loadTexture(g.texture[5][:], "pics/mossy.png")
	loadTexture(g.texture[6][:], "pics/wood.png")
	loadTexture(g.texture[7][:], "pics/colorstone.png")

	// load some sprite textures
	loadTexture(g.texture[8][:], "pics/barrel.png")
	loadTexture(g.texture[9][:], "pics/pillar.png")
	loadTexture(g.texture[10][:], "pics/greenlight.png")

	return nil
}

func (g *Graphic) CleanUp() {
	g.window.Destroy()
	sdl.Quit()
}

func (g *Graphic) Done() bool {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			fmt.Println("\nQuitting...")
			return true
		case *sdl.KeyboardEvent:
			if t.Keysym.Sym == sdl.K_ESCAPE {
				return true
			}
		}
	}
	return false
}

func (g *Graphic) HandleKeys() {
	// speed modifiers
	moveSpeed := g.frameTime * 5.0 // the constant value is in squares/second
	rotSpeed := g.frameTime * 3.0  // the constant value is in radians/second

	sdl.PumpEvents()
	keystate := sdl.GetKeyboardState()

	// move forward if no wall in front of you
	if keystate[sdl.SCANCODE_UP] != 0 {
		if worldMap[int(g.pos.x+g.dir.x*moveSpeed)][int(g.pos.y)] == 0 {
			g.pos.x += g.dir.x * moveSpeed
		}

		if worldMap[int(g.pos.x)][int(g.pos.y+g.dir.y*moveSpeed)] == 0 {
			g.pos.y += g.dir.y * moveSpeed
		}
	}

	// move backwards if no wall behind you
	if keystate[sdl.SCANCODE_DOWN] != 0 {
		if worldMap[int(g.pos.x-g.dir.x*moveSpeed)][int(g.pos.y)] == 0 {
			g.pos.x -= g.dir.x * moveSpeed
		}

		if worldMap[int(g.pos.x)][int(g.pos.y-g.dir.y*moveSpeed)] == 0 {
			g.pos.y -= g.dir.y * moveSpeed
		}
	}

	// rotate to the right
	if keystate[sdl.SCANCODE_RIGHT] != 0 {
		// both camera direction and camera plane must be rotated
		oldDir := g.dir
		g.dir.x = g.dir.x*math.Cos(-rotSpeed) - g.dir.y*math.Sin(-rotSpeed)
		g.dir.y = oldDir.x*math.Sin(-rotSpeed) + g.dir.y*math.Cos(-rotSpeed)
		oldPlane := g.plane
		g.plane.x = g.plane.x*math.Cos(-rotSpeed) - g.plane.y*math.Sin(-rotSpeed)
		g.plane.y = oldPlane.x*math.Sin(-rotSpeed) + g.plane.y*math.Cos(-rotSpeed)
	}

	// rotate to the left
	if keystate[sdl.SCANCODE_LEFT] != 0 {
		// both camera direction and camera plane must be rotated
		oldDir := g.dir
		g.dir.x = g.dir.x*math.Cos(rotSpeed) - g.dir.y*math.Sin(rotSpeed)
		g.dir.y = oldDir.x*math.Sin(rotSpeed) + g.dir.y*math.Cos(rotSpeed)
		oldPlane := g.plane
		g.plane.x = g.plane.x*math.Cos(rotSpeed) - g.plane.y*math.Sin(rotSpeed)
		g.plane.y = oldPlane.x*math.Sin(rotSpeed) + g.plane.y*math.Cos(rotSpeed)
	}
}

func (g *Graphic) SetFrameTime() {
	// timing for input and FPS counter
	g.oldTime = g.time
	g.time = sdl.GetTicks64()
	g.frameTime = float64(g.time-g.oldTime) / 1000.0 // frametime is the time this frame has taken, in seconds
}

func (g *Graphic) Calculate() {
	// FLOOR CASTING
	if FLOOR_HORIZONTAL {
		for y := SCREEN_HEIGHT/2 + 1; y < SCREEN_HEIGHT; y++ {
			// rayDir for leftmost ray (x = 0) and rightmost ray (x = w)
			rayDir0 := g.dir.sub(g.plane)
			rayDir1 := g.dir.add(g.plane)

			// Current y position compared to the center of the screen (the horizon)
			p := y - SCREEN_HEIGHT/2

			// Vertical position of the camera.
			posZ := 0.5 * SCREEN_HEIGHT

			// Horizontal distance from the camera to the floor for the current row.
			// 0.5 is the z position exactly in the middle between floor and ceiling.
			rowDistance := posZ / float64(p)

			// calculate the real world step vector we have to add for each x (parallel to camera plane)
			// adding step by step avoids multiplications with a weight in the inner loop
			floorStep := rayDir1.sub(rayDir0).muln(rowDistance).divn(SCREEN_WIDTH)

			// real world coordinates of the leftmost column. This will be updated as we step to the right.
			floor := g.pos.add(rayDir0.muln(rowDistance))

			for x := range SCREEN_WIDTH {
				// the cell coord is simply got from the integer parts of floorX and floorY
				cell := vec2i{int(floor.x), int(floor.y)}

				// get the texture coordinate from the fractional part
				tx := int(TEX_WIDTH*(floor.x-float64(cell.x))) & (TEX_WIDTH - 1)
				ty := int(TEX_HEIGHT*(floor.y-float64(cell.y))) & (TEX_HEIGHT - 1)

				floor = floor.add(floorStep)

				// choose texture and draw the pixel
				// floorTexture := 3
				floorTexture := 4
				checkerBoardPattern := int(cell.x+cell.y) & 1
				if checkerBoardPattern == 0 {
					floorTexture = 3
				}

				// floor
				c := g.texture[floorTexture][TEX_WIDTH*ty+tx]
				c = (c >> 1) & 8355711 // make a bit darker
				g.buffer[y][x] = c

				//ceiling (symmetrical, at screenHeight - y - 1 instead of y)
				ceilingTexture := 6
				c = g.texture[ceilingTexture][TEX_WIDTH*ty+tx]
				c = (c >> 1) & 8355711 // make a bit darker
				g.buffer[SCREEN_HEIGHT-y-1][x] = c
			}
		}
	}

	// WALL CASTING
	for x := range SCREEN_WIDTH {

		// calculate ray position and direction
		cameraX := 2*float64(x)/SCREEN_WIDTH - 1 // x-coordinate in camera space

		rayDir := g.dir.add(g.plane.muln(cameraX))

		// which box of the map we're in
		ipos := vec2i{int(g.pos.x), int(g.pos.y)}

		// length of ray from one x or y-side to next x or y-side
		deltaDist := vec2f{}
		if rayDir.x == 0 {
			deltaDist.x = 1e30
		} else {
			deltaDist.x = math.Abs(1 / rayDir.x)
		}

		if rayDir.y == 0 {
			deltaDist.y = 1e30
		} else {
			deltaDist.y = math.Abs(1 / rayDir.y)
		}

		// length of ray from current position to next x or y-side
		sideDist := vec2f{}
		// what direction to step in x or y-direction (either +1 or -1)
		step := vec2f{}

		// calculate step and initial sideDist
		if rayDir.x < 0 {
			step.x = -1
			sideDist.x = (g.pos.x - float64(ipos.x)) * deltaDist.x
		} else {
			step.x = 1
			sideDist.x = (float64(ipos.x) + 1.0 - g.pos.x) * deltaDist.x
		}

		if rayDir.y < 0 {
			step.y = -1
			sideDist.y = (g.pos.y - float64(ipos.y)) * deltaDist.y
		} else {
			step.y = 1
			sideDist.y = (float64(ipos.y) + 1.0 - g.pos.y) * deltaDist.y
		}

		hit := false  // was there a wall hit?
		side := false // was a NS or a EW wall hit?

		// perform DDA
		for !hit {
			// jump to next map square, either in x-direction, or in y-direction
			if sideDist.x < sideDist.y {
				sideDist.x += deltaDist.x
				ipos.x = int(float64(ipos.x) + step.x)
				side = false
			} else {
				sideDist.y += deltaDist.y
				ipos.y = int(float64(ipos.y) + step.y)
				side = true
			}

			// Check if ray has hit a wall
			if worldMap[ipos.x][ipos.y] > 0 {
				hit = true
			}
		}

		perpWallDist := 0.0
		// Calculate distance of perpendicular ray (Euclidean distance would give fisheye effect!)
		if !side {
			perpWallDist = sideDist.x - deltaDist.x
		} else {
			perpWallDist = sideDist.y - deltaDist.y
		}

		// Calculate height of line to draw on screen
		lineHeight := int(SCREEN_HEIGHT / perpWallDist)

		// calculate lowest and highest pixel to fill in current stripe
		drawStart := -lineHeight/2 + SCREEN_HEIGHT/2
		if drawStart < 0 {
			drawStart = 0
		}
		drawEnd := lineHeight/2 + SCREEN_HEIGHT/2
		if drawEnd >= SCREEN_HEIGHT {
			drawEnd = SCREEN_HEIGHT - 1
		}

		// texturing calculations
		textureNum := worldMap[ipos.x][ipos.y] - 1 // 1 subtracted from it so that texture 0 can be used!

		// calculate value of wallX
		wallX := 0.0 // where exactly the wall was hit
		if !side {
			wallX = g.pos.y + perpWallDist*rayDir.y
		} else {
			wallX = g.pos.x + perpWallDist*rayDir.x
		}
		wallX -= math.Floor(wallX)

		// x coordinate on the texture
		texX := int(wallX * TEX_WIDTH)
		if !side && rayDir.x > 0 {
			texX = TEX_WIDTH - texX - 1
		}
		if side && rayDir.y < 0 {
			texX = TEX_WIDTH - texX - 1
		}

		// How much to increase the texture coordinate per screen pixel
		texStep := 1.0 * TEX_HEIGHT / float64(lineHeight)
		// Starting texture coordinate
		texPos := (float64(drawStart) - SCREEN_HEIGHT/2 + float64(lineHeight)/2) * texStep
		for y := drawStart; y < drawEnd; y++ {
			// Cast the texture coordinate to integer, and mask with (texHeight - 1) in case of overflow
			textureY := int(texPos) & (TEX_HEIGHT - 1)
			texPos += texStep
			c := g.texture[textureNum][TEX_HEIGHT*textureY+texX]
			// make color darker for y-sides: R, G and B byte each divided through two with a "shift" and an "and"
			if side {
				c = (c >> 1) & 8355711
			}
			g.buffer[y][x] = c
		}

		// SET THE ZBUFFER FOR THE SPRITE CASTING
		g.zbuffer[x] = perpWallDist // perpendicular distance is used

		// FLOOR CASTING (vertical version, directly after drawing the vertical wall stripe for the current x)
		if !FLOOR_HORIZONTAL {
			floorWall := vec2f{} // x, y position of the floor texel at the bottom of the wall

			// 4 different wall directions possible
			if !side && rayDir.x > 0 {
				floorWall.x = float64(ipos.x)
				floorWall.y = float64(ipos.y) + wallX
			} else if !side && rayDir.x < 0 {
				floorWall.x = float64(ipos.x) + 1.0
				floorWall.y = float64(ipos.y) + wallX
			} else if side && rayDir.y > 0 {
				floorWall.x = float64(ipos.x) + wallX
				floorWall.y = float64(ipos.y)
			} else {
				floorWall.x = float64(ipos.x) + wallX
				floorWall.y = float64(ipos.y) + 1.0
			}

			distWall := perpWallDist
			distPlayer := 0.0

			if drawEnd < 0 {
				drawEnd = SCREEN_HEIGHT // becomes < 0 when the integer overflows
			}

			// draw the floor from drawEnd to the bottom of the screen
			for y := drawEnd + 1; y < SCREEN_HEIGHT; y++ {
				currentDist := SCREEN_HEIGHT / (2.0*float64(y) - SCREEN_HEIGHT)

				weight := (currentDist - distPlayer) / (distWall - distPlayer)

				currentFloor := floorWall.muln(weight).add(g.pos.muln(1.0 - weight))

				floorTexX := int(currentFloor.x*TEX_WIDTH) & (TEX_WIDTH - 1)
				floorTexY := int(currentFloor.y*TEX_HEIGHT) & (TEX_HEIGHT - 1)

				floorTexture := 4
				checkerBoardPattern := (int(currentFloor.x) + int(currentFloor.y)) & 1
				if checkerBoardPattern == 0 {
					floorTexture = 3
				}

				// floor
				g.buffer[y][x] = (g.texture[floorTexture][TEX_WIDTH*floorTexY+floorTexX] >> 1) & 8355711
				// ceiling (symmetrical!)
				g.buffer[SCREEN_HEIGHT-y][x] = g.texture[6][TEX_WIDTH*floorTexY+floorTexX]
			}
		}
	}

	// SPRITE CASTING
	// sort sprites from far to close
	for i := range NUM_SPRITES {
		g.spriteOrder[i] = i
		g.spriteDistance[i] = (g.pos.x-sprite[i].x)*(g.pos.x-sprite[i].x) + (g.pos.y-sprite[i].y)*(g.pos.y-sprite[i].y)
	}
	sortSprites(g.spriteOrder[:], g.spriteDistance[:], NUM_SPRITES)

	// after sorting the sprites, do the projection and draw them
	for i := range NUM_SPRITES {
		// translate sprite position to relative to camera
		spriteX := sprite[g.spriteOrder[i]].x - g.pos.x
		spriteY := sprite[g.spriteOrder[i]].y - g.pos.y

		//transform sprite with the inverse camera matrix
		// [ planeX   dirX ] -1                                       [ dirY      -dirX ]
		// [               ]       =  1/(planeX*dirY-dirX*planeY) *   [                 ]
		// [ planeY   dirY ]                                          [ -planeY  planeX ]
		invDet := 1.0 / (g.plane.x*g.dir.y - g.dir.x*g.plane.y) // required for correct matrix multiplication

		transformX := invDet * (g.dir.y*spriteX - g.dir.x*spriteY)
		transformY := invDet * (-g.plane.y*spriteX + g.plane.x*spriteY) //this is actually the depth inside the screen, that what Z is in 3D

		spriteScreenX := int((SCREEN_WIDTH / 2) * (1 + transformX/transformY))

		const (
			uDiv  = 1
			vDiv  = 1
			vMove = 0.0
		)

		vMoveScreen := int(vMove / transformY)

		// calculate height of the sprite on screen
		spriteHeight := int(math.Abs(float64(int(SCREEN_HEIGHT / transformY)))) // using 'transformY' instead of the real distance prevents fisheye
		// calculate lowest and highest pixel to fill in current stripe
		drawStartY := max(-spriteHeight/2+SCREEN_HEIGHT/2+vMoveScreen, 0)
		drawEndY := spriteHeight/2 + SCREEN_HEIGHT/2 + vMoveScreen
		if drawEndY >= SCREEN_HEIGHT {
			drawEndY = SCREEN_HEIGHT - 1
		}

		// calculate width of the sprite
		spriteWidth := int(math.Abs(float64(int(SCREEN_HEIGHT / transformY))))
		drawStartX := max(-spriteWidth/2+spriteScreenX, 0)
		drawEndX := spriteWidth/2 + spriteScreenX
		if drawEndX >= SCREEN_WIDTH {
			drawEndX = SCREEN_WIDTH - 1
		}

		// loop through every vertical stripe of the sprite on screen
		for stripe := drawStartX; stripe < drawEndX; stripe++ {
			texX := int(256*(stripe-(-spriteWidth/2+spriteScreenX))*TEX_WIDTH/spriteWidth) / 256
			// the conditions in the if are:
			// 1) it's in front of camera plane so you don't see things behind you
			// 2) it's on the screen (left)
			// 3) it's on the screen (right)
			// 4) ZBuffer, with perpendicular distance
			if transformY > 0 && stripe > 0 && stripe < SCREEN_WIDTH && transformY < g.zbuffer[stripe] {
				for y := drawStartY; y < drawEndY; y++ { // for every pixel of the current stripe
					d := (y-vMoveScreen)*256 - SCREEN_HEIGHT*128 + spriteHeight*128 // 256 and 128 factors to avoid floats
					texY := ((d * TEX_HEIGHT) / spriteHeight) / 256
					c := g.texture[sprite[g.spriteOrder[i]].texture][TEX_WIDTH*texY+texX] // get current color from the texture

					if (c & 0x00FFFFFF) != 0 { // paint pixel if it isn't black, black is the invisible color
						g.buffer[y][stripe] = c
					}
				}
			}
		}

	}

}

func (g *Graphic) Render() {
	offsetX := (WINDOW_WIDTH - SCREEN_WIDTH) / 2
	offsetY := (WINDOW_HEIGHT - SCREEN_HEIGHT) / 2
	g.surface.Lock()
	defer g.surface.Unlock()
	pixels := g.surface.Pixels()
	for y := range int(SCREEN_HEIGHT) {
		for x := range int(SCREEN_WIDTH) {
			index := (y+offsetY)*int(g.surface.Pitch) + (x+offsetX)*4
			binary.LittleEndian.PutUint32(pixels[index:index+4], g.buffer[y][x])
		}
	}

	for y := range SCREEN_HEIGHT {
		for x := range SCREEN_WIDTH {
			g.buffer[y][x] = 0
		}
	}

	g.window.UpdateSurface()
}

func loadTexture(texture []uint32, filename string) {
	img := openImg(filename)
	pixels := readColorFromImg(img)
	for x := range uint32(TEX_WIDTH) {
		for y := range uint32(TEX_HEIGHT) {
			pixel := pixels[TEX_WIDTH*y+x]
			c, _ := color.RGBAModel.Convert(pixel).(color.RGBA)
			texture[TEX_WIDTH*y+x] = rgbaToUint32(c)
		}
	}
}

func openImg(filename string) image.Image {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		panic(err)
	}

	return img
}

func readColorFromImg(img image.Image) []color.Color {
	var pixels []color.Color

	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixels = append(pixels, img.At(x, y))
		}
	}

	return pixels
}

func rgbaToUint32(rgba color.RGBA) uint32 {
	return uint32(rgba.R)<<16 | uint32(rgba.G)<<8 | uint32(rgba.B) | uint32(rgba.A)<<24
}

type pair[T any, E any] struct {
	first  T
	second E
}

func sortSprites(order []int, dist []float64, amount int) {
	sprites := make([]pair[float64, int], amount)

	for i := range amount {
		sprites[i].first = dist[i]
		sprites[i].second = order[i]
	}

	sort.Slice(sprites, func(i, j int) bool {
		return sprites[i].first < sprites[j].first
	})

	for i := range amount {
		dist[i] = sprites[amount-i-1].first
		order[i] = sprites[amount-i-1].second
	}
}
