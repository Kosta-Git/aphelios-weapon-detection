package pkg

import (
	"context"
	"github.com/kbinani/screenshot"
	"gocv.io/x/gocv"
	"image"
	"log"
)

type weaponImage struct {
	WeaponName     string
	MainImage      gocv.Mat
	SecondaryImage gocv.Mat
	NextImage      gocv.Mat
}

func getMinimalRect() image.Rectangle {
	// TODO: Make this configurable based on resolution
	return image.Rect(
		1000,
		1305,
		1230,
		1390,
	)
}

func getScreen(displayIndex int) (gocv.Mat, error) {
	bounds := screenshot.GetDisplayBounds(displayIndex)
	minRect := getMinimalRect()
	bounds.Min = image.Point{X: minRect.Min.X, Y: minRect.Min.Y}
	bounds.Max = image.Point{X: minRect.Max.X, Y: minRect.Max.Y}
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		panic(err)
	}
	return gocv.ImageToMatRGB(img)
}

func isMatching(toScan gocv.Mat, template gocv.Mat) bool {
	result := gocv.NewMat()
	defer result.Close()
	gocv.MatchTemplate(toScan, template, &result, gocv.TmCcoeffNormed, gocv.NewMat())
	_, maxVal, _, _ := gocv.MinMaxLoc(result)
	return maxVal > 0.9
}

func loadImages(basePath string) []weaponImage {
	return []weaponImage{
		{
			Calibrum,
			gocv.IMRead(basePath+"primary/calibrum.png", gocv.IMReadColor),
			gocv.IMRead(basePath+"secondary/calibrum.png", gocv.IMReadColor),
			gocv.IMRead(basePath+"next/calibrum.png", gocv.IMReadColor),
		},
		{
			Severum,
			gocv.IMRead(basePath+"primary/severum.png", gocv.IMReadColor),
			gocv.IMRead(basePath+"secondary/severum.png", gocv.IMReadColor),
			gocv.IMRead(basePath+"next/severum.png", gocv.IMReadColor),
		},
		{
			Gravitum,
			gocv.IMRead(basePath+"primary/gravitum.png", gocv.IMReadColor),
			gocv.IMRead(basePath+"secondary/gravitum.png", gocv.IMReadColor),
			gocv.IMRead(basePath+"next/gravitum.png", gocv.IMReadColor),
		},
		{
			Infernum,
			gocv.IMRead(basePath+"primary/infernum.png", gocv.IMReadColor),
			gocv.IMRead(basePath+"secondary/infernum.png", gocv.IMReadColor),
			gocv.IMRead(basePath+"next/infernum.png", gocv.IMReadColor),
		},
		{
			Crescendum,
			gocv.IMRead(basePath+"primary/crescendum.png", gocv.IMReadColor),
			gocv.IMRead(basePath+"secondary/crescendum.png", gocv.IMReadColor),
			gocv.IMRead(basePath+"next/crescendum.png", gocv.IMReadColor),
		},
	}
}

func unloadImages(weapons []weaponImage) {
	for _, weapon := range weapons {
		err := weapon.MainImage.Close()
		if err != nil {
			log.Fatalf("Error closing weapon weapons: %v", err.Error())
		}
		err = weapon.SecondaryImage.Close()
		if err != nil {
			log.Fatalf("Error closing weapon weapons: %v", err.Error())
		}
		err = weapon.NextImage.Close()
		if err != nil {
			log.Fatalf("Error closing weapon weapons: %v", err.Error())
		}
	}
}

func WatchApheliosWeapons(basePath string, displayIndex int, ctx context.Context) <-chan Weapons {
	weaponChan := make(chan Weapons, 1)
	go func() {
		loadout := Weapons{
			First:  Calibrum,
			Second: Severum,
			Third:  Gravitum,
			Fourth: Infernum,
			Fifth:  Crescendum,
		}
		weaponChan <- loadout // This is to ensure that the first loadout is sent
		weapons := loadImages(basePath)
		defer unloadImages(weapons)
		for {
			select {
			case <-ctx.Done():
				close(weaponChan)
				return
			default:
				screen, err := getScreen(displayIndex)
				if err != nil {
					log.Fatalf("Error getting screen: %v", err.Error())
				}
				changed := false
				for _, weapon := range weapons {
					if isMatching(screen, weapon.MainImage) && loadout.First != weapon.WeaponName {
						loadout.First = weapon.WeaponName
						changed = true
					}
					if isMatching(screen, weapon.SecondaryImage) && loadout.Second != weapon.WeaponName {
						loadout.Second = weapon.WeaponName
						changed = true
					}
					if isMatching(screen, weapon.NextImage) && loadout.Third != weapon.WeaponName {
						loadout.Third = weapon.WeaponName
						changed = true
					}
				}
				if changed {
					weaponChan <- loadout
				}
			}
		}
	}()
	return weaponChan
}
