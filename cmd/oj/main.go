package main

import "oj/pkg"

func main() {
	pkg.ExecBuild("./assets/success/build/main.exe")
	pkg.ExecBuild("./assets/timeout/build/main.exe")
	pkg.ExecBuild("./assets/output-limit-exceeded/build/main.exe")
}
