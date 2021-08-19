package main

/*
#include <windows.h>
int KeyDown(int key) {
    // 数据兼容：因为 GetKeyState() 不接受小写字母
    if (key > 96 && key < 123)  key -= 32;

    // 获取按下的键的状态，返回 0 则表示没按，其他情况表示按了
    return (GetKeyState(key) < 0) ? 1 : 0;
}
*/
import "C"
import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sync"
	"time"
)

/*
1.创建一个全局格子
2.定义多个图形方块
3.方块移动方法
4.方块旋转方法
5.数组边界判断
6.数据到底部进行重新分配方块    for 循环
*/
type APP struct {
	s           sync.Mutex
	Location    [][]int
	GlobalSlice [30][15]int
	factory     map[int]Diamonds
}

var TetrisAPP = &APP{
	factory: map[int]Diamonds{
		0: &FourSquares{},
		1: &Strip{},
	},
}

//func TestTetris(t *testing.T) {
func main() {
	//格子初始化
	for i, v := range TetrisAPP.GlobalSlice {
		for ii := range v {
			switch {
			case i == 0 || i == len(TetrisAPP.GlobalSlice)-1:
				TetrisAPP.GlobalSlice[i][ii] = 1

			case ii == 0 || ii == len(v)-1:
				TetrisAPP.GlobalSlice[i][ii] = 1

			default:
				TetrisAPP.GlobalSlice[i][ii] = 0
			}
		}
	}
	//程序运行
	TetrisAPP.Run() //程序运行 输出格子
}
func (app *APP) Run() {
	flag := false
	app.Println()                   //输出格子
	if err := clear(); err != nil { //清屏函数
		panic(err)
	}

	//循环监听键盘操作
	wg := sync.WaitGroup{}
	//ch := make(chan interface{})
	wg.Add(1)
	/*
	   	1.执行初始化，和默认方块下落的时候，对定位数组与全局数组进行加锁，锁释放之后，通道释放
	   2.键盘操控键位捕获释放的通道，进行相应的操作
	*/
	//ch <- 1
	go func(wg *sync.WaitGroup) {
		for {
			if !flag {
				app.Location = TetrisAPP.factory[rand.Intn(len(TetrisAPP.factory))].Produce() //[1,1,1,1]
				if err := app.initialize(); err != nil {
					return
				}
			}
			if err := clear(); err != nil { //清屏函数
				panic(err)
			}
			app.Println() //输出格子
			app.s.Lock()
			flag = app.Decline("s")
			app.s.Unlock()
			time.Sleep(time.Second * 1)
			if err := clear(); err != nil { //清屏函数
				panic(err)
			}
		}
	}(&wg)
	defer wg.Wait() //最终等待子线程运行完毕
	//wg.Add(1)

	for {

		switch 1 {
		// 监听按键 ，按下时返回 1，没按下时返回 0
		case int(C.KeyDown('w')): //键盘上  旋转
			if err := clear(); err != nil { //清屏函数
				panic(err)
			}
			app.Decline("w")
			app.Println() //输出格子
		case int(C.KeyDown('a')): //左侧移动
			if err := clear(); err != nil { //清屏函数
				panic(err)
			}
			app.Decline("a")
			app.Println() //输出格子
		case int(C.KeyDown('s')): //加速下落
			if err := clear(); err != nil { //清屏函数
				panic(err)
			}
			flag = app.Decline("s")
			app.Println() //输出格子
		case int(C.KeyDown('d')): //右侧移动
			if err := clear(); err != nil { //清屏函数
				panic(err)
			}
			app.Decline("d")
			app.Println() //输出格子
		}

		// 延迟 500毫秒 后再接收，避免接收太快

		time.Sleep(time.Millisecond * 100)
	}

}

func (app *APP) Println() {
	for _, v := range TetrisAPP.GlobalSlice {
		var vvv string
		for _, vv := range v {
			if vv != 0 {
				vvv += "■"
				continue
			}
			vvv += "  "
		}
		fmt.Println(vvv)
	}
}

func clear() error {
	cmd := exec.Command("cmd", "/c", "clear")
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return errors.New(fmt.Sprintf("程序出现未知错误:%v", err))
	}
	return nil
}

