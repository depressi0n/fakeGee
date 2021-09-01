package gee

import "strings"

type node struct {
	pattern  string // 待匹配路由
	part     string // 路由中的一部分
	children []*node
	isWild   bool // 精确匹配，当part含有:或*是为true
}

func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}
func (n *node) insert(pattern string,parts []string,height int)  {
	if len(parts)==height{
		n.pattern=pattern
		return
	}
	part:=parts[height]
	child:=n.matchChild(part)
	if child==nil{
		child=&node{
			part:     part,
			isWild:   part[0]==':' || part[0]=='*',
		}
		n.children=append(n.children,child)
	}
	child.insert(pattern,parts,height+1)
}
func (n *node) search(parts []string,height int) *node {
	// 如果遇到了*的匹配则也说明已经是最后了
	if len(parts)==height || strings.HasPrefix(n.part,"*"){
		// 只有最后一层节点才会设置pattern，即如果匹配到的是中间层节点，可以利用pattern是否为空来判断匹配是否成功
		if n.pattern==""{
			return nil
		}
		return n
	}
	part:=parts[height]
	children:=n.matchChildren(part)
	for _,child:=range children{
		result:=child.search(parts,height+1)
		if result!=nil{
			return result
		}
	}
	return nil
}