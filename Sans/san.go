package Sans

import (
	Bonecolumn "Flysans/BoneColumn"
	"fmt"
	"github.com/faiface/pixel"
)

type SanPhys struct {
	Rect     pixel.Rect
	Vel      pixel.Vec
	Gravity  float64
	FlySpeed float64
	Sp       *pixel.Sprite
	Frame    pixel.Rect
	Touch    bool

	OriPos pixel.Rect
}

func (S *SanPhys) Update(dt float64, ctrl pixel.Vec, Bcol []*Bonecolumn.Bone) {
	//重力：
	S.Vel.Y += S.Gravity * dt

	//矩阵移动
	S.Rect = S.Rect.Moved(S.Vel.Scaled(dt))

	//升降控制
	if ctrl.Y > 0 {
		S.Vel.Y = S.FlySpeed
	}

	//碰撞判定
	S.Touch = false
	for _, b := range Bcol {
		if S.Rect.Min.X > b.RectUp.Min.X-S.Rect.W()+10 && S.Rect.Max.X < b.RectUp.Max.X+S.Rect.W()-10 {
			if S.Rect.Max.Y-10 > b.RectUp.Min.Y {
				S.Touch = true
				fmt.Println("Lose")
			}
			if S.Rect.Min.Y+10 < b.RectDown.Max.Y {
				S.Touch = true
				fmt.Println("Lose")
			}
		}
	}

}

func (S *SanPhys) Drawit(w pixel.Target, pic pixel.Picture) {
	//碰撞判定
	if S.Touch {
		S.Rect = S.OriPos
		S.Vel.Y = 0
	}
	S.Sp = pixel.NewSprite(pic, S.Frame)
	S.Sp.Draw(w, pixel.IM.ScaledXY(pixel.ZV, pixel.V(
		S.Rect.W()/S.Frame.W(),
		S.Rect.H()/S.Frame.H())).Moved(S.Rect.Center()))
}