//格子默认下落   修改  //左右下  ads键位
func (app *APP) Decline(param string) bool {
	var (
		l    = len(app.Location)
		flag = true
	)
	switch param {
	case "w":

	case "a": //左移动，判断location最左边的全部数据移动一位是否会产生碰撞 右移动相反
		fmt.Println(app.Location)
		for i, v := range app.Location {
			var ff = true
			for _, vv := range v {
				if vv != 0 && ff {
					ff = false
					if app.GlobalSlice[i][vv-1] != 0 {
						flag = false
					}
				}
			}
		}
		fmt.Println(flag)
		if flag {
			for i, v := range app.Location {
				for ii, vv := range v {
					if vv != 0 {
						app.GlobalSlice[i][vv] = 0
						app.Location[i][ii] = vv - 1

						app.GlobalSlice[i][vv-1] = 1
					}
				}
			}
		}

	case "d": //右移动
		fmt.Println(app.Location)
		for i, v := range app.Location {
			var ff = true
			for ii := len(v) - 1; ii >= 0; ii-- {
				if v[ii] != 0 && ff {
					ff = false
					if app.GlobalSlice[i][v[ii]+1] != 0 {
						flag = false
					}
				}
			}
		}

		fmt.Println(flag)
		if flag {

			//
			for i, v := range app.Location {
				for ii, vv := range v {
					if vv != 0 {
						app.Location[i][ii] = vv + 1
						app.GlobalSlice[i][vv] = 0
						app.GlobalSlice[i][vv+1] = 1
					}
				}
			}
		}

	default:
		for _, v := range app.Location[l-1] {
			//如果当前数组的值不为0，判断当前全局数组++之后的位置有没有值，没有值的话，对当前location数组进行++操作，循环完毕之后删除原全局数组的location位置，并将新的location赋值到全局数组
			if app.GlobalSlice[l][v] != 0 { //如果下落过程中遇到实体方块，停止下落
				//fmt.Println(app.GlobalSlice[i+1], v, "调试第 %v 行", i+2, "第 ", v, "个参数")
				flag = false
			}
		}

		if flag {
			var f = false
			for i, v := range app.Location {
				for _, vv := range v {
					//如果当前数组的值不为0，判断当前全局数组++之后的位置有没有值，循环完毕之后删除原全局数组的location位置
					//如果下落过程中遇到实体方块，停止下落
					if vv != 0 {
						f = true
						app.GlobalSlice[i][vv] = 0 //清空原来的位置
					}
				}
			}

			for i, v := range app.Location {
				for _, vv := range v {
					//如果当前数组的值不为0，判断当前全局数组++之后的位置有没有值，循环完毕之后删除原全局数组的location位置
					//如果下落过程中遇到实体方块，停止下落
					if vv != 0 {
						app.GlobalSlice[i+1][vv] = 1 //全局数组方块下落
					}
				}
			}

			if f {
				app.Location = append([][]int{{0, 0, 0, 0}}, app.Location...) //当前方块变更
			}
		}
	}

	return flag
}

func (app *APP) initialize() (err error) { //第一次生成
	app.s.Lock()
	defer app.s.Unlock()
	if is := app.FallJudgment(); !is {
		return errors.New("程序结束:新生成方块位置不足！")
	}
	return
}

func (app *APP) FallJudgment() bool { //落位判断 如果当前位置可以落位就给模块数组赋值为全局数组的位置
	var (
		lens     = len(app.GlobalSlice[0])
		location int
	)
	for {
		l := rand.Intn(lens)
		//随机位置出现当前的模块，排除数组第一位和最后一位，以及模块数据在出现位置存放不下的情况
		if l != 0 && l < lens-len(app.Location[0]) {
			location = l
			break
		}
	}
	//给当前方块赋值，并且删除原全局数组的方块，进行新一次的赋值
	for i, v := range app.Location {
		for ii, vv := range v {
			if vv != 0 {
				app.Location[i][ii] = location + ii
				app.GlobalSlice[i+1][ii+location] = 1
			}
		}
	}
	app.Location = append([][]int{{0, 0, 0, 0}}, app.Location...)
	return true
}

type Diamonds interface {
	Produce() [][]int    //生产格子   多种类型的方块
	Rotate() interface{} //格子旋转
}

type FourSquares struct {
	Value [][]int
}

func (f *FourSquares) Produce() [][]int {
	f.Value = [][]int{{1, 1}, {1, 1}}
	return f.Value
}

func (f *FourSquares) Rotate() interface{} {
	return f.Value
}

type Strip struct {
	Value [][]int
}

func (s *Strip) Produce() [][]int {
	s.Value = [][]int{{1, 1, 1, 1}}
	return s.Value
}

func (s *Strip) Rotate() interface{} {
	return s.Value
}
