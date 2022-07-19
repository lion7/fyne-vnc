module github.com/lion7/fyne-vnc

go 1.13

require (
	fyne.io/fyne/v2 v2.2.3
	github.com/amitbet/vnc2video v0.0.0-20190616012314-9d50b9dab1d9
	github.com/icza/mjpeg v0.0.0-20210726201846-5ff75d3c479f // indirect
)

replace github.com/amitbet/vnc2video => github.com/deluan/vnc2video v0.0.0-20210101045232-81fb4f50aef5
