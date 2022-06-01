package Home

import (
	Bonecolumn "Flysans/BoneColumn"
	"Flysans/LoadPic"
	"Flysans/Sans"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"math/rand"
	"time"
)

type Bg struct {
	Pic   pixel.Picture
	Frame pixel.Rect
	Rect  pixel.Rect
	Sp    *pixel.Sprite
}

func Run() {
	//载入图片
	bcpic, err := Loadpic.LoadPic("img/wall.png")
	bgpic, err := Loadpic.LoadPic("img/bg.png")
	sansPic, err := Loadpic.LoadPic("img/sans.png")
	if err != nil {
		panic(err)
	}

	//窗口设置
	cfg := pixelgl.WindowConfig{
		Title:     "Flysans",
		Bounds:    pixel.R(0, 0, 1020, 768),
		Icon:      []pixel.Picture{sansPic},
		Resizable: true,
		VSync:     true,
	}
	win, err := pixelgl.NewWindow(cfg)
	win.Smooth()
	if err != nil {
		panic(err)
	}

	//角色对象
	sans := Sans.SanPhys{
		Rect:     pixel.R(-50, -50, 50, 50),
		Gravity:  -1200,
		FlySpeed: 260,
		Frame:    sansPic.Bounds(),
	}
	sans.Rect = sans.Rect.Moved(win.Bounds().Center())
	sans.OriPos = sans.Rect //保存原始位置

	//背景
	bg := []*Bg{
		{Pic: bgpic,
			Frame: bgpic.Bounds(),
			Rect:  win.Bounds(),
			Sp:    pixel.NewSprite(bgpic, bgpic.Bounds()),
		},
		{Pic: bgpic,
			Frame: bgpic.Bounds(),
			Rect:  win.Bounds(),
			Sp:    pixel.NewSprite(bgpic, bgpic.Bounds()),
		},
	}
	bg[1].Rect = bg[1].Rect.Moved(pixel.V(bg[1].Rect.W(), 0))

	//柱子位置
	frame := bcpic.Bounds() //原图大小
	bs := []*Bonecolumn.Bone{
		{
			RectUp: pixel.R(win.Bounds().Max.X, win.Bounds().Max.Y-144,
				win.Bounds().Max.X+170, win.Bounds().Max.Y),
			RectDown: pixel.R(win.Bounds().Max.X, win.Bounds().Min.Y,
				win.Bounds().Max.X+170, win.Bounds().Min.Y+432),
		},
		{
			RectUp: pixel.R(win.Bounds().Max.X, win.Bounds().Max.Y-432,
				win.Bounds().Max.X+170, win.Bounds().Max.Y),

			RectDown: pixel.R(win.Bounds().Max.X, win.Bounds().Min.Y,
				win.Bounds().Max.X+170, win.Bounds().Min.Y+144),
		},
		{
			RectUp: pixel.R(win.Bounds().Max.X, win.Bounds().Max.Y-288,
				win.Bounds().Max.X+170, win.Bounds().Max.Y),
			RectDown: pixel.R(win.Bounds().Max.X, win.Bounds().Min.Y,
				win.Bounds().Max.X+170, win.Bounds().Min.Y+288),
		},
	}

	x := rand.Intn(3)
	//柱子矩阵array
	bcs := []*Bonecolumn.Bone{
		{
			RectUp:   bs[x].RectUp,
			RectDown: bs[x].RectDown,
		},
		{
			RectUp:   bs[(x+1)/2].RectUp,
			RectDown: bs[(x+1)/2].RectDown,
		},
	}
	//第二个出现的柱子与第一个相隔170*3
	bcs[1].RectUp = bcs[1].RectUp.Moved(pixel.V(170*3, 0))
	bcs[1].RectDown = bcs[1].RectDown.Moved(pixel.V(170*3, 0))

	//记录原始位置
	bcs[0].OriPosUp = bcs[0].RectUp
	bcs[0].OriPosDown = bcs[0].RectDown
	bcs[1].OriPosUp = bcs[1].RectUp
	bcs[1].OriPosDown = bcs[1].RectDown

	////柱子sprite
	bcs[0].Sp = pixel.NewSprite(bcpic, frame)
	bcs[1].Sp = pixel.NewSprite(bcpic, frame)

	last := time.Now()
	for !win.Closed() {
		//每帧时间
		dt := time.Since(last).Seconds()
		last = time.Now()

		//背景
		win.Clear(colornames.Black)
		bg[0].BgMove(win, dt)
		bg[1].BgMove(win, dt)

		//移动柱子矩阵
		bcs[0].Moveit(win, win.Bounds(), frame, dt, bs, sans.Touch)
		bcs[1].Moveit(win, win.Bounds(), frame, dt, bs, sans.Touch)

		//角色移动
		ctrl := pixel.ZV
		if win.Pressed(pixelgl.MouseButtonLeft) {
			ctrl.Y = 1
		}
		sans.Update(dt, ctrl, bcs)
		sans.Drawit(win, sansPic)
		win.Update()
	}

}

func (bg *Bg) BgMove(w pixel.Target, dt float64) {
	bg.Sp.Draw(w, pixel.IM.ScaledXY(pixel.ZV,
		pixel.V(bg.Rect.W()/bg.Frame.W(),
			bg.Rect.H()/bg.Frame.H())).
		Moved(bg.Rect.Center()),
	)

	if bg.Rect.Max.X > 0 {
		bg.Rect = bg.Rect.Moved(pixel.V(-200*dt, 0))
	} else {
		bg.Rect = bg.Rect.Moved(pixel.V(bg.Rect.W()*2, 0))
	}

}
