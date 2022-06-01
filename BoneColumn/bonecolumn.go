package Bonecolumn

import (
	"github.com/faiface/pixel"
	"math/rand"
)

//柱子信息
type Bone struct {
	RectUp     pixel.Rect
	RectDown   pixel.Rect
	Vel        pixel.Vec
	Sp         *pixel.Sprite
	OriPosUp   pixel.Rect
	OriPosDown pixel.Rect
	MoveSpeed  float64
}

//画出柱子
func (b *Bone) Drawit(w pixel.Target, bSp *pixel.Sprite, bRect pixel.Rect, picframe pixel.Rect) {
	bSp.Draw(w, pixel.IM.ScaledXY(pixel.ZV,
		pixel.V(bRect.W()/picframe.W(), bRect.H()/picframe.H())).Moved(bRect.Center()))
}

//移动柱子
func (b *Bone) Moveit(win pixel.Target, winB pixel.Rect, frame pixel.Rect, dt float64, bs []*Bone, th bool) {
	n := rand.Intn(3)
	//移动
	if b.RectUp.Min.X >= winB.Min.X {
		b.Drawit(win, b.Sp, b.RectUp, frame)
		b.Drawit(win, b.Sp, b.RectDown, frame)
		//向左移动
		b.RectUp = b.RectUp.Moved(pixel.V(-200*dt, 0))
		b.RectDown = b.RectDown.Moved(pixel.V(-200*dt, 0))
	} else {
		//重置位置
		b.RectUp = bs[n].RectUp //从障碍物矩阵数组中随机赋值
		b.RectDown = bs[n].RectDown
	}

	//碰撞判断
	if th {
		b.RectUp = bs[n].RectUp
		b.RectDown = bs[n].RectDown

		b.RectUp.Min.X = b.OriPosUp.Min.X
		b.RectUp.Max.X = b.OriPosUp.Max.X
		b.RectDown.Min.X = b.OriPosDown.Min.X
		b.RectDown.Max.X = b.OriPosDown.Max.X
	}
}
