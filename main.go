package main

import (
	"image/color"
	_ "image/png"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	WIDTH  = 480
	HEIGHT = 600
)

const (
	PLAY = iota
	GAME_OVER
)

type player struct {
	img  *ebiten.Image
	x, y float64
}

func (p *player) update(shoot *bool) {
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) && p.x < WIDTH-99 {
		p.x += 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) && p.x > 0 {
		p.x -= 5
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		*shoot = true
	}
}

type bullet struct {
	img  *ebiten.Image
	x, y float64
}

type meteor struct {
	img    *ebiten.Image
	x, y   float64
	active bool
}

type Game struct {
	player
	bullet
	meteor
	backgroundImg *ebiten.Image
	shoot         bool
	pressedKeys   []ebiten.Key
	meteors       []meteor
	score         int
	f             font.Face
	scene         int
}


func (g *Game) Update() error {
	g.player.update(&g.shoot)

	if g.shoot {
		g.bullet.y -= 15
	}

	if g.bullet.y < 0 {
		g.shoot = false
	}

	if !g.shoot {
		g.bullet.x, g.bullet.y = g.player.x+43, g.player.y

	}

	for len(g.meteors) < 3 {
		rand.Seed(time.Now().UnixNano())
		g.meteors = append(g.meteors, meteor{
			g.meteor.img,
			float64(rand.Intn((WIDTH - 28 + 1) + 28)),
			-20,
			true,
		})
	}

	for i := 0; i < len(g.meteors); i++ {
		g.meteors[i].y += 3
		if g.meteors[i].y > HEIGHT || !g.meteors[i].active {
			g.meteors = append(g.meteors[:i], g.meteors[i+1:]...)
		}
	}

	// check collision peluru dengan meteor
	bw, bh := g.bullet.img.Size()
	for i := 0; i < len(g.meteors); i++ {
		mw, mh := g.meteors[i].img.Size()
		if g.bullet.x+float64(bw) >= g.meteors[i].x &&
			g.bullet.x <= g.meteors[i].x+float64(mw) &&
			g.bullet.y+float64(bh) >= g.meteors[i].y &&
			g.bullet.y <= g.meteors[i].y+float64(mh) && g.shoot {
			g.shoot = false
			g.meteors[i].active = false
			g.score += 1
		}

	}

	// check collision player dengan meteor
	pw, ph := g.player.img.Size()
	for i := 0; i < len(g.meteors); i++ {
		mw, mh := g.meteors[i].img.Size()
		if g.player.x+float64(pw) >= g.meteors[i].x &&
			g.player.x <= g.meteors[i].x+float64(mw) &&
			g.player.y+float64(ph) >= g.meteors[i].y &&
			g.player.y <= g.meteors[i].y+float64(mh) {
			g.scene = 1
		}

	}

	if ebiten.IsKeyPressed(ebiten.KeyR) && g.scene == 1 {
		g.scene = 0
		g.score = 0
		g.meteors = nil
		g.player.x = WIDTH / 2
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// draw background
	screen.DrawImage(g.backgroundImg, &ebiten.DrawImageOptions{})
	switch g.scene {
	case PLAY:

		// draw bullet
		if g.shoot {
			bp := &ebiten.DrawImageOptions{}
			bp.GeoM.Translate(g.bullet.x, g.bullet.y)
			screen.DrawImage(g.bullet.img, bp)
		}

		// draw meteor
		for i := 0; i < len(g.meteors); i++ {
			if g.meteors[i].active {
				mp := &ebiten.DrawImageOptions{}
				mp.GeoM.Translate(g.meteors[i].x, g.meteors[i].y)
				screen.DrawImage(g.meteors[i].img, mp)
			}
		}

		// draw player
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(g.player.x, g.player.y)
		screen.DrawImage(g.player.img, op)
		// draw score text
		score_text := "Score: " + strconv.Itoa(g.score)
		text.Draw(screen, score_text, g.f, 10, 40, color.White)
	case GAME_OVER:
		// draw game over text
		text.Draw(screen, "Game Over", g.f, WIDTH/2-70, HEIGHT/4, color.White)
		text.Draw(screen, "Tekan \"R\" untuk restart", g.f, 80, HEIGHT/2+30, color.White)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WIDTH, HEIGHT
}

func main() {
	ebiten.SetWindowSize(WIDTH, HEIGHT)
	ebiten.SetWindowTitle("Space ship")

	// load font
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	// load assets
	img, _, _ := ebitenutil.NewImageFromFile("./assets/playerShip1_orange.png")
	bg, _, _ := ebitenutil.NewImageFromFile("./assets/bg.png")
	bullet_img, _, _ := ebitenutil.NewImageFromFile("./assets/laserBlue05.png")
	meteor_img, _, _ := ebitenutil.NewImageFromFile("./assets/meteorBrown_small1.png")

	game := &Game{}
	// player
	game.player.img = img
	game.player.x = WIDTH / 2
	game.player.y = HEIGHT - 100

	// bullet
	game.bullet.img = bullet_img

	// meteor
	game.meteor.img = meteor_img

	// etc
	game.shoot = false
	game.backgroundImg = bg
	game.f, _ = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    30,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
