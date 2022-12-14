package photo

import (
  "fmt"
  "time"

  "github.com/dsoprea/go-exif/v3"
  "github.com/dsoprea/go-jpeg-image-structure/v2"
  "github.com/dsoprea/go-png-image-structure/v2"
  "github.com/dsoprea/go-tiff-image-structure/v2"
  "github.com/dsoprea/go-heic-exif-extractor/v2"
  // log "github.com/sirupsen/logrus"
)

type MediaContext interface {
	Exif() (rootIfd *exif.Ifd, data []byte, err error)
}


func (fi *Image) SetExifInfo() error {
  if fi.filetype == "" {
    fi.SetMagicInfo()
  }

  var mediaContext MediaContext
  var err error
  switch fi.filetype {
  case "image/jpeg":
    jmp := jpegstructure.NewJpegMediaParser()
    mediaContext, err = jmp.ParseBytes(fi.fileBytes)
  case "image/png", "image/vnd.mozilla.apng":
    png := pngstructure.NewPngMediaParser()
    mediaContext, err = png.ParseBytes(fi.fileBytes)
  case "image/tiff":
    tiff := tiffstructure.NewTiffMediaParser()
    mediaContext, err = tiff.ParseBytes(fi.fileBytes)
  case "image/heic", "image/heic-sequence":
    heic := heicexif.NewHeicExifMediaParser()
    mediaContext, err = heic.ParseBytes(fi.fileBytes)
  default:
    err = fmt.Errorf("Unknown filetype %v", fi.filetype)
  }
  if err != nil {
    return err
  }

  rootIfd, _, err := mediaContext.Exif()
  if err != nil {
    return fmt.Errorf("Couldn't get root exif %v", err)
  }


  var datetime string
  var offset string
  cb := func(ifd *exif.Ifd, ite *exif.IfdTagEntry) error {
    var err error
    // fmt.Printf("NAME=[%s] VALUE=[%v]\n", , )
    switch ite.TagName() {
    case "DateTime":
      datetime, err = ite.FormatFirst()
    case "OffsetTime":
      offset, err = ite.FormatFirst()
    }
    if err != nil {
      return err
    }

  	return nil
  }
  err = rootIfd.EnumerateTagsRecursively(cb)
  if err != nil {
    return fmt.Errorf("Couldn't get exifs %v", err)
  }

  t, err := time.Parse("2006:01:02 15:04:05 -07:00", datetime + " " + offset)
  if err != nil {
    return fmt.Errorf("Error parsing date %v", err)
  }
  fi.taken = t
  return nil
}
