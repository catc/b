module github.com/catc/b

go 1.15

require (
	github.com/AlecAivazis/survey/v2 v2.1.1
	github.com/eiannone/keyboard v0.0.0-20200508000154-caf4b762e807
	github.com/mgutz/ansi v0.0.0-20170206155736-9520e82c474b
	github.com/stretchr/testify v1.2.2 // indirect
	golang.org/x/sys v0.0.0-20210113181707-4bcb84eeeb78 // indirect
)

replace github.com/AlecAivazis/survey/v2 v2.1.1 => github.com/catc/survey/v2 v2.1.2-0.20201010180609-6f3d9b653792
