package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
)

func main() {

	runCmd()
	fmt.Println("hello go<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
	runCmdWithPipe()
}

func runCmdWithPipe() {
	fmt.Println("Run command `ps aux | grep apipe`")
	cmd1 := exec.Command("ps", "aux")
	cmd2 := exec.Command("grep", "apipe")
	var outPutBuf1 bytes.Buffer
	cmd1.Stdout = &outPutBuf1

	if err := cmd1.Start(); err != nil {
		fmt.Printf("Error: The first cammand can not be startup %s\n", err)
		return
	}

	if err := cmd1.Wait(); err != nil {
		fmt.Printf("Error: Couldn't wait for the first command: %s\n", err)
		return
	}
	//将命令1的输出作为命令2的输入,两条命令之间是阻塞的
	cmd2.Stdin = &outPutBuf1
	var outputBuf2 bytes.Buffer
	cmd2.Stdout = &outputBuf2

	if err := cmd2.Start(); err != nil {
		fmt.Printf("Error: The second command can not be startup: %s\n", err)
		return
	}

	if err := cmd2.Wait(); err != nil {
		fmt.Printf("Error: Couldn't wait for the second command: %s\n", err)
		return
	}
	fmt.Printf("%s\n", outputBuf2.Bytes())
}

func runCmd() {
	useBufferedIO := false
	fmt.Println("Run command `echo -n \" My first command comes from golang.\"`:")
	cmd0 := exec.Command("echo", "-n", "My first command comes from golang.")
	stdout0, err := cmd0.StdoutPipe()
	if err != nil {
		fmt.Printf("Error: Couldn't obtain the stdout pipe for command Mo.0: %s\n", err)
		return
	}

	if err := cmd0.Start(); err != nil {
		fmt.Printf("Error: The command Mo.0 can not be startup: %s\n", err)
		return
	}

	if !useBufferedIO {
		//通过循环读取的方式将命令输出到一个缓存变量中
		var outputBuf0 bytes.Buffer

		i := 1
		for {
			tempOutput := make([]byte, 5)
			n, err := stdout0.Read(tempOutput)
			if err != nil {
				if err == io.EOF {
					break
				} else {
					fmt.Printf("Error: Couldn't read data from the pipe: %s\n", err)
					return
				}
			}

			if n > 0 {
				fmt.Printf("write in buffer %d \n", i)
				i++
				outputBuf0.Write(tempOutput[:n])
			}
		}
		fmt.Printf("%s\n", outputBuf0.String())

	} else {
		//使用带缓冲的读取器
		outputBuf0 := bufio.NewReader(stdout0)
		output0, _, err := outputBuf0.ReadLine()
		if err != nil {
			fmt.Printf("Error: Couldn't read data from the pipe: %s\n", err)
			return
		}

		fmt.Printf("%s\n", string(output0))
	}

}