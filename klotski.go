package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

var prime [6]int = [6]int{2, 3, 5, 7, 11, 13}
var Path map[string]int = map[string]int{}
var Best []string = []string{}

const STOP = 5

func main() {
	start := time.Now().UnixNano() / int64(time.Millisecond)
	//生成布局，去除全横与全竖布局。
	//GenerateLayout("12335.txt")
	//RemoveSymmetric("12335.txt", "10415.txt")
	//SetSpaces("10415.txt", "44094.txt")
	//RemoveSymmetric("44094.txt", "44014.txt")
	//RemoveNoneSolution("44014.txt", "32095.txt")
	//去除5步内的变换形式 。722秒
	//RemoveTransformation("32095.txt", "8107.txt")
	//根据横将的数量排序，需要手动在sublime中排序，并去除首字母供后续使用。
	//SortLayout("8107.txt", "sorted8107.txt")
	//GeneratHtmlLink("last.txt", "link.txt")
	//BatchSolve("last.txt", "result.txt")
	stop := time.Now().UnixNano() / int64(time.Millisecond)
	fmt.Println("[Time: ", (stop - start), " ms]")
}

func BatchSolve(input, output string) {
	fr, err := os.OpenFile(input, os.O_RDONLY, 0644)
	defer fr.Close()
	if err != nil {
		println(err)
	}
	fw, err := os.OpenFile(output, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	defer fw.Close()
	if err != nil {
		println(err)
	}
	reader := bufio.NewReader(fr)
	id := 1
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			println(err)
		}
		line = strings.TrimSpace(line)
		fmt.Println("==[", id, "]==")
		Solve(line)
		fw.WriteString(line)
		fw.WriteString("#")
		fw.WriteString(generateSolution())
		fw.WriteString("\n")
		id++
	}
}

func Solve(game string) {
	//fmt.Println(game)
	Path = map[string]int{}
	extend(game)
	solution, _ := updatePath(game)
	trace(solution)
	//fmt.Println(generateSolution())
}

func HasSolution(game string) bool {
	Path = map[string]int{}
	extend(game)
	for k, _ := range Path {
		if k[17] == k[18] && k[17] == 52 {
			return true
		}
	}
	return false
}

func trace(layout string) {
	if layout == "" {
		panic("No solution.")
	}
	Best = []string{}
	Best = append(Best, layout)
	step := Path[layout]
	for step > 0 {
		items := moveForward(layout)
		for _, v := range items {
			if Path[v] == step-1 {
				Best = append(Best, v)
				layout = v
				step--
				break
			}
		}
	}
}

func generateSolution() string {
	solution := "["
	for i := len(Best) - 1; i > 0; i-- {
		j := i - 1
		rc, d := showMotion(Best[i], Best[j])
		solution += fmt.Sprintf("[%d],%d,", rc, d)
	}
	solution = solution[0:len(solution)-1] + "]"
	return solution
}

func updatePath(game string) (string, error) {
	Path[game] = 0
	for stop, step := false, 0; !stop; step++ {
		stop = true
		for c, s := range Path {
			if s == step {
				stop = false
				nxs := moveForward(c)
				for _, nx := range nxs {
					if Path[nx] > s+1 {
						Path[nx] = s + 1
						if nx[17] == 52 && nx[18] == 52 {
							fmt.Println("Solution:", s+1)
							return nx, nil
						}
					}
				}
			}
		}
	}
	return "", errors.New("No solution.")
}

func showMotion(a, b string) (int, int) {
	c := []byte(a)
	d := []byte(b)
	diff := []int{}
	for i, _ := range c {
		if c[i] != d[i] {
			diff = append(diff, i)
		}
	}
	if len(diff) == 2 {
		if diff[0]+1 == diff[1] {
			if c[diff[0]] == 48 {
				return diff[1], 0
			} else {
				return diff[0], 2
			}
		} else if diff[0]+4 == diff[1] {
			if c[diff[0]] == 48 {
				return diff[1], 1
			} else {
				return diff[0], 3
			}
		} else if diff[0]+2 == diff[1] {
			if c[diff[0]] == 48 {
				return diff[0] + 1, 0
			} else {
				return diff[0], 2
			}
		} else if diff[0]+8 == diff[1] {
			if c[diff[0]] == 48 {
				return diff[0] + 4, 1
			} else {
				return diff[0], 3
			}
		}
	} else {
		if diff[0]+1 == diff[1] && diff[0]%4 < 3 {
			if c[diff[0]] == c[diff[1]] {
				if c[diff[0]] == 48 {
					return diff[0] + 4, 1
				} else {
					return diff[0], 3
				}
			} else {
				if c[diff[0]] == 48 {
					return diff[0] + 1, 0
				} else {
					return diff[0], 2
				}
			}
		} else if diff[0]+2 == diff[1] {
			if c[diff[0]] == 48 {
				return diff[0] + 1, 0
			} else {
				return diff[0], 2
			}
		} else if diff[0]+8 == diff[2] {
			if c[diff[0]] == 48 {
				return diff[0] + 4, 1
			} else {
				return diff[0], 3
			}
		}
	}
	return -1, -1
}

