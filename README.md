# Aphelios Weapon Detector

## Requirements

GoLang >= 1.21.0

GoCV: https://gocv.io/getting-started/

## Setup

```go 
pkg/scanner.go
func getMinimalRect() image.Rectangle {
	return image.Rect(
		1000,
		1305,
		1230,
		1390,
	)
}
```

You will need to update this function to match with your resolution and hud scaling.
You can take a look at the assets to see how the rectangle should be placed. All the pictures should be contained in the rectangle.

Have fun!