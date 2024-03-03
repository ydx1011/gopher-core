package appcontext

import (
	"bytes"
	"fmt"
	"github.com/xfali/xlog"
	"io"
	"os"
)

const (
	GopheBanner = `
  ________              .__                  
 /  _____/  ____ ______ |  |__   ___________ 
/   \  ___ /  _ \\____ \|  |  \_/ __ \_  __ \
\    \_\  (  <_> )  |_> >   Y  \  ___/|  | \/
 \______  /\____/|   __/|___|  /\___  >__|   
        \/       |__|        \/     \/ 
`
	GopheBannerVersion = `================================`
)

func printGopherInfo(version string, bannerPath string, banner bool) {
	w := selectWriter()
	buf := bytes.NewBuffer(nil)
	buf.Grow(len(GopheBanner) + len(GopheBannerVersion) + 2)
	buf.WriteByte('\n')
	if banner {
		buf.WriteString(bannerString(bannerPath))
	}
	buf.WriteString(versionString(version, banner))
	buf.WriteByte('\n')

	w.Write(buf.Bytes())
}

func selectWriter() io.Writer {
	for i := xlog.INFO; i <= xlog.DEBUG; i++ {
		w := xlog.GetOutputBySeverity(i)
		if w != nil {
			return w
		}
	}
	return os.Stdout
}

func bannerString(bannerPath string) string {
	output := []byte(GopheBanner)
	if bannerPath != "" {
		f, err := os.Open(bannerPath)
		if err == nil {
			buf := bytes.NewBuffer(nil)
			_, err := io.Copy(buf, f)
			if err == nil {
				if buf.Bytes()[buf.Len()-1] != '\n' {
					buf.WriteByte('\n')
				}
				output = buf.Bytes()
			}
		}
	}
	return string(output)
}

func versionString(version string, banner bool) string {
	if banner {
		size := len(version)
		bs := len(GopheBannerVersion)
		if size == 0 || size > bs-3 {
			return GopheBannerVersion
		}
		return fmt.Sprintf("%s (%s)\n", GopheBannerVersion[:bs-size-3], version)
	} else {
		return fmt.Sprintf("=== neve === (%s)\n", version)
	}
}