func extend(child string) {
	if _, exist := Path[child]; !exist {
		Path[child] = 100000
		if child[17] == 52 && child[18] == 52 {
			return
		} else {
			x := moveForward(child)
			for _, v := range x {
				extend(v)
			}
		}
	}
}

func moveForward(child string) []string {
	result := []string{}
	for pos, c := range child {
		if c == 48 {
			if pos%4 != 0 {
				if child[pos-1] == 49 {
					layout := []byte(child)
					layout[pos] = 49
					layout[pos-1] = 48
					result = append(result, string(layout))
				}
				if child[pos-1] == 50 {
					if pos-1 == findPoint(child, pos-1) && child[pos+4] == 48 {
						layout := []byte(child)
						layout[pos] = 50
						layout[pos+4] = 50
						layout[pos-1] = 48
						layout[pos+3] = 48
						result = append(result, string(layout))
					}
				}
				if child[pos-1] == 51 {
					layout := []byte(child)
					layout[pos] = 51
					layout[pos-2] = 48
					result = append(result, string(layout))
				}
				if child[pos-1] == 52 {
					if pos-2 == findPoint(child, pos-1) && child[pos+4] == 48 {
						layout := []byte(child)
						layout[pos] = 52
						layout[pos+4] = 52
						layout[pos-2] = 48
						layout[pos+2] = 48
						result = append(result, string(layout))
					}
				}
			}
			if pos%4 != 3 && pos < 19 {
				if child[pos+1] == 49 {
					layout := []byte(child)
					layout[pos] = 49
					layout[pos+1] = 48
					result = append(result, string(layout))
				}
				if child[pos+1] == 50 {
					if pos+1 == findPoint(child, pos+1) && child[pos+4] == 48 {
						layout := []byte(child)
						layout[pos] = 50
						layout[pos+4] = 50
						layout[pos+1] = 48
						layout[pos+5] = 48
						result = append(result, string(layout))
					}
				}
				if child[pos+1] == 51 {
					layout := []byte(child)
					layout[pos] = 51
					layout[pos+2] = 48
					result = append(result, string(layout))
				}
				if child[pos+1] == 52 {
					if pos+1 == findPoint(child, pos+1) && child[pos+4] == 48 {
						layout := []byte(child)
						layout[pos] = 52
						layout[pos+4] = 52
						layout[pos+2] = 48
						layout[pos+6] = 48
						result = append(result, string(layout))
					}
				}
			}
			if pos >= 4 {
				if child[pos-4] == 49 {
					layout := []byte(child)
					layout[pos] = 49
					layout[pos-4] = 48
					result = append(result, string(layout))
				}
				if child[pos-4] == 50 {
					layout := []byte(child)
					layout[pos] = 50
					layout[pos-8] = 48
					result = append(result, string(layout))
				}
				if child[pos-4] == 51 {
					if pos-4 == findPoint(child, pos-4) && child[pos+1] == 48 {
						layout := []byte(child)
						layout[pos] = 51
						layout[pos+1] = 51
						layout[pos-4] = 48
						layout[pos-3] = 48
						result = append(result, string(layout))
					}
				}
				if child[pos-4] == 52 {
					if pos-8 == findPoint(child, pos-4) && child[pos+1] == 48 {
						layout := []byte(child)
						layout[pos] = 52
						layout[pos+1] = 52
						layout[pos-8] = 48
						layout[pos-7] = 48
						result = append(result, string(layout))
					}
				}
			}
			if pos <= 15 {
				if child[pos+4] == 49 {
					layout := []byte(child)
					layout[pos] = 49
					layout[pos+4] = 48
					result = append(result, string(layout))
				}
				if child[pos+4] == 50 {
					layout := []byte(child)
					layout[pos] = 50
					layout[pos+8] = 48
					result = append(result, string(layout))
				}
				if child[pos+4] == 51 {
					if pos+4 == findPoint(child, pos+4) && child[pos+1] == 48 {
						layout := []byte(child)
						layout[pos] = 51
						layout[pos+1] = 51
						layout[pos+4] = 48
						layout[pos+5] = 48
						result = append(result, string(layout))
					}
				}
				if child[pos+4] == 52 {
					if pos+4 == findPoint(child, pos+4) && child[pos+1] == 48 {
						layout := []byte(child)
						layout[pos] = 52
						layout[pos+1] = 52
						layout[pos+8] = 48
						layout[pos+9] = 48
						result = append(result, string(layout))
					}
				}
			}

		}
	}
	return result
}

