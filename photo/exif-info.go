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

func jpegExifDump(data []byte) ([]exif.ExifTag, error) {
  jmp := jpegstructure.NewJpegMediaParser()

  intfc, err := jmp.ParseBytes(data)
  if intfc == nil {
    return nil, fmt.Errorf("Error reading jpeg exif data %v", err)
  }

  sl := intfc.(*jpegstructure.SegmentList)
  _, _, et, err := sl.DumpExif()
  if err != nil {
    // if err == exif.ErrNoExif {
    //   log.Infof("No EXIF %v (%v)", filename, err)
    //   return
    // }
    return nil, fmt.Errorf("Error reading jpeg exif data %v", err)
  }

  return et, nil
}
func pngExifDump(data []byte) ([]exif.ExifTag, error) {
  png := pngstructure.NewPngMediaParser()

  intfc, err := png.ParseBytes(data)
  if intfc == nil {
    return nil, fmt.Errorf("Error reading png exif data %v", err)
  }

  sl := intfc.(*jpegstructure.SegmentList)
  _, _, et, err := sl.DumpExif()
  if err != nil {
    // if err == exif.ErrNoExif {
    //   log.Infof("No EXIF %v (%v)", filename, err)
    //   return
    // }
    return nil, fmt.Errorf("Error reading png exif data %v", err)
  }

  return et, nil
}
func tiffExifDump(data []byte) ([]exif.ExifTag, error) {
  tiff := tiffstructure.NewTiffMediaParser()

  intfc, err := tiff.ParseBytes(data)
  if intfc == nil {
    return nil, fmt.Errorf("Error reading tiff exif data %v", err)
  }

  sl := intfc.(*jpegstructure.SegmentList)
  _, _, et, err := sl.DumpExif()
  if err != nil {
    // if err == exif.ErrNoExif {
    //   log.Infof("No EXIF %v (%v)", filename, err)
    //   return
    // }
    return nil, fmt.Errorf("Error reading tiff exif data %v", err)
  }

  return et, nil
}
func heicExifDump(data []byte) ([]exif.ExifTag, error) {
  heic := heicexif.NewHeicExifMediaParser()

  intfc, err := heic.ParseBytes(data)
  if intfc == nil {
    return nil, fmt.Errorf("Error reading tiff exif data %v", err)
  }

  sl := intfc.(*jpegstructure.SegmentList)
  _, _, et, err := sl.DumpExif()
  if err != nil {
    // if err == exif.ErrNoExif {
    //   log.Infof("No EXIF %v (%v)", filename, err)
    //   return
    // }
    return nil, fmt.Errorf("Error reading tiff exif data %v", err)
  }

  return et, nil
}

func (fi *Image) SetExifInfo() error {
  if fi.filetype == "" {
    fi.SetMagicInfo()
  }

  var et []exif.ExifTag
  var err error
  switch fi.filetype {
  case "image/jpeg":
    et, err = jpegExifDump(fi.fileBytes)
  case "image/png", "image/vnd.mozilla.apng":
    et, err = pngExifDump(fi.fileBytes)
  case "image/tiff":
    et, err = tiffExifDump(fi.fileBytes)
  case "image/heic", "image/heic-sequence":
    et, err = heicExifDump(fi.fileBytes)
  default:
    err = fmt.Errorf("Unknown filetype %v", fi.filetype)
  }

  if err != nil {
    return err
  }

  var datetime string
  var offset string
  for _, tag := range et {
      switch tag.TagName {
      case "DateTime":
        datetime = tag.FormattedFirst
      case "OffsetTime":
        offset = tag.FormattedFirst
      }
  }
  t, err := time.Parse("2006:01:02 15:04:05 -07:00", datetime + " " + offset)
  if err != nil {
    return fmt.Errorf("Error parsing date %v", err)
  }
  fi.taken = t
  return nil
  // for i, tag := range et {
  //   // Since we dump the complete value, the thumbnails introduce
  //   // too much noise.
  //   if (tag.TagId == exif.ThumbnailOffsetTagId || tag.TagId == exif.ThumbnailSizeTagId) && tag.IfdPath == exif.ThumbnailFqIfdPath {
  //     continue
  //   }
  //
  //   if tag.TagName == "DateTime" || tag.TagName == "OffsetTime" {
  //     fmt.Printf("%2d: IFD-PATH=[%s] ID=(0x%04x) NAME=[%s] TYPE=(%d):[%s] VALUE=[%v]", i, tag.IfdPath, tag.TagId, tag.TagName, tag.TagTypeId, tag.TagTypeName, tag.FormattedFirst)
  //
  //     if tag.ChildIfdPath != "" {
  //       fmt.Printf(" CHILD-IFD-PATH=[%s]", tag.ChildIfdPath)
  //     }
  //
  //     fmt.Printf("\n")
  //   }
  // }
}
