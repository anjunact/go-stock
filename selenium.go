package main

import (
	"fmt"
	"github.com/tebeka/selenium"
	//"strings"
)

var code = `
package main
import "fmt"

func main() {
	fmt.Println("Hello WebDriver!\n")
}
`

// Errors are ignored for brevity.

func main() {
	// FireFox driver without specific version
	// *** Add gecko driver here if necessary (see notes above.) ***
	caps := selenium.Capabilities{"browserName": "chrome"}
	wd, err := selenium.NewRemote(caps, "")
	if err != nil {
		panic(err)
	}
	defer wd.Quit()

	// Get simple playground interface
	wd.Get("http://finance.sina.com.cn/data/#stock")

	// Enter code in textarea
	elems, _ := wd.FindElements(selenium.ByCSSSelector, "table#block_1.tr")
	for elem := range elems{
		fmt.Printf("%v",elem)
		//text,_ := elem.Find.Text()
		//ar := strings.Fields(text)
		//fmt.Printf("%v",ar)
	}

	//fmt.Printf("%s",elem.Text())
	//elem.Clear()
	//elem.SendKeys(code)
	//
	//// Click the run button
	//btn, _ := wd.FindElement(selenium.ByCSSSelector, "#run")
	//btn.Click()
	//
	//// Get the result
	//div, _ := wd.FindElement(selenium.ByCSSSelector, "#output")
	//
	//output := ""
	//// Wait for run to finish
	//for {
	//	output, _ = div.Text()
	//	if output != "Waiting for remote server..." {
	//		break
	//	}
	//	time.Sleep(time.Millisecond * 100)
	//}
	//
	//fmt.Printf("Got: %s\n", output)
}