func findPoint(temp string, pos int) int {
	bType := temp[pos]
	if bType == 49 {
		return pos
	}
	if bType == 50 {
		if pos < 4 {
			return pos
		}
		if pos > 15 {
			return pos - 4
		}
		if 50 != temp[pos+4] {
			return pos - 4
		}
		if 50 != temp[pos-4] {
			return pos
		}
		if pos > 3 && pos < 8 {
			return pos - 4
		}
		if pos > 11 && pos < 16 {
			return pos
		}
		if temp[pos-8] == 50 {
			return pos
		} else {
			return pos - 4
		}
	} else if bType == 51 {
		if pos%4 == 0 {
			return pos
		}
		if pos%4 == 3 {
			return pos - 1
		}
		if 51 != temp[pos-1] {
			return pos
		}
		if 51 != temp[pos+1] {
			return pos - 1
		}
		return pos - (pos % 2)
	} else if bType == 52 {
		if pos%4 > 0 && 52 == temp[pos-1] {
			if pos > 3 && 52 == temp[pos-5] {
				return pos - 5
			} else {
				return pos - 1
			}
		} else {
			if pos > 3 && 52 == temp[pos-4] {
				return pos - 4
			} else {
				return pos
			}
		}
	}
	return -1
}

// 递归调用横将或竖将填充
func pave(tmp [20]int, one, vtwo, htwo int, fw *os.File) {
	if one == 0 && one == vtwo && vtwo == htwo {
		var layout bytes.Buffer
		for _, v := range tmp {
			layout.WriteByte(byte(v + 48))
		}
		layout.WriteString("\n")
		fw.WriteString(layout.String())
	}
	pos := 0
	for pos = 0; pos < 20; pos++ {
		if tmp[pos] == 0 {
			break
		}
	}
	if one > 0 {
		temp := tmp
		temp[pos] = 1
		pave(temp, one-1, vtwo, htwo, fw)
	}
	if vtwo > 0 && pos < 16 && tmp[pos+4] == 0 {
		temp := tmp
		temp[pos] = 2
		temp[pos+4] = 2
		pave(temp, one, vtwo-1, htwo, fw)
	}
	if htwo > 0 && pos%4 < 3 && tmp[pos+1] == 0 {
		temp := tmp
		temp[pos] = 3
		temp[pos+1] = 3
		pave(temp, one, vtwo, htwo-1, fw)
	}
}

// 生成布局，考虑到对称情况，只保留曹操在左侧的布局。
func GenerateLayout(output string) {
	fw, err := os.OpenFile(output, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return
	}
	initMap := [7][20]int{
		{4, 4, 0, 0, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 4, 4, 0, 0, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 4, 4, 0, 0, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 4, 4, 0, 0, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 4, 4, 0, 0, 4, 4, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 4, 0, 0, 4, 4, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 4, 0, 0, 4, 4, 0, 0}}
	for i := 0; i < 7; i++ {
		//pave(initMap[i], 6, 5, 0, fw)
		pave(initMap[i], 6, 4, 1, fw)
		pave(initMap[i], 6, 3, 2, fw)
		pave(initMap[i], 6, 2, 3, fw)
		pave(initMap[i], 6, 1, 4, fw)
		//pave(initMap[i], 6, 0, 5, fw)
	}
}

