### FlySans

### 一、对象：

#### 障碍物-骨柱：

1. 基本属性

   1. 宽128，高{（288，288），（144，432），（432，144）}；
   2. 上下空隙：192，左右间隙：256；
   3. 建立x,y两个矩阵存储骨柱，让x与y的间隙等于256；建立数组存储两个柱子;
   4. 当首次新建柱子时，将其位置放在前一个柱子+256个像素值的位置；

2. 移动：

   1. 逻辑：存储在数组中的障碍物循环出现；当超出窗口范围时，重置其位置；

   2. 移动和重置：

      1. 移动距离为200*dt，每次移动则重新绘图；

      2. 位置重置，即将矩阵位置变为初始位置；需要建立存储所有不同长度障碍物的矩阵数组，且重置时随机赋值；

         
   
      ```go
      n := rand.Intn(3)//随机数
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
      ```

      

   3. 碰撞判断：

      1. 当碰撞参数为true即游戏失败，重置所有矩阵位置；与移动位置重置不同，因为首次建立障碍物两个柱子的位置相隔256，所以需要存储首次生成时的X轴位置，以方便赋值；
   
         
   
      ```go
      //碰撞判断
      if th {
      	b.RectUp = bs[n].RectUp
      	b.RectDown = bs[n].RectDown
      
      	b.RectUp.Min.X = b.OriPosUp.Min.X
      	b.RectUp.Max.X = b.OriPosUp.Max.X
   	b.RectDown.Min.X = b.OriPosDown.Min.X
      	b.RectDown.Max.X = b.OriPosDown.Max.X
      }
      ```
   
      
   
   完整代码：[bonecolumn.go](D:\go game\selfMake-flappySans\BoneColumn\bonecolumn.go)

#### 角色：

1. 属性：位置矩阵，移动向量，重力，移动速度，sprite，图片矩阵，碰撞判断，原始位置；

   ```go
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
   ```

   

2. 绘制：

   ```go
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
   ```

   

3. 重力和移动：

   ```go
   	//重力：
   	S.Vel.Y += S.Gravity * dt
   
   	//矩阵移动
   	S.Rect = S.Rect.Moved(S.Vel.Scaled(dt))
   
   	//升降控制
   	if ctrl.Y > 0 {
   		S.Vel.Y = S.FlySpeed
   	}
   ```

   

4. 碰撞判定：

   1. 碰撞判定范围：当角色X的最大值大于障碍物的最小值或角色X的最小值小于障碍物的最大值，且在柱子的Y轴范围内，则判断为碰撞；

      ![image-20220328150417251](C:\Users\25777\AppData\Roaming\Typora\typora-user-images\image-20220328150417251.png)

   2. 注意：判断边界尽量模糊化，即可适当增减判断的范围大小，以免过于精确造成游玩不适；

   ```go
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
   ```

   完整代码：[sans.go](D:\go game\selfMake-flappySans\Sans\san.go)



#### 窗口及背景：

1. 背景属性：图片，大小，位置，sprite

   ```go
   type Bg struct {
   	Pic   pixel.Picture
   	Frame pixel.Rect
   	Rect  pixel.Rect
   	Sp    *pixel.Sprite
   }
   ```

2. 背景移动逻辑：建立大小为2的数组存储背景元素，第二个背景位置在第一个之后；每帧向左移动固定距离，当背景矩阵X轴的最大值小于窗口最小值，则向右移动窗口（背景）宽度*2的距离，即回到初始位置；

   ```go
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
   		
   ```

   ```go
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
   
   ```

3. 窗口基本设置及对象初始化：

   ​	完整代码：[home.go](D:\go game\selfMake-flappySans\Home\home.go)

#### Tip&Exp：

- 载入图片的大小与实际显示的对象的大小差距不宜过大；
- 对象的属性设置，最好和窗口的大小对应，以便在窗口调整时，对象也随之调整；
- 一些对象的基本数据，如大小，最好提前确定 并单独用文件存放以便修改；

#### 项目位置：[Flysans](D:\go game\selfMake-flappySans)

