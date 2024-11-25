package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Label compatible with ToolbarItem interface to insert into Toolbar.
type ToolbarLabel struct {
	widget.Label
}

func NewToolbarLabel(text string) *ToolbarLabel {
	var l = &ToolbarLabel{
		Label: widget.Label{
			Text:      text,
			Alignment: fyne.TextAlignLeading,
			TextStyle: fyne.TextStyle{},
		},
	}
	l.ExtendBaseWidget(l)
	return l
}

func (tl *ToolbarLabel) ToolbarObject() fyne.CanvasObject {
	tl.Label.Importance = widget.MediumImportance
	return tl
}

// Layout that fits the images to whole space and cuts edges if it needs.
type ImageFitLayout struct {
}

// Declare conformity with Layout interface
var _ fyne.Layout = (*ImageFitLayout)(nil)

func NewImageFit(objects ...fyne.CanvasObject) *fyne.Container {
	return container.New(ImageFitLayout{}, objects...)
}

func (l ImageFitLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	var ratiofit = size.Width / size.Height
	for _, child := range objects {
		var newsize = size
		var pos = fyne.NewPos(0, 0)
		if img, ok := child.(*canvas.Image); ok {
			var ratioimg = img.Aspect()
			if ratiofit > ratioimg {
				newsize.Height = size.Width / ratioimg
				pos.Y = (size.Height - newsize.Height) / 2
			} else {
				newsize.Width = size.Height * ratioimg
				pos.Y = (size.Width - newsize.Width) / 2
			}
		}
		child.Resize(newsize)
		child.Move(pos)
	}
}

func (l ImageFitLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	var minSize = fyne.NewSize(0, 0)
	for _, child := range objects {
		if !child.Visible() {
			continue
		}
		minSize = minSize.Max(child.MinSize())
	}
	return minSize
}