// 去除对称局。曹操在中间的布局中仍有对称布局。
func RemoveSymmetric(input, output string) {
	fr, err := os.OpenFile(input, os.O_RDONLY, 0644)
	defer fr.Close()
	if err != nil {
		println(err)
	}
	fw, err := os.OpenFile(output, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	defer fw.Close()
	if err != nil {
		println(err)
	}
	var list []string
	reader := bufio.NewReader(fr)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			println(err)
		}
		list = append(list, strings.TrimSpace(line))

	}
	fmt.Println("Symmetrical layouts: (keep the left only.)")
	cnt := 0
	for entry, original := range list {
		target := []byte(original)
		for r := 0; r < 5; r++ {
			tempa := target[r*4]
			tempb := target[r*4+1]
			target[r*4] = target[r*4+3]
			target[r*4+1] = target[r*4+2]
			target[r*4+2] = tempb
			target[r*4+3] = tempa
		}
		symmetic := string(target)
		find := false
		for x := entry + 1; x < len(list); x++ {
			if symmetic == list[x] {
				fmt.Println(original, "<==>", list[x])
				find = true
				cnt++
				break
			}
		}
		if !find {
			fw.WriteString(original + "\n")
		}
	}
	fmt.Println("[Total:", cnt, "]")
}

// 设定空格
func SetSpaces(input, output string) {
	fr, err := os.OpenFile(input, os.O_RDONLY, 0644)
	defer fr.Close()
	if err != nil {
		println(err)
	}
	fw, err := os.OpenFile(output, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	defer fw.Close()
	if err != nil {
		println(err)
	}
	reader := bufio.NewReader(fr)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			println(err)
		}
		line = strings.TrimSpace(line)
		pos := [6]int{}
		idx := 0
		for i, v := range line {
			if v == 49 {
				pos[idx] = i
				idx++
			}
		}
		chooseSpaces(line, pos, fw)
	}
}

// 选择空格位置，排除相连空格产生的不同组合
func chooseSpaces(line string, pos [6]int, fw *os.File) {
	g1, g2, g3 := 1, 1, 1
	for i := 0; i < 5; i++ {
		if pos[i]%4 < 3 && pos[i+1] == pos[i]+1 {
			if g1 == 1 {
				g1 *= prime[i] * prime[i+1]
			} else if g1%prime[i] == 0 || g1%prime[i+1] == 0 {
				if g1%prime[i] == 0 {
					g1 *= prime[i+1]
				} else {
					g1 *= prime[i]
				}
				if g2%prime[i] == 0 || g2%prime[i+1] == 0 {
					g1 *= g2
					g2 = 1
				}
			} else if g2 == 1 {
				g2 *= prime[i] * prime[i+1]
			} else if g2%prime[i] == 0 || g2%prime[i+1] == 0 {
				if g2%prime[i] == 0 {
					g2 *= prime[i+1]
				} else {
					g2 *= prime[i]
				}
				if g3%prime[i] == 0 || g3%prime[i+1] == 0 {
					g2 *= g3
					g3 = 1
				}

			} else {
				g3 *= prime[i] * prime[i+1]
			}
		}
		for j := i + 1; j < 6; j++ {
			if pos[j] > pos[i]+4 {
				break
			}
			if pos[j] == pos[i]+4 {
				if g1 == 1 {
					g1 *= prime[i] * prime[j]
				} else if g1%prime[i] == 0 || g1%prime[j] == 0 {
					if g1%prime[i] == 0 {
						g1 *= prime[j]
					} else {
						g1 *= prime[i]
					}
					if g2%prime[i] == 0 || g2%prime[j] == 0 {
						g1 *= g2
						g2 = 1
					}
				} else if g2 == 1 {
					g2 *= prime[i] * prime[j]
				} else if g2%prime[i] == 0 || g2%prime[j] == 0 {
					if g2%prime[i] == 0 {
						g2 *= prime[j]
					} else {
						g2 *= prime[i]
					}
					if g3%prime[i] == 0 || g3%prime[j] == 0 {
						g2 *= g3
						g3 = 1
					}
				} else {
					g3 *= prime[i] * prime[j]
				}
			}
		}
	}
	var group1, group2, group3, group4 []int
	for x := 0; x < 6; x++ {
		if g1%prime[x] == 0 {
			group1 = append(group1, x)
		} else if g2%prime[x] == 0 {
			group2 = append(group2, x)
		} else if g3%prime[x] == 0 {
			group3 = append(group3, x)
		} else {
			group4 = append(group4, x)
		}
	}
	size := len(group4)
	if size >= 2 {
		for m := 0; m < size-1; m++ {
			for n := m + 1; n < size; n++ {
				layout := []byte(line)
				layout[pos[group4[n]]] = 48
				layout[pos[group4[m]]] = 48
				fw.WriteString(string(layout) + "\n")
			}
		}
	}
	var gArr [][]int = [][]int{group1, group2, group3, group4}
	for s := 0; s < 3; s++ {
		if len(gArr[s]) >= 2 {
			layout := []byte(line)
			layout[pos[gArr[s][0]]] = 48
			layout[pos[gArr[s][1]]] = 48
			fw.WriteString(string(layout) + "\n")
			for t := s + 1; t < 3; t++ {
				if len(gArr[t]) >= 2 {
					layout1 := []byte(line)
					layout1[pos[gArr[s][0]]] = 48
					layout1[pos[gArr[t][0]]] = 48
					fw.WriteString(string(layout1) + "\n")
				}
			}
			for m := 0; m < size; m++ {
				layout2 := []byte(line)
				layout2[pos[gArr[s][0]]] = 48
				layout2[pos[group4[m]]] = 48
				fw.WriteString(string(layout2) + "\n")
			}
		}
	}
}

// 布局排序，根据横将从0~5
func SortLayout(input, output string) {
	fr, err := os.OpenFile(input, os.O_RDONLY, 0644)
	defer fr.Close()
	if err != nil {
		println(err)
	}
	fw, err := os.OpenFile(output, os.O_CREATE|os.O_RDWR, 0644)
	defer fw.Close()
	if err != nil {
		println(err)
	}
	reader := bufio.NewReader(fr)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			println(err)
		}
		line = strings.TrimSpace(line)
		var three int = 0
		for i, v := range line {
			if i >= 20 {
				break
			}
			if v == 51 {
				three++
			}
		}
		fw.WriteString(fmt.Sprintf("%x%s\n", three/2+10, line))
	}
}

// 去除无解布局
func RemoveNoneSolution(input, output string) {
	var solved_list map[string]int = make(map[string]int, 10000)
	fr, err := os.OpenFile(input, os.O_RDONLY, 0644)
	defer fr.Close()
	if err != nil {
		println(err)
	}
	fw, err := os.OpenFile(output, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	defer fw.Close()
	if err != nil {
		println(err)
	}
	reader := bufio.NewReader(fr)
	c := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			println(err)
		}
		line = strings.TrimSpace(line)
		if _, exist := solved_list[line]; exist {
			fw.WriteString(line)
			fw.WriteString("\n")
			continue
		}
		if HasSolution(line) {
			for k, _ := range Path {
				solved_list[k] = 0
			}
			fw.WriteString(line)
			fw.WriteString("\n")
		} else {
			c++
		}
	}
	fmt.Println("[Remove for none solution: ", c, "]")
}

// 去除变换
func RemoveTransformation(input, output string) {
	var layout_list []string
	fr, err := os.OpenFile(input, os.O_RDONLY, 0644)
	defer fr.Close()
	if err != nil {
		println(err)
	}
	fw, err := os.OpenFile(output, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	defer fw.Close()
	if err != nil {
		println(err)
	}
	reader := bufio.NewReader(fr)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			println(err)
		}
		line = strings.TrimSpace(line)
		layout_list = append(layout_list, line)
	}
	c := 0
	for i, v := range layout_list {
		if v == "0" {
			continue
		}
		remove_list := []string{v}
		step := 0
		for step < STOP {
			step++
			for _, key := range remove_list {
				remove_list = append(remove_list, moveForward(key)...)
			}
		}
		for _, key := range remove_list {
			for j := i + 1; j < len(layout_list); j++ {
				if layout_list[j] == "0" {
					continue
				}
				if key == layout_list[j] {
					c++
					layout_list[j] = "0"
				}
			}
		}
	}
	for _, v := range layout_list {
		if v == "0" {
			continue
		}
		fw.WriteString(v)
		fw.WriteString("\n")
	}
	fmt.Println("[Remove: ", c, " layouts]")
}

// 生成超链接，发布使用。
func GeneratHtmlLink(input, output string) {
	fr, err := os.OpenFile(input, os.O_RDONLY, 0644)
	defer fr.Close()
	if err != nil {
		println(err)
	}
	fw, err := os.OpenFile(output, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	defer fw.Close()
	if err != nil {
		println(err)
	}
	reader := bufio.NewReader(fr)
	id := 0
	for {
		id++
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			println(err)
		}
		line = strings.TrimSpace(line)
		fw.WriteString(fmt.Sprintf("%d#:<a href=\"customize.html?%s\">%s</a><br>\n", id, line, line))
	}
}
